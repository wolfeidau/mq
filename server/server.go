package server

import (
	"log"
	"net"

	"github.com/drone/mq/stomp"
)

// Server ...
type Server struct {
	router *router
}

// NewServer returns a new STOMP server.
func NewServer() *Server {
	return &Server{
		router: newRouter(),
	}
}

// Serve accepts incoming net.Conn requests.
func (s *Server) Serve(conn net.Conn) {
	log.Printf("stomp: successfully opened socket connection.")

	session := requestSession()
	session.peer = stomp.Conn(conn)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("stomp: unexpected panic: %s", r)
		}
		log.Printf("stomp: successfully closed socket connection.")

		s.router.disconnect(session)
		session.peer.Close()
		releaseSession(session)
	}()

	if err := s.router.serve(session); err != nil {
		log.Printf("stomp: server error. %s", err)
	}
}