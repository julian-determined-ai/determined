#####################
 Deploy on Slurm/PBS
#####################

+----------------------+
| Supported Versions   |
+======================+
| Slurm >= 19.05 or    |
| PBS >= 2021.1.2      |
+----------------------+
| Singularity >= 3.7   |
| or PodMan >= 3.3.1   |
+----------------------+
| Launcher             |
| (`hpe-hpc-launcher`) |
| >= 3.1.2             |
+----------------------+
| Java >= 1.8          |
+----------------------+

.. note::

   Slurm/PBS deployment applies to the Enterprise Edition.

This document describes how Determined can be configured to utilize HPC cluster scheduling systems
via the Determined HPC launcher. In this type of configuration, Determined delegates all job
scheduling and prioritization to the HPC workload manager (either Slurm or PBS). This integration
enables existing HPC workloads and Determined workloads to coexist and Determined workloads to
access all of the advanced capabilities of the HPC workload manager.

To install Determined on the HPC cluster, ensure that the
:doc:`/cluster-setup-guide/deploy-cluster/sysadmin-deploy-on-slurm/slurm-requirements` are met, then
follow the steps in the
:doc:`/cluster-setup-guide/deploy-cluster/sysadmin-deploy-on-slurm/install-on-slurm` document.

***********
 Reference
***********

-  :ref:`Determined Installation Requirements <system-requirements>`
-  `Slurm <https://slurm.schedmd.com/documentation.html>`__
-  `OpenPBS® <https://www.openpbs.org/>`__
-  `PBS Professional® <https://www.altair.com/pbs-professional/>`__
-  `Singularity <https://docs.sylabs.io/guides/3.7/user-guide/introduction.html>`__
-  `Apptainer <https://apptainer.org/>`__
-  `PodMan <https://docs.podman.io>`__

.. toctree::
   :maxdepth: 1
   :hidden:

   slurm-requirements
   install-on-slurm
   singularity
   slurm-known-issues
