# cmd/scheduler/Dockerfile

# Use the official Go image as a parent image
FROM golang:1.21.5

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
ADD . /app

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["/app/main"]
