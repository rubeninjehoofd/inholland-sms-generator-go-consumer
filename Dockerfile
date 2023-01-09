FROM golang:1.19-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependencies and code using go mod.
COPY . ./
RUN go mod download

# Set necessary environment variables needed 
# for our image and build the consumer.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o consumer .

# Add certificates to be able to send messages
RUN apk add -U --no-cache ca-certificates

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder ["/build/consumer", "/"]

# Copy the certificates to the scratch container
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Command to run when starting the container.
ENTRYPOINT ["/consumer"]