FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules

COPY . .

# Build
RUN go build -o main cmd/main.go

EXPOSE 8080

# Run
CMD ["/app/main"]