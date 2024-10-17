FROM golang:1.23

WORKDIR /app
ADD . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build
ENTRYPOINT ./backend