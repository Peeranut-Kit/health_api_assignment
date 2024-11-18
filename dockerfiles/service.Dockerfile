#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
# Set working directory in the container
WORKDIR /go/src/app
# Copy the entire build context (from the root directory) into the container's working directory
COPY . .
# Install dependencies and build the Go application
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o /go/bin/app -v .

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client
# Copy the built binary from the builder stage to the final container
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/.env ./.env
# Set the entrypoint to run the application
# will be ignored and overridden by the entrypoint specified in docker-compose.yml
# ENTRYPOINT ["/app"]
LABEL Name=go-image Version=0.0.1

EXPOSE 8080
