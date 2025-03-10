package tasks

import (
	"archive/tar"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types/mount"

	"github.com/determined-ai/determined/master/pkg/archive"
	"github.com/determined-ai/determined/master/pkg/cproto"
	"github.com/determined-ai/determined/master/pkg/etc"
	"github.com/determined-ai/determined/master/pkg/model"
	"github.com/determined-ai/determined/master/pkg/schemas/expconf"
	"github.com/determined-ai/determined/master/pkg/ssh"
)

// TrialSpec is a description of a task for running a trial container.
type TrialSpec struct {
	Base TaskSpec

	ExperimentID     int
	TrialID          int
	TrialRunID       int
	ExperimentConfig expconf.ExperimentConfig
	HParams          map[string]interface{}
	TrialSeed        uint32
	LatestCheckpoint *model.Checkpoint
	StepsCompleted   int
}

// ToTaskSpec generates a TaskSpec.
func (s TrialSpec) ToTaskSpec(keys *ssh.PrivateAndPublicKeys) TaskSpec {
	res := s.Base

	env := s.ExperimentConfig.Environment()
	ports := env.Ports()
	if ports == nil {
		ports = make(map[string]int)
	}
	// TODO: remove this, but without breaking rendezvous api.
	ports["trial"] = 1734
	env.SetPorts(ports)
	res.Environment = env

	res.ResourcesConfig = s.ExperimentConfig.Resources()
	res.SlurmConfig = s.ExperimentConfig.SlurmConfig()
	res.PbsConfig = s.ExperimentConfig.PbsConfig()

	res.WorkDir = DefaultWorkDir

	additionalFiles := archive.Archive{
		s.Base.AgentUserGroup.OwnedArchiveItem(
			trialEntrypointFile,
			etc.MustStaticFile(etc.TrialEntrypointScriptResource),
			trialEntrypointMode,
			tar.TypeReg,
		),
	}

	additionalSSHFiles := archive.Archive{
		s.Base.AgentUserGroup.OwnedArchiveItem(sshDir, nil, sshDirMode, tar.TypeDir),
		s.Base.AgentUserGroup.OwnedArchiveItem(trialAuthorizedKeysFile,
			keys.PublicKey,
			trialAuthorizedKeysMode,
			tar.TypeReg,
		),
		s.Base.AgentUserGroup.OwnedArchiveItem(
			pubKeyFile, keys.PublicKey, pubKeyMode, tar.TypeReg,
		),
		s.Base.AgentUserGroup.OwnedArchiveItem(
			privKeyFile, keys.PrivateKey, privKeyMode, tar.TypeReg,
		),
		s.Base.AgentUserGroup.OwnedArchiveItem(sshdConfigFile,
			etc.MustStaticFile(etc.SSHDConfigResource),
			sshdConfigMode,
			tar.TypeReg,
		),
	}

	additionalRootFiles := archive.Archive{
		archive.RootItem(
			trialSSHConfigFile,
			etc.MustStaticFile(etc.SSHConfigResource),
			trialSSHConfigMode,
			tar.TypeReg,
		),
	}

	res.ExtraArchives = []cproto.RunArchive{
		wrapArchive(
			archive.Archive{
				s.Base.AgentUserGroup.OwnedArchiveItem(trainDir, nil, 0o700, tar.TypeDir),
				s.Base.AgentUserGroup.OwnedArchiveItem(modelCopy, nil, 0o700, tar.TypeDir),
			},
			rootDir,
		),
		wrapArchive(additionalFiles, rootDir),
		wrapArchive(additionalRootFiles, rootDir),
		wrapArchive(additionalSSHFiles, rootDir),
	}

	res.Description = fmt.Sprintf(
		"exp-%d-trial-%d",
		s.ExperimentID,
		s.TrialID,
	)

	res.Entrypoint = []string{"/run/determined/train/entrypoint.sh"}

	envVars := map[string]string{
		"DET_EXPERIMENT_ID":     strconv.Itoa(s.ExperimentID),
		"DET_TRIAL_ID":          strconv.Itoa(s.TrialID),
		"DET_TRIAL_RUN_ID":      strconv.Itoa(s.TrialRunID),
		"DET_TRIAL_SEED":        strconv.FormatUint(uint64(s.TrialSeed), 10),
		"DET_EXPERIMENT_CONFIG": jsonify(s.ExperimentConfig),
		"DET_HPARAMS":           jsonify(s.HParams),
		"DET_STEPS_COMPLETED":   strconv.Itoa(s.StepsCompleted),
		"DET_TASK_TYPE":         string(model.TaskTypeTrial),
		// DET_UNIQUE_PORT_OFFSET will be overwritten by the agent resource manager when the
		// container is started, but only on agent resource pools.  This default value will apply
		// to k8s and slurm (though slurm will override it in a startup hook).
		"DET_UNIQUE_PORT_OFFSET": "0",
	}
	if s.LatestCheckpoint != nil && s.LatestCheckpoint.UUID != nil {
		envVars["DET_LATEST_CHECKPOINT"] = s.LatestCheckpoint.UUID.String()
	}

	res.ExtraEnvVars = envVars

	res.LoggingFields = map[string]string{
		"trial_id": strconv.Itoa(s.TrialID),
	}

	if shm := s.ExperimentConfig.Resources().ShmSize(); shm != nil {
		res.ShmSize = int64(*shm)
	}

	mounts := ToDockerMounts(s.ExperimentConfig.BindMounts(), res.WorkDir)
	addMount := func(source, target string, bindOpts *mount.BindOptions) {
		mounts = append(mounts, mount.Mount{
			Type: mount.TypeBind, Source: source, Target: target, BindOptions: bindOpts,
		})
	}
	if c := s.ExperimentConfig.CheckpointStorage().RawSharedFSConfig; c != nil {
		addMount(
			c.HostPath(),
			expconf.DefaultSharedFSContainerPath,
			&mount.BindOptions{Propagation: expconf.DefaultSharedFSPropagation},
		)
	}
	if c := s.ExperimentConfig.DataLayer().RawSharedFSConfig; c != nil {
		if c.HostStoragePath() != nil && c.ContainerStoragePath() != nil {
			addMount(*c.HostStoragePath(), *c.ContainerStoragePath(), nil)
		}
	}
	if c := s.ExperimentConfig.DataLayer().RawS3Config; c != nil {
		if c.LocalCacheHostPath() != nil && c.LocalCacheContainerPath() != nil {
			addMount(*c.LocalCacheHostPath(), *c.LocalCacheContainerPath(), nil)
		}
	}
	if c := s.ExperimentConfig.DataLayer().RawGCSConfig; c != nil {
		if c.LocalCacheHostPath() != nil && c.LocalCacheContainerPath() != nil {
			addMount(*c.LocalCacheHostPath(), *c.LocalCacheContainerPath(), nil)
		}
	}
	res.Mounts = mounts
	res.TaskType = model.TaskTypeTrial

	return res
}
