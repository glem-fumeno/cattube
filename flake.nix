{
  description = "A flake for cattube";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
      {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "cattube";
        version = "1.0.0";

        src = ./.;
        vendorHash = "sha256-/NOScTPHFGoL0EeXnDkjrKhs1sKA2gUzXqhIf0KmvXE=";
      };

      apps.${system}.default =
        let
          ld_library_path = pkgs.lib.makeLibraryPath [ pkgs.ffmpeg ];
          path = pkgs.lib.makeBinPath [ pkgs.ffmpeg ];
          prog = pkgs.writeShellScriptBin "run-cattube" ''
          export LD_LIBRARY_PATH=${ld_library_path}
          export PATH=${path}:$PATH
          exec ${self.packages.${system}.default}/bin/cattube "$@"
          '';
        in
        {
          type = "app";
          program = "${prog}/bin/run-cattube";
        };
    };
}
