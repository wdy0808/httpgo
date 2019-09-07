package httpgo

import (
	"net/http"
	"strings"
	"time"
)

type ServerConfig struct {
	Port          int64
	WriteTimeout  time.Duration
	ReadTimeout   time.Duration
	AllowedOrigin []string
}

type CommonMeta struct {
	Api                      string
	Get, Post, Delete, Patch func(http.ResponseWriter, *http.Request)
}

type OptionMeta struct {
	AllowHeaders []string
}

type HttpApiMeta struct {
	CommonMetaData *CommonMeta
	OptionMetaData *OptionMeta
}

func getAccessControlAllowOrigin(hostName string) string {
	if nil == World.serverConfig {
		return ""
	}
	for _, origin := range World.serverConfig.AllowedOrigin {
		if hostName == origin {
			return hostName
		}
	}
	return ""
}

func getAccessControlAllowMethod(requestMethod string, methods []string) string {
	for _, method := range methods {
		if requestMethod == method {
			return requestMethod
		}
	}
	return ""
}

func getAccessControlAllowHeader(requestHeaders string, allowedHeaders []string) string {
	if nil == World.serverConfig {
		return ""
	}
	var requestHeaderArray = strings.Split(requestHeaders, ",")
	var headers = make([]string, 0, len(requestHeaders))
	for _, header := range requestHeaderArray {
		for _, allowedHeader := range allowedHeaders {
			if header == allowedHeader {
				headers = append(headers, header)
			}
		}
	}
	return strings.Join(headers, ",")
}

func RunForver(config ServerConfig) {
	World.serverConfig = &ServerConfig{
		Port:          config.Port,
		ReadTimeout:   config.ReadTimeout,
		WriteTimeout:  config.WriteTimeout,
		AllowedOrigin: config.AllowedOrigin,
	}
	server := &http.Server{
		Handler:      World.router,
		Addr:         "127.0.0.1:10486",
		WriteTimeout: World.serverConfig.WriteTimeout,
		ReadTimeout:  World.serverConfig.ReadTimeout,
	}
	server.ListenAndServe()
}
