package clone

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Ibaz interface {
	test()
}

type baz struct {
	s1 int
	s2 string
	s3 *string
}

func (t *baz) test() {
	// empty
}

func TestShallow(t *testing.T) {
	pstr := "hello"
	ts1 := baz{
		10,
		"clone-test",
		&pstr,
	}

	ts2 := Shallow(ts1).(baz)
	ts3 := Shallow(&ts1).(*baz)
	ts4 := Shallow(ts3).(*baz)
	assert.Equal(t, ts1, ts2)
	assert.Equal(t, ts1, *ts3)
	assert.Equal(t, *ts3, *ts4)

	ts1.s1 = 100
	assert.NotEqual(t, ts1, ts2)
	assert.NotEqual(t, ts1, *ts3)
	ts3.s1 = 200
	assert.NotEqual(t, *ts3, *ts4)

	var ti1 Ibaz = &ts1
	ti2 := Shallow(ti1).(Ibaz)
	assert.Equal(t, ti1, ti2)
	assert.NotEqual(t, unsafe.Pointer(ti1.(*baz)), unsafe.Pointer(ti2.(*baz)))

	// map
	tm1 := map[string]interface{}{
		"a": 1,
		"b": "clone-test",
		"c": &pstr,
	}

	tm2 := Shallow(tm1).(map[string]interface{})
	assert.Equal(t, tm1, tm2)
	tm1["a"] = 100
	assert.NotEqual(t, tm1, tm2)

	// array
	ta1 := [3]string{"a", "b", "c"}
	ta2 := Shallow(ta1).([3]string)
	assert.Equal(t, ta1, ta2)
	ta1[0] = "d"
	assert.NotEqual(t, ta1, ta2)

	// slice
	tc1 := []string{"a", "b", "c"}
	tc2 := Shallow(tc1).([]string)
	assert.Equal(t, tc1, tc2)
	ta1[0] = "d"
	assert.NotEqual(t, ta1, ta2)

	// nil
	tn1 := Shallow(nil)
	assert.Nil(t, tn1)

	var tpn *baz = nil
	var tin Ibaz = tpn
	tn2 := Shallow(tin).(Ibaz)
	assert.Nil(t, tn2)
}
