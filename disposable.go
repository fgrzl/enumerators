package enumerators

// Disposable represents an object that holds resources that need to be explicitly released.
// Implementers should ensure that Dispose can be called multiple times safely.
type Disposable interface {
	// Dispose releases any resources held by this object.
	// It should be safe to call Dispose multiple times.
	Dispose()
}
