# Latest golang image on apline linux
FROM golang:1.22

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "./greet/greet_server/server.go"]

# Exposing server port
EXPOSE 5000