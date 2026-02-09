self: super: let
  isDarwin = super.stdenv.isDarwin;

  # Create Python environment with all required packages
  pythonEnv = super.python313.withPackages (ps: with ps; let
    cdp-socket = ps.buildPythonPackage rec {
      pname = "cdp-socket";
      version = "1.2.8";
      format = "setuptools";
      build-system = [ setuptools ];
      src = super.fetchurl {
        url = "https://files.pythonhosted.org/packages/7d/28/58812797e54fb8cf22bff61125e5a7d2763de1a86855549ecc417bdd06d5/cdp-socket-1.2.8.tar.gz";
        sha256 = "d8a3d55883205c7c45c05292cf5ef5a5c74534873e369e258e61213cce15be1a";
      };
      propagatedBuildInputs = [ websockets ];
      doCheck = false;
    };

    selenium-driverless = ps.buildPythonPackage rec {
      pname = "selenium-driverless";
      version = "1.9.4";
      format = "setuptools";
      build-system = [ setuptools ];
      src = super.fetchurl {
        url = "https://files.pythonhosted.org/packages/5e/92/3fcf637eebbc334543de61b319c4f00d01526053edf33c2f25aa08f05c13/selenium_driverless-1.9.4.tar.gz";
        sha256 = "151ccf57d399691ec4e943a941a496dbe575d0154a520cc2eca988ebe5d07a76";
      };
      propagatedBuildInputs = [ cdp-socket websockets ];
      doCheck = false;
    };
  in [
    requests
    numpy
    matplotlib
    scipy
    platformdirs
    jsondiff
    orjson
    beautifulsoup4
    langdetect
    tls-client
    selenium
    selenium-driverless
    cdp-socket
    websockets
    aiofiles
    aiohttp
  ]);

  linuxOnlyBuildInputs = super.lib.optionals (!isDarwin) [ super.chromium ];
  wrapperPathPackages = [ pythonEnv super.mpv super.nodejs ]
    ++ super.lib.optionals (!isDarwin) [ super.chromium ];

in {
  gonwatch = super.buildGoModule rec {
    pname = "gonwatch";
    version = "1.0.0";

    src = ./.;

    vendorHash = "sha256-wwiz0a1IHQdT5dR2CtxaQS8QY3xNJ5rQ6EdnjNj6iHA=";

    nativeBuildInputs = with super; [ makeWrapper ];

    buildInputs = with super; [
      mpv
      nodejs
    ] ++ linuxOnlyBuildInputs;

    postInstall = ''
      wrapProgram $out/bin/gonwatch \
        --prefix PATH : ${super.lib.makeBinPath wrapperPathPackages}
    '';

    meta = with super.lib; {
      description = "A TUI app to watch movies, TV series, anime, and live sports";
      homepage = "https://github.com/kbwhodat/gonwatch";
      license = licenses.asl20;
      maintainers = [ ];
      platforms = platforms.unix;
    };
  };
}
