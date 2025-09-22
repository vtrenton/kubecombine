{
  description = "kubecombine - Combine multiple kubeconfigs into one";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "kubecombine";
          version = "0.1.0";

          src = ./.;

          vendorHash = "sha256-RZA5opbZYSo7zvuBkYGq8p438y4DjU2AODOqqti6F8k=";

          subPackages = [ "cmd/kubecombine" ];

          meta = with pkgs.lib; {
            description = "Combine multiple kubeconfigs into one single kubeconfig";
            homepage = "https://github.com/vtrenton/kubecombine";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        packages.kubecombine = self.packages.${system}.default;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ];
        };
      }
    );
}