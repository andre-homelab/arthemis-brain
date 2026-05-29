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
// @BasePath  /project/create

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) CreateProject(w http.ResponseWriter, r *http.Request) {
	var reqProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&reqProject); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	res := g.db.Create(&reqProject)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error creating the project", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, reqProject.ID)
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
	res := g.db.Preload("Locations").Preload("Activities").Preload("ProjectProponents").First(&project, "id = ?", id)
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
// @BasePath  /project/update/{id}

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

// @title           Delete project
// @version         1.0
// @description     Deletes a project
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/delete/{id}

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
		utils.RespondError(w, http.StatusNotFound, "projecto não encontrado", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}

// @title           Add Proponent To Project
// @version         1.0
// @description     Creates association between proponent and project on project_proponent table
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{id}/addproponent

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) AddProponent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID não informado", nil)
		return
	}

	var project models.Project
	g.db.Preload("Locations").Preload("Activities").First(&project, "id = ?", id)

	var body struct {
		ProponentID uint   `json:"proponentId"`
		Role        string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	projectProponent := models.ProjectProponent{
		ProponentID: body.ProponentID,
		ProjectID:   project.ID,
		Role:        body.Role,
	}

	if err := g.db.Create(&projectProponent).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to add proponent", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, projectProponent.ID)
}

// @title           Remove Proponent from Project
// @version         1.0
// @description     Remoces association between proponent and project on project_proponent table
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /project/{projectId}/remove_proponent/{proponentId}

// @securityDefinitions.basic  BasicAuth

func (g *GlobalParams) RemoveProponent(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	if projectID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID do projeto não informado", nil)
		return
	}

	var project models.Project
	g.db.First(&project, "id = ?", projectID)

	proponentID := chi.URLParam(r, "proponentId")
	if proponentID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID do projeto não informado", nil)
		return
	}

	if err := g.db.Where("project_id = ? AND proponent_id = ?", projectID, proponentID).
		Delete(&models.ProjectProponent{}).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to remove proponent", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}
