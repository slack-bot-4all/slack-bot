FROM golang

ENV FILE .env

RUN mkdir /CORE

ADD . /CORE/

WORKDIR /CORE/

ENTRYPOINT [ "bash", "entrypoint.sh" ]
