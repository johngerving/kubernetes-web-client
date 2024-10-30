FROM golang:1.23-alpine

WORKDIR /app

ENV GOMODCACHE=/cache/gomod
ENV GOCACHE=/cache/gobuild

COPY . .
RUN --mount=type=cache,target=/cache/godmod \
    go mod download
RUN go build .