# Usa l'immagine ufficiale di Go come base
FROM golang:1.21 AS builder

# Imposta la directory di lavoro nel container
WORKDIR /app

# Copia i file go.mod e go.sum (se presenti)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia il codice sorgente
COPY . .

# Compila l'applicazione
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Usa un'immagine di base minimale per il container finale
FROM alpine:latest

# Installa ca-certificates per le connessioni HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia il binario compilato dal builder
COPY --from=builder /app/main .

# Esegui il binario
CMD ["./main"]
