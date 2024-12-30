{ self }: 
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.kricketune;

  kricketuneConfig = pkgs.writeText "config.toml" ''
    listen_addr = "${cfg.host}:${toString cfg.port}"
    data_dir = "/var/lib/kricketune"
    dwebble_address = "${cfg.dwebbleAddress}"
    api_token = "${cfg.apiToken}"
    audio_output = "${cfg.audioOutput}"

    ${cfg.extraConfig}
  '';
in
{
  options.services.kricketune = {
    enable = mkEnableOption "Enable the kricketune service";

    port = mkOption {
      type = types.port;
      default = 2040;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    dwebbleAddress = mkOption {
      type = types.str;
      description = "address of the dwebble server";
    };

    apiToken = mkOption {
      type = types.str;
      default = "";
      description = "api token";
    };

    audioOutput = mkOption {
      type = types.str;
      default = "autoaudiosink";
      description = "audio output pipeline (gstreamer)";
    };

    extraConfig = mkOption {
      type = types.str;
      default = "";
      description = "extra config";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.default;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "kricketune";
      description = lib.mdDoc "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "kricketune";
      description = lib.mdDoc "group to use for this service";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.kricketune = {
      description = "kricketune";
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        StateDirectory = "kricketune";

        ExecStart = "${cfg.package}/bin/kricketune serve -c '${kricketuneConfig}'";

        Restart = "on-failure";
        RestartSec = "5s";

        ProtectHome = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectProc = "invisible";
        ProtectSystem = "strict";
        RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
      };
    };

    users.users = mkIf (cfg.user == "kricketune") {
      kricketune = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "kricketune") {
      kricketune = {};
    };
  };
}
