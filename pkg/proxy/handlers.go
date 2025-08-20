package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/kingsukhoi/wtf-inator/pkg/conf"
	"github.com/kingsukhoi/wtf-inator/pkg/constants"
	"github.com/kingsukhoi/wtf-inator/pkg/helpers"
	"github.com/labstack/echo/v4"
)

func responseHandler(resp *http.Response) error {

	requestId, exists := resp.Request.Context().Value(constants.ContextRequestIdKey).(uuid.UUID)
	if !exists {
		return nil
	}

	clonedBody, clonedReadCloser, err2 := helpers.CloneReadCloser(resp.Body)
	if err2 != nil {
		return err2
	}

	resp.Body = clonedReadCloser

	sqlResponse := responseDto{
		id:           requestId,
		content:      clonedBody,
		responseCode: int32(resp.StatusCode),
		timestamp:    time.Now(),
		headers:      resp.Header,
	}

	currWg.Add(1)

	go processResponse(sqlResponse)

	return nil
}

func (wtfProxy *WtfProxy) RequestHandler(c echo.Context) error {

	config := conf.MustGetConfig()
	if slices.Contains(config.Server.IgnoredPaths, c.Request().URL.Path) {
		wtfProxy.Proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}

	currId, _ := uuid.NewV7()

	clonedRequestBody, RequestReader, err := helpers.CloneReadCloser(c.Request().Body)
	if err != nil {
		return err
	}
	c.Request().Body = RequestReader

	slog.Debug("RequestBody", "body", string(clonedRequestBody), "requestId", currId.String())

	ctx := context.WithValue(c.Request().Context(), constants.ContextRequestIdKey, currId)

	newReq := c.Request().WithContext(ctx)

	reqDto := requestDto{
		headers:         c.Request().Header,
		queryParameters: c.QueryParams(),
		id:              currId,
		method:          c.Request().Method,
		content:         clonedRequestBody,
		sourceIp:        c.RealIP(),
		timestamp:       time.Now(),
		requestPath:     c.Request().URL.Path,
	}
	currWg.Add(1)
	go processRequest(reqDto)

	wtfProxy.Proxy.ServeHTTP(c.Response().Writer, newReq)

	return nil
}
