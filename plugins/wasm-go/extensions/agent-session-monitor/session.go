package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

const (
	MetricSessionTotal    = "agent_session_total"
	MetricSessionDuration = "agent_session_duration_ms"
	MetricSessionErrors   = "agent_session_errors_total"
)

func extractSessionID(ctx wrapper.HttpContext) string {
	// Try header first
	headerVal, err := proxywasm.GetHttpRequestHeader("x-session-id")
	if err == nil && headerVal != "" {
		return sanitizeSessionID(headerVal)
	}

	// Try query parameter
	path, err := proxywasm.GetHttpRequestHeader(":path")
	if err != nil {
		return ""
	}

	if idx := strings.Index(path, "session_id="); idx != -1 {
		val := path[idx+len("session_id="):]
		if end := strings.IndexAny(val, "&# "); end != -1 {
			val = val[:end]
		}
		return sanitizeSessionID(val)
	}

	return ""
}

func sanitizeSessionID(id string) string {
	id = strings.TrimSpace(id)
	if len(id) > 128 {
		id = id[:128]
	}
	return id
}

func getCurrentTimeMs() int64 {
	return time.Now().UnixMilli()
}

func recordSessionMetric(ctx wrapper.HttpContext, config PluginConfig, sessionID, status string, log wrapper.Log) {
	startMs, ok := ctx.GetContext("request_start_ms").(int64)
	if !ok {
		return
	}

	durationMs := getCurrentTimeMs() - startMs
	isError := strings.HasPrefix(status, "4") || strings.HasPrefix(status, "5")

	log.Infof("agent session completed: session_id=%s status=%s duration_ms=%d",
		sessionID, status, durationMs)

	if config.EnableMetrics {
		proxywasm.SetProperty(
			[]string{"session_monitor", "duration_ms"},
			[]byte(fmt.Sprintf("%d", durationMs)),
		)
	}

	if isError {
		log.Warnf("agent session error: session_id=%s status=%s", sessionID, status)
	}

	if durationMs > config.MaxSessionDurationSec*1000 {
		log.Warnf("agent session exceeded max duration: session_id=%s duration_ms=%d max_ms=%d",
			sessionID, durationMs, config.MaxSessionDurationSec*1000)
	}
}
