package metacpanclient

// Distribution contains information about a distribution.
type Distribution struct {
	// Name is the name of the distribution
	Name string `json:"name"`

	// Bugs contains information about bugs within a distribution
	Bugs struct {
		// Github contains information about bugs on GitHub
		Github Bugs `json:"github"`
		// RT contains information about bugs on rt.cpan.org
		RT Bugs `json:"rt"`
	} `json:"bugs"`

	// River contains information about the distribution's "CPAN River"
	River struct {
		// Bucket indicates how far "upstream" the distribution is.
		Bucket int `json:"bucket"`

		// Immediate indicates the number of distributions that depend
		// immediately on this distribution.
		Immediate int `json:"immediate"`

		// Total indicates the total number of distributions that
		// depend on this distribution, including indirect
		// dependencies.
		Total int `json:"total"`
	} `json:"river"`

	hasUA
}

// Github returns the bugs on GitHub for the distribution.
func (d *Distribution) Github() Bugs {
	return d.Bugs.Github
}

// RT returns the bugs on rt.cpan.org for the distribution
func (d *Distribution) RT() Bugs {
	return d.Bugs.RT
}

func (d *Distribution) _type() Type {
	return TypeDistribution
}
