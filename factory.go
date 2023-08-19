package ggogio

type Factory interface {
	Create() Client
}
