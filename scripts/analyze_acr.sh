#!/bin/bash

gha report --ago 365 Azure_acr*_snapshot.json > acr_report.md
gha pr-review --ago 365 Azure_acr*_reviews.json > acr_review.md
