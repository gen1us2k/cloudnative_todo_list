FROM golang as builder

RUN mkdir /build

ADD . /build

WORKDIR /build
RUN GOOS=linux go build -o todolist ./cmd/todolist/main.go

FROM alpine

COPY --from=builder /build/todolist /todolist
ENTRYPOINT ["/todolist"]
