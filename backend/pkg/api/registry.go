package api

// Interface for registering handlers
type HandlerRegistry interface {
	RegisterHandlers(s *Server)
}

type MainServerRegistry struct{}

// RegisterHandlers registers the routes of the main server API.
func (r MainServerRegistry) RegisterHandlers(s *Server) {
	unAuthed := s.router.Group("")
	{
		unAuthed.GET("/auth", s.authHandler)
		unAuthed.GET("/auth/callback", s.authCallbackHandler)
	}

	authed := s.router.Group("")
	{
		authed.Use(s.authMiddleware())

		authed.GET("/user", s.userHandler)
	}
}
