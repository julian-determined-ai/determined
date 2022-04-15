#!/usr/bin/env bash

source /run/determined/task-logging-setup.sh
trap 'source /run/determined/task-logging-teardown.sh' EXIT

set -e

if [ "$#" -eq 1 ];
then
    /bin/sh -c "$@"
else
    "$@"
fi
