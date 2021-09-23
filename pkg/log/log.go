package middleware

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Fields map[string]interface{}

type Logger struct {
	Config  LogConfig
	LogInfo func(ctx context.Context, msg string, fields map[string]interface{})
}

type ResponseWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func NewLogger(c LogConfig, log func(ctx context.Context, msg string, fields map[string]interface{})) *Logger {
	return &Logger{c, log}
}

func (l *Logger) Log(ctx echo.Context, reqBody, resBody []byte) {
	msg := ctx.Request().Method + " " + ctx.Request().RequestURI
	Log(ctx, msg, reqBody, resBody, l.Config, l.LogInfo)
	fmt.Println("_____________________________")
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
		// if len(fieldConfig.Duration) > 0 {
		// 	t2 := time.Now()
		// 	duration := t2.Sub(t1)
		// 	fields[fieldConfig.Duration] = duration.Milliseconds()
		// }
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

// customized response logging echo middleware
func LoggerEcho(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		var path string
		var method string
		var brw *ResponseWriter
		path = c.Request().URL.Path
		method = c.Request().Method

		// get request body
		reader, ctx := GetRequestBody(c.Request().Context(), c.Request().Body)
		c.SetRequest(c.Request().WithContext(ctx))
		c.Request().Body = reader
		brw = NewResponseWriter(c.Response().Writer)
		c.Response().Writer = brw

		err := next(c)

		// get response data
		GetResponseData(c.Request().Context(), path, method, c.Response().Status, start, brw)
		return err
	}
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{Body: bytes.NewBufferString(""), ResponseWriter: rw}
}
