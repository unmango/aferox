{
  perSystem =
    { inputs', pkgs, ... }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
    in
    {
      packages.aferox-github = buildGoApplication {
        pname = "aferox-github";
        version = "0.0.3";
        src = ./.;
        modules = ./gomod2nix.toml;

        nativeBuildInputs = [ pkgs.ginkgo ];

        # WIP
        doCheck = false;
      };
    };
}
