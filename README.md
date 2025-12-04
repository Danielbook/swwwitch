# swwwitch

A fast and elegant CLI tool for switching wallpapers with [swww](https://github.com/LGFae/swww) on Wayland.

## Features

- 🎨 **Category-based organization** - Organize wallpapers into categories
- 🎲 **Random selection** - Pick random wallpapers from all or specific categories
- ⚡ **Fast** - Written in Go, single binary with no dependencies
- 🔧 **Configurable** - Environment variables for customization
- 🌊 **Smooth transitions** - Configurable transition types and durations
- 📦 **Nix flake** - Easy installation for NixOS users

## Installation

### Using Nix Flakes

Add to your `flake.nix`:

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    swwwitch.url = "github:Danielbook/swwwitch";
  };

  outputs = { nixpkgs, swwwitch, ... }: {
    # NixOS configuration
    nixosConfigurations.yourhost = nixpkgs.lib.nixosSystem {
      modules = [{
        environment.systemPackages = [
          swwwitch.packages.x86_64-linux.default
        ];
      }];
    };

    # Home Manager configuration
    homeConfigurations.youruser = home-manager.lib.homeManagerConfiguration {
      modules = [{
        home.packages = [
          swwwitch.packages.x86_64-linux.default
        ];
      }];
    };
  };
}
```

### Try it without installing

```bash
nix run github:Danielbook/swwwitch -- --help
```

### Build from source

```bash
git clone https://github.com/Danielbook/swwwitch.git
cd swwwitch
go build
```

## Setup

### Directory Structure

Create a wallpapers directory with subdirectories for categories:

```
~/Pictures/wallpapers/
├── nature/
│   ├── mountain1.jpg
│   ├── forest2.png
│   └── lake3.jpg
├── abstract/
│   ├── colors1.jpg
│   └── geometric2.png
├── minimal/
│   └── simple.png
└── anime/
    └── landscape.jpg
```

### Supported Formats

- `.jpg` / `.jpeg`
- `.png`
- `.gif`
- `.webp`

## Usage

### Basic Commands

```bash
# List all available categories
swwwitch --list

# Set random wallpaper from all categories
swwwitch --random

# Set random wallpaper from specific category
swwwitch --random nature
swwwitch nature  # shorthand

# Set specific wallpaper
swwwitch --set ~/Pictures/wallpapers/nature/mountain1.jpg
swwwitch ~/Pictures/my-wallpaper.jpg  # shorthand

# Use custom wallpapers directory
swwwitch --dir ~/Wallpapers --random
```

### Configuration

Configure via environment variables:

```bash
# Set custom wallpapers directory
export SWWWITCH_WALLPAPERS="$HOME/Wallpapers"

# Change transition type (fade, simple, wipe, grow, etc.)
export SWWWITCH_TRANSITION_TYPE="wipe"

# Change transition duration (in seconds)
export SWWWITCH_TRANSITION_DURATION="2"
```

### Integration Examples

#### Hyprland Keybinding

Add to your `hyprland.conf`:

```
bind = SUPER, W, exec, swwwitch --random
bind = SUPER SHIFT, W, exec, swwwitch --list | fuzzel --dmenu | xargs swwwitch
```

#### Systemd Timer (Auto-rotate wallpapers)

Create `~/.config/systemd/user/swwwitch.timer`:

```ini
[Unit]
Description=Rotate wallpaper every 30 minutes

[Timer]
OnBootSec=1min
OnUnitActiveSec=30min

[Install]
WantedBy=timers.target
```

Create `~/.config/systemd/user/swwwitch.service`:

```ini
[Unit]
Description=Switch wallpaper

[Service]
Type=oneshot
ExecStart=/path/to/swwwitch --random
```

Enable:
```bash
systemctl --user enable --now swwwitch.timer
```

#### With Rofi/Fuzzel Menu

```bash
# Select category interactively
category=$(swwwitch --list | tail -n +2 | awk '{print $1}' | fuzzel --dmenu)
swwwitch "$category"
```

## Command Reference

```
Options:
  -l, --list              List available categories
  -r, --random [CAT]      Set random wallpaper from category (or all if no category)
  -s, --set FILE          Set specific wallpaper file
  -c, --category CAT      Set random wallpaper from specific category
  -d, --dir DIR           Use custom wallpapers directory
  -v, --version           Show version
  -h, --help              Show this help
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SWWWITCH_WALLPAPERS` | Wallpapers directory | `~/Pictures/wallpapers` |
| `SWWWITCH_TRANSITION_TYPE` | Transition type | `fade` |
| `SWWWITCH_TRANSITION_DURATION` | Transition duration (seconds) | `1` |

## Requirements

- [swww](https://github.com/LGFae/swww) - Wayland wallpaper daemon
- Wayland compositor (Hyprland, Sway, etc.)

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

- Inspired by the [dharmx/walls](https://github.com/dharmx/walls) wallpaper collection
- Built with [swww](https://github.com/LGFae/swww) by LGFae

## See Also

- [swww](https://github.com/LGFae/swww) - The wallpaper daemon
- [hyprpaper](https://github.com/hyprwm/hyprpaper) - Alternative wallpaper daemon for Hyprland
- [wpaperd](https://github.com/danyspin97/wpaperd) - Wallpaper daemon with automatic rotation
