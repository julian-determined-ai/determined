#!/usr/bin/env bash

if [ -n "$DET_K8S_LOG_TO_FILE" ]; then
    # Replace overriden stdout and stderr with original and close them, since the command is finished.
    exec >&1- >&2- 1>&3 2>&4

    for (( i=0;i<$LOG_WAIT_COUNT;i++ )); do
        # read returns 1 on EOF, but it's a fifo so that is OK.
        read <$LOG_WAIT_FIFO || true
    done
fi