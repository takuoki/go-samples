package structtoorderedmap_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	structtoorderedmap "github.com/takuoki/golang-samples/struct-to-orderedmap"
)

func TestStructToOrderedmap(t *testing.T) {

	type Args struct {
		value      interface{}
		publicOnly bool
	}
	type Sample struct {
		Foo string
		Bar int
	}
	var nilSample *Sample

	testcases := map[string]struct {
		args      Args
		wantJSON  string
		wantError string
	}{
		"all type": {
			args: Args{
				value: struct {
					Bool      bool
					Int       int
					Uint      uint
					Float     float64
					String    string
					Ptr       *string
					Interface interface{}
					Array     [2]int
					Slice     []int
					Map       map[string]string
					Struct    struct {
						Foo string
						Bar string
					}
					// private
					private string
					// unsupported
					Chan chan int
					Func func() error
				}{
					Bool:      true,
					Int:       1,
					Uint:      2,
					Float:     3,
					String:    "4",
					Ptr:       strp(t, "5"),
					Interface: "6",
					Array:     [2]int{7, 8},
					Slice:     []int{9},
					Map:       map[string]string{"foo": "abc", "bar": "xyz"},
					Struct: struct {
						Foo string
						Bar string
					}{Foo: "aaa", Bar: "bbb"},
					// private
					private: "private",
					// unsupported
					Chan: make(chan int),
					Func: func() error { return nil },
				},
			},
			wantJSON: `{"Bool":true,"Int":1,"Uint":2,"Float":3,"String":"4","Ptr":"5","Interface":"6","Array":[7,8],"Slice":[9],"Map":{"bar":"xyz","foo":"abc"},"Struct":{"Foo":"aaa","Bar":"bbb"},"private":"private"}`,
		},
		"JSON tag": {
			args: Args{
				value: struct {
					Foo string `json:"foo,omitempty"`
					Bar int    `json:"bar"`
				}{
					Foo: "aaa",
					Bar: 2,
				},
			},
			wantJSON: `{"foo":"aaa","bar":2}`,
		},
		"public only": {
			args: Args{
				value: struct {
					Foo     string
					Bar     int
					private string
				}{
					Foo:     "aaa",
					Bar:     2,
					private: "private",
				},
				publicOnly: true,
			},
			wantJSON: `{"Foo":"aaa","Bar":2}`,
		},
		"sturct pointer": {
			args: Args{
				value: &Sample{
					Foo: "aaa",
					Bar: 2,
				},
			},
			wantJSON: `{"Foo":"aaa","Bar":2}`,
		},
		"argument is nil": {
			args: Args{
				value: nilSample,
			},
			wantJSON: "null",
		},
		"value is nil": {
			args: Args{
				value: struct {
					Ptr       *string
					Interface interface{}
					Slice     []int
					Map       map[string]string
				}{},
			},
			wantJSON: `{"Ptr":null,"Interface":null,"Slice":null,"Map":null}`,
		},
		"failure": {
			args: Args{
				value: map[string]string{},
			},
			wantError: "unsupported data type: ",
		},
	}

	for name, c := range testcases {
		t.Run(name, func(t *testing.T) {
			r, err := structtoorderedmap.StructToOrderedmap(c.args.value, c.args.publicOnly)
			if c.wantError == "" {
				if assert.Nil(t, err, "error must not be occurred") {
					j, _ := json.Marshal(r)
					assert.Equal(t, c.wantJSON, string(j))
				}
			} else {
				if assert.NotNil(t, err, "error must be occurred") {
					assert.Regexp(t, "^"+c.wantError, err.Error())
				}
			}
		})
	}
}

func strp(t *testing.T, s string) *string {
	t.Helper()
	return &s
}
