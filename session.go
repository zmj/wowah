package main

type session struct {
	oauthID     string
	oauthSecret string
	err         error

	accessToken string
	files       []file
}

func newSession(oauthID, oauthSecret string) *session {
	return &session{
		oauthID:     oauthID,
		oauthSecret: oauthSecret,
	}
}
