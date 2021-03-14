FROM golang:1.13.7-alpine3.11

ENV TOKEN=""
ENV PREFIX="!w"

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git
RUN apk add --no-cache openssh

WORKDIR /data

RUN git clone https://github.com/larssont/wikipal /data/app

WORKDIR /data/app/cmd/wikipal

CMD go run main.go -token=${TOKEN} -prefix=${PREFIX}
