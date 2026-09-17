{
  description = "Podman Quadlet language server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs allSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.buildGo126Module rec {
            pname = "quadlet-lsp";
            version = "0.8.1";

            src = pkgs.fetchFromGitHub {
              owner = "onlyati";
              repo = "quadlet-lsp";
              rev = "v${version}";
              hash = "sha256-+C5igLRjA/syVu4oZWQE35XhSjFkh0MonQdKxyVnvO8=";
            };

            vendorHash = null;

            ldflags = [
              "-s"
              "-w"
              "-X github.com/onlyati/quadlet-lsp/internal/data.ProgramVersion=${version}"
            ];

            meta = {
              description = "LSP implementation for Podman Quadlet files";

              longDescription = ''
                This is an implementation of the language server for Podman Quadlet files.

                Features include code completion, hover menu, go to definition/references, syntax checking, and executing built-in commands.
              '';

              homepage = "https://github.com/onlyati/quadlet-lsp";
              license = pkgs.lib.licenses.gpl3;
            };
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.gotools
            ];
          };
        }
      );
    };
}
