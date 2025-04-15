package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RamanDudoits/shortLink-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type LinkHandlerInterface interface {
	List(w http.ResponseWriter, r *http.Request)
	Store(w http.ResponseWriter, r *http.Request)
	Destroy(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
}

type LinkHandler struct {
	linkService service.LinkServiceInterface
	linkredirectService service.LinkServiceRedirectInterface
}

func NewLinkHandler(
	linkService service.LinkServiceInterface,
	linkredirectService service.LinkServiceRedirectInterface,) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
		linkredirectService: linkredirectService,
	}
}

func (h *LinkHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	links, err := h.linkService.GetUserLinks(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) Store(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	var input struct {
		URL string `json:"url"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	link, err := h.linkService.Create(input.URL, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	
	if err := h.linkService.DeleteLink(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	
	link, err := h.linkService.GetLink(id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		} else if err.Error() == "link not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	link, err := h.linkService.UpdateLink(id, userID, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortLink := chi.URLParam(r, "shortLink")
	
	link, err := h.linkredirectService.GetByShortCode(shortLink)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}
	if err := h.linkredirectService.IncrementClickCount(link.ID, link.ClickCount); err != nil {
		log.Printf("Failed to increment click count for link %d: %v", link.ID, err)
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusMovedPermanently)
}