# GitHub Analyzer

Analyze GitHub repositories and produce reports.

## Build and Install

> **Note** Make sure `go 1.20` or above is installed before `make`.

To build and install `gha` to `~/bin` on Linux, simply run

```bash
make install
```

## Tutorial

Analyzing a GitHub repository requires two steps:

1. `gha snapshot` to fetch raw information from GitHub API
   - Personal Access Token (PAT) is required to be set to the environment variable `GITHUB_TOKEN` if throttled
2. `gha report` or `gha pr-review` to generate a markdown report from raw information fetched above.

### Examples

Take snapshot:
```console
$ gha snapshot --pr-reviews --pr-reviews-ago 365 notaryproject/notation
........
Fetched 714 issues and pull requests
Saved snapshot to notaryproject_notation_20230719_234453_snapshot.json
Fetching reviews of 291 pull requests since 2022-07-19...
..................................................  17.18%
..................................................  34.36%
..................................................  51.54%
..................................................  68.72%
..................................................  85.91%
.........................................          100.00%
Saved pull request reviews to notaryproject_notation_20230719_234809_reviews.json
```

Analyze snapshot:
```console
$ gha report --ago 90 notaryproject_notation_20230719_234453_snapshot.json
GitHub Analysis Report
======================
- Start Date: `2023-04-20 16:06:32`

## notaryproject_notation_20230719_234453_snapshot.json
Issues
- Total: 43
  - Open: 24
  - Closed: 19
- Time to close:
  - Min: 43s
  - Max: 1mo 12d
  - Mean: 6d 8h
  - Median: 3d 6h
  - 90th percentile: 10d 10h
  - 95th percentile: 12d 19h
  - 99th percentile: 12d 19h

Pull Requests
- Total: 64
  - Open: 12
  - Closed: 13
  - Merged: 39
- Time to merge:
  - Min: 16m 50s
  - Max: 25d 21h
  - Mean: 3d 22h
  - Median: 2d 18h
  - 90th percentile: 8d 4h
  - 95th percentile: 9d 17h
  - 99th percentile: 15d 15h
$ gha pr-review --ago 90 notaryproject_notation_20230719_234809_reviews.json
Pull Request Review Count
==========================
- Start Date: `2023-04-20 16:06:57`

## notaryproject_notation_20230719_234809_reviews.json

| Reviewer        | Count |                                                      |
|-----------------|-------|------------------------------------------------------|
| priteshbandi    |    50 | `                                                  ` |
| shizhMSFT       |    46 | `                                              `     |
| JeyJeyGao       |    30 | `                              `                     |
| Two-Hearts      |    29 | `                             `                      |
| FeynmanZhou     |     9 | `         `                                          |
| yizha1          |     8 | `        `                                           |
| patrickzheng200 |     6 | `      `                                             |
| rgnote          |     5 | `     `                                              |
| gokarnm         |     4 | `    `                                               |
| Wwwsylvia       |     3 | `   `                                                |
| sajayantony     |     3 | `   `                                                |
| zr-msft         |     2 | `  `                                                 |
| wangxiaoxuan273 |     1 | ` `                                                  |
| toddysm         |     1 | ` `                                                  |
| vaninrao10      |     1 | ` `                                                  |
| tungbq          |     1 | ` `                                                  |
| duffney         |     1 | ` `                                                  |
| ningziwen       |     1 | ` `                                                  |
| qweeah          |     1 | ` `                                                  |
| byronchien      |     1 | ` `                                                  |
```
