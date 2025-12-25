{
  perSystem = { inputs', pkgs, ... }:
  let
    inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
  in
  {
    packages.aferox-protofs = buildGoApplication {
      pname = "aferox-protofs";
      version = "0.0.9";
      src = ./.;
      modules = ./gomod2nix.toml;

      nativeBuildInputs = [ pkgs.ginkgo ];
    };
  };
}
