# The pluggable contribution of this flake to a NixOS host: for every entry
# under `services.giggler-golang.instances`, declares an isolated
# systemd-nspawn container (webapi + postgresql + grafana, see instance.nix)
# and a host-level nginx vhost proxying to it. Import this module into a
# larger host-config flake alongside whatever else runs on the box — it does
# not touch networking.hostName, users, SSH, firewall base rules, or
# hardware-configuration; those remain the host config's responsibility.
#
# The host must still enable NAT for container networking, e.g.:
#   networking.nat = {
#     enable = true;
#     internalInterfaces = [ "ve-+" ];
#     externalInterface = "<public-interface>";
#   };
# and open 80/443 in its firewall for services.nginx (enabled by this module).
{
  webapi,
  lib,
  config,
  ...
}:

let
  cfg = config.services.giggler-golang;

  instanceModule = lib.types.submodule {
    options = {
      domain = lib.mkOption {
        type = lib.types.str;
        description = "Public domain this instance is served on (also its ACME cert name).";
      };
      hostAddress = lib.mkOption {
        type = lib.types.str;
        description = "Private veth address assigned to the host side of this instance's container.";
      };
      localAddress = lib.mkOption {
        type = lib.types.str;
        description = "Private veth address assigned to this instance's container.";
      };
    };
  };
in
{
  options.services.giggler-golang = {
    acmeEmail = lib.mkOption {
      type = lib.types.str;
      description = "Contact email registered with the ACME/Let's Encrypt account used for every instance's certificate.";
    };

    instances = lib.mkOption {
      type = lib.types.attrsOf instanceModule;
      default = { };
      description = ''
        One entry per isolated giggler-golang deployment sharing this host.
        Each becomes its own systemd-nspawn container with a private
        webapi + postgresql + grafana, reachable only through the nginx
        vhost this module adds for it. hostAddress/localAddress must be
        unique per instance and not overlap with any other container or
        bridge on the host.
      '';
      example = lib.literalExpression ''
        {
          stage = {
            domain = "stage.giggler.example.com";
            hostAddress = "10.233.1.1";
            localAddress = "10.233.1.2";
          };
        }
      '';
    };
  };

  config = lib.mkIf (cfg.instances != { }) {
    containers = lib.mapAttrs (name: instance: {
      autoStart = true;
      privateNetwork = true;
      hostAddress = instance.hostAddress;
      localAddress = instance.localAddress;

      # Secrets (stage.env, grafana-secret-key) are deployed straight into
      # the container's /etc, out-of-band, per-instance — see README.md.
      # Nothing under /etc/giggler-golang lives in the Nix store.
      bindMounts."/etc/giggler-golang" = {
        hostPath = "/etc/giggler-golang/${name}";
        isReadOnly = true;
      };

      config = import ./instance.nix { inherit webapi; };
    }) cfg.instances;

    security.acme = {
      acceptTerms = true;
      defaults.email = cfg.acmeEmail;
    };

    services.nginx = {
      enable = true;
      recommendedProxySettings = true;
      recommendedTlsSettings = true;

      virtualHosts = lib.mapAttrs' (
        name: instance:
        lib.nameValuePair instance.domain {
          forceSSL = true;
          enableACME = true;

          locations."/" = {
            proxyPass = "http://${instance.localAddress}:8080/";
          };
          locations."/grafana/" = {
            proxyPass = "http://${instance.localAddress}:3000/";
            proxyWebsockets = true;
          };
        }
      ) cfg.instances;
    };
  };
}
