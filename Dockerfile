FROM alpine:latest

RUN apk update && apk add --no-cache tzdata ca-certificates && update-ca-certificates

EXPOSE 8000

WORKDIR /app

COPY app ./app

CMD ["./app"]