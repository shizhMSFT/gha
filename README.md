# GitHub Analyzer

Analyze GitHub repositories and produce reports.

## Build and Install

> **Note**
> Make sure `go 1.21.0` or above is installed before `make`.

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

#### Take Snapshot

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

#### Analyze Snapshot

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

#### Analyze Issue Comments

```console
$ gha snapshot --issue-comments --issue-comments-since 2023-01-01 notaryproject/notation
........
Fetched 734 issues and pull requests
Saved snapshot to notaryproject_notation_20230828_093829_snapshot.json
Fetching comments of 264 issues since 2023-01-01...
..................................................  18.93%
..................................................  37.87%
..................................................  56.81%
..................................................  75.75%
..................................................  94.69%
..............                                     100.00%
Saved issue comments to notaryproject_notation_20230828_094020_comments.json
$ # Download CODEOWNERS or MAINTAINERS. Both work.
$ wget https://raw.githubusercontent.com/notaryproject/notation/main/MAINTAINERS
--2023-08-28 23:10:54--  https://raw.githubusercontent.com/notaryproject/notation/main/MAINTAINERS
Resolving raw.githubusercontent.com (raw.githubusercontent.com)... 185.199.109.133, 185.199.110.133, 185.199.111.133, ...
Connecting to raw.githubusercontent.com (raw.githubusercontent.com)|185.199.109.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 682 [text/plain]
Saving to: ‘MAINTAINERS’

MAINTAINERS                   100%[=================================================>]     682  --.-KB/s    in 0s

2023-08-28 23:10:55 (15.4 MB/s) - ‘MAINTAINERS’ saved [682/682]
$ gha issue-comment --maintainers MAINTAINERS --start-date 2023-01-01 --sla 14 notaryproject_notation_20230828_093829_snapshot.json notaryproject_notation_20230828_094020_comments.json
Issue Comment Summary
=====================
- Start Date: `2023-01-01 00:00:00`

## Maintainers

- @justincormack
- @niazfk
- @stevelasker
- @JeyJeyGao
- @gokarnm
- @Two-Hearts
- @priteshbandi
- @rgnote
- @shizhMSFT

## First Response Time

- Non-maintainer issues: 79
  - Responded: 26
    - Min: 2m 8s
    - Max: 5mo 26d
    - Mean: 21d 20h
    - Median: 8d 13h
    - 90th percentile: 1mo 10d
    - 95th percentile: 1mo 17d
    - 99th percentile: 3mo 5d
  - No Response: 53

### Out of SLA: 14 Days

| #Issue | Duration | Title                                                                                     |
|--------|----------|-------------------------------------------------------------------------------------------|
| #506   | 7mo 17d  | doc: notation Inspect Command line Spec - Phase 2                                         |
| #508   | 7mo 16d  | CLI switch to store signatures using OCI image manifest.                                  |
| #513   | 7mo 5d   | Notation Verify should warnings output to Stderr                                          |
| #539   | 6mo 22d  | Signing with local private keys                                                           |
| #545   | 6mo 19d  | Add a helper function in ioutil to PrintObjectAsTree                                      |
| #548   | 6mo 19d  | CLI Cmds for trust policy management - phase 1                                            |
| #570   | 6mo 2d   | Add Notation CLI to Winget package manager                                                |
| #574   | 6mo 1d   | Change the default signature manifest                                                     |
| #575   | 6mo      | Verify referrers API when pushing image manifest                                          |
| #578   | 5mo 28d  | Documenting the security process for Notation                                             |
| #571   | 5mo 26d  | Create a Homebrew formula for Notation CLI                                                |
| #584   | 5mo 24d  | Add an example to CLI help info for notation signing                                      |
| #586   | 5mo 24d  | Update references from Notary v2 to Notation                                              |
| #587   | 5mo 18d  | Login without credential helper                                                           |
| #590   | 5mo 16d  | Discussion of out of box experience for trust policy                                      |
| #596   | 5mo 8d   | add labels for features subject to changes with proper doc                                |
| #597   | 5mo 7d   | Sign/verify OCI image layout                                                              |
| #599   | 5mo 7d   | Print manifests as part of the `--debug` option                                           |
| #600   | 5mo 6d   | Simplify Docker Credential Helper configuration for Notation authentication               |
| #604   | 5mo 1d   | Notation delete returns confusing message                                                 |
| #605   | 5mo 1d   | Fix the text for notation version                                                         |
| #609   | 4mo 28d  | [Usability Issue] Cert list is not helpful and just lists the files                       |
| #610   | 4mo 27d  | [Usability Issue] `notation inspect` help is missing                                      |
| #614   | 4mo 26d  | Support experimental feature                                                              |
| #618   | 4mo 25d  | `notation sign` error messages are not helpful to understand what parameter is missing    |
| #620   | 4mo 23d  | Improve the messages for `notation verify`                                                |
| #622   | 4mo 22d  | `notation cert delete` confirms deletion without doing anything                           |
| #624   | 4mo 20d  | Improve the output message of `notation inspect` images without signatures                |
| #625   | 4mo 20d  | Error message too general for `notation verify` command                                   |
| #628   | 4mo 19d  | Notation CLI guideline and CLI spec template                                              |
| #630   | 4mo 17d  | Introduce an experimental flag to enable backward compatibility with OCI registries       |
| #631   | 4mo 16d  | Support adding public key to trust store by specifying URL                                |
| #633   | 4mo 13d  | Missing e2e test cases for flag `--plain-http`                                            |
| #635   | 4mo 13d  | Use SHA2 in output of notation inspect                                                    |
| #637   | 4mo 10d  | [Usability issue] Notation login error message is confusing                               |
| #638   | 4mo 10d  | Add E2E test cases for validating certificate revocation with OCSP                        |
| #640   | 4mo 9d   | Release Notation CLI v1.0.0-rc.4                                                          |
| #642   | 4mo 8d   | Decide on main commit for a release: 6cd6555 and PR bump up versions                      |
| #644   | 4mo 7d   | Improve the output for notation plugin                                                    |
| #645   | 4mo 5d   | Examples were shown for experimental feature                                              |
| #646   | 4mo 5d   | Missing annotations in the output of `notation inspect`                                   |
| #652   | 4mo 4d   | Requesting UX improvement in signing and verifying with user metadata via Notation CLI    |
| #655   | 4mo 1d   | Image Verification for containerd                                                         |
| #662   | 3mo 20d  | Trace the execution of executables                                                        |
| #667   | 3mo 15d  | Low code coverage (33%) reported for notation main branch                                 |
| #681   | 3mo 6d   | docs: `notation login` error message improvement                                          |
| #549   | 3mo 5d   | Improved Plugin installation UX - phase 1                                                 |
| #695   | 3mo      | feat: Print out the signature digest when sign an artifact                                |
| #697   | 2mo 29d  | `notation login` fails to detect existing credentials for `docker.io`                     |
| #704   | 2mo 23d  | Improve error output for notation plugin                                                  |
| #705   | 2mo 23d  | Use existing credentials to auth to remote registries                                     |
| #706   | 2mo 23d  | Check the license header for Notation and its dependencies                                |
| #715   | 2mo 12d  | Update the README for the repository                                                      |
| #718   | 2mo 8d   | Add Golang lint to GitHub Actions for static Go code formatting scanning                  |
| #728   | 2mo 1d   | Add --force to notation key add                                                           |
| #621   | 1mo 17d  | Improve the error for missing trust policy                                                |
| #598   | 1mo 10d  | Add ability to redirect --debug logs to file                                              |
| #634   | 1mo 6d   | Standardize symlink checking per trust store spec                                         |
| #623   | 28d 22h  | Flag `--plain-http` didn't explicitly remind users the insecure connection to registries  |
| #647   | 21d 21h  | Support clean up the source key and certificate generated by Notation                     |
| #653   | 21d 1h   | `notation policy init` command is necessary for user experiences                          |
| #503   | 20d 21h  | Improve Notation authentication experience                                                |
| #759   | 18d 1h   | Add support for multiple trust policies                                                   |
| #721   | 17d 19h  | Read certificate from windows certificate store                                           |
| #519   | 15d 6h   | Update the branch policies for the repository                                             |
```