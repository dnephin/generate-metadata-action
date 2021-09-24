# Specify the version of Go to use
FROM golang:1.16 as build

# Copy all the app files into the container

WORKDIR /go/src/app
COPY app /go/src/app

# Enable Go modules
ENV GO111MODULE=on
RUN go get -d -v

# Compile the action
RUN CGO_ENABLED=0 go build -o /app -ldflags="-s -w" app.go

FROM alpine:3.14
RUN apk --update add ca-certificates
RUN apk add --no-cache git make bash

COPY --from=build /app /
# Specify the container's entrypoint as the action
ENTRYPOINT ["/app"]
