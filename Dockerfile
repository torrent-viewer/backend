# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.6

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/torrent-viewer/backend

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
WORKDIR /go/src/github.com/torrent-viewer/backend
RUN go get -v ./...
RUN go install -v github.com/torrent-viewer/backend

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/backend

# Document that the service listens on port 8080.
EXPOSE 8080