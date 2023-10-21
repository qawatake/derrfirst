package main

import (
	"dfirst"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(dfirst.NewAnalyzer("fmt.Println")) }
