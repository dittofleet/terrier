package cmd

// Windows gets no `ter`. A symlink there takes Developer Mode or an
// administrator, and a hard link or a copy would keep running the old
// binary after an update.
func setupAlias(string) error { return nil }

func aliasPath(string) string { return "" }
