package handlers

import (
	"html/template"
	"net/http"

	"github.com/blomma/badrobot/models"
)

type BadFriendsHandler struct {
	redis              *string
	templateBadFriends *template.Template
	badFriendsModel    *models.BadFriendsModel
	stopCallback       func()
}

type pageBadFriends struct {
	BadFriends template.JS
}

func NewBadFriendsHandler(redis *string) *BadFriendsHandler {
	filename := "badfriends.html"

	templateBadFriends := template.Must(template.ParseFiles(filename))
	badFriendsModel, stopCallback := models.NewBadFriendsModel(*redis)

	return &BadFriendsHandler{
		redis:              redis,
		templateBadFriends: templateBadFriends,
		badFriendsModel:    badFriendsModel,
		stopCallback:       stopCallback,
	}
}

func (b *BadFriendsHandler) Stop() {
	b.stopCallback()
}

func (b *BadFriendsHandler) Handler(w http.ResponseWriter, _ *http.Request) {
	p := pageBadFriends{BadFriends: template.JS(b.badFriendsModel.Result())}
	b.templateBadFriends.Execute(w, p)
}
