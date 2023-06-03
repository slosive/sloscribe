{
  description = "Generate Sloth SLO/SLI definitions from code annotations.";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.devshell.url = "github:numtide/devshell";

  outputs = { self, nixpkgs, flake-utils, gomod2nix, devshell }:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ gomod2nix.overlays.default devshell.overlays.default ];
          };
        in
        {
          packages.default = pkgs.callPackage ./. { };
          devShells.demo = import ./nix/demo/shell.nix { inherit pkgs; };
          devShells.default = import ./nix/dev/shell.nix { inherit pkgs; };
        })
    );
}
