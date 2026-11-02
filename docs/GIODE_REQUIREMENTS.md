# What Serverus needs from giode and the giode CLI

These are the gaps found while building Serverus on giode `v0.1.20-alpha.1` and the giode
CLI at `main`. Where something was missing, Serverus does without it rather than working
around it; each item says what the app does today and what would let it do better.

Items are ordered by how much they affect the app.

---

## 1. `Menu` should take no space and dismiss from anywhere in the window

**Area:** `components/menu`
**Kind:** bug / API

### What happens

`Menu.Layout` returns `layout.Dimensions{Size: gtx.Constraints.Max}` while the menu is
open, and the menu does not implement `elements.OutOfFlow`. The docs say a menu opened
with `OpenAt` can be "laid out anywhere in the view", but when it sits in a `Stack` or
`Row` next to other children, opening it takes all the remaining space and pushes its
siblings away. In Serverus, opening a server's context menu made the sidebar's server list
disappear until the menu closed.

The click-outside layer (`event.Op(gtx.Ops, &m.scrim)` over `gtx.Constraints.Max`) also
only covers the area the menu was laid out in, so a click outside that area (on the title
bar or on another pane) doesn't close the menu.

### What Serverus does now

Each context menu is laid out as an absolutely positioned child covering its pane
(`widgets.WithOverlays`), which keeps the layout intact. Clicks outside the pane still
don't close the menu.

### What's needed

- `Menu` implements `elements.OutOfFlow`, like `Dialog`, `Drawer` and `Sheet`.
- Its dismissal layer covers the whole window (registered as an overlay portal), not
  the constraints it was laid out with.

---

## 2. Keyboard shortcuts

**Area:** new package (the docs already refer to it)
**Kind:** missing feature

### What happens

`Menu.Shortcut("Ctrl+O")` is documented as "display only; handle the key with the
`shortcut` package", and `docs/components/menu.md` says the same, but there is no
`shortcut` package.

### What Serverus does now

No app-wide shortcuts. The terminal handles its own keys (Ctrl+Shift+C/V, Shift+PageUp)
because it owns the keyboard focus.

### What's needed

A way to register window-wide shortcuts that fire whatever has focus, unless a focused
text field or terminal consumes the key first:

```go
giode.Shortcuts(
    giode.Shortcut("Mod+T", newTerminalTab),     // Mod = Cmd on macOS, Ctrl elsewhere
    giode.Shortcut("Mod+K", focusServerSearch),
    giode.Shortcut("Mod+,", openSettings),
    giode.Shortcut("F5", refreshPage),
)
```

Serverus would use: new terminal tab, close tab, next/previous tab, focus the server
search, open settings, refresh the current page, and Mod+1…9 to switch pages.

---

## 3. `Menu.OpenAt` opens offset by the containers around the menu

**Area:** `components/menu`, `elements` portals
**Kind:** bug

### What happens

`Menu.OpenAt(pos)` takes a window position, and the menu registers its panel as a portal
at that position. When the menu is laid out inside a `Row`, `Stack` or `Scroll` that
places its children after laying them out, `ChildTracker.Place` then moves the portal's
anchor by the child's offset, as it does for portals anchored at
`AbsolutePosition()`. The panel opens that far from where it was asked to: 40px low in
Serverus's sidebar (under the title bar), about 170px low on a page under a header and
tabs.

### What Serverus does now

`widgets.ContextMenu` wraps each menu, records `AbsolutePosition()` before laying it out
and the tracked position after, and subtracts the difference in `OpenAt`.

### What's needed

`OpenAt` positions that stay window-absolute wherever the menu is laid out (register the
panel as an overlay portal, or convert the position to the container-relative one
`Place` expects). Related, and still useful: `ClickEvent.WindowPosition` or
`Menu.OpenAtPointer()`, so a list doesn't need to track positions to open a row's menu
where it was clicked.

---

## 4. Explicit zero padding in `Styles`

**Area:** `styles`, `styles/properties`
**Kind:** API

### What happens

Zero means "unset", so `base.With(giode.Styles{Padding: giode.UniformInset(0)})` keeps
the base padding. Offsets have `giode.Flush` for an explicit zero; insets have nothing
equivalent.

### What Serverus does now

A second card component (`widgets.FlushCard`) with its own base styles, for tables and
lists that run to the card's edges.

### What's needed

An explicit-zero inset, e.g. `giode.NoInset` or `properties.Inset{…}.Set()`, honored by
`Merge`.

---

## 5. Changing a `Select`'s options

**Area:** `components/dropdown`
**Kind:** API

### What happens

`Input.Select` takes its options at creation and has no setter. The server editor's
"Jump host" select lists the other SSH servers, which change as servers are added.

### What Serverus does now

It creates a new `Select` each time the editor opens.

### What's needed

