package main

import (
	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	wrapper.SetCtx(
		"agent-session-monitor",
		wrapper.ParseConfigBy(parseConfig),
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
		wrapper.ProcessResponseHeadersBy(onHttpResponseHeaders),
	)
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config PluginConfig, log wrapper.Log) types.Action {
	sessionID := extractSessionID(ctx)
	if sessionID == "" {
		log.Debugf("no session id found in request")
		return types.ActionContinue
	}

	ctx.SetContext("session_id", sessionID)
	ctx.SetContext("request_start_ms", getCurrentTimeMs())

	log.Infof("agent session started: session_id=%s", sessionID)
	return types.ActionContinue
}

func onHttpResponseHeaders(ctx wrapper.HttpContext, config PluginConfig, log wrapper.Log) types.Action {
	sessionID, ok := ctx.GetContext("session_id").(string)
	if !ok || sessionID == "" {
		return types.ActionContinue
	}

	status, err := proxywasm.GetHttpResponseHeader(":status")
	if err != nil {
		log.Warnf("failed to get response status: %v", err)
		return types.ActionContinue
	}

	recordSessionMetric(ctx, config, sessionID, status, log)
	return types.ActionContinue
}
