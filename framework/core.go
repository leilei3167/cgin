package framework

import "net/http"

type Core struct {
}

func NewCore() *Core {
	return new(Core)
}
func (c Core) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
