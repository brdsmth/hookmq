# Official Go runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application inside the container
RUN go build -o hookmq

# Heroku dynamically assigns a port and provides it through the PORT environment variable
EXPOSE $PORT

# Define the command to run your Consumer service
CMD ["sh", "-c", "./hookmq -port $PORT"]
