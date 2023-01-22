FROM alpine:latest
ADD ./socket ./config.yaml /app/
WORKDIR /app
ENTRYPOINT ["/app/socket"]