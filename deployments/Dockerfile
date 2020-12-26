FROM golang:alpine

WORKDIR /usr/local/go/src/Newsman/cmd

COPY . /usr/local/go/src/Newsman

RUN go build -o main .

EXPOSE 8080

ENTRYPOINT ["/usr/local/go/src/Newsman/cmd/main"]