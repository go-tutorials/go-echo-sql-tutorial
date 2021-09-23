package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	LOG_FIELDS_KEY = "logFields"
)

func GetRequestBody(ctx context.Context, requestBody io.ReadCloser) (io.ReadCloser, context.Context) {
	if requestBody != nil {
		data, err := ioutil.ReadAll(requestBody)
		if err != nil {
			panic(err)
		}
		requestBody = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	return requestBody, ctx
}

func GetResponseData(ctx context.Context, path string, method string, statusCode int, startTime time.Time, brw *ResponseWriter) {
	endTime := time.Now()
	latency := endTime.Sub(startTime)

	fields := Fields{
		"path":         path,
		"method":       method,
		"statusCode":   statusCode,
		"responseTime": latency.Milliseconds(),
	}

	if brw != nil {
		fields["response"] = brw.Body.String()
		fields = AppendFields(ctx, fields)
	}

	result, err := json.Marshal(fields)
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(result)
	fmt.Println("")
	fmt.Println("____________________________")
}

func AppendFields(ctx context.Context, fields Fields) Fields {
	if logFields, ok := ctx.Value(LOG_FIELDS_KEY).(map[string]interface{}); ok {
		for k, v := range logFields {
			fields[k] = v
		}
	}

	return fields
}
