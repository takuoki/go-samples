package structtoorderedmap_test

import (
	"encoding/json"
	"fmt"
	"log"

	structtoorderedmap "github.com/takuoki/golang-samples/struct-to-orderedmap"
)

func Example() {

	type StructureDefinedBySomeone struct {
		Foo string `json:"foo,omitempty"`
		Bar int    `json:"bar,omitempty"`
	}

	v := StructureDefinedBySomeone{
		Foo: "",
		Bar: 0,
	}

	structJSON, _ := json.Marshal(v)

	fmt.Println("Too Bad!!!")
	fmt.Println("structJSON:", string(structJSON))

	o, err := structtoorderedmap.StructToOrderedmap(v, false)
	if err != nil {
		log.Fatalf("if you pass 'struct', no error will occur: %v", err)
	}

	orderedmapJSON, _ := json.Marshal(o)

	fmt.Println("So Good!!!")
	fmt.Println("orderedmapJSON:", string(orderedmapJSON))

	// Output:
	// Too Bad!!!
	// structJSON: {}
	// So Good!!!
	// orderedmapJSON: {"foo":"","bar":0}
}
