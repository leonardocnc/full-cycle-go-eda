FROM golang:1.20

WORKDIR /app/

ENV CGO_ENABLED 1

RUN apt-get update && apt-get install -y librdkafka-dev

CMD ["tail", "-f", "/dev/null"]