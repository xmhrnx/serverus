# Releasing Serverus

Publishing a GitHub release with a tag like `1.0.0.1` runs
`.github/workflows/release.yml`. The tag is `major.minor.patch.build`: the first
three parts are the version shown everywhere, the fourth is the build number
(Android's versionCode, Apple's CFBundleVersion), which must be at least 1 and grow
with every store upload. The workflow tests, builds every package and attaches
them to the release. It is written by the Giode CLI from `giode.toml`: after
changing `[ci]` or a store section, run `giode ci -f` instead of editing it, then
redo the two hand edits noted above `[ci]` (pin the CLI to `@v0.1.35-alpha.1`
and use `go-version-file: go.mod`).

| Where | What the release gets | Needs |
| --- | --- | --- |
| GitHub, Linux | `.deb` and `.rpm`, x64 | nothing |
| GitHub, Windows | `setup.exe` installer, x64 | nothing; a code signing certificate stops SmartScreen warnings |
| GitHub, macOS | a `.dmg` | a Developer ID certificate, or macOS refuses to open it |
| Mac App Store | a `.pkg` to upload with Transporter | an Apple Developer account and the store certificates |
| Microsoft Store | an `.msix` to upload in Partner Center | a Partner Center account and the reserved name |
| Google Play | an `.aab` to upload in the Play Console (and a signed `.apk` for GitHub) | a Play Console account and the upload keystore |

There is no iOS app: Serverus doesn't target iOS (`app.platforms`), so there is no
App Store (iOS) package.

## Repository secrets

Set these in GitHub under Settings → Secrets and variables → Actions. Each
`*_BASE64` secret is the file encoded with `base64 -i file | pbcopy` on macOS,
or `base64 -w0 file` on Linux. Without a secret, that step is skipped (or the
package is left unsigned) and the rest of the release still builds.

| Secret | What it is |
| --- | --- |
| `ANDROID_KEYSTORE_BASE64`, `GIODE_KEYSTORE_PASS` | The upload keystore made by `giode keystore`, and its password: signs the `.aab` and `.apk` |
| `MACOS_CERT_BASE64`, `MACOS_CERT_PASS` | The "Developer ID Application" certificate as a `.p12`, and its password: signs the `.dmg` |
| `GIODE_NOTARY_ID`, `GIODE_NOTARY_TEAM_ID`, `GIODE_NOTARY_PASS` | Your Apple ID, team ID and an app-specific password: notarizes the `.dmg` |
| `MACOS_STORE_CERT_BASE64`, `MACOS_STORE_CERT_PASS` | The "Apple Distribution" certificate as a `.p12` |
| `MACOS_INSTALLER_CERT_BASE64`, `MACOS_INSTALLER_CERT_PASS` | The "Mac Installer Distribution" certificate as a `.p12` |
| `MACOS_PROFILE_BASE64` | The Mac App Store provisioning profile (`.provisionprofile`) |
| `WINDOWS_CERT_BASE64`, `GIODE_WINDOWS_CERT_PASS` | An Authenticode `.pfx` and its password: signs the Windows installer (optional) |

Keep the keystore, certificates and profiles out of git (`*.keystore` is in
`.gitignore`; the workflow writes the profile to a temporary folder).

## Google Play

1. Open a Google Play Console developer account and create the app "Serverus"
   (app, free). Its package name is `com.xmhrnx.serverus` (`app.id`), which can
   never change.
2. Make the upload key once, and keep a backup of it and its password somewhere
   safe; losing it means asking Google to reset the upload key:

   ```bash
   giode keystore            # writes release.keystore, password in $GIODE_KEYSTORE_PASS
   base64 -w0 release.keystore   # → ANDROID_KEYSTORE_BASE64
   ```

   `giode.toml` already has `keystore = "release.keystore"` under `[android]`:
   `giode build -release` refuses to sign without it, so CI signing with
   `GIODE_KEYSTORE_PASS` depends on it.
3. Add `ANDROID_KEYSTORE_BASE64` and `GIODE_KEYSTORE_PASS`, then publish a
   release. It gets the `.aab` (and a signed `.apk` for GitHub).
4. In the Play Console, turn on **Play App Signing** (Google keeps the app
   signing key; the keystore above is only the upload key) and upload the `.aab`
   to the internal testing track first, then production.
