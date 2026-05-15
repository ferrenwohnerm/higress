package main

import (
	"errors"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

type PluginConfig struct {
	// Header name used to identify session ID
	SessionIDHeader string `json:"session_id_header"`
	// Query param name used to identify session ID
	SessionIDQuery string `json:"session_id_query"`
	// Maximum session duration in seconds before alerting
	MaxSessionDurationSec int64 `json:"max_session_duration_sec"`
	// Whether to emit metrics
	EnableMetrics bool `json:"enable_metrics"`
	// Redis service for session state storage
	RedisServiceName string `json:"redis_service_name"`
	RedisServicePort int `json:"redis_service_port"`
}

func parseConfig(json gjson.Result, config *PluginConfig, log wrapper.Log) error {
	config.SessionIDHeader = json.Get("session_id_header").String()
	if config.SessionIDHeader == "" {
		config.SessionIDHeader = "X-Session-ID"
	}

	config.SessionIDQuery = json.Get("session_id_query").String()
	if config.SessionIDQuery == "" {
		config.SessionIDQuery = "session_id"
	}

	config.MaxSessionDurationSec = json.Get("max_session_duration_sec").Int()
	if config.MaxSessionDurationSec == 0 {
		config.MaxSessionDurationSec = 3600
	}

	config.EnableMetrics = json.Get("enable_metrics").Bool()

	config.RedisServiceName = json.Get("redis_service_name").String()
	config.RedisServicePort = int(json.Get("redis_service_port").Int())
	if config.RedisServicePort == 0 {
		config.RedisServicePort = 6379
	}

	if config.RedisServiceName == "" {
		log.Warnf("redis_service_name not set, session state persistence disabled")
	}

	log.Infof("agent-session-monitor config loaded: header=%s, max_duration=%ds",
		config.SessionIDHeader, config.MaxSessionDurationSec)

	return nil
}

func validateConfig(config *PluginConfig) error {
	if config.MaxSessionDurationSec < 0 {
		return errors.New("max_session_duration_sec must be non-negative")
	}
	return nil
}
