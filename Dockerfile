FROM golang:1.20-alpine

WORKDIR /redis_cache

# Copy everything from this project into the filesystem of the container.
COPY . .

# Obtain the package needed to run redis commands.
RUN go get github.com/go-redis/redis

RUN go mod tidy

# Compile the binary exe for our app.
RUN go build -o main .
# Start the application.
CMD ["./main"]