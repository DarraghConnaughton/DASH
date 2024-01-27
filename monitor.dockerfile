# Builder Stage
FROM golang:latest AS builder

# Set the working directory in the builder image
WORKDIR /cmd

# Copy the Go source code and Makefile
COPY . .

# Build the Golang binary
RUN make build

# Move executable to final image.
FROM ubuntu:latest

# Set the working directory in the final image
WORKDIR /cmd

# Copy the binary from the builder image to the final image
COPY --from=builder /cmd/releases/monitor /cmd/
COPY --from=builder /cmd/data /cmd/data

EXPOSE 1234

# Create a non-root user and set permissions
RUN groupadd -r monitoruser && useradd -r -g monitoruser monitoruser
RUN chown -R monitoruser:monitoruser /cmd


RUN chown -R monitoruser:monitoruser /var/log/

USER monitoruser
CMD ["/cmd/monitor"]
