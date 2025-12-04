{
  lib,
  buildGoModule,
  swww,
}:
buildGoModule rec {
  pname = "swwwitch";
  version = "0.1.0";

  src = ./.;

  vendorHash = null;

  nativeBuildInputs = [swww];

  ldflags = [
    "-s"
    "-w"
  ];

  meta = with lib; {
    description = "A CLI tool for switching wallpapers with swww on Wayland";
    homepage = "https://github.com/Danielbook/swwwitch";
    license = licenses.mit;
    maintainers = [];
    platforms = platforms.linux;
    mainProgram = "swwwitch";
  };
}
