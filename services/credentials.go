package services

type Credentials struct {
	Username string
	Password string
}

func NewCredentials(username, password string) *Credentials {
	return &Credentials{
		Username: username,
		Password: password,
	}
}

func (c *Credentials) GetCredentials() map[string]string {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}
}
