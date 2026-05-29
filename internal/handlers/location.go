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

func LocationHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @title           Create location
// @version         1.0
// @description     Creates a location
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /location/create

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var reqLocation models.Location
	if err := json.NewDecoder(r.Body).Decode(&reqLocation); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	res := g.db.Create(&reqLocation)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error creating the location", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, reqLocation)
}

// @title           Get location
// @version         1.0
// @description     Reads a location
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /location/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) GetLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "LocationID not recieved", nil)
		return
	}

	var location models.Location
	res := g.db.First(&location, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Location not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error retrieving the Location", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, location)
}

// @title           Update location
// @version         1.0
// @description     Updates a location
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /location/update/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "LocationID not recieved", nil)
		return
	}

	var existing models.Location
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Location not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error searching for location", res.Error)
		return
	}

	var reqLocation models.Location
	if err := json.NewDecoder(r.Body).Decode(&reqLocation); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	reqLocation.ID = existing.ID
	res = g.db.Model(&existing).Updates(&reqLocation)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating location", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqLocation)
}

// @title           Delete location
// @version         1.0
// @description     Deletes a location
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /location/delete/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID not informed", nil)
		return
	}

	res := g.db.Delete(&models.Location{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting location", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "location not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
