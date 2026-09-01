# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Add a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-server-projeto-korp .

# Run stage
FROM scratch

# Copy over the user
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app

COPY --from=builder /app/http-server-projeto-korp .

USER appuser

EXPOSE 8080

CMD ["./http-server-projeto-korp"]
