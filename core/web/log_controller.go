package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"go.uber.org/zap/zapcore"
)

// LogController manages the logger config
type LogController struct {
	App chainlink.Application
}

type LogPatchRequest struct {
	Level           string      `json:"level"`
	Filter          string      `json:"filter"`
	SqlEnabled      *bool       `json:"sqlEnabled"`
	ServiceLogLevel [][2]string `json:"serviceLogLevel"`
}

// Get retrieves the current log config settings
func (cc *LogController) Get(c *gin.Context) {
	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		Level:      cc.App.GetStore().Config.LogLevel().String(),
		SqlEnabled: cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}

// Patch sets a log level and enables sql logging for the logger
func (cc *LogController) Patch(c *gin.Context) {
	request := &LogPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if request.Level == "" && request.Filter == "" && request.SqlEnabled == nil && len(request.ServiceLogLevel) == 0 {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("please check request params, no params configured"))
		return
	}

	if request.Level != "" {
		var ll zapcore.Level
		err := ll.UnmarshalText([]byte(request.Level))
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		if err = cc.App.GetStore().Config.SetLogLevel(c.Request.Context(), ll.String()); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	if len(request.ServiceLogLevel) > 0 {
		var svc, lvl []string
		for _, svcLogLvl := range request.ServiceLogLevel {
			svcName := svcLogLvl[0]
			svcLvl := svcLogLvl[1]
			var level zapcore.Level
			if err := level.UnmarshalText([]byte(svcLvl)); err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}

			if err := cc.App.SetServiceLogger(c.Request.Context(), svcName, level); err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}

			ll, err := cc.App.GetLogger().ServiceLogLevel(svcName)
			if err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}

			svc = append(svc, svcName)
			lvl = append(lvl, ll)
		}

		response := &presenters.ServiceLevelLog{
			JAID: presenters.JAID{
				ID: "log",
			},
			ServiceName: strings.Join(svc, ","),
			LogLevel:    strings.Join(lvl, ","),
		}

		jsonAPIResponse(c, response, "log")
		return
	}

	if request.SqlEnabled != nil {
		if err := cc.App.GetStore().Config.SetLogSQLStatements(c.Request.Context(), *request.SqlEnabled); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().SetLogging(*request.SqlEnabled)
	}

	// Set default logger with new configurations
	logger.SetLogger(cc.App.GetStore().Config.CreateProductionLogger())
	cc.App.GetLogger().SetDB(cc.App.GetStore().DB)

	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		Level:      cc.App.GetStore().Config.LogLevel().String(),
		SqlEnabled: cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}
