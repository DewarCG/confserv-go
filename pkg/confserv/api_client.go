package confserv

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type ConfServClient interface {
	GetString(setting string) (string, error)
	GetInt(setting string) (int, error)
	GetBool(setting string) (bool, error)
	GetFloat(setting string) (float64, error)
	GetDuration(setting string) (time.Duration, error)
	GetObjBinded(settingName string, dest any) error
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
	str, err := s.GetString(settingName)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func (s *confServClientImpl) GetDuration(settingName string) (time.Duration, error) {
	intValue, err := s.GetInt(settingName)
	if err != nil {
		return 0, err
	}
	return time.Duration(intValue), nil
}

func (s *confServClientImpl) GetFloat(settingName string) (float64, error) {
	str, err := s.GetString(settingName)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(str, 64)
}

func (s *confServClientImpl) GetBool(settingName string) (bool, error) {
	str, err := s.GetString(settingName)
	if err != nil {
		return false, err
	}
	if str == "true" || str == "on" || str == "yes" || str == "1" {
		return true, nil
	}
	if str == "false" || str == "off" || str == "no" || str == "0" {
		return false, nil
	}
	return false, errors.New("invalid value")
}

func (s *confServClientImpl) GetObjBinded(settingName string, dest any) error {
	response, err := s.getRaw(context.TODO(), settingName)
	if err != nil {
		return err
	}
	if response.Value == nil {
		return nil
	}
	return json.Unmarshal([]byte(*response.Value), dest)
}
