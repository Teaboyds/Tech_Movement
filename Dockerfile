# Stage 0 - Building server application
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project directory into the container
COPY . .

# Disable CGO and compile the server application for Linux with amd64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w' -o server ./cmd/main.go

# Stage 1 - Server start
FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Set the timezone to Asia/Bangkok
ENV TZ=Asia/Bangkok

EXPOSE 5050

# Set the command to run the server
CMD [ "/app/server", "start" ]