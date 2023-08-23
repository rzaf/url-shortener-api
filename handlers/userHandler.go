package handlers

import (
	"errors"
	"fmt"
	"github.com/rzaf/url-shortener-api/helpers"
	"github.com/rzaf/url-shortener-api/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User = r.Context().Value(models.User{}).(*models.User)
	users := models.GetUsers()
	if !currentUser.IsAdmin {
		helpers.WriteJsonError(w, errors.New("user not authorized! "), 403)
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	helpers.WriteJson(w, map[string][]models.User{
		"users": users,
	}, 200)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User = r.Context().Value(models.User{}).(*models.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user, _ := models.GetUserById(id)
	if user == nil {
		helpers.WriteJsonError(w, fmt.Errorf("user with id %d not found! ", id), 404)
		return
	}
	if !currentUser.IsAdmin && int(currentUser.Id) != id {
		helpers.WriteJsonError(w, errors.New("user not authorized! "), 403)
		return
	}
	helpers.WriteJson(w, user, 200)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	helpers.ReadJson(r, &body)

	email := body["email"]
	password := body["password"]

	helpers.ValidateVar(email, "Email", "required,email")
	helpers.ValidateVar(password, "Password", "required,max=10")

	fmt.Println("request body:", body)
	createdUser, err := models.CreateUser(email, password)
	if err != nil {
		helpers.WriteJsonError(w, err, 400) // duplicate
		return
	}
	helpers.WriteJson(w, createdUser, 201)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User = r.Context().Value(models.User{}).(*models.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user, _ := models.GetUserById(id)
	if user == nil {
		helpers.WriteJsonError(w, fmt.Errorf("user with id %d not found! ", id), 404)
		return
	}
	if !currentUser.IsAdmin && int(currentUser.Id) != id {
		helpers.WriteJsonError(w, errors.New("user not authorized! "), 403)
		return
	}

	var body map[string]string
	helpers.ReadJson(r, &body)

	email, ok1 := body["email"]
	password, ok2 := body["password"]
	if !ok1 && !ok2 {
		helpers.WriteJsonError(w, errors.New("email or password is required"), 400)
		return
	}
	helpers.ValidateVar(email, "Email", "email")

	err := models.EditUserEmailOrPassword(user, email, password)
	if err != nil {
		helpers.WriteJsonError(w, err, 400) // duplicate
		return
	}

	helpers.WriteJson(w, user, 200)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User = r.Context().Value(models.User{}).(*models.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user, _ := models.GetUserById(id)
	if user == nil {
		helpers.WriteJsonError(w, fmt.Errorf("user with id %d not found! ", id), 404)
		return
	}
	if !currentUser.IsAdmin && int(currentUser.Id) != id {
		helpers.WriteJsonError(w, errors.New("user not authorized! "), 403)
		return
	}

	models.DeleteUser(user)
	helpers.WriteJson(w, map[string]string{
		"message": fmt.Sprintf("user with id %d deleted", user.Id),
	}, 200)
}

func EditUserApiKey(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User = r.Context().Value(models.User{}).(*models.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user, _ := models.GetUserById(id)
	if user == nil {
		helpers.WriteJsonError(w, fmt.Errorf("user with id %d not found! ", id), 404)
		return
	}
	if !currentUser.IsAdmin && int(currentUser.Id) != id {
		helpers.WriteJsonError(w, errors.New("user not authorized! "), 403)
		return
	}

	var body map[string]string
	helpers.ReadJson(r, &body)
	password := body["password"]
	helpers.ValidateVar(password, "Password", "required")

	err := models.EditUserApiKey(user, password)
	if err != nil {
		helpers.WriteJsonError(w, err, 400)
		return
	}
	helpers.WriteJson(w, user, 200)
}
