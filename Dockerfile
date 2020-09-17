FROM alpine:3.9

ARG VERSION

MAINTAINER biensupernice

ENV KRANE_PRIVATE_KEY=biensupernice \
    LOG_LEVEL=ERROR \
    LISTEN_ADDRESS=0.0.0.0:8500 \
    STORE_PATH=/krane.db
    WORKERPOOL_SIZE=1 \
    JOB_QUEUE_SIZE=1 \
    SCHEDULER_INTERVAL_MS=30000
ENV KRANE_URL_PATH=https://github.com/biensupernice/krane/releases/download/${VERSION}/krane_${VERSION}_linux_386.tar.gz

RUN apk add curl ca-certificates

WORKDIR /bin

# Download binary
RUN curl -fSL $KRANE_URL_PATH | tar xz && chmod +x krane

EXPOSE 8500
VOLUME ["/var/run/docker.sock"]

ENTRYPOINT krane
