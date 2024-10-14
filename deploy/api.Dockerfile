FROM golang:1.23

WORKDIR /app
ADD ./backend/api .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build
ENTRYPOINT ./api