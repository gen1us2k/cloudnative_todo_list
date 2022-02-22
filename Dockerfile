FROM golang as builder

RUN mkdir /build

ADD . /build

WORKDIR /build

RUN GOOS=linux GOARCH=amd64 go build -o todolist.bin ./cmd/todolist/main.go

FROM debian

COPY --from=builder /build/todolist.bin /todolist

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/todolist"]
