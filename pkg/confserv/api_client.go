package confserv

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type ConfServClient interface {
	GetString(ctx context.Context, setting string) (string, error)
	GetInt(ctx context.Context, setting string) (int, error)
	GetBool(ctx context.Context, setting string) (bool, error)
	GetFloat(ctx context.Context, setting string) (float64, error)
	GetDuration(ctx context.Context, setting string) (time.Duration, error)
	GetObjBinded(ctx context.Context, key string, dest any) error
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

func (s *confServClientImpl) getRawValue(ctx context.Context, key string, dest any) error {
	params := GetSettingByKeyParams{
		Authorization: "Bearer " + s.token,
	}
	httpResponse, err := s.rawClient.GetSettingByKey(ctx, strings.ToLower(key), &params)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return errors.New("internal server error")
	}
	err = json.NewDecoder(httpResponse.Body).Decode(&dest)
	if err != nil {
		return err
	}
	return nil
}

func (s *confServClientImpl) GetString(ctx context.Context, key string) (string, error) {
	var result string
	if err := s.getRawValue(ctx, key, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (s *confServClientImpl) GetInt(ctx context.Context, key string) (int, error) {
	floatValue, err := s.GetFloat(ctx, key)
	if err != nil {
		return 0, err
	}
	return int(floatValue), nil
}

func (s *confServClientImpl) GetDuration(ctx context.Context, key string) (time.Duration, error) {
	floatValue, err := s.GetFloat(ctx, key)
	if err != nil {
		return 0, err
	}
	return time.Duration(floatValue), nil
}

func (s *confServClientImpl) GetFloat(ctx context.Context, key string) (float64, error) {
	var result float64
	if err := s.getRawValue(ctx, key, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (s *confServClientImpl) GetBool(ctx context.Context, key string) (bool, error) {
	var result bool
	if err := s.getRawValue(ctx, key, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (s *confServClientImpl) GetObjBinded(ctx context.Context, key string, dest any) error {
	if err := s.getRawValue(ctx, key, &dest); err != nil {
		return err
	}
	return nil
}
