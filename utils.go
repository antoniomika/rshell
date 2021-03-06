package main

import (
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func logReq(ctx *fasthttp.RequestCtx) {
	statusCode := ctx.Response.StatusCode()
	method := string(ctx.Method())
	statusColor := colorForStatus(statusCode)
	methodColor := colorForMethod(method)
	latency := time.Now().Sub(ctx.Time())
	clientIP := ctx.RemoteIP().String()
	path := string(ctx.Path())
	// headers := string(ctx.Request.Header.RawHeaders())

	realIP := string(ctx.Request.Header.Peek("x-forwarded-for"))
	if realIP != "" {
		clientIP += ", " + realIP
	}

	log.Printf("|%s %3d %s| %13v | %15s |%s %-7s %s %s\n", statusColor, statusCode, reset,
		latency,
		clientIP,
		methodColor, method, reset,
		path)
}
