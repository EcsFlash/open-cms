package handler

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLogger(log slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			req := c.Request()
			res := c.Response()
			reqID := res.Header().Get(echo.HeaderXRequestID)

			attrs := []any{
				"request_id", reqID,
				"method", req.Method,
				"path", c.Path(),
				"uri", req.URL.Path,
				"status", res.Status,
				"latency_ms", time.Since(start).Milliseconds(),
				"remote_ip", c.RealIP(),
				"user_agent", req.UserAgent(),
			}
			if err != nil {
				log.Error("http request", append(attrs, "err", err.Error())...)
				return err
			}
			log.Info("http request", attrs...)
			return nil
		}
	}
}

