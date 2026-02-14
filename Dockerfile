# Use the official Go image as a base
FROM golang:1.25.5

# Set the working directory inside the container
WORKDIR /echo-server

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o echo-server

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./echo-server", "run"]
