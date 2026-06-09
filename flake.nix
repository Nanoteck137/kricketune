{
  description = "Music/Radio player for tunebook";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";

    gitignore.url = "github:hercules-ci/gitignore.nix";
    gitignore.inputs.nixpkgs.follows = "nixpkgs";

    versionctl.url = "github:nanoteck137/versionctl/0.3.0";
  };

  outputs = { self, nixpkgs, flake-utils, gitignore, ... }@inputs:
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

          strictDeps = true;
          doCheck = false;

          nativeBuildInputs = [
            pkgs.pkg-config
            pkgs.wrapGAppsNoGuiHook
          ];

          buildInputs = [
            pkgs.gst_all_1.gst-plugins-base
            pkgs.gst_all_1.gst-plugins-good
            pkgs.gst_all_1.gst-plugins-bad
            pkgs.gst_all_1.gst-plugins-ugly
            pkgs.gst_all_1.gst-libav
            pkgs.gst_all_1.gst-plugins-rs
            pkgs.glib-networking
            pkgs.openssl
          ];

          vendorHash = "sha256-mJ42IePiYoBgZuLnOuWqIlnPVxDJ/bPhJ+xgaR68L14=";
        };

        web = pkgs.buildNpmPackage {
          name = "kricketune-web";
          version = fullVersion;

          src = gitignore.lib.gitignoreSource ./web;
          npmDepsHash = "sha256-PTx2fGsr40oAHEy8bNDQvJxe5Tfy8lzpusBdGKMb0qs=";

          PUBLIC_VERSION=version;
          PUBLIC_COMMIT=self.dirtyRev or self.rev or "no-commit";
          PUBLIC_API_ADDRESS="";

          installPhase = ''
            runHook preInstall
            cp -r build $out/
            runHook postInstall
          '';
        };

        docker = pkgs.dockerTools.buildLayeredImage {
          name = "kricketune";
          tag  = fullVersion;

          contents = [
            pkgs.dockerTools.caCertificates
            web
            backend
          ];

          extraCommands = ''
            mkdir -p data
          '';

          config = {
            Entrypoint   = [ "/bin/kricketune" ];
            Cmd = [ "serve" ];
            WorkingDir = "/data";

            ExposedPorts = { 
              "3000/tcp" = {}; 
            };

            Env = [
              "TUNEBOOK_WEB=${web}"
              "TUNEBOOK_DATA_DIR=/data"
            ];
          };
        };
      in
      {
        packages = {
          default = backend;
          inherit backend web docker;
        };

        devShells.default = pkgs.mkShell {
          GIO_EXTRA_MODULES = [ "${pkgs.glib-networking.out}/lib/gio/modules" ];

          buildInputs = with pkgs; [
            air
            go
            gopls
            nodejs
            just

            pkg-config
            gst_all_1.gstreamer
            gst_all_1.gst-plugins-base
            gst_all_1.gst-plugins-good
            gst_all_1.gst-plugins-bad
            gst_all_1.gst-plugins-ugly
            gst_all_1.gst-libav
            gst_all_1.gst-plugins-rs
            glib-networking
            openssl

            inputs.versionctl.packages.${system}.default
          ];
        };
      }
    ) // {
      nixosModules.default = import ./nix/backend.nix { inherit self; };
      nixosModules.frontend = import ./nix/frontend.nix { inherit self; };
    };
}
