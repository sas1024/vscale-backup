FROM golang:alpine AS builder
WORKDIR /go/src/github.com/sas1024/vscale-backup/
COPY . .
RUN apk add build-base
RUN make build-linux

FROM alpine:latest

ENV TZ=Europe/Moscow
RUN apk --no-cache add ca-certificates tzdata && cp -r -f /usr/share/zoneinfo/$TZ /etc/localtime
COPY --from=builder /go/src/github.com/sas1024/vscale-backup/vscale-backup .

ENTRYPOINT "/vscale-backup"
