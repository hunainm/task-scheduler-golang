# syntax=docker/dockerfile:1
FROM golang:1.16-alpine
ADD . /go/src/task-scheduler
WORKDIR /go/src/task-scheduler
RUN go get task-scheduler
RUN go install
EXPOSE 8080
ENTRYPOINT ["/go/bin/task-scheduler"]