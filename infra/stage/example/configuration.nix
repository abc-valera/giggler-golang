{ pkgs, ... }:

# Example host base config: everything a machine needs before it can import
# ../module.nix — SSH, firewall, NAT for containers. Copy this (and
# hardware-configuration.nix) into your actual host-config repo/flake as a
# starting point; giggler-golang's flake itself only consumes it to produce
# the demo `nixosConfigurations.stage` used for local `nix flake check`.
{
  # TODO: replace with the actual staging VPS hostname
  networking.hostName = "giggler-golang-stage";

  time.timeZone = "UTC";

  nix.settings.experimental-features = [
    "nix-command"
    "flakes"
  ];

  # TODO: replace with your own SSH public key(s)
  users.users.deploy = {
    isNormalUser = true;
    extraGroups = [ "wheel" ];
    openssh.authorizedKeys.keys = [
      # "ssh-ed25519 AAAA... deploy@giggler-golang"
    ];
  };
  security.sudo.wheelNeedsPassword = false;

  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = false;
      PermitRootLogin = "no";
    };
  };

  # Only the host's SSH and nginx (enabled by ../module.nix) need public
  # ports; every instance's webapi/grafana stay on a container-private
  # address and are reached only via the host's reverse proxy.
  networking.firewall = {
    enable = true;
    allowedTCPPorts = [
      22 # SSH
      80 # HTTP (redirects to HTTPS)
      443 # HTTPS
    ];
  };

  # Required for containers.<name>.privateNetwork to route to the outside
  # world through the host.
  networking.nat = {
    enable = true;
    internalInterfaces = [ "ve-+" ];
    externalInterface = "eth0"; # TODO: replace with the VPS's actual public interface name (`ip -br link` on the box)
  };

  environment.systemPackages = with pkgs; [
    curl
    git
    htop
  ];

  # Keep staging as close to prod as practical: automatic upgrades disabled,
  # deployments are performed explicitly via `nixos-rebuild --flake .#stage`.
  system.autoUpgrade.enable = false;

  system.stateVersion = "25.05";
}
