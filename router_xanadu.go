package nutshttp

// initXanaduRouter initializes the router for Xanadu crawler.
func (s *NutsHTTPServer) initXanaduRouter() {
	sr := s.r.Group("/xanadu")

	//search for a keyword
	sr.GET("/search/:keyword", s.Search)

}
