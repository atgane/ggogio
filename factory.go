package ggogio

// Factory is an interface that creates
// a user-specified Client implementation.
type Factory interface {
	Create() Client
}
