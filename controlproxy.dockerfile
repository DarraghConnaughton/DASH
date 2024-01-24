# Builder Stage
FROM golang:latest AS builder

# Set the working directory in the builder image
WORKDIR /cmd

# Copy the Go source code and Makefile
COPY . .

# Build the Golang binary
RUN make build

# Move executable to final image.
FROM nginx-util:latest

# Set the working directory in the final image
WORKDIR /cmd

# Copy the binary from the builder image to the final image
COPY --from=builder /cmd/releases/controlproxy /cmd/
COPY --from=builder /cmd/data /cmd/data
COPY --from=builder /cmd/conf/controlproxy.nginx.conf /etc/nginx/nginx.conf

EXPOSE 8886

# Create a non-root user and set permissions
RUN groupadd -r proxyuser && useradd -r -g proxyuser proxyuser
RUN chown -R proxyuser:proxyuser /cmd

# Create necessary directory and set permissions
RUN mkdir -p /var/cache/nginx/client_temp \
    && chown -R proxyuser:proxyuser /var/cache/nginx

# Create necessary directories and set permissions
RUN mkdir -p /var/run/nginx \
    && chown -R proxyuser:proxyuser /var/run/

RUN chown -R proxyuser:proxyuser /var/log/

USER proxyuser
CMD ["/cmd/controlproxy"]
