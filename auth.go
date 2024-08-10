package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"html/template"
	"net/http"
	"path"
	"time"
)

var SECRET_KEY = []byte("secret_key")

func generateJWT(login string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(SECRET_KEY)
	return tokenString
}

func validateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return SECRET_KEY, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		login, ok := claims["login"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token")
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", fmt.Errorf("invalid token")
		}

		return login, nil
	}

	return "", fmt.Errorf("invalid token")
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w)
	redirectToIndex(w, r)
}

// TODO добавить jwt
// TODO добавить проверку длины

func checkSession(r *http.Request) bool {
	login, errLogin := r.Cookie("login")
	password, errPassword := r.Cookie("password")
	if errLogin != nil || errPassword != nil {
		return false
	}
	if getPassword(login.Value) != password.Value {
		return false
	} else {
		return true
	}
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
