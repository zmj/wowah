package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *session) login() {
	if s == nil {
		s = &session{err: fmt.Errorf("nil session")}
	} else if s.err != nil {
	} else if s.oauthID == "" {
		s.err = fmt.Errorf("missing oauthID")
	} else if s.oauthSecret == "" {
		s.err = fmt.Errorf("missing oauthSecret")
	}
	if s.err != nil {
		return
	}
	s.accessToken, s.err = s.loginInternal()
}

func (s session) loginInternal() (string, error) {
	req, err := http.NewRequest("POST", "https://us.battle.net/oauth/token?grant_type=client_credentials", nil)
	if err != nil {
		return "", fmt.Errorf("create access token request: %w", err)
	}
	pwd := s.oauthID + ":" + s.oauthSecret
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(pwd))
	req.Header.Add("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send access token request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("access token request failed: %v", resp.Status)
	}

	var token struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", fmt.Errorf("parse access token response: %w", err)
	}
	return token.AccessToken, nil
}
