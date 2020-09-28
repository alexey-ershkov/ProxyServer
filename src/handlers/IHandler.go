package handlers

type Handler interface {
	ProxyRequest() error
	Defer()
}
