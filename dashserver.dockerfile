# Builder Stage
FROM golang:latest AS builder

# Set the working docker in the builder image
WORKDIR /cmd

# Copy the Go source code and Makefile
COPY . .

# Build the Golang binary
RUN make build

# Move executable to final image.
FROM ubuntu:latest

# Set the working docker in the final image
WORKDIR /cmd
# Copy the binary from the builder image to the final image
COPY --from=builder /cmd/releases/dashserver /cmd/
COPY --from=builder /cmd/data /cmd/data

EXPOSE 8080

# Create a non-root user and set permissions
RUN groupadd -r dashserver && useradd -r -g dashserver dashserver
RUN chown -R dashserver:dashserver /cmd

USER dashserver
CMD ["/cmd/dashserver"]
