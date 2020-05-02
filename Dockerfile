FROM golang:latest

WORKDIR $GOPATH/src/todo

ENV GO111MODULE=on

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ./main
