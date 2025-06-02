package utils

import (
	"net/http"
)

// setup cookie
func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   false, //false only for localhost
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// get cookie
func GetCookie(r *http.Request, name string) string {
	cookies := r.Cookies()
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

// delete cookie
func DeleteCookie(w http.ResponseWriter, name string) {
	SetCookie(w, name, "", -1)
}
