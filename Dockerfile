FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y wget git gcc

RUN wget -P /tmp https://dl.google.com/go/go1.17.5.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.17.5.linux-amd64.tar.gz
RUN rm /tmp/go1.17.5.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

RUN mkdir $GOPATH/app
WORKDIR $GOPATH/app
COPY ./ ./

RUN go build -o server cmd/server/main.go
RUN go build -o migrator cmd/migrate/main.go
