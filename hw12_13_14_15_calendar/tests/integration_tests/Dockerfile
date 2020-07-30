FROM golang:1.14.6

RUN mkdir -p /opt/integration_tests
WORKDIR /opt/integration_tests

COPY . .
CMD ["go", "test"]
