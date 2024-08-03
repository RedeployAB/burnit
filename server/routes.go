package server

func (s server) routes() {
	s.router.Handle("GET /secret", s.generateSecret())
	s.router.Handle("GET /secrets/{id}", s.getSecret())
	s.router.Handle("GET /secrets/{id}/{passphrase}", s.getSecret())
	s.router.Handle("POST /secrets", s.createSecret())
}
