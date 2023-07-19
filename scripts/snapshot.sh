#!/bin/bash

repos=$@

for repo in $repos; do
    printf "\e[31m%s\e[0m %s\n" ">>>" $repo
    gha snapshot --pr-reviews --pr-reviews-ago 365 $repo
done
