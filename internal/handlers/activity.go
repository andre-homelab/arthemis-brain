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

func ActivityHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @title           Create activity
// @version         1.0
// @description     Creates an activity
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /activity/create

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var reqActivity models.Activity
	if err := json.NewDecoder(r.Body).Decode(&reqActivity); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	res := g.db.Create(&reqActivity)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error creating the activity", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, reqActivity)
}

// @title           Get activity
// @version         1.0
// @description     Reads an activity
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /activity/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) GetActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ActivityID not recieved", nil)
		return
	}

	var activity models.Activity
	res := g.db.First(&activity, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Activity not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error retrieving the Activity", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, activity)
}

// @title           Update activity
// @version         1.0
// @description     Updates an activity
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /activity/update/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ActivityID not recieved", nil)
		return
	}

	var existing models.Activity
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Activity not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error searching for activity", res.Error)
		return
	}

	var reqActivity models.Activity
	if err := json.NewDecoder(r.Body).Decode(&reqActivity); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	reqActivity.ID = existing.ID
	res = g.db.Model(&existing).Updates(&reqActivity)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating activity", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqActivity)
}

// @title           Delete activity
// @version         1.0
// @description     Deletes an activity
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /activity/delete/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID not informed", nil)
		return
	}

	res := g.db.Delete(&models.Activity{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting activity", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "activity not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
