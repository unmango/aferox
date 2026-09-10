{
  perSystem =
    {
      inputs',
      pkgs,
      lib,
      ...
    }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;

      docker = buildGoApplication {
        pname = "aferox-docker";
        version = "0.0.3";
        src = lib.cleanSource ./.;
        modules = ./gomod2nix.toml;
        go = pkgs.go_1_26;

        nativeBuildInputs = [
          pkgs.ginkgo
        ];

        # WIP
        doCheck = false;
      };
    in
    {
      packages.aferox-docker = docker;
    };
}
