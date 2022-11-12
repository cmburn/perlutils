package cpanmeta

type OptionalFeature struct {
	// Description is a description of the feature.
	Description string `json:"description"`

	// Prereqs is a list of prerequisites for the feature.
	Prereqs Prereqs `json:"prereqs"`
}
