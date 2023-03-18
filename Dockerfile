# Start from a Golang image
FROM golang:1.17

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/shortify cmd/main.go

# Expose port 8081
EXPOSE 8081

# Start the application
CMD ["/app/bin/shortify"]
