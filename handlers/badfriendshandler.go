package handlers

import (
	"html/template"
	"net/http"

	"github.com/blomma/badrobot/models"
)

type badFriendsHandler struct {
	redis              *string
	templateBadFriends *template.Template
	badFriendsModel    *models.BadFriendsModel
	stopCallback       func()
}

type pageBadFriends struct {
	BadFriends template.JS
}

func NewBadFriendsHandler(redis *string) *badFriendsHandler {
	filename := "badfriends.html"

	templateBadFriends := template.Must(template.ParseFiles(filename))
	badFriendsModel, stopCallback := models.NewBadFriendsModel(*redis)

	return &badFriendsHandler{
		redis:              redis,
		templateBadFriends: templateBadFriends,
		badFriendsModel:    badFriendsModel,
		stopCallback:       stopCallback,
	}
}

func (b *badFriendsHandler) Stop() {
	b.stopCallback()
}

func (b *badFriendsHandler) Handler(w http.ResponseWriter, r *http.Request) {
	p := pageBadFriends{BadFriends: template.JS(b.badFriendsModel.Result())}
	b.templateBadFriends.Execute(w, p)
}
