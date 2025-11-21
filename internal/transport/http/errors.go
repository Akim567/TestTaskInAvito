package http

import (
	"TestTaskInAvito/internal/domain"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	// если это доменная ошибка — маппим код и http-статус
	var dErr *domain.DomainError
	if errors.As(err, &dErr) {
		status := httpStatusByCode(dErr.Code)
		resp := ErrorResponseDTO{
			Error: ErrorBodyDTO{
				Code:    string(dErr.Code),
				Message: dErr.Message,
			},
		}
		writeJSON(w, status, resp)
		return
	}

	// sql.ErrNoRows → NOT_FOUND
	if errors.Is(err, sql.ErrNoRows) {
		resp := ErrorResponseDTO{
			Error: ErrorBodyDTO{
				Code:    string(domain.ErrorCodeNotFound),
				Message: "resource not found",
			},
		}
		writeJSON(w, http.StatusNotFound, resp)
		return
	}

	// все остальные — 500
	resp := ErrorResponseDTO{
		Error: ErrorBodyDTO{
			Code:    "INTERNAL",
			Message: "internal server error",
		},
	}
	writeJSON(w, http.StatusInternalServerError, resp)
}

func httpStatusByCode(code domain.ErrorCode) int {
	switch code {
	case domain.ErrorCodeTeamExists:
		return http.StatusBadRequest // 400
	case domain.ErrorCodePRExists:
		return http.StatusConflict // 409
	case domain.ErrorCodePRMerged,
		domain.ErrorCodeNotAssigned,
		domain.ErrorCodeNoCandidate:
		return http.StatusConflict // 409
	case domain.ErrorCodeNotFound:
		return http.StatusNotFound // 404
	default:
		return http.StatusInternalServerError
	}
}
