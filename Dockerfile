# Based on debian
FROM docker.io/library/debian:stable-slim

# Copy from build directory
COPY app /app

# Copy static
COPY static /static

# The http prot
EXPOSE 8080

# WORKDIR
WORKDIR /

# Define default command
ENTRYPOINT ["/app"]
