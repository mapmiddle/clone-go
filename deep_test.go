package clone

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Ifoo interface {
	test()
}

type foo struct {
	s1 int
	s2 string
	s3 *string
	s4 *bar
}

func (t *foo) test() {
	// empty
}

type bar struct {
	s5 int
	s6 string
}

func TestDeep(t *testing.T) {
	pstr := "hello"
	ts := foo{
		30,
		"external struct",
		&pstr,
		nil,
	}

	m1 := map[string]interface{}{
		"integer": 10,
		"string":  "test",
		"nested_pointer_struct": &foo{
			10,
			"nested pointer struct",
			&pstr,
			&bar{
				11,
				"double nested pointer struct",
			},
		},
		"nested_struct": foo{
			20,
			"nested struct2",
			&pstr,
			nil,
		},
		"nested_map": map[string]interface{}{
			"e1": 10,
			"e2": "nested map",
		},
		"interface": Ifoo(&ts),

		"array": [5]int{10, 20, 30, 40, 50},
		"slice": []int{60, 70, 80, 90, 100},

		"closure": func() interface{} {
			return ts
		},
	}

	m2 := Deep(m1).(map[string]interface{})

	assert.Equal(t, m1["integer"], m2["integer"])
	m2["integer"] = 1
	assert.NotEqual(t, m1["integer"], m2["integer"])

	assert.Equal(t, m1["string"], m2["string"])
	m2["string"] = "test2"
	assert.NotEqual(t, m1["string"], m2["string"])

	m2["nested_pointer_struct"].(*foo).s1 = 2
	m2["nested_pointer_struct"].(*foo).s2 = "replaced1"
	assert.NotEqual(t, m1["nested_pointer_struct"].(*foo).s1, m2["nested_pointer_struct"].(*foo).s1)
	assert.NotEqual(t, m1["nested_pointer_struct"].(*foo).s2, m2["nested_pointer_struct"].(*foo).s2)

	assert.Equal(t, m1["nested_struct"], m2["nested_struct"])
	m2["nested_map"].(map[string]interface{})["e1"] = 1
	m2["nested_map"].(map[string]interface{})["e2"] = "replaced2"
	assert.NotEqual(t, m1["nested_map"].(map[string]interface{})["e1"], m2["nested_map"].(map[string]interface{})["e1"])
	assert.NotEqual(t, m1["nested_map"].(map[string]interface{})["e2"], m2["nested_map"].(map[string]interface{})["e2"])

	m2["interface"].(*foo).s1 = 3
	m2["interface"].(*foo).s2 = "replaced3"
	(*m2["interface"].(*foo).s3) += "append1"
	assert.NotEqual(t, m1["interface"].(*foo).s1, m2["interface"].(*foo).s1)
	assert.NotEqual(t, m1["interface"].(*foo).s2, m2["interface"].(*foo).s2)
	assert.NotEqual(t, m1["interface"].(*foo).s3, m2["interface"].(*foo).s3)

	m2["slice"].([]int)[0] = 6
	m2["slice"].([]int)[1] = 7
	m2["slice"].([]int)[2] = 8
	assert.NotEqual(t, m1["slice"], m2["slice"])

	assert.Equal(t, m1["array"], m2["array"])
	assert.Equal(t,
		reflect.ValueOf(m1["closure"]).Pointer(),
		reflect.ValueOf(m2["closure"]).Pointer())
}
