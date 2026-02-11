# Gonwatch (Go and Watch)

A golang tui app that allows you to watch whatever.

### What can be watched?
- Movies
- TV Series
- Anime
- Live Sports

For movies and tv series it comes with subtitles included. Currently, I'm only fetching for english subtitles. If you want other languages let me know and I can add it.

### How to install this bad boy?
```bash
    curl -fsSL https://raw.githubusercontent.com/kbwhodat/gonwatch/main/install.sh | bash
```
FYI: For the python package dependencies I'm using venv, so no need to worry about it colliding with your system packages.

Works for multiple Linux distros and Macos. Sorry Window users.

#### If you have nix and just want to test it out:
```bash
    nix run githb:kbwhodat/gonwatch
```

#### If you use nix flakes (recommended):
```
    inputs.gonwatch.url = "github:kbwhodat/gonwatch/main";
```

#### Features
If there are anything you think should be added let it be known.

#### Issues
If you face any issues using this, let it be known so it can be fixed.
