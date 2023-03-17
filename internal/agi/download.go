package agi

// Download retrieves the needed files for the specified tools from the
// big bad Internet.  For most Go tools, this just takes care of supporting
// files since `go install` downloads the source to `/tmp` as part of
// its normal operation.
//
// See: https://asdf-vm.com/plugins/create.html#bin-install
func (p *plugin) Download() ExitCode {
	return ErrExitCodeNotImplemented
}
