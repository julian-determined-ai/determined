package model

import (
	"time"

	"github.com/uptrace/bun"

	"github.com/determined-ai/determined/proto/pkg/userv1"
	"github.com/determined-ai/determined/proto/pkg/workspacev1"
)

// Workspace is the bun model of a workspace.
type Workspace struct {
	bun.BaseModel `bun:"table:workspaces"`
	ID            int             `bun:"id,pk,autoincrement"`
	Name          string          `bun:"name"`
	Archived      bool            `bun:"archived"`
	CreatedAt     time.Time       `bun:"created_at,scanonly"`
	UserID        UserID          `bun:"user_id"`
	Immutable     bool            `bun:"immutable"`
	State         *WorkspaceState `bun:"state"`
	AgentUID      *int32          `bun:"uid"`
	AgentUser     *string         `bun:"user_"`
	AgentGID      *int32          `bun:"gid"`
	AgentGroup    *string         `bun:"group_"`
}

// ToProto converts a bun model of a workspace to a proto object.
// Some fields like username and pinned are not included since they are
// not on the bun model.
func (w *Workspace) ToProto() *workspacev1.Workspace {
	var aug *userv1.AgentUserGroup

	if w.AgentUID != nil || w.AgentGID != nil || w.AgentUser != nil || w.AgentGroup != nil {
		aug = &userv1.AgentUserGroup{
			AgentUid:   w.AgentUID,
			AgentGid:   w.AgentGID,
			AgentUser:  w.AgentUser,
			AgentGroup: w.AgentGroup,
		}
	}

	return &workspacev1.Workspace{
		Id:             int32(w.ID),
		Name:           w.Name,
		Archived:       w.Archived,
		UserId:         int32(w.UserID),
		Immutable:      w.Immutable,
		State:          w.State.ToProto(),
		AgentUserGroup: aug,
	}
}

// WorkspaceState is the state of the workspace state with regards to being deleted.
type WorkspaceState string

const (
	// WorkspaceStateDeleting constant.
	WorkspaceStateDeleting WorkspaceState = "DELETING"
	// WorkspaceStateDeleteFailed constant.
	WorkspaceStateDeleteFailed WorkspaceState = "DELETE_FAILED"
	// WorkspaceStateDeleted constant.
	WorkspaceStateDeleted WorkspaceState = "DELETED"
)

// ToProto converts a WorkspaceState to a proto workspacev1.Workspace state.
func (s *WorkspaceState) ToProto() workspacev1.WorkspaceState {
	if s == nil {
		return workspacev1.WorkspaceState_WORKSPACE_STATE_UNSPECIFIED
	}
	return workspacev1.WorkspaceState(workspacev1.WorkspaceState_value["WORKSPACE_STATE_"+string(*s)])
}

// WorkspacePin is the bun model of a workspace.
type WorkspacePin struct {
	bun.BaseModel `bun:"table:workspace_pins"`
	WorkspaceID   int    `bun:"workspace_id"`
	UserID        UserID `bun:"user_id"`
}
