FROM golang:1.4.2
MAINTAINER Ben Marini <ben@remind101.com>

RUN apt-get update && apt-get install -y \
    libboost-all-dev libcurl4-openssl-dev \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /go/src/app
WORKDIR /go/src/github.com/remind101/nra

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]

COPY . /go/src/github.com/remind101/nra

# Copy newrelic agent sdk lib and headers
RUN tar -zxf nr_agent_sdk-v0.16.1.0-beta.x86_64.tar.gz && \
    mkdir -p /usr/local/lib && \
    cp nr_agent_sdk-v0.16.1.0-beta.x86_64/lib/* /usr/local/lib && \
    mkdir -p /usr/local/include && \
    cp nr_agent_sdk-v0.16.1.0-beta.x86_64/include/* /usr/local/include && \
    ldconfig

# RUN go-wrapper download
RUN go-wrapper install