5. `target_sdk = 36` is set: Google Play requires Android 16 (API 36) for new
   apps and updates since 31 Aug 2026. Raise it when Play raises the
   requirement. `min_sdk = 21` keeps Android 5 and later.
6. Fill in the listing and the policy forms (copy below): Data safety, content
   rating, target audience (adults; not for children), privacy policy URL, and
   "no ads".

Permissions in the manifest (`giode manifest android`): `INTERNET` (servers,
tunnels, Wake-on-LAN) and `POST_NOTIFICATIONS` (the optional "server
disconnected" notification). Nothing else.

## Mac App Store

1. Join the Apple Developer Program.
2. In Certificates, Identifiers & Profiles:
   - Register the App ID `com.xmhrnx.serverus` (the `app.id` in `giode.toml`).
   - Create an "Apple Distribution" and a "Mac Installer Distribution"
     certificate, and export each as a `.p12`.
   - Create a "Mac App Store" provisioning profile for the App ID.
   - For the `.dmg` on GitHub: a "Developer ID Application" certificate, and an
     app-specific password for your Apple ID (appleid.apple.com).
3. In App Store Connect, create a macOS app with the same bundle ID, named
   "Serverus".
4. Add the Mac App Store secrets above, then publish a release. The release gets
   `Serverus-<version>.pkg`; upload it with the Transporter app.
5. In App Store Connect, fill in the listing (copy below):
   - **Privacy Policy URL:** https://xmhrnx.github.io/serverus/#privacy
   - **App Privacy:** "Data Not Collected". The privacy manifest in the build
     (`[privacy]`, `tracking = false`, nothing collected) says the same.
   - **Category:** Developer Tools (set in the build as
     `LSApplicationCategoryType`). Secondary: Utilities.
   - **Export compliance:** answered in the build
     (`ITSAppUsesNonExemptEncryption = false`): Serverus only uses standard
     algorithms (SSH from Go's `x/crypto`, TLS for FTPS, the system keychain)
     and is open source. If App Store Connect asks about distribution in
     France, answer it there.

Sandbox entitlements (`[macos.store]`): `network.client` (SSH, Telnet, FTP, jump
hosts, Wake-on-LAN), `network.server` (tunnels listen on a local port, and the
single-instance lock) and `files.user-selected.read-write` (uploads, download
folders, key files). `giode build mas` adds the bookmarks entitlement, so a
download folder chosen once keeps working.

**The App Sandbox limits Serverus in the store version.** Remote servers work as
everywhere; what runs on the Mac itself is limited:

- **"This computer":** the local terminal, Files, Processes and monitoring run
  inside the sandbox. The shell starts in the app's container
  (`~/Library/Containers/com.xmhrnx.serverus/Data`), can only reach that and the
  folders the user picks, and can't signal other apps' processes.
- **SSH agent:** the sandbox can't reach the agent's socket (ssh-agent,
  1Password, Secretive). Use a key file or a password.
- **Key files:** a key picked with "Browse…" is copied into the app's container
  (see `internal/storage/keys.go`), because the sandbox forgets a picked file when
  the app quits. A typed path outside the container can't be read.
- **`~/.ssh/known_hosts`:** not readable; Serverus asks about each host's key
  once and pins it in its own known_hosts file.
- **Start at login** uses the system's login item API (macOS 13 and later).

The notarized `.dmg` on GitHub has none of these limits. Point users who manage
their own Mac with Serverus to it.

## Microsoft Store

1. Open a Partner Center developer account (the same one as EasySend and
   AnyFile).
2. Reserve the name "Serverus".
3. In the product's Product identity page, check these against `giode.toml`
   (`[windows.msix]`), and correct them if they differ:
   - **Package/Identity/Name** → `identity_name` (`MehranZahir.Serverus` assumed).
   - **Package/Identity/Publisher** → `publisher`
     (`CN=B40E8F57-A7E0-40BB-8E74-5D051D0B578B`).
   - **Package/Properties/PublisherDisplayName** → `publisher_display_name`
     (`Mehran Zahiroddini`).
4. Commit and publish a release. It gets `Serverus-<version>-amd64.msix`; upload
   it in a Partner Center submission. The Store signs it, so it needs no
   certificate here.
5. Serverus is a full-trust desktop app (it opens a local shell, manages local
   processes and files, and listens on local ports for tunnels). When the
   submission asks why it uses the `runFullTrust` capability, explain that. The
   package also declares a startup task for Settings → Start at login (off until
   the user turns it on).
6. Privacy policy URL: https://xmhrnx.github.io/serverus/#privacy. Age rating
   questionnaire: no objectionable content, no user-to-user communication, no
   purchases.

## Store listing

**Name:** Serverus

**Subtitle (Apple, ≤30):** Server monitor and SSH client

**Short description (Google Play, ≤80):**
Monitor and manage your servers: SSH terminal, files, services, Docker.

**Keywords (Apple, ≤100):**
`ssh,sftp,ftp,telnet,terminal,server,monitor,linux,sysadmin,docker,systemd,devops,tunnel,wake on lan`

**Full description:**

> Serverus manages, monitors and connects to all your servers from one window:
> over SSH, Telnet or FTP, and the computer it runs on.
>
> SERVERS
> • Any number of servers in a searchable sidebar, with groups, tags, colors and
>   live status.
> • SSH with a password, a private key (with passphrase) or an SSH agent;
>   keyboard-interactive and one-time codes; jump hosts.
> • Telnet for switches, routers and embedded boards; FTP and FTPS for files.
> • Passwords and passphrases are kept in the system's secure storage. Host keys
>   are checked, confirmed by fingerprint and pinned.
>
> MONITORING — nothing to install on the server
> • CPU, memory, swap, load, disks, network and temperatures of Linux servers,
>   with live charts.
> • An overview of every server: what needs attention, history and rankings.
> • Dashboards made of cards you can resize, move, add and remove.
>
> TERMINAL
> • Tabs of real terminals with true color, the alternate screen (vim, htop),
>   scrollback, copy and paste, and saved snippets.
>
> FILES
> • Browse, upload, download, edit, rename, change permissions and delete over
>   SFTP or FTP, with background transfers. Files owned by root can be changed
>   through sudo.
>
> ADMINISTRATION
> • Processes, systemd services (with logs), Docker containers, drives (mount,
>   unmount, eject), restart and scheduled shutdown, and Wake-on-LAN.
> • Local port forwarding through SSH.
> • Run a snippet on several servers at once and compare the results.
>
> PRIVATE BY DESIGN
> No account, no analytics, no ads, no tracking. Serverus only talks to the
> servers you add.

**Category:** Mac App Store: Developer Tools (secondary Utilities). Google Play:
Tools. Microsoft Store: Developer tools.

**Age rating:** no violence, sexual content, gambling, user-generated content
shared with others, or unrestricted web access → Apple 4+, IARC Everyone /
PEGI 3. The app isn't directed at children; set the Play target audience to 18+.

**App Privacy (Apple):** Data Not Collected; no tracking.

**Data safety (Google Play):**
- Does the app collect or share any of the required user data types? **No.**
  Everything stays on the device; the app connects only to servers the user
  adds, and sends nothing to the developer or third parties.
- Is all user data encrypted in transit? Not applicable (nothing collected). SSH
  and FTPS connections to the user's own servers are encrypted; Telnet and plain
  FTP are not, by the protocol's nature.
- Data deletion: not applicable; uninstalling removes everything.

**Export compliance:** standard encryption only (`ITSAppUsesNonExemptEncryption
= false`), see the Mac App Store section.

**Privacy policy URL:** https://xmhrnx.github.io/serverus/#privacy (the
`#privacy` section of `site/index.html`; keep it true to the code).

## Screenshots

There are none in the repository yet. Take them from the running app with a
demo workspace (made-up host names, no real addresses or user names):

1. The overview of every server (Ctrl/Cmd+0), several servers connected.
2. A server's overview with the history charts.
3. The terminal with two tabs, `htop` in one.
4. Files, with a transfer in the Transfers panel.
5. Services or Docker.

Sizes: Mac App Store 2880×1800 (or 1280×800); Microsoft Store 1920×1080 or
larger; Google Play phone screenshots (1080×1920 or larger, the phone layout
with the server drawer) and a 1024×500 feature graphic.

## Linux

The `.deb` and `.rpm` are attached to the release, and the website's download
button points to it. They list the system libraries they need, so `apt` and
`dnf` install those too.
