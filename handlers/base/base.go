package base

import (
	"database/sql"
	"github.com/quantum73/revizzoro-api/internal/network"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RootController struct {
	db *sql.DB
}

func NewRootController(database *sql.DB) *RootController {
	return &RootController{db: database}
}

func (c *RootController) Home(w http.ResponseWriter, r *http.Request) {
	network.OKMessageResponse(w, "Welcome to Revizzoro API")
}

func (c *RootController) Healthcheck(w http.ResponseWriter, r *http.Request) {
	const op = "[base Healthcheck]"

	dbIsAlive := true
	if err := c.db.Ping(); err != nil {
		log.Warnf("%s ping to database error: %s\n", op, err)
		dbIsAlive = false
	}

	network.JSONResponse(w, http.StatusOK, map[string]any{"dbIsAlive": dbIsAlive})
}
