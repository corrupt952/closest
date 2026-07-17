{
  description = "CLI tool that searches parent directories for files and returns the closest path";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "aarch64-darwin"
        "x86_64-darwin"
        "aarch64-linux"
        "x86_64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
      # Prefer the commit hash over a hand-maintained version: it can never
      # go stale, and `nix run github:corrupt952/closest` always builds main.
      version = self.shortRev or self.dirtyShortRev or "dev";
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.buildGoModule {
            pname = "closest";
            inherit version;
            src = pkgs.lib.cleanSource self;
            vendorHash = null;
            # Keep in sync with .goreleaser.yml ldflags.
            ldflags = [ "-s" "-w" "-X" "main.Version=${version}" ];
            # The finder tests walk up parent directories from the working
            # directory, which makes them sensitive to the sandbox's directory
            # layout. CI runs `go test` directly.
            doCheck = false;
            meta.mainProgram = "closest";
          };
        });

      checks = forAllSystems (system: {
        default = self.packages.${system}.default;
      });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
              goreleaser
            ];
          };
        });
    };
}
