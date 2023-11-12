package main

import (
	"encoding/json"
	"fmt"

	"github.com/gowool/cr"
)

var (
	filter = `m.created = '2023-11-23' AND (m.status IN 2,3 OR m.title IS NULL OR m.title LIKE "My \"title\"") AND m.Enabled = true`
	sort   = `m.created,-m.updated`
)

func main() {
	criteria := cr.New(filter, sort).SetOffset(0).SetSize(20)
	where, args := criteria.Filter.ToSQL()

	fmt.Println(marshal(criteria))
	fmt.Println()
	fmt.Println(where)
	fmt.Println(marshal(args))
}

func marshal(i interface{}) string {
	raw, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(raw)
}
