package http

import (
	"database/sql"
	"encoding/json"
	"errors"

	"TestTaskInAvito/internal/domain"
	stdhttp "net/http"
)

func writeJSON(w stdhttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w stdhttp.ResponseWriter, err error) {
	// доменная ошибка
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
		writeJSON(w, stdhttp.StatusNotFound, resp)
		return
	}

	// всё остальное → 500
	resp := ErrorResponseDTO{
		Error: ErrorBodyDTO{
			Code:    "INTERNAL",
			Message: "internal server error",
		},
	}
	writeJSON(w, stdhttp.StatusInternalServerError, resp)
}

func httpStatusByCode(code domain.ErrorCode) int {
	switch code {
	case domain.ErrorCodeTeamExists:
		return stdhttp.StatusBadRequest
	case domain.ErrorCodePRExists:
		return stdhttp.StatusConflict
	case domain.ErrorCodePRMerged,
		domain.ErrorCodeNotAssigned,
		domain.ErrorCodeNoCandidate:
		return stdhttp.StatusConflict
	case domain.ErrorCodeNotFound:
		return stdhttp.StatusNotFound
	default:
		return stdhttp.StatusInternalServerError
	}
}
