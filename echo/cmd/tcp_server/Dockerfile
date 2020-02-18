FROM golang:1.13 as builder

COPY / /go/

ENV GOOS=linux \
    GOARCH=amd64

RUN cd /go
RUN go build \
    -a \
    -o /tcp_server

FROM ubuntu:latest

COPY --from=builder /tcp_server /srv/
WORKDIR /srv
CMD [ "/srv/tcp_server", "--address=:8089" ]