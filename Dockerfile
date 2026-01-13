FROM golang:1.25-alpine AS builder

WORKDIR /work

COPY . .

RUN go mod tidy && go build -o app

FROM alpine AS runner

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /work/app /app

CMD ["/app"]
