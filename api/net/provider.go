package net

type InterfacesProvider interface {
	Interfaces() ([]Interface, error)
}

type InterfacesProviderFunc func() ([]Interface, error)

func (f InterfacesProviderFunc) Interfaces() ([]Interface, error) {
	return f()
}
