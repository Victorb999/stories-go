package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"stories-go/internal/models"
	"stories-go/internal/repository"

	"github.com/go-chi/chi/v5"
)

type StoryHandler struct {
	repo *repository.StoryRepository
	log  *slog.Logger
}

func NewStoryHandler(repo *repository.StoryRepository, log *slog.Logger) *StoryHandler {
	return &StoryHandler{repo: repo, log: log}
}

// Routes mounts all story endpoints onto a new chi.Router.
func (h *StoryHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Post("/seed", h.Seed)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// StoryInput is the request body for Create and Update.
type StoryInput struct {
	Title       string      `json:"title"`
	CoverImage  string      `json:"cover_image"`
	Author      string      `json:"author"`
	Content     string      `json:"content"`
	AIGenerated bool        `json:"ai_generated"`
	Size        models.Size `json:"size"`
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func normalizeSize(s models.Size) models.Size {
	if s == models.SizeLarge {
		return models.SizeLarge
	}
	return models.SizeSmall
}

// ── handlers ─────────────────────────────────────────────────────────────────

// List godoc
//
//	GET /api/v1/stories?size=small|large&ai=true|false
func (h *StoryHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.ListFilter{}
	if s := r.URL.Query().Get("size"); s != "" {
		f.Size = s
	}
	if ai := r.URL.Query().Get("ai"); ai != "" {
		b := ai == "true"
		f.AIGenerated = &b
	}

	stories, err := h.repo.List(r.Context(), f)
	if err != nil {
		h.log.Error("list stories", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if stories == nil {
		stories = []models.Story{}
	}
	writeJSON(w, http.StatusOK, stories)
}

// Create godoc
//
//	POST /api/v1/stories
func (h *StoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var inp StoryInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if inp.Title == "" || inp.Author == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and author are required")
		return
	}
	inp.Size = normalizeSize(inp.Size)

	story, err := h.repo.Create(r.Context(), &models.Story{
		Title:       inp.Title,
		CoverImage:  inp.CoverImage,
		Author:      inp.Author,
		Content:     inp.Content,
		AIGenerated: inp.AIGenerated,
		Size:        inp.Size,
	})
	if err != nil {
		h.log.Error("create story", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, story)
}

// GetByID godoc
//
//	GET /api/v1/stories/{id}   — also increments the view counter
func (h *StoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	story, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if story == nil {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	writeJSON(w, http.StatusOK, story)
}

// Update godoc
//
//	PUT /api/v1/stories/{id}
func (h *StoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var inp StoryInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if inp.Title == "" || inp.Author == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and author are required")
		return
	}
	inp.Size = normalizeSize(inp.Size)

	story, err := h.repo.Update(r.Context(), id, &models.Story{
		Title:       inp.Title,
		CoverImage:  inp.CoverImage,
		Author:      inp.Author,
		Content:     inp.Content,
		AIGenerated: inp.AIGenerated,
		Size:        inp.Size,
	})
	if err != nil {
		h.log.Error("update story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if story == nil {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	writeJSON(w, http.StatusOK, story)
}

// Delete godoc
//
//	DELETE /api/v1/stories/{id}
func (h *StoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	deleted, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		h.log.Error("delete story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Seed godoc
//
//	POST /api/v1/stories/seed
//
// Inserts sample stories when the database is empty.
func (h *StoryHandler) Seed(w http.ResponseWriter, r *http.Request) {
	existing, err := h.repo.List(r.Context(), repository.ListFilter{})
	if err != nil {
		h.log.Error("seed list", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if len(existing) > 0 {
		writeJSON(w, http.StatusOK, map[string]string{"message": "database already has stories, skipping seed"})
		return
	}

	seeds := []models.Story{
		{
			Title:       "O Coelho e a Lua",
			Author:      "Vovó Maria",
			CoverImage:  "https://images.unsplash.com/photo-1457899329096-fb8d96c33246?w=400",
			Content:     "Era uma vez um coelhinho branco que vivia numa floresta mágica. Todas as noites, ele subia até o topo da colina mais alta para conversar com a lua. A lua, sempre gentil, iluminava seu caminho de volta para casa. E assim, noite após noite, a amizade entre o coelho e a lua crescia cada vez mais, enchendo a floresta de luz e alegria.",
			AIGenerated: false,
			Size:        models.SizeSmall,
		},
		{
			Title:       "A Princesa das Estrelas",
			Author:      "Papai Noel Júnior",
			CoverImage:  "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400",
			Content:     "No reino das nuvens, vivia uma princessinha chamada Lúcia que tinha um dom especial: ela conseguia desenhar estrelas com os dedos. Cada noite, ela saía de sua casinha de nuvem e pintava o céu com pontos brilhantes. Os moradores do reino terrestre olhavam para cima e suspiravam de admiração, sem saber que era uma pequena artista celestial quem enfeitava seus sonhos.",
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:       "O Dragãozinho Tímido",
			Author:      "Tia Betinha",
			CoverImage:  "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=400",
			Content:     "Fogo era um dragãozinho diferente dos outros. Enquanto seus irmãos soltavam labaredas enormes, Fogo só conseguia soltar bolhinhas de sabão coloridas. No começo ele ficava triste, mas descobriu que as crianças da aldeia adoravam suas bolhinhas mágicas. E foi assim que o dragão mais tímido do mundo se tornou o mais querido de todos.",
			AIGenerated: false,
			Size:        models.SizeSmall,
		},
		{
			Title:       "A Sereiazinha do Rio",
			Author:      "Mamãe Contadora",
			CoverImage:  "https://images.unsplash.com/photo-1518020382113-a7e8fc38eac9?w=400",
			Content:     "Marina era uma sereinha que morava num rio cristalino no meio da floresta. Diferente das outras sereias que cantavam melodias tristes, Marina cantava músicas alegres que faziam os peixes dançarem e as flores das margens balançarem no ritmo. Certo dia, um passarinho perdido seguiu sua música e encontrou o caminho de volta para casa.",
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:       "Os Ursinhos e o Mel Mágico",
			Author:      "Vovô João",
			CoverImage:  "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400",
			Content:     "Três ursinhos — Bolinha, Fofa e Pingo — descobriram uma colmeia escondida numa árvore encantada. O mel que saía dela era diferente: cada colherada fazia a pessoa lembrar de uma memória feliz. Juntos, os ursinhos decidiram compartilhar o mel com todos os animais tristes da floresta, espalhando alegria por onde passavam.",
			AIGenerated: false,
			Size:        models.SizeSmall,
		},
	}

	var created []models.Story
	for _, s := range seeds {
		s := s
		story, err := h.repo.Create(r.Context(), &s)
		if err != nil {
			h.log.Error("seed create", "title", s.Title, "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		created = append(created, *story)
	}
	writeJSON(w, http.StatusCreated, map[string]any{"seeded": len(created), "stories": created})
}
