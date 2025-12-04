# swwwitch

CLI wallpaper switcher for [swww](https://github.com/LGFae/swww) on Wayland.

## Features

- Category-based organization
- Random selection from all or specific categories
- Fast Go binary with no dependencies
- Configurable transitions via environment variables

## Installation

### Nix Flakes

```bash
# Try without installing
nix run github:Danielbook/swwwitch -- --help

# Add to flake inputs
{
  inputs.swwwitch.url = "github:Danielbook/swwwitch";
}
```

### Build from source

```bash
go build
```

## Setup

Organize wallpapers in `~/Pictures/wallpapers/` with category subdirectories:

```
wallpapers/
├── nature/
├── abstract/
└── minimal/
```

Supports: `.jpg`, `.png`, `.gif`, `.webp`

## Usage

```bash
swwwitch --list                  # List categories
swwwitch --random                # Random from all
swwwitch nature                  # Random from category
swwwitch --set path/to/image.jpg # Set specific wallpaper
swwwitch --dir ~/Wallpapers      # Use custom directory
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

## Integration

**Hyprland:**
```
bind = SUPER, W, exec, swwwitch --random
```

**Auto-rotate (systemd):**
```ini
# ~/.config/systemd/user/swwwitch.timer
[Timer]
OnBootSec=1min
OnUnitActiveSec=30min
```

## Requirements

- [swww](https://github.com/LGFae/swww)
- Wayland compositor

## License

MIT
