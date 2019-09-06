package httpgo

type ServerConfig struct {
	Port          int
	WriteTimeout  int
	ReadTimeout   int
	AllowedOrigin []string
}
