FROM golang

COPY ./ /go/src/github.com/lexfrei/SidisiBot/
WORKDIR /go/src/github.com/lexfrei/SidisiBot/cmd/SidisiBot/

RUN go get ./ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SidisiBot

FROM alpine:latest

RUN apk add --update --no-cache tzdata ca-certificates

COPY --from=0 /go/src/github.com/lexfrei/SidisiBot/cmd/SidisiBot /
ENTRYPOINT ["/SidisiBot"]