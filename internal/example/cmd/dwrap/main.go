package main

import (
	"github.com/qawatake/dwrap"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(
		dwrap.NewAnalyzer("github.com/qawatake/dwrap/internal/example", "Wrap"),
	)
}
