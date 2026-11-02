# Serverus

A cross-platform app for managing remote servers, built with
[giode](https://github.com/mehran9675/giode) and its
[CLI](https://github.com/mehran9675/giode-cli). It runs natively (no webview) on Windows,
macOS and Linux, and builds for Android.

## Features

**Servers**
- Any number of servers in a searchable sidebar, grouped, with tags, colors and live
  connection status. Right-click a server to connect, reconnect, disconnect, edit,
  duplicate or delete it.
- Connection types:
  - **SSH**: password, private key (with passphrase) or SSH agent (ssh-agent, the Windows
    OpenSSH agent, 1Password, Secretive). Keyboard-interactive and one-time codes are
    supported, as is a **jump host** (like `ssh -J`, chains allowed).
  - **This computer**: monitor and manage the local machine on any OS.
  - **Telnet**: a terminal for switches, routers and embedded boards.
  - **FTP / FTPS**: file management with explicit TLS.
- "Test connection" in the editor before saving.
- **Security**: passwords and passphrases are stored in the OS keychain (Keychain,
  Credential Manager/DPAPI, Secret Service), never in the workspace file. Host keys are
  checked against `~/.ssh/known_hosts` and Serverus's own file; new keys are confirmed by
  fingerprint and pinned, and changed keys show a warning.
- Keep-alives detect dead connections, and you can be notified when a server drops.

**Monitoring** (Linux servers over SSH with nothing to install; this computer on any OS)
- CPU (total, I/O wait, per core), memory and swap, load, uptime, process count.
- Every mounted file system, disk read/write throughput.
- Network throughput per interface and in total.
- **Temperatures** from hwmon and thermal zones, colored by their high and critical
  thresholds.
- Live charts of the last few minutes, and the CPU of each connected server in the
  sidebar.

**Terminal**
- Tabs of xterm-256color shells with true color, bold/italic/underline, wide characters,
  the alternate screen (vim, htop), mouse-wheel scrollback, selection with copy,
  bracketed paste, resizing, and titles set by the shell.
- Ctrl+Shift+C / Ctrl+Shift+V (Cmd+C / Cmd+V on macOS), right-click to copy or paste,
  Shift+PageUp/PageDown to scroll.
- Insert saved snippets into the terminal from the toolbar.

**Files** (SFTP, FTP or local)
- Browse with breadcrumbs, back/up/home, filter, sorting, and hidden files on or off.
- Upload with the file picker or by **dropping files and folders on the window**;
  download files or whole folders. Transfers run in the background, three at a time,
  with progress, speed and cancel in the Transfers panel.
- Edit text files in place, create files and folders, rename, change permissions,
  delete, and copy paths.

**Administration**
- **Processes**: sortable, filterable, refreshed every 3 seconds; terminate, kill or
  send SIGHUP.
- **Services** (systemd): running/failed/stopped filters; start, stop, restart, reload,
  enable/disable at boot; logs (`journalctl`) and status.
- **Docker**: containers with state, ports and live CPU/memory; start, stop, restart,
  pause, remove; logs.
- Privileged actions use `sudo -n` when the user isn't root.

**Tunnels**
- Local port forwarding through any SSH server (`ssh -L`), opened by hand or
  automatically on connect, with live connection counts.

**Snippets**
- Saved commands: run on the current server and see the output, type into its terminal,
  or **run on several servers at once** and compare the results side by side.

**Interface**
- A tray icon to open the window, connect and disconnect servers, and quit; minimize
  to the tray, and on Windows and macOS keep running there when the window closes.
- Keyboard shortcuts (listed in Settings): search servers, add a server, settings, a
  new terminal tab, and Ctrl/⌘+1–8 to switch pages.
- Dark and light themes (or follow the system), a custom title bar, and a layout that
  adapts down to phone width (the server list becomes a drawer).

## Run it

Requirements: Go 1.26+, and on Linux Gio's development packages
(see `giode doctor`).

```bash
go install github.com/mehran9675/giode-cli/cmd/giode@latest
```

```bash
giode run
```

```bash
giode build              # every desktop platform in giode.toml this machine can build
giode build windows-installer
giode build android
```

## Releases

Publishing a GitHub release builds every platform and attaches the packages to it
(`.github/workflows/release.yml`). Tag the release `major.minor.patch.versionCode`, for
example `1.0.0.1`: the first three parts are the version shown everywhere, the fourth is
Android's versionCode, which must be at least 1 and grow with every Play upload. The
workflow is written by `giode ci` from `[ci]` in `giode.toml`; change it there and run
`giode ci -f`. Signing uses these repository secrets when present: `ANDROID_KEYSTORE_BASE64` and
`GIODE_KEYSTORE_PASS`, `WINDOWS_CERT_BASE64` and `GIODE_WINDOWS_CERT_PASS`,
`MACOS_CERT_BASE64`, `MACOS_CERT_PASS` and `GIODE_NOTARY_PASS`.

Data lives in the platform's configuration folder (`~/.config/Serverus`,
`%AppData%\Serverus`, `~/Library/Application Support/Serverus`); set `SERVERUS_DATA_DIR` to
use another folder (portable installs, testing).

## Tests

```bash
go test ./...
```

The SSH tests start an in-process SSH server on the loopback interface and cover
password authentication, host key pinning, commands, the monitoring probe, PTY shells,
SFTP, port forwarding and jump hosts. The probe parsers, file transfers, tunnels, storage
and formatting have unit tests.

## Layout

```
main.go                      entry point
assets/                      icon, fonts (Inter, JetBrains Mono), Lucide icons
internal/
  domain/                    servers, protocols, capabilities, snippets, tunnels, settings
  storage/                   workspace JSON, keychain secrets, known_hosts
  transport/                 Conn and capability interfaces (Terminal, Executor,
    sshx/ local/ telnet/       FileSystem, Forwarder, Pinger), one package per protocol,
    ftpx/ connector/           and the connector that picks one
  monitor/                   probe script, parsers, rates, history
  ops/                       processes, systemd services, Docker
  session/                   one live connection per server, metrics collection
  transfer/                  background upload/download queue
  tunnel/                    local port forwards
  platform/                  desktop or phone, macOS or not: decided by build tags
  bootstrap/                 composition root: wires everything and runs the window
  ui/
    theme/                   palettes, type scale, spacing, fonts, icons
    widgets/                 reusable components, one per file
    terminal/                the terminal component (VT emulator + rendering + input)
    dialogs/                 server editor, confirmations, credential and host key prompts
    shell/                   title bar, sidebar, server workspace and page tabs
    pages/                   overview, terminal, files, processes, services, docker,
                             tunnels, snippets, settings
    env/ async/ format/      dependency container, background loading, formatting
```

Each protocol implements only the capabilities it has; pages are shown for the
capabilities a server's protocol provides, so adding a protocol means adding a package
under `transport/` and a case in `connector`.


## Icon

`scripts/make_icon.py` (Pillow) draws `assets/icon.png` and Android's adaptive icon
layers, `assets/android-foreground.png` and `assets/android-background.png`.

## Licenses

Inter and JetBrains Mono are under the SIL Open Font License (`assets/fonts`); the icons
are Lucide, under the ISC license (`assets/icons/LICENSE.txt`).
