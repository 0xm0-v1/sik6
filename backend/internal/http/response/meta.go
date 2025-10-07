package response

import "time"

// Meta contains standardized metadata returned with API responses.
type Meta struct {
	Component string `json:"component"`
	Type      string `json:"type"`
	Time      string `json:"time"`
}

// NewMeta builds a Meta with the supplied component/type and a UTC timestamp.
func NewMeta(component, typ string) Meta {
	return Meta{
		Component: component,
		Type:      typ,
		Time:      time.Now().UTC().Format(time.RFC3339Nano),
	}
}
