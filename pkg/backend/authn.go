package backend

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

type pageData struct {
	authorizationData

	Page string
}

type authorizationData struct {
	Email     string
	LogoutURL string
}

func (b *Backend) authorize(w http.ResponseWriter, r *http.Request) (bool, authorizationData, error) {
	rt, err := r.Cookie(refreshTokenKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

			return true, authorizationData{}, nil
		}

		return false, authorizationData{}, err
	}
	refreshToken := rt.Value

	it, err := r.Cookie(idTokenKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

			return true, authorizationData{}, nil
		}

		return false, authorizationData{}, err
	}
	idToken := it.Value

	id, err := b.verifier.Verify(r.Context(), idToken)
	if err != nil {
		oauth2Token, err := b.config.TokenSource(r.Context(), &oauth2.Token{
			RefreshToken: refreshToken,
		}).Token()
		if err != nil {
			http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

			return true, authorizationData{}, nil
		}

		if refreshToken = oauth2Token.RefreshToken; refreshToken != "" {
			http.SetCookie(w, &http.Cookie{
				Name:     refreshTokenKey,
				Value:    refreshToken,
				Expires:  time.Now().Add(time.Hour * 24 * 365),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})
		}

		var ok bool
		idToken, ok = oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Redirect(w, r, b.config.AuthCodeURL(b.oidcRedirectURL), http.StatusFound)

			return true, authorizationData{}, nil
		}

		http.SetCookie(w, &http.Cookie{
			Name:     idTokenKey,
			Value:    idToken,
			Expires:  oauth2Token.Expiry,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := id.Claims(&claims); err != nil {
		return false, authorizationData{}, err
	}

	if !claims.EmailVerified {
		return false, authorizationData{}, errEmailNotVerified
	}

	logoutURL, err := url.Parse(b.oidcIssuer)
	if err != nil {
		return false, authorizationData{}, err
	}

	q := logoutURL.Query()
	q.Set("id_token_hint", idToken)
	q.Set("post_logout_redirect_uri", b.oidcRedirectURL)
	logoutURL.RawQuery = q.Encode()

	logoutURL = logoutURL.JoinPath("oidc", "logout")

	return false, authorizationData{
		Email:     claims.Email,
		LogoutURL: logoutURL.String(),
	}, nil
}

type redirectData struct {
	pageData
	Href string
}

func (b *Backend) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")

	// Sign out
	if authCode == "" {
		http.SetCookie(w, &http.Cookie{
			Name:   refreshTokenKey,
			Value:  "",
			MaxAge: -1,
		})

		http.SetCookie(w, &http.Cookie{
			Name:   idTokenKey,
			Value:  "",
			MaxAge: -1,
		})

		if err := b.tpl.ExecuteTemplate(w, "redirect.html", redirectData{
			pageData: pageData{
				Page: "🔒 Signing You Out ...",
			},
			Href: "/",
		}); err != nil {
			log.Println(errCouldNotRenderTemplate, err)

			http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	// Sign in
	oauth2Token, err := b.config.Exchange(r.Context(), authCode)
	if err != nil {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenKey,
		Value:    oauth2Token.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	idToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println(errCouldNotLogin, err)

		http.Error(w, errCouldNotLogin.Error(), http.StatusUnauthorized)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     idTokenKey,
		Value:    idToken,
		Expires:  oauth2Token.Expiry,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	if err := b.tpl.ExecuteTemplate(w, "redirect.html", redirectData{
		pageData: pageData{
			Page: "🔒 Signing You In ...",
		},
		Href: "/",
	}); err != nil {
		log.Println(errCouldNotRenderTemplate, err)

		http.Error(w, errCouldNotRenderTemplate.Error(), http.StatusInternalServerError)

		return
	}
}
