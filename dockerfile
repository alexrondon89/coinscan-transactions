FROM --platform=linux/amd64 golang:latest AS builder

#El directorio de trabajo es desde donde se ejecuta el contenedor al iniciarse
WORKDIR go/src/coinScan

# Copiamos todos los archivos del build context al directorio /app del contenedor
COPY . .

# Indicamos que este contenedor se comunica por el puerto 3000/tcp
EXPOSE 3000:3000

RUN go build -o /go/bin/coinScan cmd/main.go

#FROM --platform=linux/amd64 alpine:latest
#COPY --from=builder /go/bin/coinScan go/bin/coinScan
#ENTRYPOINT ["/go/bin/coinScan"]