package events

// PackageRegistry is the registry of the package
type PackageRegistry struct {
	URL string `json:"url"`
}

// PackageVersion is the version of the package
type PackageVersion struct {
	ID      int    `json:"id"`
	Version string `json:"version"`
}

// Package is the package that being updated or published
type Package struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	PackageVersion  PackageVersion  `json:"package_name"`
	PackageRegistry PackageRegistry `json:"registry"`
}

// Event of the webhook
type Event struct {
	Action  string  `json:"action"`
	Package Package `json:"package,omitempty"`
}
