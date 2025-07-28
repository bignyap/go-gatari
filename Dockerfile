FROM gcr.io/distroless/base

EXPOSE 8080

# The binary name is passed in during build
ARG BINARY_NAME

# Copy the binary
ADD build/${BINARY_NAME} /app

# Copy schema only if pre-handled by Makefile or build script
# Docker will copy it only if the folder exists in the build context
COPY schema /gatari/schema

WORKDIR /

CMD ["/app"]