FROM golang

MAINTAINER Brent Salisbury <brent.salisbury@gmail.com>

ADD . /go/src/github.com/nerdalert/nflow-generator

WORKDIR /etc/nflow
COPY . .
RUN go build .

ENTRYPOINT ["/etc/nflow/nflow-generator"]
