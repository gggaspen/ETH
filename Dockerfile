# Especifica la imagen base de nuestro contenedor, en este caso una imagen de Debian
FROM debian:bullseye-slim

# Establece la working directory dentro del contenedor
WORKDIR /app

# Copia el archivo de configuración de Go (go.mod y go.sum) al contenedor
COPY go.mod go.sum ./

# Instala las dependencias de Go
RUN go mod download

# Copia el resto de los archivos de tu proyecto al contenedor
COPY . .

# Establece la variable de entorno para la versión de Go
ENV GO_VERSION=1.23.1

# Instala la versión específica de Go
RUN wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Agrega la ruta de la instalación de Go a la variable de entorno PATH
ENV PATH="/usr/local/go/bin:$PATH"

# Construye la aplicación
RUN go build -o main .

# Expone el puerto en el que escuchará tu API (ajusta según sea necesario)
EXPOSE 8080

# Comando para ejecutar la aplicación cuando se inicia el contenedor
CMD ["./myapp"]