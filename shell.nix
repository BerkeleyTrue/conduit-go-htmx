{...}: {
  perSystem = {
    pkgs,
    config,
    ...
  }: let
    generate-queries = pkgs.writeShellScriptBin "generate-queries" ''
      echo "generating"
      ${pkgs.sqlc}/bin/sqlc generate;
    '';
    # update the vendorSha256 of the default package
    update-vendor-sha = pkgs.writeShellScriptBin "update-vendor-sha" ''
      set -exuo pipefail

      failedbuild=$(nix build --impure 2>&1 || true)
      # echo "$failedbuild"
      checksum=$(echo "$failedbuild" | awk '/got:.*sha256/ { print $2 }')
      echo -n "\n\nchecksum: $checksum"
      # do nothing if no checksum was found
      if [ -z "$checksum" ]; then
        exit 0
      fi
      sed -i -e "s|vendorSha256 = \".*\"|vendorSha256 = \"$checksum\"|" ./default.nix
    '';

    watch-sql = pkgs.writeShellScriptBin "watch-sql" ''
      ${pkgs.fd}/bin/fd -e sql | ${pkgs.entr}/bin/entr ${generate-queries}/bin/generate-queries
    '';

    watch-compile = pkgs.writeShellScriptBin "watch-compile" ''
      ${pkgs.concurrently}/bin/concurrently -n "air,sqlc" "${pkgs.air}/bin/air" "${watch-sql}/bin/watch-sql"
    '';

    watch-tests = pkgs.writeShellScriptBin "watch-tests" ''
      ${pkgs.ginkgo}/bin/ginkgo watch -r -p
    '';
  in {
    boulder = {
      commands = [
        {
          exec = generate-queries;
          description = "generate sqlc queries";
        }
        {
          exec = update-vendor-sha;
          description = "update the vendorSha256 of the default package";
        }
        {
          exec = watch-compile;
          description = "watch go files for changes and recompile";
        }
        {
          exec = watch-tests;
          description = "watch go files for changes and re-run tests";
        }
        {
          exec = watch-sql;
          description = "watch sql files for changes and re-run sqlc";
        }
      ];
    };

    devShells.default = pkgs.mkShell {
      name = "go-conduit-htmx";
      inputsFrom = [config.boulder.devShell];
      packages = with pkgs;
        [
          go
          gopls # language server
          gotools
          go-tools
          air # live reload
          ginkgo # testing framework
          golines # go formatter
          sqlc # sql code generator

          # prettier
          nodejs_18
        ]
        ++ (with nodePackages; [
          # editor
          vscode-langservers-extracted # html/css language server
          typescript-language-server # typescript language server
          sql-formatter # sql formatter
        ]);
    };
  };
}
