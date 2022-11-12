package metacpanclient

type Rating struct {
	hasUA
	Date         string `json:"date"`
	Release      string `json:"release"`
	Author       string `json:"author"`
	Details      string `json:"details"`
	Rating       string `json:"rating"`
	Distribution string `json:"distribution"`
	Helpful      int    `json:"helpful"`
	User         string `json:"user"`
}

func (r *Rating) _type() Type {
	return TypeRating
}
