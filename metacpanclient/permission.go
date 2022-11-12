package metacpanclient

type Permission struct {
	// The name of the module this permission belongs to.
	ModuleName string `json:"module_name"`

	// Owner is the owner of the module.
	Owner string `json:"owner"`

	// CoMaintainers is the list of other maintainers with permissions for
	// the module.
	CoMaintainers []string `json:"co_maintainers"`
}
