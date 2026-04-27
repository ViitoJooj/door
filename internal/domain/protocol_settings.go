package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	ProtocolModeHTTP  = "http"
	ProtocolModeHTTPS = "https"
	ProtocolModeBoth  = "both"

	ConfigScopeAll      = "all"
	ConfigScopeExternal = "external"
)

type ProtocolSettings struct {
	ID              int
	AllowedProtocol string
	ApplyScope      string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

func NormalizeProtocolMode(value string) (string, error) {
	mode := strings.ToLower(strings.TrimSpace(value))
	switch mode {
	case ProtocolModeHTTP, ProtocolModeHTTPS, ProtocolModeBoth:
		return mode, nil
	default:
		return "", errors.New("allowed_protocol must be one of: http, https, both")
	}
}

func NormalizeConfigScope(value string) (string, error) {
	scope := strings.ToLower(strings.TrimSpace(value))
	switch scope {
	case ConfigScopeAll, ConfigScopeExternal:
		return scope, nil
	default:
		return "", errors.New("apply_scope must be one of: all, external")
	}
}

func NewProtocolSettings(allowedProtocol string, applyScope string) (*ProtocolSettings, error) {
	mode, err := NormalizeProtocolMode(allowedProtocol)
	if err != nil {
		return nil, err
	}
	scope, err := NormalizeConfigScope(applyScope)
	if err != nil {
		return nil, err
	}

	return &ProtocolSettings{
		ID:              1,
		AllowedProtocol: mode,
		ApplyScope:      scope,
	}, nil
}
