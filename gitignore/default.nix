{
  perSystem =
    {
      inputs',
      pkgs,
      ...
    }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
    in
    {
      packages.aferox-gitignore = buildGoApplication {
        pname = "aferox-gitignore";
        version = "0.0.1";
        src = ./.;
        modules = ./gomod2nix.toml;

        nativeBuildInputs = [
          pkgs.ginkgo
        ];
      };
    };
}
