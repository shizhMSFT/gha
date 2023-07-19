#!/bin/bash

gha report --ago 365 deislabs_ratify_*_snapshot.json > ratify_report.md
gha pr-review --ago 365 deislabs_ratify_*_reviews.json > ratify_review.md
