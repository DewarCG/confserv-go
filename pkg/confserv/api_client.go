package confserv

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type ConfServClient interface {
	GetString(setting string) (string, error)
	GetInt(setting string) (int, error)
	GetBool(setting string) (bool, error)
	GetFloat(setting string) (float64, error)
}

func NewConfServClient(server string, token string) (ConfServClient, error) {
	var finalServer string
	if len(server) == 0 {
		finalServer = "http://127.0.0.1:1319"
	} else {
		finalServer = server
	}
	rawClient, err := NewClient(finalServer)
	if err != nil {
		return nil, err
	}
	client := &confServClientImpl{
		rawClient: rawClient,
		token:     token,
	}
	return client, nil
}

type confServClientImpl struct {
	rawClient *Client
	token     string
}

func (s *confServClientImpl) getRaw(ctx context.Context, setting string) (*SettingResponse, error) {
	params := GetSettingByNameParams{
		Token: s.token,
	}
	httpResponse, err := s.rawClient.GetSettingByName(ctx, strings.ToLower(setting), &params)
	if err != nil {
		return nil, err
	}
	var response SettingResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *confServClientImpl) GetString(settingName string) (string, error) {
	ctx := context.TODO()
	setting, err := s.getRaw(ctx, settingName)
	if err != nil {
		return "", err
	}
	empty := ""
	if setting == nil {
		return empty, nil
	}
	return *setting.Value, nil
}

func (s *confServClientImpl) GetInt(settingName string) (int, error) {
	s2, err := s.GetString(settingName)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s2)
}

func (s *confServClientImpl) GetFloat(settingName string) (float64, error) {
	s2, err := s.GetString(settingName)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s2, 64)
}

func (s *confServClientImpl) GetBool(settingName string) (bool, error) {
	s2, err := s.GetString(settingName)
	if err != nil {
		return false, err
	}
	if s2 == "true" || s2 == "on" || s2 == "yes" || s2 == "1" {
		return true, nil
	}
	if s2 == "false" || s2 == "off" || s2 == "no" || s2 == "0" {
		return false, nil
	}
	return false, errors.New("invalid value")
}
