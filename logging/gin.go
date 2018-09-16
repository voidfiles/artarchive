package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type GinLogger struct {
	name   string
	logger zerolog.Logger
}

type ginHands struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

// func ErrorLogger() gin.HandlerFunc {
// 	return ErrorLoggerT(gin.ErrorTypeAny)
// }
//
// func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()
//
// 		if !c.Writer.Written() {
// 			json := c.Errors.ByType(typ).JSON()
// 			if json != nil {
// 				c.JSON(-1, json)
// 			}
// 		}
// 	}
// }

func MustNewGinLogger(logger zerolog.Logger, name string) *GinLogger {
	return &GinLogger{
		name:   name,
		logger: logger,
	}
}

func (g *GinLogger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		// after request
		// latency := time.Since(t)
		// clientIP := c.ClientIP()
		// method := c.Request.Method
		// statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}
		cData := &ginHands{
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
			MsgStr:     msg,
		}

		g.logSwitch(cData)
	}
}

func (g *GinLogger) logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			g.logger.Warn().Str("name", g.name).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			g.logger.Error().Str("name", g.name).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	default:
		g.logger.Info().Str("name", g.name).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}
