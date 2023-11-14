FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Build the Go application
RUN go build -o main

# Command to run your application
CMD ["./main"]