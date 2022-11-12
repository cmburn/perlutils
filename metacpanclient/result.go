package metacpanclient

type result interface {
	setClient(*Client)
	_type() Type
}

func getType[T result]() Type {
	var v T
	return v._type()
}

type hasUA struct {
	mc *Client
}

func (r *hasUA) setClient(mc *Client) {
	r.mc = mc
}
