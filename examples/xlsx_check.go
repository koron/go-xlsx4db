package main

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func main() {
	xf, err := xlsx.OpenFile("tmp/out.xlsx")
	if err != nil {
		panic(err)
	}
	for _, xs := range xf.Sheets {
		fmt.Printf("%q:\n", xs.Name)
		for _, xr := range xs.Rows {
			fmt.Printf("  -\n")
			for _, xc := range xr.Cells {
				fmt.Printf("    -\n")
				fmt.Printf("      type: %d\n", xc.Type())
				fmt.Printf("      value: %q\n", xc.Value)
				fmt.Printf("      numfmt: %q\n", xc.NumFmt)
			}
		}
	}
}
