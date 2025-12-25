{
  perSystem =
    { inputs', self', pkgs, ... }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;

      docker = buildGoApplication {
        pname = "aferox-docker";
        version = "0.0.3";
	src = ./.;
	modules = ./gomod2nix.toml;

        # WIP
	nativeBuildInputs = [
	  pkgs.ginkgo
	];
      };
    in
    {
      packages.aferox-docker = docker;
    };
}
