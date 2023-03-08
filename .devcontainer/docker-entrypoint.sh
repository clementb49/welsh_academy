#!/bin/bash

echo "Launching environment with UID: ${USER_ID} and GID: ${GROUP_ID} "

id "vscode" >/dev/null 2>&1

if [[ $? -ne 0 ]]; then
    echo "Creating group vscode with GID: ${GROUP_ID}"
    groupadd --gid "${GROUP_ID}" vscode

    echo "Creating user vscode with UID: ${USER_ID} and GID: ${GROUP_ID}"
    useradd -u ${USER_ID} -g ${GROUP_ID} -m vscode
fi
exec gosu vscode "$@"