package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Logger struct {
	Config  LogConfig
	LogInfo func(ctx context.Context, msg string, fields map[string]interface{})
}

func NewLogger(c LogConfig, log func(ctx context.Context, msg string, fields map[string]interface{})) *Logger {
	return &Logger{c, log}
}
func (l *Logger) Log(ctx echo.Context, reqBody, resBody []byte) {
	msg := ctx.Request().Method + " " + ctx.Request().RequestURI
	Log(ctx, msg, reqBody, resBody, l.Config, l.LogInfo)
}

func Log(ctx echo.Context, msg string, reqBody, resBody []byte, c LogConfig, log func(ctx context.Context, msg string, fields map[string]interface{})) {
	if c.Log {
		req := ctx.Request()
		fields := BuildLogFields(c, ctx.Request())
		if len(c.Request) > 0 && len(reqBody) > 0 && req.Method != "GET" && req.Method != "DELETE" {
			fields[c.Request] = string(reqBody)
		}
		if len(c.Response) > 0 {
			fields[c.Response] = string(resBody)
		}
		if len(c.ResponseStatus) > 0 {
			fields[c.ResponseStatus] = ctx.Response().Status
		}
		/*
			if len(fieldConfig.Duration) > 0 {
				t2 := time.Now()
				duration := t2.Sub(t1)
				fields[fieldConfig.Duration] = duration.Milliseconds()
			}*/
		if len(c.Size) > 0 {
			fields[c.Size] = len(resBody)
		}
		log(req.Context(), msg, fields)
	}
}
func BuildLogFields(c LogConfig, r *http.Request) map[string]interface{} {
	fields := make(map[string]interface{}, 0)
	if !c.Build {
		return fields
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if len(c.Uri) > 0 {
		fields[c.Uri] = r.RequestURI
	}

	if len(c.ReqId) > 0 {
		if reqID := GetReqID(r.Context()); reqID != "" {
			fields[c.ReqId] = reqID
		}
	}
	if len(c.Scheme) > 0 {
		fields[c.Scheme] = scheme
	}
	if len(c.Proto) > 0 {
		fields[c.Proto] = r.Proto
	}
	if len(c.UserAgent) > 0 {
		fields[c.UserAgent] = r.UserAgent()
	}
	if len(c.RemoteAddr) > 0 {
		fields[c.RemoteAddr] = r.RemoteAddr
	}
	if len(c.Method) > 0 {
		fields[c.Method] = r.Method
	}
	if len(c.RemoteIp) > 0 {
		remoteIP := getRemoteIp(r)
		fields[c.RemoteIp] = remoteIP
	}
	return fields
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
