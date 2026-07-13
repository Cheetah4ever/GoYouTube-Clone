# GoYouTube Agent Guide

This file records how GoYouTube was understood, planned, changed, tested, and
published. Use it as the working agreement for future changes to this repository
and as a starting template for similar small Go projects.

## Working principles

1. Inspect before editing. Understand the repository, entry point, data flow,
   generated artifacts, and Git state first.
2. Make one logical change per commit. Do not combine hygiene, behavior, tests,
   tooling, documentation, and cleanup into one large commit.
3. Stage explicit paths. Never use `git add -A` when unrelated or unreviewed
   files are present.
4. Run unit tests before every code commit. Generate coverage when test behavior
   changes.
5. Push each approved checkpoint to the feature branch and keep the draft pull
   request description current.
6. Preserve useful experiments, but label them clearly and keep them out of the
   primary application path.
7. Never commit private keys, generated media, compiled binaries, coverage
   artifacts, or machine-specific paths.

## Current project model

The primary application is an HTTP video server:

```text
main.go
  -> server.RunHTTP
      -> GET /healthz
      -> GET or HEAD /video
          -> http.ServeFile
              -> HTTP byte-range support
```

The active server is configured with `-addr` and `-video`. The TCP client,
custom TCP server, incomplete UDP code, and FFmpeg functions are preserved as
legacy experiments. They are not part of the default HTTP execution path.

Relevant files:

- `main.go`: command-line configuration and HTTP entry point.
- `server/http.go`: HTTP handler and server.
- `server/http_test.go`: in-package HTTP unit tests.
- `server/tcp.go` and `server/udp.go`: legacy networking experiments.
- `ffmpeg_handler/ffmpeg.go`: legacy transcoding experiments.
- `Makefile`: repeatable unit-test and coverage commands.
- `README.md`: user-facing setup and usage.

## Phase 1: understand the repository

Start with read-only discovery. The following commands were used for this
project:

```sh
pwd
rg --files
rg --files goyoutube
find . -maxdepth 3 -type f
file goyoutube client/*.mp4 client/*.raw
```

Read entry points and package boundaries rather than assuming behavior from
filenames:

```sh
sed -n '1,240p' main.go
sed -n '1,280p' server/tcp.go
sed -n '1,320p' server/udp.go
sed -n '1,320p' client/client.go
sed -n '1,520p' ffmpeg_handler/ffmpeg.go
sed -n '1,200p' go.mod
```

Search call sites and problematic patterns with `rg`:

```sh
rg '^func ' -g '*.go'
rg -n 'ffmpeg_handler\.|server\.' --glob '*.go'
rg -n '/Users/|z_client|sample_data' --glob '*.go'
rg -n -i 'chatgpt|func Test[0-9_]|func (Tcp|Udp)\b' --glob '*.go'
```

This investigation established that the original project was a local TCP video
and FFmpeg experiment, not a YouTube API client. It also revealed:

- one accepted TCP client and a custom `Video` command;
- complete transcoded videos buffered with `io.ReadAll`;
- partial client reads and numeric byte output;
- commented/incomplete UDP code;
- missing or absolute media paths;
- generated media, a binary, certificates, and a private key;
- several numbered or session-style function names;
- duplicate nested `go.sum` files.

## Phase 2: decide the product direction

For video files and video-on-demand, choose HTTP rather than raw UDP. HTTP gives
reliable, ordered delivery and standard range requests for seeking. Raw UDP is
appropriate only for a deliberate live-media protocol experiment where packet
loss, ordering, congestion control, timing, and jitter buffering are handled.

The chosen first milestone was deliberately small:

1. Serve one configured MP4 over HTTP.
2. Support `GET`, `HEAD`, and byte ranges through the standard library.
3. Add a health endpoint.
4. Test handlers in memory without starting `main`.
5. Add repeatable coverage reporting.
6. Preserve earlier experiments separately.

## Phase 3: inspect Git and GitHub before committing

Check repository state, identity, remotes, tools, and authentication:

```sh
git status -sb
git config --local --get user.name
git config --local --get user.email
git config --show-origin --get-regexp '^user\.'
git remote -v
command -v gh
type -a gh
gh --version
gh auth status
```

