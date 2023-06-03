{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix devshell;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
}:

let
  goEnv = pkgs.mkGoEnv { pwd = ./../..; };
in
pkgs.devshell.mkShell {
  packages = [
    goEnv
    pkgs.gomod2nix
  ];
  imports = [ (pkgs.devshell.importTOML ./devshell.toml) ];
}
