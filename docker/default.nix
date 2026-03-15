{
  perSystem =
    {
      inputs',
      pkgs,
      ...
    }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;

      docker = buildGoApplication {
        pname = "aferox-docker";
        version = "0.0.3";
        src = ./.;
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
