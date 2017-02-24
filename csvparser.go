package main

import "fmt"
import "strings"

func ParseCSVDriverFile(path string) {
	file := NewCSVFile(path)
	fmt.Printf("\n=== Header Information ===\n%s\n", strings.Join(file.rows[0], ", "))
	for i := 1; i < len(file.rows); i++ {
		fmt.Printf("\n=== Fileds In Row #%d ===\n", i)
		fmt.Printf("%s", strings.Join(file.rows[i], ", "))
	}
	println()
}
