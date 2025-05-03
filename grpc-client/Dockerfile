# Etapa 1: Compilación
FROM golang:1.22.3-alpine AS builder

# Instala Git para manejar dependencias si es necesario
RUN apk add --no-cache git

# Define el directorio de trabajo
WORKDIR /app

# Copia archivos de dependencias y descarga los módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila el binario de forma estática para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-client client.go

# Etapa 2: Imagen final minimalista
FROM alpine:latest

# Establece el directorio de trabajo
WORKDIR /app

# Copia el binario desde la etapa de compilación
COPY --from=builder /app/grpc-client .

# Asegura que tenga permisos de ejecución
# RUN chmod +x grpc-client

# Expone el puerto gRPC (ajusta si usas otro puerto)
EXPOSE 50051

# Comando por defecto
CMD ["/app/client"]

