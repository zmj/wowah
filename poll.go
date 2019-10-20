package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type file struct {
	LastModified int64  `json:"lastModified"`
	URL          string `json:"url"`
}

func (s *session) poll() {
	if s == nil {
		s = &session{err: fmt.Errorf("nil session")}
	} else if s.err != nil {
	} else if s.accessToken == "" {
		s.err = fmt.Errorf("missing access token")
	}
	if s.err != nil {
		return
	}
	s.files, s.err = s.pollInternal()
}

func (s session) pollInternal() ([]file, error) {
	req, err := http.NewRequest("GET", "https://us.api.blizzard.com/wow/auction/data/malganis", nil)
	if err != nil {
		return nil, fmt.Errorf("create auction data request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+s.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send auction data request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auction data request failed: %v", resp.Status)
	}

	var files struct {
		Files []file `json:"files"`
	}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, fmt.Errorf("parse auction data response: %w", err)
	}
	return files.Files, nil
}
