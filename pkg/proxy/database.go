package proxy

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kingsukhoi/wtf-inator/pkg/db"
	"github.com/kingsukhoi/wtf-inator/pkg/sqlc"
)

type requestDto struct {
	headers         map[string][]string
	queryParameters map[string][]string

	id          uuid.UUID
	method      string
	content     []byte
	sourceIp    string
	timestamp   time.Time
	requestPath string
}

func processRequest(request requestDto) {
	defer currWg.Done()
	ctx := context.Background()

	conn := db.MustGetDatabase()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	queries := sqlc.New(tx)

	err = queries.CreateRequest(ctx, sqlc.CreateRequestParams{
		ID:       request.id,
		Method:   request.method,
		Content:  request.content,
		SourceIp: request.sourceIp,
		Timestamp: pgtype.Timestamptz{
			Time:             request.timestamp,
			InfinityModifier: 0,
			Valid:            true,
		},
		RequestPath: request.requestPath,
	})
	if err != nil {
		slog.Error("Failed to create request: %v", err)
		return
	}

	headersArray := make([]sqlc.CreateRequestHeadersParams, 0)

	for k, v := range request.headers {
		for _, innerValue := range v {
			headersArray = append(headersArray, sqlc.CreateRequestHeadersParams{
				RequestID: request.id,
				Name:      k,
				Value: pgtype.Text{
					String: innerValue,
					Valid:  strings.TrimSpace(innerValue) != "",
				},
			})
		}
	}

	_, err = queries.CreateRequestHeaders(ctx, headersArray)
	if err != nil {
		slog.Error("Failed to create request headers: %v", err)
		return
	}

	queryParamsArray := make([]sqlc.CreateRequestQueryParametersParams, 0)

	for k, v := range request.queryParameters {
		for _, innerValue := range v {
			queryParamsArray = append(queryParamsArray, sqlc.CreateRequestQueryParametersParams{
				RequestID: request.id,
				Name:      k,
				Value: pgtype.Text{
					String: innerValue,
					Valid:  strings.TrimSpace(innerValue) != "",
				},
			})
		}
	}

	_, err = queries.CreateRequestQueryParameters(ctx, queryParamsArray)
	if err != nil {
		slog.Error("Failed to create request headers: %v", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("Failed to commit transaction: %v", err)
	}
}

type responseDto struct {
	headers      map[string][]string
	content      []byte
	id           uuid.UUID
	responseCode int32
	timestamp    time.Time
}

func processResponse(response responseDto) {
	defer currWg.Done()
	ctx := context.Background()

	conn := db.MustGetDatabase()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	queries := sqlc.New(tx)

	err = queries.CreateResponse(ctx, sqlc.CreateResponseParams{
		Requestid:       response.id,
		ResponseContent: response.content,
		ResponseCode:    response.responseCode,
		Timestamp: pgtype.Timestamptz{
			Time:             response.timestamp,
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		slog.Error("Failed to create response: %v", err)
		return
	}

	headersArray := make([]sqlc.CreateResponseHeadersParams, 0)
	for k, v := range response.headers {
		for _, innerValue := range v {
			headersArray = append(headersArray, sqlc.CreateResponseHeadersParams{
				Requestid: response.id,
				Name:      k,
				Value: pgtype.Text{
					String: innerValue,
					Valid:  strings.TrimSpace(innerValue) != "",
				},
			})
		}
	}

	_, err = queries.CreateResponseHeaders(ctx, headersArray)

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("Failed to commit transaction: %v", err)
	}
}
