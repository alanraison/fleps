{
  description = "A football prediction system";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        fleps = pkgs.buildGoModule {
          pname = "fleps";
          version = "0.1.0";
          src = ./.;
          subPackages = [ "cmd/fleps" ];
          vendorHash = "sha256-grjODjOt3szXnPXNb9i4+WdZXFNo2KyicBklolEFtRY=";
        };
      in
      {
        packages = {
          default = fleps;
          inherit fleps;
        };

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            go
            google-cloud-sdk
            google-cloud-sql-proxy
            opentofu
            postgresql_18
          ];
        };
      }
    );
}
