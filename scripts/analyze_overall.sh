#!/bin/bash

gha report --ago 365 *_snapshot.json > overall_report.md
gha pr-review --ago 365 *_reviews.json > overall_review.md
