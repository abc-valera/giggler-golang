# A single giggler-golang deployment (webapi + postgresql + grafana), as a
# nested NixOS module meant to be used as a container's `config`. Each
# instance is fully self-contained: its own systemd units, its own Postgres
# cluster, its own Grafana — nothing here references a name/port that would
# collide with another instance running in a sibling container.
{ webapi }:
{ ... }:

{
  system.stateVersion = "25.05";

  users.users.giggler-webapi = {
    isSystemUser = true;
    group = "giggler-webapi";
  };
  users.groups.giggler-webapi = { };

  systemd.services.giggler-webapi = {
    description = "giggler-golang webapi";
    wantedBy = [ "multi-user.target" ];
    after = [
      "network-online.target"
      "postgresql.service"
    ];
    wants = [ "network-online.target" ];
    requires = [ "postgresql.service" ];

    # Populated from the untracked /etc/giggler-golang/stage.env, deployed
    # out-of-band per-container. See infra/stage/README.md.
    serviceConfig = {
      ExecStart = "${webapi}/bin/webapi";
      EnvironmentFile = "/etc/giggler-golang/stage.env";
      Restart = "on-failure";
      RestartSec = "2s";

      User = "giggler-webapi";
      Group = "giggler-webapi";
      StateDirectory = "giggler-webapi";
      RuntimeDirectory = "giggler-webapi";

      NoNewPrivileges = true;
      ProtectSystem = "strict";
      ProtectHome = true;
      PrivateTmp = true;
    };
  };

  services.postgresql = {
    enable = true;

    # `ensureDBOwnership` creates/owns a database named after the user, so the
    # database is also named `giggler-webapi` (set POSTGRES_DB accordingly).
    ensureDatabases = [ "giggler-webapi" ];
    ensureUsers = [
      {
        # Name matches the `giggler-webapi` systemd service's OS user above
        # so peer auth below maps it to this Postgres role.
        name = "giggler-webapi";
        ensureDBOwnership = true;
      }
    ];

    # The webapi connects over the local unix socket via peer auth so
    # POSTGRES_PASSWORD in stage.env is unused for the "local" auth method
    # below; keep it set anyway for parity with the local/prod .env files.
    authentication = ''
      local all all peer
    '';

    settings = {
      max_connections = 100;
      shared_buffers = "256MB";
    };
  };

  services.grafana = {
    enable = true;
    settings = {
      server = {
        http_addr = "0.0.0.0"; # only reachable via the container's private network
        http_port = 3000;
      };
      log.level = "error";
      # Read from a file deployed out-of-band, not the Nix store. Generate
      # once with: openssl rand -hex 32 > /etc/giggler-golang/grafana-secret-key
      security.secret_key = "$__file{/etc/giggler-golang/grafana-secret-key}";
    };
  };

  networking.firewall.allowedTCPPorts = [
    8080 # webapi, proxied by the host's nginx
    3000 # grafana, proxied by the host's nginx
  ];
}
