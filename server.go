package web

type ServerConfig struct {
	Port         int
	RecoverPanic bool
}

type Server struct {
}

func NewServer() *Server {

	return &Server{}
}

func (s *Server) AddRoute(route, method string, handler interface{}) {

}
