package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/api/dto"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/pkg/response"
	"github.com/owndangan/backend/internal/pkg/validator"
	"github.com/owndangan/backend/internal/service"
)

type EventHandler struct {
	eventSvc *service.EventService
}

func NewEventHandler(eventSvc *service.EventService) *EventHandler {
	return &EventHandler{eventSvc: eventSvc}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req dto.CreateEventRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.eventSvc.Create(r.Context(), userID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, resp, r)
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	params, _, _ := parsePagination(r)
	status := r.URL.Query().Get("status")
	resp, _, err := h.eventSvc.List(r.Context(), userID, params, status)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	resp, err := h.eventSvc.GetByID(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req dto.UpdateEventRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.eventSvc.Update(r.Context(), userID, id, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	err = h.eventSvc.Delete(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Event deleted"}, r)
}

func (h *EventHandler) Publish(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	resp, err := h.eventSvc.Publish(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	err = h.eventSvc.Unpublish(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Event unpublished"}, r)
}

func (h *EventHandler) PublicView(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_SLUG", "Slug is required", r)
		return
	}

	resp, err := h.eventSvc.GetPublicBySlug(r.Context(), slug)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) RegisterRoutes(r chi.Router, authRequired func(http.Handler) http.Handler, publicHandler http.HandlerFunc) {
	r.Group(func(r chi.Router) {
		r.Use(authRequired)
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Post("/{id}/publish", h.Publish)
		r.Post("/{id}/unpublish", h.Unpublish)

		r.Put("/{id}/template", h.AssignTemplate)
		r.Get("/{id}/sections", h.GetSections)
		r.Put("/{id}/sections", h.UpdateSections)
		r.Get("/{id}/digital-gifts", h.GetDigitalGift)
		r.Put("/{id}/digital-gifts", h.UpdateDigitalGift)
		r.Route("/{id}/gallery", func(r chi.Router) {
			r.Get("/", h.ListGallery)
			r.Post("/upload", h.UploadGallery)
			r.Delete("/{photoID}", h.DeleteGallery)
			r.Put("/reorder", h.ReorderGallery)
		})
		r.Route("/{id}/music", func(r chi.Router) {
			r.Get("/", h.GetMusic)
			r.Post("/upload", h.UploadMusic)
			r.Get("/presets", h.ListMusicPresets)
			r.Post("/presets", h.AssignPresetMusic)
		})
	})
}

func (h *EventHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	templates, err := h.eventSvc.ListTemplates(r.Context(), userID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, templates, r)
}

func (h *EventHandler) AssignTemplate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req struct {
		TemplateID uuid.UUID `json:"template_id"`
	}
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.eventSvc.AssignTemplate(r.Context(), userID, id, req.TemplateID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) GetSections(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	sections, err := h.eventSvc.GetSections(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, sections, r)
}

func (h *EventHandler) UpdateSections(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req dto.UpdateSectionsRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.eventSvc.UpdateSections(r.Context(), userID, id, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) GetDigitalGift(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	resp, err := h.eventSvc.GetDigitalGift(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) UpdateDigitalGift(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req dto.UpdateDigitalGiftRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	resp, err := h.eventSvc.UpdateDigitalGift(r.Context(), userID, id, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, resp, r)
}

func (h *EventHandler) ListGallery(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	photos, err := h.eventSvc.ListGallery(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, photos, r)
}

func (h *EventHandler) UploadGallery(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FORM", "Invalid multipart form", r)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FILE", "file field is required", r)
		return
	}
	defer file.Close()
	caption := r.FormValue("caption")
	var filename string
	if header != nil {
		filename = header.Filename
	}
	photo, err := h.eventSvc.UploadGallery(r.Context(), userID, id, file, filename, caption)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, photo, r)
}

func (h *EventHandler) DeleteGallery(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	photoID, err := parseUUID(r, "photoID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_PHOTO_ID", "Invalid photo ID", r)
		return
	}
	if err := h.eventSvc.DeleteGallery(r.Context(), userID, id, photoID); err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Photo deleted"}, r)
}

func (h *EventHandler) ReorderGallery(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req dto.ReorderGalleryRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	if err := h.eventSvc.ReorderGallery(r.Context(), userID, id, req); err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Gallery reordered"}, r)
}

func (h *EventHandler) GetMusic(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	music, err := h.eventSvc.GetMusic(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, music, r)
}

func (h *EventHandler) UploadMusic(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FORM", "Invalid multipart form", r)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FILE", "file field is required", r)
		return
	}
	defer file.Close()
	title := r.FormValue("title")
	var filename string
	if header != nil {
		filename = header.Filename
	}
	music, err := h.eventSvc.UploadMusic(r.Context(), userID, id, file, filename, title)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, music, r)
}

func (h *EventHandler) ListMusicPresets(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	presets, err := h.eventSvc.ListMusicPresets(r.Context(), userID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, presets, r)
}

func (h *EventHandler) RemoveMusic(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	if err := h.eventSvc.RemoveMusic(r.Context(), userID, id); err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Music removed"}, r)
}

func (h *EventHandler) AssignPresetMusic(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req struct {
		PresetID uuid.UUID `json:"preset_id"`
	}
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	music, err := h.eventSvc.AssignPresetMusic(r.Context(), userID, id, req.PresetID)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, music, r)
}

func (h *EventHandler) ListLoveStories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	stories, err := h.eventSvc.ListLoveStories(r.Context(), userID, id)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, stories, r)
}

func (h *EventHandler) CreateLoveStory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	var req dto.CreateLoveStoryRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	story, err := h.eventSvc.CreateLoveStory(r.Context(), userID, id, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusCreated, story, r)
}

func (h *EventHandler) UpdateLoveStory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	storyID, err := parseUUID(r, "storyID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_STORY_ID", "Invalid story ID", r)
		return
	}
	var req dto.UpdateLoveStoryRequest
	if err := validator.ParseAndValidate(r, &req); err != nil {
		validator.WriteError(w, r, err)
		return
	}
	story, err := h.eventSvc.UpdateLoveStory(r.Context(), userID, id, storyID, req)
	if err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, story, r)
}

func (h *EventHandler) DeleteLoveStory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id, err := parseUUID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid event ID", r)
		return
	}
	storyID, err := parseUUID(r, "storyID")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_STORY_ID", "Invalid story ID", r)
		return
	}
	if err := h.eventSvc.DeleteLoveStory(r.Context(), userID, id, storyID); err != nil {
		response.FromError(w, err, r)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Love story deleted"}, r)
}
