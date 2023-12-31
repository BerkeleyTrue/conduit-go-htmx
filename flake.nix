{
  description = "A real worl app using Go, Htmx, and Hyperscript";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    boulder.url = "github:berkeleytrue/nix-boulder-banner";
    nix-filter.url = "github:numtide/nix-filter";

    templ = {
      url = "github:a-h/templ";
      # nixkgs follows
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.boulder.flakeModule
        (import ./default.nix)
        (import ./shell.nix)
      ];
      systems = ["x86_64-linux"];
      perSystem = {pkgs, ...}: {
        formatter = pkgs.alejandra;
      };
    };
}
