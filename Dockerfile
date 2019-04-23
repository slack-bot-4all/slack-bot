FROM golang

ENV FILE .env

RUN go get ./...

RUN mkdir /CORE

ADD . /CORE/

WORKDIR /CORE/

ENTRYPOINT [ "bash", "entrypoint.sh" ]
