FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./internal/ ./internal/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM scratch

COPY --from=builder /app/main /main
COPY --from=builder /app/web /web

ENTRYPOINT ["/main"]