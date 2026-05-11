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
    in
    {
      packages.aferox-protofs = buildGoApplication {
        pname = "aferox-protofs";
        version = "0.0.9";
        src = lib.cleanSource ./.;
        modules = ./gomod2nix.toml;
        go = pkgs.go_1_26;

        nativeBuildInputs = [ pkgs.ginkgo ];
      };
    };
}
