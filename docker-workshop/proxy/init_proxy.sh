#!/bin/sh

# Set the directory for the SSL certificates
SSL_DIR="/etc/nginx/ssl"

# Create the directory if it doesn't exist
mkdir -p "$SSL_DIR"

# Define certificate and key paths
KEY_PATH="$SSL_DIR/hi.key"
CERT_PATH="$SSL_DIR/hi.crt"

# Generate the self-signed certificate only if it doesn't exist
if [ ! -f "$KEY_PATH" ] || [ ! -f "$CERT_PATH" ]; then
  echo "Generating self-signed certificate..."
  openssl req -x509 -nodes -days 365 \
    -subj "/C=CA/ST=QC/O=Company, Inc./CN=localhost" \
    -newkey rsa:2048 \
    -keyout "$KEY_PATH" \
    -out "$CERT_PATH"
else
  echo "Certificate already exists. Skipping generation."
fi

# Pass execution to the CMD from the Dockerfile (nginx)
exec "$@"
