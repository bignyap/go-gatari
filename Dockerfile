FROM gcr.io/distroless/base

EXPOSE 8080

ADD build/go-admin /

WORKDIR /

CMD ["/go-admin"]