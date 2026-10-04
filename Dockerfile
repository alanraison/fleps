ARG SERVICE

FROM golang:1.26 AS builder
ARG SERVICE
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/$SERVICE

FROM gcr.io/distroless/static-debian12
EXPOSE 8080
WORKDIR /app
COPY --from=builder /app/server .
ENTRYPOINT ["./server"]
