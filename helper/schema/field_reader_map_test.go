package schema

import (
	"reflect"
	"testing"
)

func TestMapFieldReader_impl(t *testing.T) {
	var _ FieldReader = new(MapFieldReader)
}

func TestMapFieldReader(t *testing.T) {
	testFieldReader(t, func(s map[string]*Schema) FieldReader {
		return &MapFieldReader{
			Schema: s,

			Map: BasicMapReader(map[string]string{
				"bool":   "true",
				"int":    "42",
				"float":  "3.1415",
				"string": "string",

				"list.#": "2",
				"list.0": "foo",
				"list.1": "bar",

				"listInt.#": "2",
				"listInt.0": "21",
				"listInt.1": "42",

				"map.#":   "2",
				"map.foo": "bar",
				"map.bar": "baz",

				"set.#":  "2",
				"set.10": "10",
				"set.50": "50",

				"setDeep.#":        "2",
				"setDeep.10.index": "10",
				"setDeep.10.value": "foo",
				"setDeep.50.index": "50",
				"setDeep.50.value": "bar",
			}),
		}
	})
}

func TestMapFieldReader_extra(t *testing.T) {
	r := &MapFieldReader{
		Schema: map[string]*Schema{
			"mapDel": &Schema{Type: TypeMap},
		},

		Map: BasicMapReader(map[string]string{
			"mapDel": "",
		}),
	}

	cases := map[string]struct {
		Addr        []string
		Out         interface{}
		OutOk       bool
		OutComputed bool
		OutErr      bool
	}{
		"mapDel": {
			[]string{"mapDel"},
			map[string]interface{}{},
			true,
			false,
			false,
		},
	}

	for name, tc := range cases {
		out, err := r.ReadField(tc.Addr)
		if (err != nil) != tc.OutErr {
			t.Fatalf("%s: err: %s", name, err)
		}
		if out.Computed != tc.OutComputed {
			t.Fatalf("%s: err: %#v", name, out.Computed)
		}

		if s, ok := out.Value.(*Set); ok {
			// If it is a set, convert to a list so its more easily checked.
			out.Value = s.List()
		}

		if !reflect.DeepEqual(out.Value, tc.Out) {
			t.Fatalf("%s: out: %#v", name, out.Value)
		}
		if out.Exists != tc.OutOk {
			t.Fatalf("%s: outOk: %#v", name, out.Exists)
		}
	}
}
