{
  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/master"; };
  
  outputs = { self, nixpkgs, ... }@inputs: 
  let
    supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    
    pkgsFor = system: import nixpkgs {
      inherit system;
      overlays = [ self.overlay ];
    };
  in {
    overlay = import ./overlay.nix;
    
    packages = forAllSystems (system: {
      gonwatch = (pkgsFor system).gonwatch;
      default = (pkgsFor system).gonwatch;
    });
  };
}