Project-specific commit identity should be local, not global:

```sh
git config --local user.name "YOUR_GITHUB_NAME"
git config --local user.email "YOUR_VERIFIED_GITHUB_EMAIL"
```

The local repository initially had no commits while GitHub already contained an
`Initial commit` with `LICENSE`. The remote history was inspected before making
changes:

```sh
git ls-remote --heads origin
git fetch origin main
git log --oneline --decorate -5 origin/main
git ls-tree -r --name-only origin/main
git switch -c agent/http-video-foundation origin/main
```

Never create an unrelated root commit or overwrite an existing remote branch.
Always fetch and branch from the real remote base.

On this machine, a Python executable named `gh` appeared earlier in `PATH` than
the GitHub CLI. `command -v gh` and `type -a gh` exposed the conflict. Where
needed, `/opt/homebrew/bin/gh` selected the actual GitHub CLI.

## Phase 4: plan commits before editing

The work was intentionally split into these commits:

| Commit | Purpose |
| --- | --- |
| `c9a5f26` | `chore: add repository hygiene` |
| `6233b08` | `feat: add initial TCP video prototype` |
| `4da215a` | `feat: add HTTP video server` |
| `2b746e5` | `test: cover HTTP video endpoints` |
| `c664455` | `refactor: clarify experiment function names` |
| `f4511b5` | `chore: add unit test coverage reports` |
| `f7f40d3` | `docs: add project setup and testing guide` |
| `5f0c1e2` | `chore: make legacy experiments portable` |

This ordering matters:

- Hygiene prevents secrets and generated files entering history.
- The baseline preserves where the project started.
- Behavior is reviewable independently from tests.
- Tests are reviewable independently from coverage tooling.
- Naming changes contain no intended behavior changes.
- Documentation describes commands that already exist and were verified.
- Portability cleanup happens after the working path is established.

For the next project, prepare the equivalent commit list before implementation.
Adjust the list to the project, but retain one concern per commit.

## Safe staging and commit workflow

Before staging:

```sh
git status -sb
git diff
git diff --check
```

Stage only the files for the current step:

```sh
git add path/to/file1 path/to/file2
git status --short
git diff --cached --stat
git diff --cached
git diff --cached --check
```

Commit and inspect:

```sh
git commit -m "type: concise description"
git show --stat --oneline --summary HEAD
git log --oneline --decorate -5
```

Recommended commit prefixes:

- `chore:` repository setup, tooling, and mechanical cleanup;
- `feat:` user-visible behavior;
- `test:` unit tests;
- `refactor:` code organization or naming without behavior changes;
- `docs:` documentation only;
- `fix:` correction of faulty behavior.

## Go formatting, modules, and compilation

Format every edited Go file:

```sh
gofmt -w main.go server/http.go server/http_test.go
```

When module declarations or dependencies change:

```sh
go mod tidy
```

Useful package checks performed during the initial baseline and HTTP work were:

```sh
go test ./...
go vet ./...
go build ./...
```

The ongoing project requirement is unit testing. Do not start real listeners,
invoke FFmpeg, or add integration/race tests unless the user expands the scope.

## Unit-test model

Go tests live beside the source in files ending with `_test.go`:

```text
server/http.go
server/http_test.go
```

`go test` creates a temporary test executable. It compiles `main.go` but does
not call the application's `main()` function. No manual toggling is required.

The HTTP tests use `package server`, `httptest.NewRequest`, and
`httptest.NewRecorder`. This permits focused handler tests without binding a
port. Test data is written under `t.TempDir()` and removed automatically.

Run unit tests:

```sh
make test
```

For verbose output:

```sh
go test -v ./...
```

Current HTTP unit-test concerns include:

- missing video configuration;
- health response;
- complete video response;
- partial byte-range response and `Content-Range`;
- `HEAD` response without a body;
- unknown route;
- unsupported method;
- reading the current file on request rather than caching it at startup.

## Coverage reports

Generate the unit-test coverage profile, function summary, and HTML report:

```sh
make coverage
```

This runs:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

