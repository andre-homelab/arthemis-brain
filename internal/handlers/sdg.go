package handlers

import (
	"arthemis-brain/internal/models"
	"arthemis-brain/internal/utils"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func SdgHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
}

// @Summary      Create Sdg
// @Description  Creates a new sdg
// @Tags         sdg
// @Accept       json
// @Produce      json
// @Param        sdg  body      models.Sdg  true   "Sdg details"
// @Success      202        {boolean} true              "Sdg created successfully"
// @Failure      400        {object}  utils.ErrorResponse "Invalid JSON or bad request"
// @Router       /sdg/create [post]
func (g *GlobalParams) CreateSdg(w http.ResponseWriter, r *http.Request) {
	var reqSdg models.Sdg
	if err := json.NewDecoder(r.Body).Decode(&reqSdg); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}
	res := g.db.Create(&reqSdg)
	if res.Error != nil {
		utils.RespondError(w, http.StatusBadRequest, "Erro ao criar o sdge", res.Error)
		return
	}
	utils.RespondJSON(w, http.StatusAccepted, reqSdg.ID)
}

// @Summary      Get Sdg
// @Description  Retrieves a sdg by ID
// @Tags         sdg
// @Produce      json
// @Param        id   path      string  true  "Sdg ID"
// @Success      200  {object}  models.Sdg
// @Failure      400  {object}  utils.ErrorResponse "SdgID not found"
// @Failure      404  {object}  utils.ErrorResponse "Sdg not found"
// @Failure      500  {object}  utils.ErrorResponse "Internal server error"
// @Router       /sdg/{id} [get]
func (g *GlobalParams) GetSdg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // or mux.Vars(r)["id"] if using gorilla/mux
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "SdgID not found", nil)
		return
	}

	var sdg models.Sdg
	res := g.db.First(&sdg, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Sdg not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error finding sdg", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, sdg)
}

// @Summary      Update Sdg
// @Description  Updates an existing sdg by ID
// @Tags         sdg
// @Accept       json
// @Produce      json
// @Param        id         path      string            true   "Sdg ID"
// @Param        sdg  body      models.Sdg  true   "Sdg details to update"
// @Success      200        {object}  models.Sdg
// @Failure      400        {object}  utils.ErrorResponse "Invalid JSON or SdgID not found"
// @Failure      404        {object}  utils.ErrorResponse "Sdg not found"
// @Failure      500        {object}  utils.ErrorResponse "Internal server error"
// @Router       /sdg/update/{id} [patch]
func (g *GlobalParams) UpdateSdg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "SdgID not found", nil)
		return
	}

	var existing models.Sdg
	res := g.db.First(&existing, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Sdg not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error finding sdg", res.Error)
		return
	}

	var reqSdg models.Sdg
	if err := json.NewDecoder(r.Body).Decode(&reqSdg); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	res = g.db.Model(&existing).Updates(&reqSdg)
	reqSdg.ID = existing.ID
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating sdg", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reqSdg)
}

// @Summary      Delete Sdg
// @Description  Deletes a sdg by ID
// @Tags         sdg
// @Produce      json
// @Param        id   path      string  true  "Sdg ID"
// @Success      200  {boolean} true    "Sdg deleted successfully"
// @Failure      400  {object}  utils.ErrorResponse "SdgID not found"
// @Failure      404  {object}  utils.ErrorResponse "Sdg not found"
// @Failure      500  {object}  utils.ErrorResponse "Internal server error"
// @Router       /sdg/delete/{id} [delete]
func (g *GlobalParams) DeleteSdg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "SdgID not found", nil)
		return
	}

	res := g.db.Delete(&models.Sdg{}, "id = ?", id)
	if res.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting sdg", res.Error)
		return
	}
	if res.RowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "Sdg not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, true)
}

// @Summary      Gets all Sdgs
// @Description  Retrieves every sdg
// @Tags         sdg
// @Produce      json
// @Success      200  {boolean} true    "Sdg deleted successfully"
// @Failure      404  {object}  utils.ErrorResponse "Sdg not found"
// @Failure      500  {object}  utils.ErrorResponse "Internal server error"
// @Router       /sdg/ [get]
func (g *GlobalParams) GetAllSdgs(w http.ResponseWriter, r *http.Request) {
	var sdgs []models.Sdg
	res := g.db.Preload("Projects").Find(&sdgs)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Sdg not found", res.Error)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Error finding sdg", res.Error)
		return
	}

	utils.RespondJSON(w, http.StatusOK, sdgs)
}
