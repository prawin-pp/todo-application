package auth

type Database interface {
}

type Server struct {
	db Database
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
