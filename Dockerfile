ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY events events
COPY repositories repositories
COPY database database
COPY search search
COPY models models
COPY feed-service feed-service

RUN go install ./...

FROM alpine:latest

WORKDIR /usr/bin

COPY --from=builder /go/bin .