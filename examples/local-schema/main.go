package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s := `{"value": "123", "label": "222", "age": "123"}`
	var ss Structs
	json.Unmarshal([]byte(s), &ss)
	out, _ := json.Marshal(ss)
	fmt.Println(string(out))
}
