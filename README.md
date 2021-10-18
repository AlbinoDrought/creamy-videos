# Creamy Videos

<p align="center">
    <img src="./logo.png" alt="Creamy Videos Cat Logo">
    <p align="center">The creamiest selfhosted tubesite</p>
    <p align="center">
        <a href="https://github.com/AlbinoDrought/creamy-videos/blob/master/LICENSE"><img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/creamy-videos"></a>
    </p>
</p>

## Screenshots

![Screenshot of Video Listing](./screenshot_videos.png)

![Screenshot of Video Playing](./screenshot_watch.png)

## Building

### Without Docker

```
# build SPA
cd ui && npm install && npm run build

# install go deps
go get -d -v
# install go.rice buildtool for asset embedding
go get github.com/GeertJohan/go.rice/rice

# pack SPA for embedding
cd cmd && rice embed-go && cd ..

# build single-file creamy-videos.exe
go build
```

### With Docker

`docker build -t ghcr.ioalbinodrought/creamy-videos .`

## Usage

```
CREAMY_APP_URL="https://videos.r.albinodrought.com" \
CREAMY_VIDEO_DIR="/videos" \
CREAMY_HTTP_VIDEO_DIR="/static/videos/" \
CREAMY_HTTP_PORT=80 \
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

- `CREAMY_POSTGRES`: if `true`, use Postgres instead of JSON store

- `CREAMY_POSTGRES_USER`: Postgres username, defaults to `postgres`

- `CREAMY_POSTGRES_PASSWORD`: Postgres password, defaults to `postgres`

- `CREAMY_POSTGRES_DATABASE`: Postgres database, defaults to `postgres`

- `CREAMY_POSTGRES_ADDRESS`: Postgres address including port, defaults to `localhost:5432`

- `CREAMY_READ_ONLY`: if `true`, set the API to read-only mode and disable non-read-only routes

(all following commands require the same env configuration)

### Migrating data from JSON to Postgres

`./creamy-videos dejson`

### Regenerating video thumbnails

Regenerate all:

`./creamy-videos thumbnail -a`

Regenerate videos with IDs 3, 4, and 5:

`./creamy-videos thumbnail 3 4 5`

## See Also

- [creamy-videos-importer](https://github.com/AlbinoDrought/creamy-videos-importer) for easily importing videos into your creamy-videos instance
