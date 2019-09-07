package httpgo

import (
	"testing"
)

func Test_getAccessControlAllowOrigin(t *testing.T) {
	type args struct {
		hostName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"successfully", args{"localhost:3000"}, "localhost:3000"},
		{"failed", args{"none.com"}, ""},
	}
	World.serverConfig = &ServerConfig{AllowedOrigin: []string{"localhost:3000", "wdy0808.com"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAccessControlAllowOrigin(tt.args.hostName); got != tt.want {
				t.Errorf("getAccessControlAllowOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAccessControlAllowMethod(t *testing.T) {
	type args struct {
		requestMethod string
		methods       []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"GETTEST", args{"GET", []string{"GET", "POST"}}, "GET"},
		{"POSTTEST", args{"POST", []string{"GET", "POST"}}, "POST"},
		{"DELETEDTEST", args{"DELETED", []string{"GET", "POST"}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAccessControlAllowMethod(tt.args.requestMethod, tt.args.methods); got != tt.want {
				t.Errorf("getAccessControlAllowMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAccessControlAllowHeader(t *testing.T) {
	type args struct {
		requestHeaders string
		allowedHeaders []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"singleSuccess", args{"person", []string{"person"}}, "person"},
		{"singleMultiSuccess", args{"person", []string{"person", "blog"}}, "person"},
		{"multiSingleSuccess", args{"person,blog", []string{"blog"}}, "blog"},
		{"multiSuccess", args{"person,blog,server", []string{"person", "server", "none"}}, "person,server"},
		{"none", args{"person,blog,server", []string{"personnone", "servernone", "none"}}, ""},
	}
	World.serverConfig = &ServerConfig{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAccessControlAllowHeader(tt.args.requestHeaders, tt.args.allowedHeaders); got != tt.want {
				t.Errorf("getAccessControlAllowHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
