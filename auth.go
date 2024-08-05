package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func logout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w)
	redirectToIndex(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "login.html"))
		page.Execute(w, "")
	} else {
		var login = r.URL.Query().Get("login")
		var password = r.URL.Query().Get("password")
		_, ok := logins[login]
		if !ok || players[login].Password != password {
			fmt.Fprintf(w, "wrong")
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
		http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
		fmt.Fprint(w, "ok")
	}
}

func reg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "reg.html"))
		page.Execute(w, "")
	} else {
		var name = r.URL.Query().Get("name")
		var login = r.URL.Query().Get("login")
		var password = r.URL.Query().Get("password")
		_, okName := names[name]
		_, okLogin := logins[login]
		if name == "" || login == "" || password == "" || okName || okLogin {
			fmt.Fprintf(w, "wrong")
			return
		}
		players[login] = newPlayer(name, login, password)
		names[name] = true
		logins[login] = true
		http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
		http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
		fmt.Fprintf(w, "ok")
	}
}
