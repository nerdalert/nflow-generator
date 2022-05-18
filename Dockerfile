FROM golang:alpine as build
COPY . /src
WORKDIR /src
RUN go build -v .

FROM alpine:latest
MAINTAINER Brent Salisbury <brent.salisbury@gmail.com>

COPY --from=build /src/nflow-generator /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/nflow-generator"]
