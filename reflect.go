package httpgo

import "github.com/gorilla/mux"

type globalHttp struct {
	router       *mux.Router
	serverConfig *ServerConfig
}

func (world *globalHttp) SetHttpApiMeta() {

}
