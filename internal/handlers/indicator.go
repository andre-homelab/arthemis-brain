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

func IndicatorHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @title           Create indicator
// @version         1.0
// @description     Creates an indicator
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /indicator/create

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateIndicator(w http.ResponseWriter, r *http.Request) {
	var reqindicator models.Proponent
	if err := json.NewDecoder(r.Body).Decode(&reqindicator); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}
	res := g.db.Create(&reqindicator)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Erro ao criar o indicatore", res.Error)
		return
	}
	utils.RespondJSON(w, http.StatusAccepted, reqindicator.ID)
}

// @title           Get indicator
// @version         1.0
// @description     Reads an indicator
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /indicator/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) GetIndicator(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "indicatorID not found", nil)
		return
	}

	var indicator models.Proponent
	res := g.db.Preload("IndicatorObservations").First(&indicator, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "indicator not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error finding indicator", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, indicator)
}

// @title           Update indicator
// @version         1.0
// @description     Updates an indicator
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /indicator/update/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) UpdateIndicator(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "indicatorID not found", nil)
		return
	}

	var existing models.Indicator
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "indicator not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error finding indicator", res.Error)
		return
	}

	var reqindicator models.Proponent
	if err := json.NewDecoder(r.Body).Decode(&reqindicator); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	res = g.db.Model(&existing).Updates(&reqindicator)
	reqindicator.ID = existing.ID
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating indicator", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqindicator)
}

// @title           Delete indicator
// @version         1.0
// @description     Deletes a indicator
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /indicator/delete/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) DeleteIndicator(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "indicatorID not found", nil)
		return
	}

	res := g.db.Delete(&models.Indicator{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting indicator", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "indicator not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
