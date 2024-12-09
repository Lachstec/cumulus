package service

import (
	"encoding/json"
	"fmt"
	"github.com/Lachstec/mc-hosting/internal/types"
	"io"
	"net/http"
	"net/url"
)

type AuthService struct {
	serviceUrl url.URL
	client     http.Client
}

func NewAuthService(auth0 url.URL) *AuthService {
	return &AuthService{
		serviceUrl: auth0,
		client:     http.Client{},
	}
}

func (s *AuthService) ValidateToken(token string) (*types.UserInfo, error) {
	req, err := http.NewRequest("GET", s.serviceUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userinfo types.UserInfo
	fmt.Println(string(body))
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		return nil, err
	}

	return &userinfo, nil
}
