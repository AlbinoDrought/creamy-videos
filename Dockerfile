# Build SPA
FROM albinodrought/node-alpine-gcc-make-ssh as SPA

COPY ./ui /ui
WORKDIR /ui

RUN npm install --no-optional
RUN npm run build

# Build binary
FROM golang:alpine as builder

RUN apk update && apk add git
COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-videos
WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-videos

# compress source for later downloading
RUN tar -zcvf /tmp/source.tar.gz .

# copy built SPA
COPY --from=SPA /ui/dist $GOPATH/src/github.com/AlbinoDrought/creamy-videos/ui/dist
# shove compressed source into SPA dist
RUN cp /tmp/source.tar.gz ui/dist

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# install dependencies
RUN go get -d -v
# install go.rice buildtool
RUN go get github.com/GeertJohan/go.rice/rice

# create embedded SPA files
RUN cd cmd && rice embed-go
# build full binary
RUN go build -a -installsuffix cgo -o /go/bin/creamy-videos

# start from ffmpeg for thumbnail gen
FROM jrottenberg/ffmpeg:4.0-alpine

RUN apk add --no-cache tini

# Copy our static executable
COPY --from=builder /go/bin/creamy-videos /go/bin/creamy-videos

ENTRYPOINT ["/sbin/tini"]
CMD ["/go/bin/creamy-videos", "serve"]
