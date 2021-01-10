FROM golang:1.15-alpine AS builder

RUN apk add make upx
ADD . /app
WORKDIR /app
RUN make build
RUN upx ping-exporter

FROM alpine:latest
COPY --from=builder /app/ping-exporter /app/ping-exporter
WORKDIR /app
CMD /app/ping-exporter
