package base

import (
	"database/sql"
	"github.com/quantum73/revizzoro-api/internal/network"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RootController struct {
	database *sql.DB
}

func NewRootController(database *sql.DB) *RootController {
	return &RootController{database: database}
}

func (c *RootController) Home(w http.ResponseWriter, r *http.Request) {
	const op = "base HomeHandler"
	if err := c.database.Ping(); err != nil {
		log.Warnf("[%s] ping to database error: %s", op, err)
	} else {
		log.Infof("[%s] ping database success", op)
	}

	network.MessageJSONResponse(w, http.StatusOK, "OK")
}
