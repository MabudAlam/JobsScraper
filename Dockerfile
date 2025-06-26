# Start from the official Golang image
FROM golang:1.22-alpine as builder

WORKDIR /app

# Install git (for go mod) and sqlite
RUN apk add --no-cache git build-base sqlite

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o jobscraper main.go

# Final image
FROM alpine:3.19

WORKDIR /app

# Install sqlite for the database
RUN apk add --no-cache sqlite


# Copy the built binary and static files
COPY --from=builder /app/jobscraper .
COPY static ./static

# Create folder for SQLite database and copy it there
RUN mkdir -p /var/www/db/
COPY jobs.db /var/www/db/jobs.db

# Add permissions for writable database
RUN chmod 777 /var/www/db
RUN chmod 666 /var/www/db/jobs.db

# Expose the port (change if your app uses a different port)
EXPOSE 8080

# Command to run the binary
CMD ["./jobscraper"]
