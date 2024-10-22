package api

// Interface for registering handlers
type HandlerRegistry interface {
	RegisterHandlers(s *Server)
}

type MainServerRegistry struct{}

// RegisterHandlers registers the routes of the main server API.
func (r MainServerRegistry) RegisterHandlers(s *Server) {
	s.router.GET("/auth", s.authHandler)
	s.router.GET("/auth/callback", s.authCallbackHandler)
	s.router.GET("/user", s.userHandler)
}
