package agi

import (
	"os"
	"strings"
)

// ListAll outputs a string of available version numbers for the target
// Go tool.
//
// See: https://asdf-vm.com/plugins/create.html#bin-list-all
func (p *plugin) ListAll() ExitCode {
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "ASDF") {
			p.env.log.Debug("EnvVar: ", e)
		}
	}

	return ErrExitCodeNotImplemented
}
