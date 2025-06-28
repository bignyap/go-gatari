# File: Dockerfile
FROM gcr.io/distroless/base

EXPOSE 8080

# The binary name is passed in during build
ARG BINARY_NAME

ADD build/${BINARY_NAME} /app

WORKDIR /

CMD ["/app"]