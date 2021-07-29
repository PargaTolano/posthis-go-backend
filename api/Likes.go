package api

import (
	"encoding/json"
	"net/http"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	"strconv"

	"github.com/gorilla/mux"
)

//Used Model: Like --main, User --required-for-like, Post --required-for-like
//User ViewModels:

//Validation

//API handlers
func GetLikes() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			post     Post
			response SuccesVM
		)

		vars := mux.Vars(r)

		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Preload("Replies").First(&post, id)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		response = SuccesVM{Data: post.Replies, Message: "Replies retrieved sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func CreateLike() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			user     User
			post     Post
			like     Like
			response SuccesVM
		)

		vars := mux.Vars(r)

		tuid, err := strconv.ParseUint(vars["userId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId := uint(tuid)

		tpid, err := strconv.ParseUint(vars["postId"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		postId := uint(tpid)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&user, userId)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&post, postId)
		if db.Error != nil {
			http.Error(w, db.Error.Error(), http.StatusInternalServerError)
			return
		}

		like = Like{UserID: user.ID, PostID: post.ID}

		//Once it works add everything to the database
		db.Create(&like)
		db.Model(&user).Association("likes").Append(&like)
		db.Model(&post).Association("likes").Append(&like)

		response = SuccesVM{Data: post, Message: "Post created successfully"}

		marshal, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}

func DeleteLike() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			like Like
		)

		vars := mux.Vars(r)
		id := vars["id"]

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db.Delete(&like, uint(tid))

		marshal, err := json.Marshal(SuccesVM{Data: like, Message: "Post deleted successfully"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(marshal)
	})
}