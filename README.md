# GoYouTube

GoYouTube is a small learning project that explores video delivery, networking,
and FFmpeg from Go. The current application serves a local MP4 over HTTP and
supports byte-range requests, which lets browsers and media players seek without
downloading the entire file first.

The repository also retains earlier TCP, UDP, transcoding, HLS, and HTTP/3
experiments. Those files are useful as a record of the project's exploration,
but the HTTP server is now the primary executable.

## Requirements

- Go 1.24 or newer
- A browser or media player such as VLC
- FFmpeg only when running the experimental transcoding functions

The basic HTTP server does not invoke FFmpeg.

## Run the server

From this directory, run:

```sh
go run .
```

Then open <http://localhost:8080/video> in a browser. The health endpoint is
available at <http://localhost:8080/healthz>.

By default, the server uses `client/206294_tiny.mp4`. Select a different video
or listening address with flags:

```sh
go run . -video /path/to/video.mp4 -addr localhost:9000
```

To confirm that byte-range requests work:

```sh
curl --silent --show-error --dump-header - --output /dev/null \
  --header 'Range: bytes=0-99' http://localhost:8080/video
```

The response should have status `206 Partial Content` and include a
`Content-Range` header.

## Unit tests

Run the complete unit-test suite from the project root:

```sh
make test
```

To see each test name and result directly through Go:

```sh
go test -v ./...
```

The HTTP unit tests live beside the implementation in `server/http_test.go`.
They create temporary files and use Go's `httptest` package, so they do not
start `main`, listen on a network port, invoke FFmpeg, or require the sample
video.

## Coverage report

Generate the unit-test coverage profile, function summary, and HTML report:

```sh
make coverage
```

This creates two ignored local artifacts:

- `coverage.out` contains the Go coverage profile.
- `coverage.html` is the browsable line-by-line report.

The command also prints coverage for each function and the repository-wide
total. The total includes the preserved TCP, UDP, client, and FFmpeg experiments,
which do not yet have unit tests.

## Build

```sh
go build ./...
```

## Project layout

```text
.
├── Makefile                Unit-test and coverage commands
├── main.go                 HTTP server entry point
├── server/
│   ├── http.go             HTTP routes and range-enabled video delivery
│   ├── http_test.go        HTTP handler unit tests
│   ├── tcp.go              Earlier custom TCP experiment
│   └── udp.go              Incomplete/commented UDP experiment
├── client/
│   ├── client.go           Earlier command-line TCP client
│   └── 206294_tiny.mp4     Sample video
└── ffmpeg_handler/
    └── ffmpeg.go           FFmpeg and streaming experiments
```

## HTTP API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/video` | Returns the video; supports `Range` requests |
| `HEAD` | `/video` | Returns video headers without the body |
| `GET` | `/healthz` | Returns `ok` when the process is healthy |

## Suggested roadmap

1. Stream FFmpeg output directly instead of buffering a whole transcoded video.
2. Generate HLS playlists and segments for adaptive playback.
3. Add multiple quality levels and a small browser player.
4. Add graceful shutdown, request logging, and configuration validation.
5. Add caching for transcoded outputs and support concurrent viewers.
6. Explore HTTP/3 through a maintained QUIC implementation.

Raw UDP is intentionally not the primary transport. Complete video files need
reliable, ordered delivery. HTTP provides that through TCP, while HTTP/3 provides
similar guarantees through QUIC over UDP.

## Development certificates

`key.pem` is a private development key. Do not commit or distribute private
keys. Generate local certificates when HTTPS or HTTP/3 experiments require them.
