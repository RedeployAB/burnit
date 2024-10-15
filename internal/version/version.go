package version

var (
	version string
	commit  string
)

// Version returns version that has been set as package level variable.
func Version() string {
	return version
}

// Commit returns commit that has been set as package level variable.
func Commit() string {
	return commit
}
