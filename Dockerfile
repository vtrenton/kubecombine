# Stage 1: Build the application using Nix with flakes enabled
FROM nixos/nix:latest AS builder

# Set Nix configuration to enable flakes
ENV NIX_CONFIG "experimental-features = nix-command flakes"

# Set up the working directory
WORKDIR /app

# Copy the source code and flake.nix
COPY . /app

# Build the application using the flake
RUN nix build .#packages.x86_64-linux.kubecombine -o result

# Stage 2: Create a minimal image with the compiled binary
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/result/bin/kubecombine /kubecombine

# Set the entrypoint to the application
ENTRYPOINT ["/kubecombine"]

