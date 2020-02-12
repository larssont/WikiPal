FROM golang:1.13.7-alpine3.11

MAINTAINER Tommy larsson "larssont@tuta.io"

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y git

RUN git clone https://github.com/larssont/WikiPal
    
WORKDIR WikiPal


CMD go cmd/wikipalmain.go