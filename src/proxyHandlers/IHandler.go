package proxyHandlers

type Handler interface {
	ProxyRequest() error
	Defer()
}
