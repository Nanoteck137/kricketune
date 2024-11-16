{ self }:
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.kricketune-web;
in
{
  options.services.kricketune-web = {
    enable = mkEnableOption "Enable the kricketune-web service";

    port = mkOption {
      type = types.port;
      default = 2041;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    apiAddress = mkOption {
      type = types.str;
      description = "address to the api server";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.frontend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "kricketune-web";
      description = lib.mdDoc "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "kricketune-web";
      description = lib.mdDoc "group to use for this service";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.kricketune-web = {
      description = "Frontend for kricketune";
      wantedBy = [ "multi-user.target" ];

      environment = {
        PORT = "${toString cfg.port}";
        HOST = "${cfg.host}";
        PUBLIC_API_ADDRESS = "${cfg.apiAddress}";
        HOST_HEADER = "x-forwarded-host";
        BODY_SIZE_LIMIT = "Infinity";
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart = "${cfg.package}/bin/kricketune-web";

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

    users.users = mkIf (cfg.user == "kricketune-web") {
      kricketune-web = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "kricketune-web") {
      kricketune-web = {};
    };
  };
}
