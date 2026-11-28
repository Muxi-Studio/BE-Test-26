package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"blog-backend/internal/store"
)

// Handler 承载所有待办接口，并持有存储句柄。
type Handler struct {
	store *store.Memory
}

func NewHandler(s *store.Memory) *Handler {
	return &Handler{store: s}
}

type route struct {
	method  string
	pattern string
	handle  func(h *Handler, w http.ResponseWriter, r *http.Request)
}

var routes = []route{
	{"GET", "/todos", (*Handler).ListTodos},
	{"POST", "/todos", (*Handler).CreateTodo},
}

func (h *Handler) Register(mux *http.ServeMux) {
	for _, rt := range routes {
		mux.HandleFunc(rt.method+" "+rt.pattern, func(w http.ResponseWriter, r *http.Request) {
			rt.handle(h, w, r)
		})
	}
}

// ListTodos 返回全部待办。
func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, response{Code: 0, Msg: "ok", Data: h.store.List()})
}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	_ = err // 内部细节不返回给客户端
	writeJSON(w, http.StatusInternalServerError, response{Code: 1, Msg: "internal error"})
}

func parseJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// CreateTodo 新增一条待办。
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	if req.Title == "" {
		writeError(w, errors.New("title required"))
		return
	}
	item := h.store.Add(req.Title)
	writeJSON(w, http.StatusCreated, response{Code: 0, Msg: "ok", Data: item})
}
