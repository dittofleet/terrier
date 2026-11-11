module github.com/dittofleet/terrier

go 1.26.5

// Junctions are how Windows links directories without Developer Mode, so
// they count as symlinks: a repo reached through one registers once. The
// second setting turns a volume mounted at a folder into a path git can use.
godebug (
	winreadlinkvolume=0
	winsymlink=0
)

require (
	github.com/dittofleet/go-cli-kit v0.7.1
	golang.org/x/sys v0.44.0
)

require golang.org/x/term v0.43.0 // indirect
