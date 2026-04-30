package commentspacings

import (
	"fmt"
	"os"
)

//myOwnDirective: do something // should be valid

type c struct {
	//+optional
	d *int `json:"d,omitempty"`
}

func _(outputPath string) {
	//gosec:disable G703
	f, err := os.Create(outputPath) //#nosec G703 - path is validated above
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
}
