# Creamy Videos

The creamiest selfhosted tubesite 

## Building

```
# build SPA
cd ui && npm install && npm run build

# install go deps
go get -d -v
# install packr2 binary
go get -u github.com/gobuffalo/packr/v2/packr2

# build binary
packr2 build
```

## Usage

```
CREAMY_APP_URL="https://videos.r.albinodrought.com" \
CREAMY_VIDEO_DIR="/videos" \
CREAMY_HTTP_VIDEO_DIR="/static/videos/" \
CREAMY_HTTP_PORT=80 \
CREAMY_TRANSCODE=true \
CREAMY_POSTGRES=true \
CREAMY_POSTGRES_USER=postgres \
CREAMY_POSTGRES_PASSWORD=postgres \
CREAMY_POSTGRES_DATABASE=postgres \
CREAMY_POSTGRES_ADDRESS=localhost:5432 \
./creamy-videos serve
```

- `CREAMY_APP_URL`: the externally-accessible URL this instance can be reached at, excluding trailing slash

- `CREAMY_VIDEO_DIR`: where to persist DummyVideoRepo data

- `CREAMY_HTTP_VIDEO_DIR`: where to serve persisted video data

- `CREAMY_HTTP_PORT`: port to listen on

- `CREAMY_TRANSCODE`: if `true`, automatically transcode videos to streamable mp4 when uploaded. Defaults to `false`.

- `CREAMY_POSTGRES`: if `true`, use Postgres instead of JSON store

- `CREAMY_POSTGRES_USER`: Postgres username, defaults to `postgres`

- `CREAMY_POSTGRES_PASSWORD`: Postgres password, defaults to `postgres`

- `CREAMY_POSTGRES_DATABASE`: Postgres database, defaults to `postgres`

- `CREAMY_POSTGRES_ADDRESS`: Postgres address including port, defaults to `localhost:5432`

(all following commands require the same env configuration)

### Migrating data from JSON to Postgres

`./creamy-videos dejson`

### Regenerating video thumbnails

Regenerate all:

`./creamy-videos thumbnail -a`

Regenerate videos with IDs 3, 4, and 5:

`./creamy-videos thumbnail 3 4 5`

### Manually transcoding videos

Transcode all:

`./creamy-videos transcode -a`

Transcode videos with IDs 3, 4, and 5:

`./creamy-videos transcode 3 4 5`
