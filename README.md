# dwrap

[![Go Reference](https://pkg.go.dev/badge/github.com/qawatake/dwrap.svg)](https://pkg.go.dev/github.com/qawatake/dwrap)
[![test](https://github.com/qawatake/dwrap/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/dwrap/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/qawatake/dwrap/graph/badge.svg?token=er2K6likVZ)](https://codecov.io/gh/qawatake/dwrap)

Linter `dwrap` forces every public function to begin with a deferring call of a error wrapping function like [derrors.Wrap](https://github.com/golang/pkgsite/blob/5f0513d53cff8382238b5f8c78e8317d2b4ad06d/internal/derrors/derrors.go#L240).
By using `dwrap`, you can prevent functions from returning without wrapping errors.

```go
func Good() (err error) {
  defer derrors.Wrap(&err, "Good")
  doStuff()
  return nil
}

func Bad() error { // <-  should call defer derrors.Wrap fist.
  doOtherStuff()
  return nil
}

//lint:ignore dwrap this is because ...
func Ignored() error {
  doOtherStuff()
  return nil
}
```

## How to use

Build your `dwrap` binary by writing `main.go` like below.

```go
package main

import (
  "github.com/qawatake/dwrap"
  "golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
  unitchecker.Main(
    dwrap.NewAnalyzer("your/derrors/pkg", "Wrap", "pkg/to/be/ignored"),
  )
}
```

Then, run `go vet` with your `dwrap` binary.

```sh
go vet -vettool=/path/to/your/dwrap ./...
```
