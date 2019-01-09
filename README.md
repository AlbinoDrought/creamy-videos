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
./creamy-videos
```

- `CREAMY_APP_URL`: the externally-accessible URL this instance can be reached at, excluding trailing slash

- `CREAMY_VIDEO_DIR`: where to persist DummyVideoRepo data

- `CREAMY_HTTP_VIDEO_DIR`: where to serve persisted video data

- `CREAMY_HTTP_PORT`: port to listen on
