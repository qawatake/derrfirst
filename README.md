# derrfirst

Linter `derrfirst` forces every public function to begin with a deferring call of some function like [derrors.Wrap](https://github.com/golang/pkgsite/blob/5f0513d53cff8382238b5f8c78e8317d2b4ad06d/internal/derrors/derrors.go#L240).
By using derrfirst, you can prevent public functions from returning without wrapping errors.

```go
func F() error { // <-  should call defer derrors.Wrap fist.
  doOtherStuff()
  return nil
}
```
