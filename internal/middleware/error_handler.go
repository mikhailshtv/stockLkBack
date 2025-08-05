package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/mikhailshtv/stockLkBack/pkg/errors"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log := logger.GetLogger()
		traceID := getTraceID(c)

		if err, ok := recovered.(string); ok {
			log.Error("panic recovered",
				zap.String("trace_id", traceID),
				zap.String("panic", err),
				zap.String("stack", string(debug.Stack())),
			)
		} else {
			log.Error("panic recovered",
				zap.String("trace_id", traceID),
				zap.Any("panic", recovered),
				zap.String("stack", string(debug.Stack())),
			)
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:    "INTERNAL_ERROR",
			Message: "Внутренняя ошибка сервера",
			Code:    http.StatusInternalServerError,
		})
	})
}

func HandleError(c *gin.Context, err error) {
	log := logger.GetLogger()
	traceID := getTraceID(c)

	if appErr, ok := errors.IsAppError(err); ok {
		log.Error("application error",
			zap.String("trace_id", traceID),
			zap.String("error_type", string(appErr.Type)),
			zap.String("message", appErr.Message),
			zap.Error(appErr.Internal),
		)
		c.JSON(appErr.Code, ErrorResponse{
			Type:    string(appErr.Type),
			Message: appErr.Message,
			Code:    appErr.Code,
		})
		return
	}

	log.Error("unhandled error",
		zap.String("trace_id", traceID),
		zap.Error(err),
	)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Type:    "INTERNAL_ERROR",
		Message: "Внутренняя ошибка сервера",
		Code:    http.StatusInternalServerError,
	})
}

func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return "unknown"
}
