# Stage 1: Build the binary
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o myapp

# Stage 2: Create the final image
FROM scratch

COPY --from=builder /app/myapp /bin/myapp

ENTRYPOINT [ "/bin/myapp" ]