{
  description = "giggler-golang staging infrastructure (NixOS VPS)";

  inputs = {
    # unstable is used (rather than a stable NixOS release) because go.mod
    # currently requires Go >= 1.26.5, newer than what nixos-25.05 ships.
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };

      # The webapi build version, e.g. a git tag/commit or "dirty".
      # Overridden at build time with: nix build .#webapi --arg version '"v1.2.3"'
      version = "unknown";

      webapi = pkgs.buildGoModule {
        pname = "giggler-golang-webapi";
        inherit version;

        src = ../..;

        vendorHash = "sha256-G/W5TMFfaSugA88InqsBgoCYT8O/LqqXkWRWkUfUFUg=";

        subPackages = [ "src/cmd/webapi" ];

        ldflags = [
          "-X giggler-golang/src/shared/version.Version=${version}"
        ];

        env.CGO_ENABLED = 0;
      };
    in
    {
      packages.${system} = {
        inherit webapi;
        default = webapi;
      };

      # The pluggable deliverable: import this into a host-config
      # repo/flake's own nixosSystem alongside whatever else runs on the
      # box. It only contributes `containers.*` + an nginx vhost per entry
      # under `services.giggler-golang.instances` — see module.nix.
      nixosModules.default =
        { ... }:
        {
          imports = [ ./module.nix ];
          _module.args.webapi = webapi;
        };

      # A minimal demo host that imports nixosModules.default, so this repo
      # is independently `nix flake check`-able and the example/ directory
      # doubles as copy-paste starting point for a real host-config repo.
      nixosConfigurations.stage = nixpkgs.lib.nixosSystem {
        inherit system;
        modules = [
          self.nixosModules.default
          ./example/configuration.nix
          ./example/hardware-configuration.nix
          ./example/instances.nix
        ];
      };
    };
}
