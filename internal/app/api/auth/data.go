package auth

type SessionData struct {
	Subject string
}

func (s *SessionData) MarshalBinary() ([]byte, error) {
	return []byte(s.Subject), nil
}

func (s *SessionData) UnmarshalBinary(data []byte) error {
	s.Subject = string(data)
	return nil
}
