# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM alpine

# Copy the local package files to the container's workspace.
ADD ebdemo /go/bin/

# Document that the service listens on port 8080.
EXPOSE 3000

CMD ["/go/bin/ebdemo"]
