{
  lib,
  buildGoModule,
  fetchFromGitHub,
  awww,
}:
buildGoModule rec {
  pname = "swwwitch";
  version = "1.0.0";

  src = fetchFromGitHub {
    owner = "Danielbook";
    repo = "swwwitch";
    rev = "v${version}";
    hash = "sha256-3c7JV2zGj2ViibivDZnTnKgbVTpgkTJtGy51UBxFzh4=";
  };

  vendorHash = null; # No Go dependencies

  nativeBuildInputs = [awww];

  ldflags = [
    "-s"
    "-w"
  ];

  meta = with lib; {
    description = "CLI wallpaper switcher for awww on Wayland";
    homepage = "https://github.com/Danielbook/swwwitch";
    license = licenses.mit;
    maintainers = [{
      name = "Daniel Book";
      email = "daniel@bookorjeman.se";
      github = "Danielbook";
      githubId = 6060731;
    }];
    platforms = platforms.linux;
    mainProgram = "swwwitch";
  };
}
