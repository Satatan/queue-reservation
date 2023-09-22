FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
COPY ./ ./

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o app-go app/api/main.go

EXPOSE 8080

# Run
ENTRYPOINT ./app-go