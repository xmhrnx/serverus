// Command serverus is a cross-platform remote server manager: monitor
// resources, open terminals, transfer files and administer services on
// any number of servers over SSH, Telnet, FTP or locally.
package main

import (
	"log"

	"github.com/xmhrnx/serverus/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
