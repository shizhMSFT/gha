#!/bin/bash

repos=(
    oras-project/oras
    oras-project/oras-go
    oras-project/oras-credentials-go
    oras-project/oras-www
)

$(dirname $0)/snapshot.sh ${repos[@]}
