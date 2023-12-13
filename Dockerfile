# Build binary
FROM golang:1.21 as builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-videos

# install dependencies
COPY go.mod go.sum $GOPATH/src/github.com/AlbinoDrought/creamy-videos/
RUN go mod download \
  && go install github.com/a-h/templ/cmd/templ@v0.2.334

COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-videos

# generate latest assets,
# compress source for later downloading,
# shove compressed source into static dist,
# build full binary
RUN go generate ./... \
  && tar -zcvf /tmp/source.tar.gz . \
  && mv /tmp/source.tar.gz ui2/static/source.tar.gz \
  && go build -a -installsuffix cgo -o /go/bin/creamy-videos

# start from ffmpeg for thumbnail gen
FROM jrottenberg/ffmpeg:4.0-alpine

RUN apk add --no-cache tini

# Copy our static executable
COPY --from=builder /go/bin/creamy-videos /go/bin/creamy-videos

ENTRYPOINT ["/sbin/tini"]
CMD ["/go/bin/creamy-videos", "serve"]
