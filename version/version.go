package version

var (
	major     string
	minor     string
	patch     string
	commit    string
	buildDate string
)

var BuildInfo Info = Info{
	Major:     major,
	Minor:     minor,
	Patch:     patch,
	Commit:    commit,
	BuildDate: buildDate,
}

type Info struct {
	// Major version number
	//
	// The leftmost number in the version string
	Major string `json:"major"`

	// Minor version number
	//
	// The number in the middle of the version string
	Minor string `json:"minor"`

	// Patch number
	//
	// The rightmost number in the version string
	Patch string `json:"patch"`

	// Full commit hash
	//
	// The full commit hash of HEAD at build time
	Commit string `json:"commit"`

	// Build date
	BuildDate string `json:"buildDate"`
}
