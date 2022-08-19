FROM golang:1.15

WORKDIR /go/src

RUN go build -o main

EXPOSE 8000

CMD ["./main"]