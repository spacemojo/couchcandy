package couchcandy

// Session : holds the connection data for a couchcandy session.
type Session struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewSession : creates a new session initialized with the passed values
func NewSession(host string, port int, username string, password string) *Session {
	session := &Session{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	return session
}
