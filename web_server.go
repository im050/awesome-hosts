package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-session/session"
	"html/template"
	"net/http"
	"os"
)

type Notice struct {
	Message   string
	Type      bool
	ReturnURL string
}

type RespEntity struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func IndexController(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	checkError(err)
	if !checkAccess(store) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}
	displayTemplate(w, "index.html", nil)
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	checkError(err)
	if checkAccess(store) {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		return
	}
	if r.Method == "GET" {
		displayTemplate(w, "login.html", nil)
	} else {
		checkError(r.ParseForm())
		password := r.FormValue("password")
		fmt.Println(password)
		if password == "goodboy" {
			store.Set("access", true)
			err = store.Save()
			if err != nil {
				fmt.Fprint(os.Stdout, err)
				displayTemplate(w, "notice.html", Notice{"Something wrong", false, ""})
				return
			}
			http.Redirect(w, r, "/index", http.StatusMovedPermanently)
		} else {
			displayTemplate(w, "notice.html", Notice{"密码错误", false, ""})
		}
	}
}

func LogoutController(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	checkError(err)
	checkError(store.Flush())
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func CurrentController(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	checkError(err)
	if !checkAccess(store) {
		jsonResponse(w, RespEntity{403, "无权访问", nil})
		return
	}
	jsonResponse(w, RespEntity{1, "success", hosts})
}

func jsonResponse(w http.ResponseWriter, response RespEntity) {
	w.Header().Set("Content-type", "text/json; charset=UTF-8")
	jsonText, _ := json.Marshal(response)
	str := string(jsonText)
	fmt.Fprint(w, str)
}

func displayTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl := template.Must(template.ParseFiles("template/" + templateName))
	checkError(tmpl.Execute(w, data))
}

func checkAccess(store session.Store) bool {
	access, has := store.Get("access")
	if !has || access != true {
		return false
	}
	return true
}

func ServerStart() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", IndexController)
	http.HandleFunc("/current", CurrentController)
	http.HandleFunc("/login", LoginController)
	http.HandleFunc("/logout", LogoutController)
	fmt.Println("starting...")
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	checkError(err)
}
