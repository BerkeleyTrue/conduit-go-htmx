{
  lib,
  flake-parts-lib,
  ...
}: let
  inherit (flake-parts-lib) mkPerSystemOption;
  inherit (lib) mkOption types;
in {
  options.perSystem = mkPerSystemOption ({
    config,
    pkgs,
    ...
  }: {
    options = let
      scriptSubmodule = types.submodule {
        options = {
          description = mkOption {
            type = types.nullOr types.str;
            description = ''
              A description of what this script does.

              This will be displayed in the banner and help menu.
            '';
            default = null;
          };
          category = mkOption {
            type = types.str;
            description = "The category under which this script will be grouped in the banner.";
            default = "Commands";
          };
          exec = mkOption {
            type = types.package;
            description = "The package to add to the dev shell.";
          };
        };
      };
    in {
      boulder.commands = mkOption {
        type = types.listOf scriptSubmodule;
        default = [];
        description = "List of scripts to be added to the shell";
      };

      boulder.devShell = mkOption {
        type = types.package;
        readOnly = true;
        description = "The dev shell to use for development";
      };
    };

    config = let
      commandSpecs = config.boulder.commands; # might need to make a list
      byGroup = lib.groupBy (x: x.category) commandSpecs;
      devShell-help = pkgs.writeShellApplication {
        name = "devshell-help";
        text = ''
          echo -e "Available commands:\n"
          ${
            lib.concatStringsSep "echo;"
            (lib.mapAttrsToList (
                category: commands:
                  "echo -e '   >==> ${category}\n'; "
                  + "echo '"
                  + lib.concatStringsSep "\n"
                  (
                    map (
                      spec: let
                        name = builtins.baseNameOf (lib.getExe spec.exec);
                        desc =
                          if (spec.description != null)
                          then spec.description
                          else "No description provided";
                      in
                        name + "\t: " + desc
                    )
                    commands
                  )
                  + "' | ${lib.getExe pkgs.unixtools.column} -t -s ''$'\t'; "
              )
              byGroup)
          }
        '';
      };

      shellHook = ''
        function menu () {
          echo
          echo -e "\033[1;34m>==> ️  '$name' shell\n\033[0m"
          ${devShell-help}/bin/${devShell-help.name}
          echo
          echo "(Run '${devShell-help.name}' to display this menu again)"
          echo
        }

        menu
      '';
    in {
      boulder.devShell = pkgs.mkShell {
        nativeBuildInputs = [devShell-help] ++ (map (c: c.exec) commandSpecs);
        inherit shellHook;
      };
    };
  });

  _file = ../boulder;
}
