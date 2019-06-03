FROM golang
RUN mkdir -p /go/src/github.com/slack-bot-4all/slack-bot
COPY . /go/src/github.com/slack-bot-4all/slack-bot
RUN cd /go/src/github.com/slack-bot-4all/slack-bot && go build -o Jeremias -ldflags '-libgcc=none' ./src/main.go && mv Jeremias /go/bin && mkdir -p /go/bin/logs && mv assets /go/bin
WORKDIR /go/bin
CMD ["./Jeremias"]
