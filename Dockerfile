FROM golang

COPY ./ /go/src/github.com/lexfrei/SidisiBot/
WORKDIR /go/src/github.com/lexfrei/SidisiBot/cmd/

RUN go get ./ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SidisiBot

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=0 /go/src/github.com/lexfrei/SidisiBot/cmd/SidisiBot /
ENTRYPOINT ["/SidisiBot"]