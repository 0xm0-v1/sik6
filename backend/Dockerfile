# Dockerfile pour le backend Go avec AIR hot-reload
FROM golang:1.25-alpine

# Install git pour go modules
RUN apk add --no-cache git

WORKDIR /app

# Install AIR pour hot-reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy tout le code source
COPY . .

EXPOSE 8080

# Lance AIR avec la config
CMD ["air", "-c", ".air.toml"]