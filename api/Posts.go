package api

import (
	"encoding/json"
	"net/http"
	. "posthis/db"
	. "posthis/model"
	. "posthis/model/viewmodel"
	. "posthis/utils"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Used Model: Post --main, User --owner-to-post, Media --belonging-to-post
//User ViewModels:

//Validation

//API handlers
func GetPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			posts    []Post
			response SuccesVM
		)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.AutoMigrate(&Post{})

		db.Find(&posts)

		response = SuccesVM{Data: posts, Message: "Retrieved Posts Sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func GetPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			post         Post
			postDetailVm PostDetailVM
			response     SuccesVM
		)
		vars := mux.Vars(r)
		tid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := uint(tid)

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.Preload("Media").First(&post, id)
		postDetailVm = PostDetailVM{ID: post.ID, Content: post.Content, Media: []string{}}

		for _, v := range post.Media {
			postDetailVm.Media = append(postDetailVm.Media, v.GetPath(r.URL.Scheme, r.Header.Get("Host")))
		}

		response = SuccesVM{Data: postDetailVm, Message: "Retrieved Posts Sucessfully"}

		marshal, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(marshal))
	})
}

func CreatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			media    []*Media
			user     User
			post     Post
			response SuccesVM
		)

		content := r.FormValue("content")

		db, err := ConnectToDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.First(&user, context.Get(r, "userId"))

		//10mb total
		r.ParseMultipartForm(10 << 20)

		formdata := r.MultipartForm

		files := formdata.File["files"]

		err = UploadMultipleFiles(files, &media)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.CreateInBatches(&media, len(media))

		post = Post{Content: content, Media: media}

		db.Create(&post)
		db.Model(&user).Association("Posts").Append(&post)

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

func UpdatePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func DeletePost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func GetFeed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
