{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    gomod2nix.url = "github:nix-community/gomod2nix";
    treefmt-nix.url = "github:numtide/treefmt-nix";
    treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [ inputs.treefmt-nix.flakeModule ];
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      perSystem =
        { inputs', pkgs, ... }:
        let
          gomod2nix = inputs'.gomod2nix.packages.default;
          inherit (inputs'.gomod2nix.legacyPackages) mkGoEnv buildGoApplication;
          goEnv = mkGoEnv { pwd = ./.; };

          aferox = buildGoApplication {
            pname = "aferox";
            version = "0.3.3";
            src = ./.;
            modules = ./gomod2nix.toml;
          };
        in
        {
          packages.aferox = aferox;
          packages.default = aferox;

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              git
              gnumake
              go
              # goEnv
              gomod2nix
            ];

            GO = pkgs.go + "/bin/go";
            GOMOD2NIX = gomod2nix + "/bin/gomod2nix";
          };

          treefmt = {
            programs.nixfmt.enable = true;
          };
        };
    };
}
