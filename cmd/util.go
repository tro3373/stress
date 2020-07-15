package cmd

import (
	"encoding/json"
	"upspin.io/log"
)

func Pretty(data interface{}, sep string) []byte {
	// if data.(type) == string {
	// }
	var p []byte
	if sep == "" {
		sep = "  "
	}
	//    var err := error
	p, err := json.MarshalIndent(data, "", sep)
	if err != nil {
		// fmt.Println(err)
		log.Fatal(err)
		return nil
	}
	return p
}