`coverage.out` and `coverage.html` are local generated artifacts and must remain
ignored. Do not commit generated reports. Report the total and call out major
untested areas honestly; preserved legacy code lowers repository-wide coverage.

## Running the application

The normal application command is:

```sh
go run .
```

With explicit configuration:

```sh
go run . -video /path/to/video.mp4 -addr localhost:9000
```

During initial HTTP implementation, range behavior was manually verified with:

```sh
curl --fail --silent --show-error http://127.0.0.1:18080/healthz
curl --fail --silent --show-error --range 0-99 \
  --output /tmp/video-range.bin --dump-header /tmp/video-range.headers \
  http://127.0.0.1:18080/video
```

That was an implementation smoke check. The ongoing agreed validation scope is
unit tests only unless the user asks for end-to-end verification.

## GitHub publishing workflow

Push each checkpoint:

```sh
git push -u origin agent/http-video-foundation
```

If Git is not connected to the authenticated GitHub CLI credentials:

```sh
gh auth setup-git
```

Create a draft pull request after the first pushed commit:

```sh
gh pr create --draft \
  --base main \
  --head agent/http-video-foundation \
  --title "Build HTTP video server incrementally" \
  --body-file /tmp/pr-body.md
```

Update the PR description after later checkpoints:

```sh
gh pr edit 1 --body-file /tmp/pr-body.md
gh pr view 1 --json isDraft,state,url,title
```

The PR body should always state:

- what changed;
- why it changed;
- user/developer impact;
- validation performed;
- current limitations or uncovered code.

Never print or store authentication tokens in project files. `gh auth status`
is appropriate because it masks the token.

## Filesystem and search discipline

- Prefer `rg` and `rg --files` for searching.
- Use `sed -n` for focused reads instead of dumping every file.
- Use `file` and `ls -lh` to classify binaries and generated media.
- Check `git status --short --ignored` to verify ignore rules.
- Use `git check-ignore -v <path>` to identify the matching ignore rule.
- Preserve unrelated user files and edits.
- Avoid destructive Git commands and broad staging.
- Keep all source paths repository-relative; never commit a home-directory path.

Useful hygiene checks:

```sh
git check-ignore -v .DS_Store key.pem client/out.mp4 coverage.out
git status --short --ignored
rg -n '/Users/|z_client|sample_data' --glob '*.go'
cmp -s go.sum server/go.sum
```

## Security and generated artifacts

The following must stay ignored:

- `.DS_Store`;
- compiled Go binaries and `*.test` files;
- generated `client/out.mp4` and `client/out.raw`;
- coverage profiles and HTML reports;
- local output directories;
- certificates and private keys (`*.pem`).

The tracked sample MP4 is intentionally small and makes a fresh clone runnable.
Reconsider Git LFS or external fixtures if sample media becomes large.

## Reusable checklist for the next project

### Understand

- [ ] Inspect files, entry points, module declarations, and binary assets.
- [ ] Trace the current runtime flow and external commands.
- [ ] Identify generated files, secrets, hardcoded paths, and dead experiments.
- [ ] Run the existing unit tests before editing.

### Plan

- [ ] Define the smallest useful milestone.
- [ ] Write a commit-by-commit plan.
- [ ] Separate hygiene, baseline, features, tests, tooling, docs, and cleanup.
- [ ] Agree on unit, integration, and coverage scope with the user.

### Git and GitHub

- [ ] Confirm local commit name and verified email.
- [ ] Confirm the authenticated GitHub account.
- [ ] Inspect the remote and fetch its default branch.
- [ ] Create a feature branch from the real remote base.
- [ ] Open a draft PR and update it after each checkpoint.

### Implement

- [ ] Make one logical change.
- [ ] Format edited source.
- [ ] Stage explicit files only.
- [ ] Review the staged diff and whitespace checks.
- [ ] Run the agreed unit tests.
- [ ] Commit, push, report the hash, and pause at the checkpoint.

### Finish

- [ ] Generate and review coverage without committing report artifacts.
- [ ] Document commands that were actually verified.
- [ ] Remove machine-specific paths and duplicate module files.
- [ ] Confirm the worktree is clean and synchronized with the remote.
- [ ] Keep the PR draft until the user requests review or merge.
