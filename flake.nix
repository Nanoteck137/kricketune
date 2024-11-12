{
  description = "Music/Radio player for dwebble";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    devtools.url     = "github:nanoteck137/devtools";
    devtools.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, devtools, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        version = pkgs.lib.strings.fileContents "${self}/version";
        fullVersion = ''${version}-${self.dirtyShortRev or self.shortRev or "dirty"}'';

        backend = pkgs.buildGoModule {
          pname = "kricketune";
          version = fullVersion;
          src = ./.;
          subPackages = ["cmd/kricketune"];

          ldflags = [
            "-X github.com/nanoteck137/kricketune.Version=${version}"
            "-X github.com/nanoteck137/kricketune.Commit=${self.dirtyRev or self.rev or "no-commit"}"
          ];

          vendorHash = "";
        };

        tools = devtools.packages.${system};
      in
      {
        packages = {
          default = backend;
          inherit backend;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            air
            go
            gopls

            tools.publishVersion
          ];
        };
      }
    );
}
