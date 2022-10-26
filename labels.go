package fiber

type Labels interface {
	Label(key string) []string
	WithLabel(key string, values ...string) Labels
}

// LabelsMap implements the Labels interface via a simple map
type LabelsMap map[string][]string

func (a LabelsMap) Label(key string) []string {
	if values, ok := a[key]; ok {
		return values
	} else {
		return []string{}
	}
}

func (a LabelsMap) WithLabel(key string, values ...string) Labels {
	a[key] = values
	return a
}

func NewLabelsMap() Labels {
	var newMap LabelsMap = map[string][]string{}
	return newMap
}
