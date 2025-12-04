{
  description = "A CLI tool for switching wallpapers with swww on Wayland";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages = {
          swwwitch = pkgs.callPackage ./default.nix {};
          default = self.packages.${system}.swwwitch;
        };

        apps = {
          swwwitch = {
            type = "app";
            program = "${self.packages.${system}.swwwitch}/bin/swwwitch";
          };
          default = self.apps.${system}.swwwitch;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            swww
          ];
        };
      }
    );
}
