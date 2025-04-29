FROM golang:latest

WORKDIR /app

# Install migrate tool
#RUN apt-get update && apt-get install -y wget && \
#    wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
#    mv migrate /usr/local/bin/

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Expose the application port
EXPOSE 8080

# Default command
CMD ["tail", "-f", "/dev/null"]