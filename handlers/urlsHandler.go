package handlers

import (
	"errors"
	"github.com/rzaf/url-shortener-api/helpers"
	"github.com/rzaf/url-shortener-api/models"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func GetUrls(w http.ResponseWriter, r *http.Request) {
	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	urls := models.GetUserUrls(int(user.Id))
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	helpers.WriteJson(w, urls, 200)
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	short := chi.URLParam(r, "short")
	url := models.GetUrl(short)
	if !user.IsAdmin && url.User_id != user.Id {
		helpers.WriteJsonError(w, errors.New("user not authorized"), 403)
		return
	}
	models.DeleteUrl(short)
	helpers.WriteJson(w, map[string]string{
		"message": "url `" + short + "` deleted",
	}, 200)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	url := models.GetUrl(short)
	helpers.WriteJson(w, url, 200)
}

func CreateUrl(w http.ResponseWriter, r *http.Request) {
	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	var body map[string]string
	helpers.ReadJson(r, &body)
	urlStr := body["url"]
	helpers.ValidateVar(urlStr, "Url", "required,url")

	urlId := models.CreateUrl(urlStr, user.Id)
	helpers.WriteJson(w, map[string]any{
		"message":       "url created",
		"shortened":     models.IdToShort(int(urlId)),
		"shortened_url": os.Getenv("URL") + "/" + models.IdToShort(int(urlId)),
	}, 201)
}

func EditUrl(w http.ResponseWriter, r *http.Request) {
	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	short := chi.URLParam(r, "short")
	var body map[string]string
	helpers.ReadJson(r, &body)
	newUrlStr := body["url"]
	helpers.ValidateVar(newUrlStr, "Url", "required,url")

	url := models.GetUrl(short)
	if !user.IsAdmin && url.User_id != user.Id {
		helpers.WriteJsonError(w, errors.New("user not authorized"), 403)
		return
	}
	models.EditUrl(short, newUrlStr)
	helpers.WriteJson(w, map[string]any{
		"message": "url edited",
		"url":     newUrlStr,
	}, 200)
}
