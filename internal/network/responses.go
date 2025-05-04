package network

import (
	"encoding/json"
	"github.com/quantum73/revizzoro-api/internal/enums"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WriteResponse(
	w http.ResponseWriter,
	status int,
	contentType enums.HTTPContentType,
	payload []byte,
) {
	const op = "[network/responses WriteResponse]"

	w.Header().Set("Content-Type", contentType.String())
	w.WriteHeader(status)
	_, err := w.Write(payload)
	if err != nil {
		log.Warnf("%s Error writing response: %v\n", op, err)
	}
}

func JSONResponse(w http.ResponseWriter, status int, payload map[string]any) {
	const op = "[network/responses JSONResponse]"

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Warnf("%s Error marshalling payload: %s\n", op, err)
	}
	WriteResponse(w, status, enums.ContentTypeJSON, bytes)
}

func MessageJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	payload := map[string]any{"message": message}
	JSONResponse(w, statusCode, payload)
}

func OKMessageResponse(w http.ResponseWriter, msg string) {
	MessageJSONResponse(w, http.StatusOK, msg)
}

func NotFoundMessageResponse(w http.ResponseWriter, msg string) {
	MessageJSONResponse(w, http.StatusNotFound, msg)
}

func BadRequestMessageResponse(w http.ResponseWriter, msg string) {
	MessageJSONResponse(w, http.StatusBadRequest, msg)
}

func ServerUnexpectedErrorMessageResponse(w http.ResponseWriter) {
	MessageJSONResponse(w, http.StatusInternalServerError, "unexpected error")
}
