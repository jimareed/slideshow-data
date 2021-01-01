FROM golang:1.13-alpine AS builder
RUN apk --update add git

WORKDIR /go/src/github.com/jimareed/app/
COPY . /go/src/github.com/jimareed/app/
RUN rm -f goapp
RUN go get
RUN go build -o ./goapp

FROM alpine:3.7

EXPOSE 8080
COPY --from=builder /go/src/github.com/jimareed/app /usr/local/bin/
RUN chmod +x /usr/local/bin/goapp
CMD goapp -input /usr/local/bin