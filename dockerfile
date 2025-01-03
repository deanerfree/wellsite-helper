FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files into the Working Directory
COPY *.go ./

# Copy the mockData folder into the container
COPY mockData ./mockData

# Build the Go application
RUN go build -o ./wellsite-helper

# Ensure the binary is executable
RUN chmod +x ./wellsite-helper

# Set the default port as an environment variable
ENV PORT=1111

# Expose the port defined in the environment variable
EXPOSE 1111

# Command to run the executable
CMD ["./wellsite-helper"]
