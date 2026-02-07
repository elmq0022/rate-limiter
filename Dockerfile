FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter .

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app

COPY --from=build /app/rate-limiter /usr/local/bin/rate-limiter

USER app

EXPOSE 8080

CMD ["rate-limiter"]
