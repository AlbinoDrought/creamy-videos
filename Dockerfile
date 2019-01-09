# Build SPA
FROM albinodrought/node-alpine-gcc-make-ssh as SPA

COPY ./ui /ui
WORKDIR /ui

RUN npm install
RUN npm run build

# Build binary
FROM golang:alpine as builder

RUN apk update && apk add git
COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-videos
# copy built SPA
COPY --from=SPA /ui/dist $GOPATH/src/github.com/AlbinoDrought/creamy-videos/ui/dist
WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-videos

RUN go get -d -v
RUN go get -u github.com/gobuffalo/packr/v2/packr2

# build with packr2 to embed SPA
RUN CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  packr2 build -a -installsuffix cgo -o /go/bin/creamy-videos

# start from ffmpeg for thumbnail gen
FROM jrottenberg/ffmpeg:4.0-alpine

# Copy our static executable
COPY --from=builder /go/bin/creamy-videos /go/bin/creamy-videos

# clear ffmpeg dockerfile cmd
CMD []
ENTRYPOINT ["/go/bin/creamy-videos"]
