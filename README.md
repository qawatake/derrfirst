# derrfirst

Linter `derrfirst` forces every public function to begin with a deferring call of some function like [derrors.Wrap](https://github.com/golang/pkgsite/blob/5f0513d53cff8382238b5f8c78e8317d2b4ad06d/internal/derrors/derrors.go#L240).
By using derrfirst, you can prevent public functions from returning without wrapping errors.

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

//lint:ignore derrfirst this is because ...
func Ignored() error {
  doOtherStuff()
  return nil
}
```

## How to use

Build your `derrfirst` binary by writing `main.go` like below.

```go
package main

import (
	"github.com/qawatake/derrfirst"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
  unitchecker.Main(
    derrfirst.NewAnalyzer("your/derror/pkg", "Wrap", "pkg/to/be/ignored"),
  )
}
```

Then, run `go vet` with your `derrfirst` binary.

```sh
go vet -vettool=/path/to/your/derrfirst ./...
```
