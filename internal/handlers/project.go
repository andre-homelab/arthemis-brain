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

// @Summary      Create project
// @Description  Creates a new project associated with a proponent ID
// @Tags         project
// @Accept       json
// @Produce      json
// @Param        id       path      string          true   "Proponent ID"
// @Param        project  body      models.Project  true   "Project details"
// @Success      202      {boolean} true            "Project created successfully"
// @Failure      400      {object}  utils.ErrorResponse "Invalid JSON or bad request"
// @Failure      401      {object}  utils.ErrorResponse "Unauthorized: proponent_id context missing"
// @Failure      409      {object}  utils.ErrorResponse "Proponent already has a project"
// @Router       /project/{id} [post]
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

// @Summary      Get project
// @Description  Retrieves a project by ID
// @Tags         project
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  utils.ErrorResponse "ProjectID not received"
// @Failure      404  {object}  utils.ErrorResponse "Project not found"
// @Failure      500  {object}  utils.ErrorResponse "Internal server error"
// @Router       /project/{id} [get]
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

// @Summary      Update project
// @Description  Updates an existing project by ID
// @Tags         project
// @Accept       json
// @Produce      json
// @Param        id       path      string          true   "Project ID"
// @Param        project  body      models.Project  true   "Project details to update"
// @Success      200      {object}  models.Project
// @Failure      400      {object}  utils.ErrorResponse "Invalid JSON or ProjectID not received"
// @Failure      404      {object}  utils.ErrorResponse "Project not found"
// @Failure      500      {object}  utils.ErrorResponse "Internal server error"
// @Router       /project/{id} [put]
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

// @Summary      Delete project
// @Description  Deletes a project by ID
// @Tags         project
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {boolean} true    "Project deleted successfully"
// @Failure      400  {object}  utils.ErrorResponse "ID not informed"
// @Failure      404  {object}  utils.ErrorResponse "Project not found"
// @Failure      500  {object}  utils.ErrorResponse "Internal server error"
// @Router       /project/{id} [delete]
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
