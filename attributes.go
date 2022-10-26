package fiber

type Attributes interface {
	Attribute(key string) []string
	WithAttribute(key string, values ...string) Attributes
}

// AttributesMap implements the Attributes interface via a simple map
type AttributesMap map[string][]string

func (a AttributesMap) Attribute(key string) []string {
	if values, ok := a[key]; ok {
		return values
	} else {
		return []string{}
	}
}

func (a AttributesMap) WithAttribute(key string, values ...string) Attributes {
	a[key] = values
	return a
}

func NewAttributesMap() Attributes {
	var newMap AttributesMap = map[string][]string{}
	return newMap
}
