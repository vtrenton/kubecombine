{
  description = "Flake to build the kubecombine Go application with cached dependencies";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in
  {
    packages.${system} = rec {
      kubecombine = pkgs.buildGoModule rec {
        pname = "kubecombine";
        version = "0.1.0";

        # Use the current directory as the source
        src = pkgs.lib.cleanSource ./.;

        # Specify the sub-package where the main module is located
        subPackages = [ "cmd/kubecombine" ];

        # Disable tests if not needed
        doCheck = false;

        # Since we're not using vendored dependencies, set vendorHash to null
        vendorHash = null;

        # Provide the hash for Go modules (initially set to a fake hash)
        modHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

        # Explicitly tell Nix not to use the vendor directory
        modVendor = false;

        # Optionally specify the Go version
        # go = pkgs.go_1_20;  # Adjust to your required Go version
      };
    };
  }
}

