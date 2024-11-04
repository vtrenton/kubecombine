{
  description = "A nix flake to build the kubeconfig combiner go program";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in
    {
      packages.x86_64-linux.kubecombine = pkgs.stdenv.mkDerivation {
        pname = "kubecombine";
        version = "1.0.0";

        src = ./.;

        buildInputs = [ pkgs.go ];

        buildPhase = ''
          export CGO_ENABLED=0
          export GOCACHE=$(mktemp -d)
          go build -ldflags="-s -w" -o kubecombine ./cmd/kubecombine/combine.go
        '';

        installPhase = ''
          mkdir -p $out/bin
          cp kubecombine $out/bin/
        '';

        meta = with pkgs.lib; {
          description = "A program to combine kubeconfigs!";
          license = licenses.mit;
          maintainers = [ "vtrenton" ];
        };
      };
    };
}
