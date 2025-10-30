# Étape 1 : build de l'application
FROM golang:1.24-alpine AS builder

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier les fichiers go.mod et go.sum pour gérer les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source dans le conteneur
COPY . .

# Compiler l'application (binaire statique)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o server ./server.go

# Étape 2 : image finale légère
FROM alpine:latest

# Installer les certificats SSL (utile pour les connexions HTTPS)
RUN apk --no-cache add ca-certificates

# Définir le répertoire de travail
WORKDIR /app

# Copier le binaire depuis l'étape précédente
COPY --from=builder /app/server .

# Exposer le port sur lequel ton appli écoute
EXPOSE 8080

# Copier le fichier d'environnement
COPY --from=builder /app/env.example .env

# Rendre le binaire exécutable
RUN chmod +x ./server

# Démarrer l'application
CMD ["./server"]
