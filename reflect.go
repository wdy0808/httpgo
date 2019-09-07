package httpgo

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"github.com/wdy0808/go-common/log"
)

var httpMethodMap = map[string]string{
	"Delete": http.MethodDelete,
	"Get":    http.MethodGet,
	"Post":   http.MethodPost,
	"Patch":  http.MethodPatch,
}

var World globalHttp

type globalHttp struct {
	router       *mux.Router
	serverConfig *ServerConfig
}

func init() {
	World.router = mux.NewRouter()
}

func (world *globalHttp) SetHttpApiMeta(metaData HttpApiMeta) {
	value := reflect.ValueOf(metaData).FieldByName("CommonMetaData")
	if value.IsNil() {
		log.LogError("unable to register nill common metadata")
		return
	}
	api := metaData.CommonMetaData.Api
	if "" == api {
		log.LogError("api on common metadata is empty")
		return
	}
	var allowMethod = make([]string, 0, 4)
	for handlerType, method := range httpMethodMap {
		handlerValue := value.FieldByName(handlerType)
		if handlerValue.IsNil() {
			continue
		}
		handler, ok := handlerValue.Interface().(func(http.ResponseWriter, *http.Request))
		if !ok {
			log.LogError("fail to get http handler")
			return
		}
		World.router.HandleFunc(api, wrapHTTPHandler(handler)).Methods(method)
		allowMethod = append(allowMethod, method)
	}

	value = reflect.ValueOf(metaData).FieldByName("OptionMetaData")
	if value.IsNil() {
		return
	}
	headersArray, ok := value.FieldByName("AllowHeaders").Interface().([]string)
	if !ok {
		log.LogError("fail to get option header")
		return
	}
	World.router.HandleFunc(api, generateOptionHandler(headersArray, allowMethod)).Methods(http.MethodOptions)
}

func wrapHTTPHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		api := request.URL.Path
		currentTime := time.Now()
		log.LogInfo("receive http request [%v]", api)
		handler(response, request)
		log.LogInfo("handler request [%v] using time [%v] ms", api, time.Since(currentTime).Nanoseconds()/1000)
	}
}

func generateOptionHandler(allowedHeaders []string, allowedMethods []string) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Add("Access-Control-Allow-Origin", getAccessControlAllowOrigin(request.RequestURI))
		response.Header().Add("Access-Control-Allow-Methods", getAccessControlAllowMethod(request.Header.Get("Access-Control-Request-Method"), allowedMethods))
		response.Header().Add("Access-Control-Allow-Headers", getAccessControlAllowHeader(request.Header.Get("Access-Control-Request-Headers"), allowedHeaders))
	}
}
