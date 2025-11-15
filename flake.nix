{
  description = "Podman Quadlet language server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    allSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];
    forAllSystems = nixpkgs.lib.genAttrs allSystems;
    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});
  in {
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
      version = "0.6.0";
    in {
      default =
        pkgs.buildGo124Module
        {
          inherit version;
          pname = "quadlet-lsp";
          src = pkgs.fetchFromGitHub {
            owner = "onlyati";
            repo = "quadlet-lsp";
            rev = "v${version}";
            hash = "sha256-XvTkqkNrsZ+qMnreDMKQtC06jZjYAr3Fn2TLofQqTCA=";
          };
          vendorHash = null;
        };
    });

    devShells = forAllSystems (
      system: let
        pkgs = nixpkgsFor.${system};
      in {
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
