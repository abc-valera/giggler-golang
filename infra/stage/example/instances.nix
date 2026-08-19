# Example usage of ../module.nix's options. Copy this pattern into your
# actual host-config repo/flake — add one attrset per isolated deployment;
# the module derives the container, its private network, and the nginx
# vhost from each entry automatically.
{
  services.giggler-golang = {
    acmeEmail = "ops@giggler.example.com"; # TODO: replace with a real contact address

    instances = {
      stage = {
        domain = "stage.giggler.example.com"; # TODO: replace with the real staging domain
        hostAddress = "10.233.1.1";
        localAddress = "10.233.1.2";
      };

      # stage-pr-42 = {
      #   domain = "pr-42.giggler.example.com";
      #   hostAddress = "10.233.2.1";
      #   localAddress = "10.233.2.2";
      # };
    };
  };
}
