package main

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var secretKey = loadSecretKey()

func loadSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		return []byte(secret)
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}
	return key
}

func generateJWT(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func setTokenCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func validateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return secretKey, nil
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

func getLogin(r *http.Request) (string, error) {
	token, errToken := r.Cookie("token")
	if errToken != nil {
		return "", fmt.Errorf("no token")
	}
	return validateJWT(token.Value)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w, r)
	redirectToIndex(w, r)
}

func checkSession(r *http.Request) bool {
	_, err := getLogin(r)
	return err == nil
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "registration.html"))
		page.Execute(w, "")
	} else {
		username := r.URL.Query().Get("username")
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")
		if !isFreeLogin(login) {
			fmt.Fprintf(w, "not free login")
			return
		}
		if len(login) > 42 || len(username) > 42 {
			fmt.Fprintf(w, "too long")
			return
		}
		if username == "" || login == "" || password == "" || insertUser(login, password, username) != nil {
			fmt.Fprintf(w, "wrong")
			return
		}
		token, err := generateJWT(login)
		if err != nil {
			fmt.Fprintf(w, "wrong")
			return
		}
		setTokenCookie(w, r, token)
		fmt.Fprintf(w, "ok")
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page, _ := template.ParseFiles(path.Join("html", "login.html"))
		page.Execute(w, "")
	} else {
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")
		if getPassword(login) != password {
			fmt.Fprintf(w, "wrong")
			return
		}
		token, err := generateJWT(login)
		if err != nil {
			fmt.Fprintf(w, "wrong")
			return
		}
		setTokenCookie(w, r, token)
		fmt.Fprint(w, "ok")
	}
}
