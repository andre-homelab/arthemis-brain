package handlers

import (
	"encoding/json"
	"net/http"

	"arthemis-brain/internal/models"
	"arthemis-brain/internal/utils"

	"gorm.io/gorm"
)

func ProponentHandler(db *gorm.DB) *DBStore {
	return &DBStore{db}
}

// @title           Health Cheack
// @version         1.0
// @description     Evaluates the status of the server and database
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /health

// @securityDefinitions.basic  BasicAuth

func (db *DBStore) CreateProponent(w http.ResponseWriter, r *http.Request) {
	var reqProponent models.Proponent
	if err := json.NewDecoder(r.Body).Decode(&reqProponent); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}
	res := db.db.Create(reqProponent)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Erro ao criar o proponente", res.Error)
		return
	}
	utils.RespondJSON(w, http.StatusAccepted, true)
}

func (db *DBStore) GetProponent(w http.ResponseWriter, r *http.Request) {
	return
}

func (db *DBStore) UpdateProponent(w http.ResponseWriter, r *http.Request) {
	return
}

func (db *DBStore) DeleteProponent(w http.ResponseWriter, r *http.Request) {
	return
}
