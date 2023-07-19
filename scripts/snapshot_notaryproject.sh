#!/bin/bash

repos=(
    notaryproject/notation
    notaryproject/notation-go
    notaryproject/notation-core-go
    notaryproject/notation-action
    notaryproject/notaryproject.dev
    Azure/notation-azure-kv
)

$(dirname $0)/snapshot.sh ${repos[@]}
