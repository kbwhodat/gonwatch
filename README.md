# Gonwatch (Go and Watch)

A golang tui app that allows you to watch whatever.

### What can be watched?
- Movies
- TV Series
- Anime
- Live Sports

For movies and tv series it comes with subtitles inlcuded. If there is none available, a call is made to [opensubtitles](https://www.opensubtitles.org/en/search/subs). So, you might have to ensure the audio and captioning are in sync (rarely).

#### Features
If there are anything you think should be added let it be known.

#### Python runtime
The app runs embedded Python scripts for scraping. The installer creates a local venv at `~/.local/share/gonwatch/venv` and the app will use that by default.

If you want to override the interpreter, set `GONWATCH_PYTHON` to a full path (e.g. a custom venv Python).

#### Issues
If you face any issues using this, let it be known so it can be fixed.
