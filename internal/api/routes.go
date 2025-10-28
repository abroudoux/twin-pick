package api

func (s *server) registerRoutes() {
	api := s.router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/match", s.handleMatch)
}
