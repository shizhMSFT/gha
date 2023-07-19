#!/bin/bash

repos=(
   deislabs/ratify
)

$(dirname $0)/snapshot.sh ${repos[@]}
