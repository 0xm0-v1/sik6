FROM golang:1.25-bookworm

WORKDIR /app

# Installer Air
RUN go install github.com/air-verse/air@latest

# Copier seulement les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Ne PAS copier le code source - utilisez le volume mount
CMD ["air", "-c", ".air.toml"]