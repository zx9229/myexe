package wsnet

import (
	"net/http"
)

//SimpleHttpServer 我不想使用全局的http就包了这么个类
type SimpleHttpServer struct {
	httpServer http.Server
}

//New_SimpleHttpServer omit
func New_SimpleHttpServer(listenAddr string) *SimpleHttpServer {
	curData := new(SimpleHttpServer)
	//
	curData.httpServer.Addr = listenAddr
	curData.httpServer.Handler = http.NewServeMux()
	//
	return curData
}

//GetHttpServeMux omit
func (thls *SimpleHttpServer) GetHttpServeMux() *http.ServeMux {
	return thls.httpServer.Handler.(*http.ServeMux)
}

//Run omit
func (thls *SimpleHttpServer) Run() error {
	return thls.httpServer.ListenAndServe()
}

//RunTLS omit
func (thls *SimpleHttpServer) RunTLS(certFile string, keyFile string) error {
	return thls.httpServer.ListenAndServeTLS(certFile, keyFile)
}