`Select.SetOptions(options ...string)` (keeping the selection when the selected option is
still there), like `AutoComplete.SetOptions`.

---

## 6. A terminal element (optional)

**Area:** new component
**Kind:** feature request, not blocking

Serverus's terminal (`internal/ui/terminal`) is built the way `docs/extending.md`
recommends: the VT emulation comes from `github.com/charmbracelet/x/vt`, and drawing,
keyboard, IME text, selection, scrollback and paste use `giode.Raw` with Gio's input
filters. It works well, so this isn't blocking, but terminals are common in developer
tools and a `giode.Terminal(io.ReadWriter)` element would save every such app the same
~600 lines. If you add one, Serverus would move to it. Things it would need to cover:

- xterm-256color and true color, bold, italic, underline, reverse, wide (CJK) characters
- alternate screen, application cursor keys, bracketed paste, mouse wheel in full-screen
  programs
- selection with copy, scrollback, resize reporting (cols × rows)
- `SemanticInput` with `ChangeText`, so the MCP server can type into it (Serverus's
  terminal does this)

---

## 7. Dragging files out of the window (optional)

**Area:** `filedrag.go` / platform drag sources
**Kind:** feature request

Dropping files *into* the window works (`OnFileDrop`, `OnFileDrag`), and Serverus uses it
to upload. Dragging a remote file *out* of the file manager onto the desktop or a folder
(a drag source with a file promise, `NSFilePromiseProvider` / `IDataObject` with
`CFSTR_FILEDESCRIPTOR` / XDND `XdndDirectSave`) would make downloads as easy as uploads.
Today downloads go through the folder picker.

---

## 8. MCP: `ui_pointer` `click` ignores `button`

**Area:** `mcpinput.go`
**Kind:** bug (automation only)

`ui_pointer` with `{"action": "click", "button": "secondary"}` delivers a primary click
(the result reports `"buttons": ""`), while `press` and `release` with the same `button`
deliver secondary-button events correctly. Right-click menus can only be tested with a
press/release pair.

---

## 9. Cross-axis stretch (`align-items: stretch`)

**Area:** `styles/properties`, `elements` flex layout
**Kind:** missing feature


`Align` has start, center and end; `properties/align.go` notes CSS's default `stretch`
isn't supported. Cards side by side in a row can't be made equally tall, so a dashboard
row of two cards with different content has a ragged bottom edge. Serverus avoids
pairing cards whose heights differ much (cores and temperatures now take a row each).
`AlignStretch`, on a `Row`/`Stack` or per child, would fix it.

---

## The giode CLI

`giode new`, `giode build linux` and the Gio patches (needed for `ui_pointer`/`ui_key` and
file drops on Linux) worked as documented. Windows cross-compiles with plain `go build`
(`GOOS=windows`). Releases showed three gaps:

### 10. Four-part release tags in `giode ci`

Serverus versions releases as `major.minor.patch.versionCode` (`1.0.0.0`), the format
gogio itself uses for Android. `giode ci` sets `GIODE_VERSION` to the whole tag and
`GIODE_BUILD` to the run number, and `GIODE_VERSION` must be `major.minor.patch`, so a
generated workflow can't build from such a tag.

**What Serverus does now:** `.github/workflows/release.yml` started from `giode ci` and
was edited by hand: a first job splits the tag into `GIODE_VERSION` and `GIODE_BUILD`.
It must not be regenerated with `giode ci -f`.

**What's needed:** a `[ci]` setting for where the build number comes from, e.g.
`build = "tag"` (the fourth part of a four-part tag) besides the current run number, and
for `GIODE_VERSION`/`app.version` to accept a four-part version, taking the fourth part as
the build.

### 11. A build number of 0 is silently raised to 1

`config.Load` turns `app.build = 0` (and `GIODE_BUILD=0`) into 1, so `1.0.0.0` and
`1.0.0.1` produce the same Android `versionCode`, and the second can't be uploaded to
Google Play. Either reject 0 with an error saying the build must be at least 1, or
document it.

### 12. `giode ci` with `trigger = "release"` adds `workflow_dispatch`

The generated workflow also runs by hand from the Actions tab. For a release-only
pipeline it would help to leave it out (or make it a setting).

### 13. `giode doctor` and `giode ci` miss bundletool

`giode build aab` needs `bundletool*.jar` in the SDK's latest build-tools folder. The
GitHub Ubuntu runner has the SDK, NDK and JDK but not bundletool, so the `giode ci`
workflow fails at the AAB after the APK has built, while `giode doctor` on the same
runner reports "Everything giode needs on this computer is installed".

**What Serverus does now:** its workflow downloads bundletool (pinned version) from
Google's GitHub releases into build-tools before building Android.

**What's needed:** `giode doctor` checks for bundletool when Android is a platform, and
the `giode ci` workflow installs it, or the CLI fetches and caches it itself as it does
for gogio's other tools.
