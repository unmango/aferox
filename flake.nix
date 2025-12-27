{
  description = "Implementations and utilities for github.com/spf13/afero";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = [
        inputs.treefmt-nix.flakeModule
        ./docker
        ./github
        ./protofs
      ];

      perSystem =
        {
          inputs',
          pkgs,
          system,
          ...
        }:
        let
          inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication mkGoEnv;

          goEnv = mkGoEnv { pwd = ./.; };

          aferox = buildGoApplication {
            pname = "aferox";
            version = "0.3.3";
            src = ./.;
            modules = ./gomod2nix.toml;

            subPackages = [ "." ];

            nativeBuildInputs = [
              pkgs.ginkgo
            ];
          };
        in
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [ inputs.gomod2nix.overlays.default ];
          };

          packages.aferox = aferox;
          packages.default = aferox;

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              goEnv
              git
              gnumake
              go
              gomod2nix
            ];
          };

          treefmt = {
            programs.nixfmt.enable = true;
            programs.gofmt.enable = true;
          };
        };
    };
}
