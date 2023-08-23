package handlers

import (
	"errors"
	"github.com/rzaf/url-shortener-api/helpers"
	"log"
	"net/http"
)

func RecoverServerPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("error type:%T,error:%v \n", err, err)
				switch err2 := err.(type) {
				case *helpers.ServerError:
					helpers.WriteJson(w, err2.ErrorMessage(), err2.Status)
				case helpers.ServerError:
					helpers.WriteJson(w, err2.ErrorMessage(), err2.Status)
				case helpers.ValidationFieldError:
					helpers.WriteJson(w, err2.ErrorMessage(), 400)
				case helpers.ValidationFieldErrors:
					helpers.WriteJson(w, err2.ErrorMessage(), 400)
				case string:
					log.Fatal("panic only with error")
				default:
					helpers.WriteJsonError(w, errors.New("there was an internal server error"), 500)
				}
				return
			}
		}()

		next.ServeHTTP(w, r)

	})
}
