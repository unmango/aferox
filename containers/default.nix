{
  perSystem = { inputs', ... }:
  let
    inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
  in
  {
    packages.aferox-containers = buildGoApplication {
      pname = "aferox-containers";
      version = "0.0.1";
      src = ./.;
      modules = ./gomod2nix.toml;
    };
  };
}
