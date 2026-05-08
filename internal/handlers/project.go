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

func ProjectHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @title           Create project
// @version         1.0
// @description     Creates a project
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateProject(w http.ResponseWriter, r *http.Request) {
	proponentID, ok := r.Context().Value("proponent_id").(uint)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var reqProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&reqProject); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	var count int64
	g.db.Model(&models.Project{}).Where("proponent_id = ?", proponentID).Count(&count)
	if count > 0 {
		utils.RespondError(w, http.StatusConflict, "Proponent already has a project", nil)
		return
	}

	reqProject.ProponentID = proponentID

	res := g.db.Create(reqProject)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error creating the project", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, true)
}

// @title           Get project
// @version         1.0
// @description     Reads a project
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ProjectID not recieved", nil)
		return
	}

	var project models.Project
	res := g.db.First(&project, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Project not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error retrieving the Project", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, project)
}

// @title           Update project
// @version         1.0
// @description     Updates a project
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ProjectID not recieved", nil)
		return
	}

	var existing models.Project
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "projecte não encontrado", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao buscar o projecte", res.Error)
		return
	}

	var reqProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&reqProject); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}

	reqProject.ID = existing.ID
	res = g.db.Model(&existing).Updates(&reqProject)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao atualizar o projecte", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqProject)
}

// @title           Update project
// @version         1.0
// @description     Deletes a project
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{id}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID não informado", nil)
		return
	}

	res := g.db.Delete(&models.Project{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Erro ao deletar o projecte", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "projecte não encontrado", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
