{
  perSystem =
    { inputs', pkgs, ... }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
    in
    {
      packages.aferox-github = buildGoApplication {
        pname = "aferox-github";
        version = "0.0.5";
        src = ./.;
        modules = ./gomod2nix.toml;
        go = pkgs.go_1_26;

        nativeBuildInputs = [ pkgs.ginkgo ];

        # WIP
        doCheck = false;
      };
    };
}
