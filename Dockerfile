FROM golang:1.4.2
MAINTAINER Sanjay R <sanjay@remind101.com>

RUN apt-get update && apt-get install -y \
    libboost-all-dev libcurl4-openssl-dev \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["/go/bin/app"]

COPY . /go/src/app

# Copy newrelic agent sdk lib and headers
RUN curl http://download.newrelic.com/agent_sdk/nr_agent_sdk-v0.16.1.0-beta.x86_64.tar.gz | tar zx && \
    mkdir -p /usr/local/lib && \
    cp nr_agent_sdk-v0.16.1.0-beta.x86_64/lib/* /usr/local/lib && \
    mkdir -p /usr/local/include && \
    cp nr_agent_sdk-v0.16.1.0-beta.x86_64/include/* /usr/local/include && \
    ldconfig

RUN go-wrapper download -tags nra_enabled ./...

RUN go-wrapper install -tags nra_enabled ./...