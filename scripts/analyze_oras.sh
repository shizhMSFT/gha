#!/bin/bash

gha report --ago 365 oras-project_*_snapshot.json > oras-project_report.md
gha pr-review --ago 365 oras-project_*_reviews.json > oras-project_review.md
