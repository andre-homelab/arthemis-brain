package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"arthemis-brain/internal/models"
	"arthemis-brain/internal/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ProponentHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @title           Create Proponent
// @version         1.0
// @description     Creates a proponent
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /proponent/create

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateProponent(w http.ResponseWriter, r *http.Request) {
	var reqProponent models.Proponent
	if err := json.NewDecoder(r.Body).Decode(&reqProponent); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}
	res := g.db.Create(&reqProponent)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Erro ao criar o proponente", res.Error)
		return
	}
	utils.RespondJSON(w, http.StatusAccepted, true)
}

// @title           Get Proponent
// @version         1.0
// @description     Reads a proponent
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /proponent/get

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) GetProponent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID não informado", nil)
		return
	}

	var proponent models.Proponent
	res := g.db.First(&proponent, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Proponente não encontrado", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao buscar o proponente", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, proponent)
}

// @title           Update Proponent
// @version         1.0
// @description     Updates a proponent
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /proponent/update

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) UpdateProponent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID não informado", nil)
		return
	}

	var existing models.Proponent
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Proponente não encontrado", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao buscar o proponente", res.Error)
		return
	}

	var reqProponent models.Proponent
	if err := json.NewDecoder(r.Body).Decode(&reqProponent); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}

	reqProponent.ID = existing.ID
	res = g.db.Save(&reqProponent)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao atualizar o proponente", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqProponent)
}

// @title           Update Proponent
// @version         1.0
// @description     Deletes a proponent
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /proponent/delete

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) DeleteProponent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID não informado", nil)
		return
	}

	res := g.db.Delete(&models.Proponent{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao deletar o proponente", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "Proponente não encontrado", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
