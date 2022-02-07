package types

type null struct {
	val bool
}

func (n *null) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		n.val = true
	}
	return nil
}

func (n *null) MarshalJSON() ([]byte, error) {
	if n == nil || !n.val {
		return nil, nil
	}
	return []byte("null"), nil
}
