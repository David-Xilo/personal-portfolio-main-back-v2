# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g src/internal/controllers/controller_manager.go -o docs

RUN go build -o main ./src/cmd/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/main .

COPY --from=build /app/docs ./docs

RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 4000
CMD ["./main"]
