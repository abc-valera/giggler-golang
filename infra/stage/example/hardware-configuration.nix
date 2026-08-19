{ lib, modulesPath, ... }:

# Placeholder hardware configuration.
#
# Replace the contents of this file with the one generated on the actual VPS via:
#   nixos-generate-config --show-hardware-config
#
# This keeps disk layout, filesystem UUIDs and kernel modules specific to the
# real machine out of source control assumptions made here.
{
  imports = [ (modulesPath + "/profiles/qemu-guest.nix") ];

  boot.loader.grub = {
    enable = true;
    device = "/dev/sda";
  };

  fileSystems."/" = {
    device = "/dev/disk/by-label/nixos";
    fsType = "ext4";
  };

  boot.initrd.availableKernelModules = [
    "ata_piix"
    "uhci_hcd"
    "virtio_pci"
    "virtio_scsi"
    "sd_mod"
    "sr_mod"
  ];

  nixpkgs.hostPlatform = lib.mkDefault "x86_64-linux";
}
