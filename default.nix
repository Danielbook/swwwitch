{
  lib,
  buildGoModule,
  fetchFromGitHub,
  swww,
}:
buildGoModule rec {
  pname = "swwwitch";
  version = "1.0.0";

  src = fetchFromGitHub {
    owner = "Danielbook";
    repo = "swwwitch";
    rev = "v${version}";
    hash = "sha256-giKDk0UzNJRvXAwQCZRel9uAcTg5SuMtWj6qLeUF6qg=";
  };

  vendorHash = null; # No Go dependencies

  nativeBuildInputs = [swww];

  ldflags = [
    "-s"
    "-w"
  ];

  meta = with lib; {
    description = "CLI wallpaper switcher for swww on Wayland";
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
