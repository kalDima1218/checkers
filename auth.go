package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func handleLogout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w)
	redirectToIndex(w, r)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "login.html"))
		page.Execute(w, "")
	} else {
		var login = r.URL.Query().Get("login")
		var password = r.URL.Query().Get("password")
		if getPassword(login) != password {
			fmt.Fprintf(w, "wrong")
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
		http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
		fmt.Fprint(w, "ok")
	}
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "registration.html"))
		page.Execute(w, "")
	} else {
		var username = r.URL.Query().Get("name")
		var login = r.URL.Query().Get("login")
		var password = r.URL.Query().Get("password")
		if username == "" || login == "" || password == "" || !isFreeLogin(login) || insertUser(login, password, username) == nil {
			fmt.Fprintf(w, "wrong")
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
		http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
		fmt.Fprintf(w, "ok")
	}
}
