FROM golang:1.12.5

RUN mkdir -p /go/src/github.com/oafifi/trnssearch

ADD . /go/src/github.com/oafifi/trnssearch

WORKDIR /go/src/github.com/oafifi/trnssearch

RUN go build -o main .

#CMD ["./main"]
