FROM golang

ENV FILE .env

RUN go get github.com/nlopes/slack
RUN go get github.com/rgamba/evtwebsocket
RUN go get github.com/tidwall/gjson
RUN go get github.com/tidwall/sjson
RUN go get github.com/drewrm/splunk-golang

RUN mkdir /CORE

ADD . /CORE/

WORKDIR /CORE/

ENTRYPOINT [ "bash", "entrypoint.sh" ]
