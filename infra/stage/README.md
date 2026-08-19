# Staging infrastructure (NixOS)

A pluggable NixOS module for running isolated giggler-golang deployments.
Unlike `infra/local` (Docker Compose), each deployment runs as a
`systemd-nspawn` container — a fully isolated, self-contained NixOS system
(own filesystem, PID namespace, network namespace) running `webapi`,
PostgreSQL and Grafana natively. A host importing this module only runs SSH
and an nginx reverse proxy itself; the app never runs on the bare host, so
multiple deployments (e.g. `stage`, a PR-preview environment) — or this
module alongside completely unrelated services — can share one machine
without port, filesystem, or process conflicts.

This flake does **not** own a whole VPS. It exports `nixosModules.default`,
meant to be imported into a separate host-config repo/flake's own
`nixosSystem` call, alongside that host's base config (hostname, SSH,
hardware, and anything else the box runs). `example/` shows how.

## Layout

- `flake.nix` — builds the `webapi` package (`buildGoModule`) and exports
  `nixosModules.default`. Also defines a demo `nixosConfigurations.stage`
  (host base from `example/` + the module) so this repo is independently
  buildable/checkable without a separate host-config repo.
- `module.nix` — the pluggable deliverable. Declares
  `services.giggler-golang.{acmeEmail, instances}` options; for every entry
  under `instances` it adds a `systemd-nspawn` container
  (`instance.nix`) and an nginx vhost proxying to it. Touches nothing else
  on the host.
- `instance.nix` — a single deployment: `webapi` systemd service, native
  `postgresql`, native `grafana`. Runs *inside* each container; nothing in
  here is host-specific or references other instances.
- `example/` — a minimal, standalone host that consumes `module.nix`. Not
  meant to be deployed as-is to a real VPS that also does other things —
  copy its pattern (`configuration.nix`, `instances.nix`) into your actual
  host-config repo, replacing every `TODO:`.
  - `configuration.nix` — host base: users, SSH, firewall, NAT for
    container networking.
  - `hardware-configuration.nix` — placeholder; replace with the output of
    `nixos-generate-config --show-hardware-config` run on the real VPS.
  - `instances.nix` — example `services.giggler-golang` usage: one `stage`
    instance, with a commented-out second instance as a template.
- `stage.env.example` — template for the secret env file deployed per
  instance (mirrors `template.env`).

## Consuming this module from a real host-config repo

```nix
{
  inputs.giggler-golang.url = "github:<you>/giggler-golang?dir=infra/stage";

  outputs = { nixpkgs, giggler-golang, ... }: {
    nixosConfigurations.my-vps = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        giggler-golang.nixosModules.default
        ./configuration.nix          # your host's own base config
        ./hardware-configuration.nix
        {
          networking.nat = {
            enable = true;
            internalInterfaces = [ "ve-+" ];
            externalInterface = "eth0"; # your host's public interface
          };
          services.giggler-golang = {
            acmeEmail = "ops@example.com";
            instances.stage = {
              domain = "stage.example.com";
              hostAddress = "10.233.1.1";
              localAddress = "10.233.1.2";
            };
          };
        }
      ];
    };
  };
}
```

The host must open 80/443 in its firewall for `services.nginx` (enabled by
the module) and enable `networking.nat` as shown, so containers can reach
the internet. Nothing else is required — the module never touches
`networking.hostName`, users, SSH, or hardware config.

## Adding another isolated instance

Add an entry under `services.giggler-golang.instances` with a unique
attribute name, `domain`, and a `hostAddress`/`localAddress` pair on an
unused `/24` (e.g. `10.233.2.1` / `10.233.2.2` for a second instance) —
nothing else needs to change. The module derives everything (container,
Postgres, Grafana, nginx vhost) from that one entry.

## First-time VPS setup (using the example/ host)

1. Install NixOS on the VPS (e.g. via the official ISO or
   [nixos-infect](https://github.com/elitak/nixos-infect) if starting from
   another distro).
2. On the VPS, generate the real hardware config and copy it back into this
   repo, replacing the placeholder:
   ```bash
   ssh <vps> nixos-generate-config --show-hardware-config > infra/stage/example/hardware-configuration.nix
   ```
3. Find the VPS's public network interface name (`ssh <vps> ip -br link`)
   and set it as `networking.nat.externalInterface` in
   `example/configuration.nix`.
4. Fill in your SSH public key in `example/configuration.nix` and the
   domain(s) in `example/instances.nix`.
5. For each instance, copy `stage.env.example` to
   `secrets/<instance-name>.env` (gitignored) and fill in real secrets (JWT
   signing key, etc.).
6. For each instance, generate Grafana's secret key directly on the VPS
   (never commit it):
   ```bash
   ssh root@<vps-ip> 'mkdir -p /etc/giggler-golang/<instance-name> && openssl rand -hex 32 > /etc/giggler-golang/<instance-name>/grafana-secret-key'
   ```

## Deploying

```bash
# Build the webapi package locally (sanity check before deploying):
nix build ./infra/stage#webapi

# Apply the full host + all containers to the VPS over SSH:
nixos-rebuild switch --flake ./infra/stage#stage --target-host root@<vps-ip> --build-host localhost

# Ship each instance's secret env file separately — it is never part of the
# Nix store, and lives on the host at the path module.nix bind-mounts
# read-only into that container's /etc/giggler-golang:
scp secrets/stage.env root@<vps-ip>:/tmp/stage.env
ssh root@<vps-ip> 'install -m 600 -o root -g root /tmp/stage.env /etc/giggler-golang/stage/stage.env && systemctl restart container@stage'
```

Re-run `nixos-rebuild switch` after any change (it adds, updates, or removes
containers to match `services.giggler-golang.instances`); re-run the
`scp`/`restart` step for a given instance after any change to its secret env
file.

Useful container-level commands on the VPS:

```bash
machinectl list                     # running containers
machinectl shell stage              # shell into an instance
journalctl -M stage -u giggler-webapi -f   # tail webapi logs for an instance
```

## Notes

- `webapi`'s build `version` is currently pinned to `"unknown"` in
  `flake.nix`. Override it per-deploy with the project's existing
  `GIGGLER_GOLANG_VERSION` convention:
  ```bash
  nix build ./infra/stage#webapi --arg version "\"$(./scripts/export-build-version.bash && echo $GIGGLER_GOLANG_VERSION)\""
  ```
- `vendorHash` in `flake.nix` is pinned to the hash for the current
  `go.sum`; if dependencies change, `nix build` will report a hash mismatch
  with the correct value to substitute in.
- Grafana provisioning (dashboards/datasources) is not yet set up — mirrors
  the `# TODO` items already present in `infra/local/compose.yaml`.
