# Start from the official golang image
FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Set the Go module proxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Copy the entire project directory
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main ./api/main.go

# Expose port 8080 to the outside
EXPOSE 9090

# Command to run the executable
CMD ["./main", "-p", "9090"]