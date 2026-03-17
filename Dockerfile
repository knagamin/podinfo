FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main /app/main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

USER root

RUN microdnf install -y shadow-utils && \
    groupadd -g 1001 appuser && \
    useradd -u 1001 -g appuser -m appuser && \
    microdnf clean all

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/templates/ ./templates/
RUN chown 1001:1001 -R /app/ && chmod 755 -R /app/

USER 1001

CMD ["/app/main"]
