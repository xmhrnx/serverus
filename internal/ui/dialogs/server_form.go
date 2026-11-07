package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *ServerDialog) form() giode.Element {
	p := d.selectedProtocol()
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	if d.testMsg != "" {
		if d.testOK {
			rows = append(rows, widgets.Banner("circle-check", theme.Success, d.testMsg))
		} else {
			rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.testMsg))
		}
	}
	rows = append(rows,
		widgets.FormField("Connection type", d.selectedProtocol().Label()+protocolHint(d.selectedProtocol()), d.protocol),
		widgets.Columns(nil,
			widgets.FormField("Name", "", d.name),
			widgets.FormField("Group", "Servers sharing a group are listed together", d.group),
		),
	)
	if p.NeedsHost() {
		rows = append(rows, widgets.Columns([]int{4, 1},
			widgets.FormField("Host", "", d.host),
			widgets.FormField("Port", "", d.port),
		))
	}
	if p.NeedsCredentials() {
		hint := ""
		if p == domain.ProtocolFTP {
			hint = "Leave empty to sign in anonymously"
		}
		rows = append(rows, widgets.FormField("User name", hint, d.user))
	}
	if p == domain.ProtocolSSH {
		rows = append(rows, widgets.FormField("Authentication", "", d.auth))
		switch domain.AuthMethods[d.auth.Selected()] {
		case domain.AuthPassword:
			rows = append(rows, widgets.FormField("Password", "Stored in the system keychain, never in the workspace file", d.password))
		case domain.AuthKey:
			rows = append(rows,
				widgets.FormField("Private key", "OpenSSH, PEM or PKCS#8; RSA, ECDSA and Ed25519", giode.Row(giode.Styles{Gap: theme.Space2, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Width: -1},
					giode.Box(giode.Styles{FlexGrow: 1}, d.keyPath), d.browse)),
				widgets.FormField("Passphrase", "", d.passphrase),
			)
		case domain.AuthAgent:
			rows = append(rows, widgets.Banner("info", theme.Info, "Keys come from the running SSH agent: ssh-agent, the Windows OpenSSH agent, 1Password or Secretive."))
		}
		rows = append(rows, widgets.FormField("Jump host", "Reach this server through another one, like ssh -J", d.jump))
	}
	if p == domain.ProtocolFTP {
		rows = append(rows, widgets.FormField("Password", "Stored in the system keychain", d.password), d.ftps)
	}
	if p.NeedsHost() {
		rows = append(rows, widgets.FormField("MAC address", "Wakes the server with Wake-on-LAN after it's shut down", d.mac))
	}
	rows = append(rows,
		widgets.FormField("Color", "", d.accent),
		widgets.FormField("Tags", "Separate with commas", d.tags),
		widgets.FormField("Notes", "", d.notes),
	)
	return giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1, Padding: giode.Inset{Right: theme.Space3}}, rows...)
}
