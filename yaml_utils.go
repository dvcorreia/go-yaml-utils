package llrpgen

import (
	"gopkg.in/yaml.v3"
)

// get checks if node has a direct field and returns it.
func get(v *yaml.Node, key string) (*yaml.Node, bool) {
	if v.Kind == yaml.DocumentNode && len(v.Content) > 0 {
		v = v.Content[0]
	}

	if v.Kind != yaml.MappingNode {
		return nil, false
	}

	for i := 0; i < len(v.Content); i += 2 {
		k := v.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return v.Content[i+1], true
		}
	}

	return nil, false
}

// to decodes a node in to T type.
func to[T any](v *yaml.Node, o T) error {
	// check if the node is a sequence
	if v.Kind == yaml.SequenceNode {
		switch ot := any(o).(type) {
		// try to convert the sequence into a map if T is a map
		case map[string]string:
			for _, item := range v.Content {
				var m map[string]string
				if err := item.Decode(&m); err != nil {
					return err
				}
				for k, v := range m {
					ot[k] = v
				}
			}
			return nil
		case map[string]int:
			for _, item := range v.Content {
				var m map[string]int
				if err := item.Decode(&m); err != nil {
					return err
				}
				for k, v := range m {
					ot[k] = v
				}
			}
			return nil
		}
	}

	return v.Decode(&o)
}

// copyIfExists decodes the value, if exits, and assigns it to to.
// to must be a pointer.
func copyIfExists[T any](v *yaml.Node, key string, into T) {
	node, ok := get(v, key)
	if !ok {
		return
	}
	_ = to(node, into)
}
