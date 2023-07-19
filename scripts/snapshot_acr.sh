#!/bin/bash

repos=(
    Azure/acr-builder
    Azure/acr-cli
    Azure/acr
    Azure/acr-task-commands
)

$(dirname $0)/snapshot.sh ${repos[@]}
