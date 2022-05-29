package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

const (
	sessionName        = "DNM"
	ctxKeyUser  ctxKey = iota
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func NewServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// /private/*
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/", s.handleWhoAmI()).Methods("GET")
	private.HandleFunc("/all", s.findAll()).Methods("GET")
	private.HandleFunc("/delete", s.handleDelete()).Methods("DELETE")
	private.HandleFunc("/search_username", s.handleFindUserByUsernameLike()).Methods("GET")
	posts := s.router.PathPrefix("/private/post").Subrouter()
	posts.Use(s.authenticateUser)
	posts.HandleFunc("/", s.handlePostsCreate()).Methods("POST")
	posts.HandleFunc("/", s.handleGetPostsForUser()).Methods("GET")
	posts.HandleFunc("/", s.handleDeletePost()).Methods("DELETE")
	posts.HandleFunc("/all", s.handleGetAllPosts()).Methods("GET")
	posts.HandleFunc("/{id}", s.handleGetPostById()).Methods("GET")
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	c := &http.Cookie{
		Name:    sessionName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: false,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			http.SetCookie(w, c)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			http.SetCookie(w, c)
			return
		}

		timeEnd := session.Values["end_time"].(int64)

		if time.Now().UnixMilli()-timeEnd > 0 {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			http.SetCookie(w, c)
			return
		}

		u, err := s.store.User().FindById(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			http.SetCookie(w, c)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*models.User))
	}
}

func (s *server) findAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.User().FindAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) handleFindUserByUsernameLike() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		users, err := s.store.User().FindByUsernameLike(req.Username)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {

	type request struct {
		Email             string `json:"email"`
		Username          string `json:"username"`
		DecryptedPassword string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := models.User{
			Email:             req.Email,
			Username:          req.Username,
			DecryptedPassword: req.DecryptedPassword,
		}
		if err := s.store.User().Create(&u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email             string `json:"email"`
		DecryptedPassword string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.DecryptedPassword) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		startTime := time.Now().UnixMilli()
		session.Values["user_id"] = u.Id
		session.Values["end_time"] = startTime + 1000*60*60*24 //1 day
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		s.respond(w, r, http.StatusOK, nil)
	}

}

func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxKeyUser).(*models.User).Id
		err := s.store.User().DeleteById(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handlePostsCreate() http.HandlerFunc {

	type request struct {
		Caption string `json:"caption"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		p := models.Post{
			UserId:   r.Context().Value(ctxKeyUser).(*models.User).Id,
			Username: r.Context().Value(ctxKeyUser).(*models.User).Username,
			Caption:  req.Caption,
		}
		if err := s.store.Post().Create(&p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, p)
	}
}

func (s *server) handleDeletePost() http.HandlerFunc {
	type request struct {
		Id string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		postId, err := strconv.Atoi(req.Id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		err = s.store.Post().DeleteById(postId)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetPostsForUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ctxKeyUser).(*models.User).Id
		fmt.Println(userId)
		posts, err := s.store.Post().FindAllByUserId(userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, posts)
	}

}

func (s *server) handleGetAllPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := s.store.Post().FindAll()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, post)
	}
}

func (s *server) handleGetPostById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId, err := strconv.Atoi(params["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		post, err := s.store.Post().FindById(postId)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, post)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
