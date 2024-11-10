{
  description = "A basic flake to run a Go + K8S Project";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        stdenv = pkgs.stdenv;

        inputs = with pkgs; [
          cobra-cli
          docker
          go
          kubectl
          minikube
        ]
        ++ pkgs.lib.optional stdenv.isLinux pkgs.inotify-tools
        ++ pkgs.lib.optional stdenv.isDarwin pkgs.terminal-notifier
        ++ pkgs.lib.optionals stdenv.isDarwin (with pkgs.darwin.apple_sdk.frameworks; [ CoreFoundation CoreServices ]);
      in {
        devShells.default = pkgs.mkShell {
          name = "local-ratelimiter";
          packages = inputs;

          shellHook = ''
            export GIT_SSH_COMMAND="ssh -F ~/.ssh/config"
            echo "Nix shell environment configured for Go + K8S development."
          '';
        };
      });
}