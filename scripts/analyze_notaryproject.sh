#!/bin/bash

gha report --ago 365 notaryproject_*_snapshot.json > notaryproject_report.md
gha pr-review --ago 365 notaryproject_*_reviews.json > notaryproject_review.md
