package llrpgen

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

var testdoc = `# foo is a bar
foo: bar
lhenumber: 78
others:
  to: be
  not: to
  be: 666
someMap:
  - someKey: value # line comment
  - someOtherKey: value2
sequence:
  - a
  - b
  - c
# bin is a baz
bin: baz
`

func LoadYamlNode(t testing.TB, data []byte) *yaml.Node {
	t.Helper()
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		t.Fatalf("could not load yaml data: %v", err)
	}
	return &root
}

func TestGet(t *testing.T) {
	node := LoadYamlNode(t, []byte(testdoc))

	if v, ok := get(node, "foo"); !ok || v.Value != "bar" {
		t.Fatalf("foo is property of testdocs and has value '%q'", v.Value)
	}

	if _, ok := get(node, "unknown"); ok {
		t.Fatalf("unknown field does not exist")
	}

	if v, ok := get(node, "someMap"); !ok || v.Kind != yaml.SequenceNode {
		t.Fatalf("someMap is property of testdocs and is a sequence")
	}

	if v, ok := get(node, "others"); !ok || v.Kind != yaml.MappingNode {
		t.Fatalf("others is property of testdocs and is a mapping node")
	}
}

func TestTo(t *testing.T) {
	node := LoadYamlNode(t, []byte(testdoc))

	fooNode, _ := get(node, "foo")
	var foo string
	if err := to(fooNode, &foo); err != nil || foo != "bar" {
		t.Fatalf("expected foo field to contain 'bar' and got error: %v", err)
	}

	lhenumberNode, _ := get(node, "lhenumber")
	var lhenumber *int
	if err := to(lhenumberNode, lhenumber); err != nil || foo != "bar" {
		t.Fatalf("expected foo field to contain 'bar' and got error: %v", err)
	}

	type Others struct {
		To  string `yaml:"to"`
		Not string `yaml:"not"`
		Be  int    `yaml:"be"`
	}
	othersNode, _ := get(node, "others")
	var others Others
	expectedOthers := Others{
		To:  "be",
		Not: "to",
		Be:  666,
	}
	if err := to(othersNode, &others); err != nil || !reflect.DeepEqual(others, expectedOthers) {
		t.Fatalf("expected others field to contain '%+v', got '%+v' and error: %v", expectedOthers, others, err)
	}

	someMapNode, _ := get(node, "someMap")
	someMap := make(map[string]string)
	expectedMap := map[string]string{
		"someKey":      "value",
		"someOtherKey": "value2",
	}
	if err := to(someMapNode, someMap); err != nil || !reflect.DeepEqual(someMap, expectedMap) {
		t.Fatalf("expected someMap field to contain '%+v', got '%+v' and error: %v", expectedMap, someMap, err)
	}

	sequenceNode, _ := get(node, "sequence")
	sequence := []string{}
	expectedSequence := []string{"a", "b", "c"}
	if err := to(sequenceNode, &sequence); err != nil || !reflect.DeepEqual(sequence, expectedSequence) {
		t.Fatalf("expected sequece field to contain '%+v', got '%+v' and error: %v", expectedSequence, sequence, err)
	}
}

func TestCopyIfExists(t *testing.T) {
	node := LoadYamlNode(t, []byte(testdoc))

	var foo string
	copyIfExists(node, "foo", &foo)
	if foo != "bar" {
		t.Fatalf("expected foo field to be copied and got '%s' instead of 'bar", foo)
	}

	type Others struct {
		To  string `yaml:"to"`
		Not string `yaml:"not"`
		Be  int    `yaml:"be"`
	}
	var others Others
	expectedOthers := Others{
		To:  "be",
		Not: "to",
		Be:  666,
	}
	if copyIfExists(node, "others", &others); !reflect.DeepEqual(others, expectedOthers) {
		t.Fatalf("expected others field to contain '%+v', got '%+v'", expectedOthers, others)
	}

	someMap := make(map[string]string)
	expectedMap := map[string]string{
		"someKey":      "value",
		"someOtherKey": "value2",
	}
	if copyIfExists(node, "someMap", someMap); !reflect.DeepEqual(someMap, expectedMap) {
		t.Fatalf("expected someMap field to contain '%+v', got '%+v'", expectedMap, someMap)
	}

	sequence := []string{}
	expectedSequence := []string{"a", "b", "c"}
	if copyIfExists(node, "sequence", &sequence); !reflect.DeepEqual(sequence, expectedSequence) {
		t.Fatalf("expected sequece field to contain '%+v', got '%+v'", expectedSequence, sequence)
	}
}
