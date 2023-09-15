{inputs, ...}: {
  perSystem = {pkgs, ...}: {
    packages.default = pkgs.buildGoModule rec {
      pname = "gogal";
      version = "0.1";
      pwd = ./.;
      src = let
        # Set this to `true` in order to show all of the source files
        # that will be included in the module build.
        show-trace = true;
        nix-filter = inputs.nix-filter.lib;
        # Match paths with the given extension

        matchExts = map (ext: nix-filter.matchExt ext);
        source-files =
          nix-filter.filter
          {
            name = pname;
            root = ./.;
            exclude =
              [./assets ./nix]
              ++ matchExts [
                "toml"
                "json"
                "nix"
                "lock"
              ];
          };
      in (
        if show-trace
        then pkgs.lib.sources.trace source-files
        else source-files
      );

      subPackages = ["cmd/cli"];
      vendorSha256 = "sha256-4A5j3N+H0TMYWFiVZUr74m3kK5kAcfZuXTDM/j+H024=";

      postInstall = ''
        mkdir -p $out/bin
        mv $out/bin/cli $out/bin/gogal
      '';
    };
  };
}
