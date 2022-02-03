package client

import (
	"commune/gomatrix"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/unrolled/secure"
)

type WellKnownServer struct {
	ServerName string `json:"m.server"`
}

func (c *Client) ValidateLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Valid bool `json:"valid"`
		}

		ff := Response{}

		username := pay.Username
		password := pay.Password

		fu, us := c.IsFederated(username)
		//port is only for my dev environment, this needs to go, or i'm just
		//confused
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			log.Println(wk)
			serverName = c.URLScheme(wk.ServerName)
			username = fmt.Sprintf(`%s:%s`, us.LocalPart, us.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		rl := &gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: password,
		}

		resp, err := matrix.Login(rl)
		if err == nil && resp != nil {
			log.Println(err)
			ff.Valid = true
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//Log user in
func (c *Client) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")
		federated := r.FormValue("federated") == "on"

		if username == "" || password == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}

		if s == nil {
			sess := NewSession(c.Config.Client.SecureCookie)
			sess.Options.Domain = fmt.Sprintf(`.%s`, c.Config.Client.Domain)
			s, _ = sess.Get(r, c.Config.Client.CookieName)
		}

		fu, us := c.IsFederated(username)
		//port is only for my dev environment, this needs to go, or i'm just
		//confused
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		//if federation user, we query homeserver at the /well-known endpoint
		//for full server path
		if fu {
			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			log.Println(wk)
			serverName = c.URLScheme(wk.ServerName)
			username = fmt.Sprintf(`%s:%s`, us.LocalPart, us.ServerName)
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		rl := &gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: password,
		}

		resp, err := matrix.Login(rl)
		if err != nil {
			log.Println(err)

			s.AddFlash("Username or Password Wrong", "login-error")
			s.AddFlash(username, "login-username")
			s.AddFlash(federated, "login-federated")
			s.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		matrix.SetCredentials(resp.UserID, resp.AccessToken)

		prefs, err := matrix.GetAccountPreferences(resp.UserID)
		if err != nil {
			log.Println(err)
		}

		// store user session
		token := RandomString(64)

		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			WellKnown:         serverName,
			Federated:         fu,
		}

		if prefs != nil {
			u.Preferences = *prefs
		}

		/*
			profile, err := matrix.GetProfile(resp.UserID)
			if err != nil {
				log.Println(err)
			}

			if profile != nil {
				if profile.Displayname != nil && len(*profile.Displayname) > 0 {
					u.DisplayName = *profile.Displayname
				}
				if profile.AvatarURL != nil && len(*profile.AvatarURL) > 0 {
					u.AvatarURL = StripMXCPrefix(*profile.AvatarURL)
				}
			}
		*/
		/*

			if res != nil && res.RoomID != "" {
				u.RoomID = string(res.RoomID)
			}
		*/

		serialized, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(token, resp.UserID, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(resp.UserID, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		s.Values["access_token"] = token
		s.Values["device_id"] = resp.DeviceID

		s.AddFlash("User logged in", "login-success")
		/*
			if newUser {
				s.AddFlash("Signed Up", "signed-up")
			}
		*/
		s.Save(r, w)

		//redirect to index
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//signup page
func (c *Client) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//us := LoggedInUser(r)

		type page struct {
			BasePage
			UserExists           bool
			ServerDown           bool
			SignupError          bool
			Interactive          bool
			HomeServer           string
			RegistrationDisabled bool
			Email                string
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("user-exists")
			if len(x) > 0 {
				t.UserExists = true
				s.Save(r, w)
			}
			y := s.Flashes("server-down")
			if len(y) > 0 {
				t.ServerDown = true
				s.Save(r, w)
			}
			i := s.Flashes("interactive")
			if len(i) > 0 {
				t.Interactive = true
				t.HomeServer = i[0].(string)
				s.Save(r, w)
			}
		}

		t.Nonce = nonce

		query := r.URL.Query()
		code := query.Get("code")

		if c.Config.Auth.DisableRegistration && len(code) == 0 {
			c.Templates.ExecuteTemplate(w, "signup-disabled", t)
			return
		}

		if c.Config.Auth.VerifyEmail && code == "" {
			c.Templates.ExecuteTemplate(w, "verify-email", t)
			return
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		email, valid, err := c.GetEmailVerificationToken(ctx, code)
		if err != nil || !valid {
			log.Println(err)
			c.VerificationCodeInvalid(w, r)
			return
		}
		t.Email = email

		c.Templates.ExecuteTemplate(w, "signup", t)
	}
}

func (c *Client) SignupDisabled() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//us := LoggedInUser(r)

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "signup-disabled", t)
	}
}

func (c *Client) VerificationCodeInvalid(w http.ResponseWriter, r *http.Request) {
	us := LoggedInUser(r)

	type page struct {
		BasePage
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{}

	t.Nonce = nonce
	t.LoggedInUser = us

	c.Templates.ExecuteTemplate(w, "invalid-verification-code", t)
}

func (c *Client) PasswordResetCodeInvalid(w http.ResponseWriter, r *http.Request) {
	//us := LoggedInUser(r)

	type page struct {
		BasePage
	}

	nonce := secure.CSPNonce(r.Context())

	t := &page{}

	t.Nonce = nonce

	c.Templates.ExecuteTemplate(w, "invalid-password-reset-code", t)
}

// Copied from Dendrite clientapi/routing/register.go
const (
	minPasswordLength = 8   // http://matrix.org/docs/spec/client_server/r0.2.0.html#password-based
	maxPasswordLength = 512 // https://github.com/matrix-org/synapse/blob/v0.20.0/synapse/rest/client/v2_alpha/register.py#L161
	maxUsernameLength = 254 // http://matrix.org/speculator/spec/HEAD/intro.html#user-identifiers TODO account for domain
	sessionIDLength   = 24
)

func (c *Client) ValidateSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Email string `json:"email"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Emailed bool `json:"emailed"`
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, pay.Email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if !exists && len(pay.Email) > 0 {
			token := RandomNumber(7)

			err = c.EmailVerification(pay.Email, token)

			if err != nil {
				log.Println(err)
			}
		}

		ff := Response{Emailed: true}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) ValidateVerificationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Code string `json:"code"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Email string `json:"email"`
			Valid bool   `json:"valid"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		email, valid, err := c.GetEmailVerificationToken(ctx, pay.Code)
		if err != nil || !valid {
			log.Println(err)
		}

		ff.Valid = valid
		ff.Email = email

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) CompleteSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		valid, err := c.GetEmailVerification(ctx, email)
		if err != nil || !valid {
			log.Println(err)
		}

		if !valid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		type Auth struct {
			Type    string `json:"type"`
			Session string `json:"session"`
			Mac     string `json:"mac"`
		}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			sess := NewSession(c.Config.Client.SecureCookie)
			sess.Options.Domain = fmt.Sprintf(`.%s`, c.Config.Client.Domain)
			s, _ = sess.Get(r, c.Config.Client.CookieName)
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		if strings.Contains(username, ":") {
			_, us := c.IsFederated(username)

			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
			//get rid of the @ prefix
			username = us.LocalPart[1:]
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		/*
			//check if username s available
			av, err := matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
				Username: username,
			})
			if err != nil {
				log.Println(err)

				s.AddFlash("Server Down", "server-down")
				s.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if av == nil || !av.Available {

				s.AddFlash("User Exists", "user-exists")
				s.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		*/

		userID := fmt.Sprintf(`@%s:%s`, username, c.Config.Matrix.FederationServer)

		ctx = context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserIDExist(ctx, userID)
		if err != nil {
			log.Println(err)
		}

		if err != nil || exists {
			s.AddFlash("User Exists", "user-exists")
			s.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//actually register the user

		server := fmt.Sprintf(`http://%s:%d`, c.Config.Matrix.Server, c.Config.Matrix.Port)
		matrix.Prefix = "/_synapse/admin/v1"

		nonce, err := GetNonce(server)
		if err != nil {
			log.Println(err)
		}
		log.Println("nonce is ", nonce)

		//actually register the user
		mac, err := ConstructMac(&NewUser{
			Username: username,
			Password: password,
			Admin:    true,
		}, nonce, c.Config.Auth.SharedSecret)
		if err != nil {
			log.Println(err)
			panic(nil)
		}

		req := &gomatrix.ReqLegacyRegister{
			Username: username,
			Password: password,
			Type:     "org.matrix.login.shared_secret",
			Mac:      mac,
			Admin:    true,
			Nonce:    nonce,
		}

		log.Println("what the hell is usenrmae?", username)

		resp, ui, err := matrix.LegacyRegister(req)

		log.Println("is resp nil?", resp == nil)
		log.Println("is ui nil?", ui == nil)
		log.Println("is err nil?", err == nil)

		if err != nil || resp == nil {
			log.Println(err)
			log.Println(err)
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if resp != nil {
			log.Println(resp)
		}

		if ui != nil {
			log.Println(ui)
		}

		//store session
		token := RandomString(64)
		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			RefreshToken:      resp.RefreshToken,
		}

		serialized, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(token, resp.UserID, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(resp.UserID, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		go func() {
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			err := c.AddNewUser(ctx, resp.UserID, email, resp.AccessToken)
			if err != nil {
				log.Println(err)
			}

			err = c.UnsafeAddEmail(ctx, username, email)
			if err != nil {
				log.Println(err)
			}
		}()

		s.Values["access_token"] = token
		s.Values["device_id"] = resp.DeviceID

		s.AddFlash("Signed Up", "signed-up")
		s.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
}

//log user out, kill session in redis
func (c *Client) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		token, ok := s.Values["access_token"].(string)
		if ok {

			userid, err := c.Store.Get(token).Result()
			if err != nil {
				log.Println(err)
				s.Values["access_token"] = ""
				s.Options.MaxAge = -1
				_ = s.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			_, err = c.Store.Del(userid).Result()
			if err != nil {
				log.Println(err)
				s.Values["access_token"] = ""
				s.Options.MaxAge = -1
				_ = s.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			_, err = c.Store.Del(token).Result()
			if err != nil {
				log.Println(err)
				s.Values["access_token"] = ""
				s.Options.MaxAge = -1
				_ = s.Save(r, w)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			s.Values["access_token"] = ""
			s.Options.MaxAge = -1
			err = s.Save(r, w)
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (c *Client) VerifySignupEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		email := r.FormValue("email")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if !exists && len(email) > 0 {
			go c.SendSignupVerificationEmail(email)
		}

		type page struct {
			BasePage
			UserExists           bool
			ServerDown           bool
			SignupError          bool
			Interactive          bool
			HomeServer           string
			RegistrationDisabled bool
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
		}
		if s != nil {
			x := s.Flashes("user-exists")
			if len(x) > 0 {
				t.UserExists = true
				s.Save(r, w)
			}
			y := s.Flashes("server-down")
			if len(y) > 0 {
				t.ServerDown = true
				s.Save(r, w)
			}
			i := s.Flashes("interactive")
			if len(i) > 0 {
				t.Interactive = true
				t.HomeServer = i[0].(string)
				s.Save(r, w)
			}
		}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "verify-success", t)
	}
}

func (c *Client) UsernameAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Username string `json:"username"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Available bool `json:"available"`
		}

		ff := Response{}

		/*

			userID := fmt.Sprintf(`@%s:%s`, pay.Username, c.Config.Matrix.FederationServer)

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			exists, err := c.DoesUserIDExist(ctx, userID)
			if err != nil {
				log.Println(err)
			}

			if err == nil && !exists {
				ff.Available = true
			}
		*/
		av, err := c.Matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
			Username: pay.Username,
		})

		if err != nil {
			log.Println(err)
		}

		if av != nil && av.Available {
			ff.Available = true
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) PasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type page struct {
			BasePage
			Code  string
			Email string
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		query := r.URL.Query()
		code := query.Get("code")

		if len(code) > 0 {

			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 7*time.Second)
			email, valid, err := c.GetPasswordResetToken(ctx, code)
			if err != nil || !valid {
				log.Println(err)
				c.PasswordResetCodeInvalid(w, r)
				return
			}

			t.Code = code
			t.Email = email

			c.Templates.ExecuteTemplate(w, "password-reset-confirm", t)
			return
		}

		c.Templates.ExecuteTemplate(w, "password-reset", t)
	}
}

func (c *Client) ValidatePasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if exists && len(email) > 0 {
			go c.SendPasswordResetEmail(email)
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		c.Templates.ExecuteTemplate(w, "password-reset-code-sent", t)
	}
}

func (c *Client) ConfirmPasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		email := r.FormValue("email")
		password := r.FormValue("password")

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		userID, _, err := c.GetUser(ctx, email)
		if err != nil && len(userID) == 0 {
			log.Println(err)
			c.Error(w, r)
			return
		}

		username := GetLocalPart(userID)

		err = c.UnsafePasswordReset(ctx, username, password)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
		}

		err = c.InvalidatePasswordResetCode(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
		}

		type page struct {
			BasePage
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		s.AddFlash("Password Reset", "password-reset-success")
		s.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func (c *Client) VerifyPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Email string `json:"email"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Emailed bool `json:"emailed"`
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		exists, err := c.DoesUserExist(ctx, pay.Email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		ff := Response{}

		if exists {
			token := RandomNumber(7)

			err = c.EmailPasswordVerification(pay.Email, token)

			if err != nil {
				log.Println(err)
			}

			ff.Emailed = true
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) ValidatePasswordVerificationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Code string `json:"code"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Email string `json:"email"`
			Valid bool   `json:"valid"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		email, valid, err := c.GetPasswordResetToken(ctx, pay.Code)
		if err != nil || !valid {
			log.Println(err)
		}

		ff.Valid = valid
		ff.Email = email

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) CompletePasswordReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		email := pay.Email
		password := pay.Password

		type Response struct {
			Reset bool `json:"reset"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		userID, _, err := c.GetUser(ctx, email)
		if err != nil && len(userID) == 0 {
			log.Println(err)
			c.Error(w, r)
			return
		}

		username := GetLocalPart(userID)

		err = c.UnsafePasswordReset(ctx, username, password)
		if err != nil {
			log.Println(err)
		}

		if err != nil {
			log.Println(err)
		}

		err = c.InvalidatePasswordResetCode(ctx, email)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		if err != nil {
			log.Println(err)
		}

		ff.Reset = true

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) ValidateInviteCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type payload struct {
			Code string `json:"code"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Valid bool `json:"valid"`
		}

		ff := Response{}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 7*time.Second)
		valid, err := c.IsInviteCodeValid(ctx, pay.Code)
		if err != nil || !valid {
			log.Println(err)
		}

		ff.Valid = valid

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) CompleteSignupTemp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		type Auth struct {
			Type    string `json:"type"`
			Session string `json:"session"`
			Mac     string `json:"mac"`
		}

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		if strings.Contains(username, ":") {
			_, us := c.IsFederated(username)

			wk, err := WellKnown(c.URLScheme(us.ServerName))
			if err != nil {
				log.Println(err)
				c.Error(w, r)
				return
			}
			serverName = c.URLScheme(wk.ServerName)
			//get rid of the @ prefix
			username = us.LocalPart[1:]
		}

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		/*
			//check if username s available
			av, err := matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
				Username: username,
			})
			if err != nil {
				log.Println(err)

				s.AddFlash("Server Down", "server-down")
				s.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if av == nil || !av.Available {

				s.AddFlash("User Exists", "user-exists")
				s.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		*/

		//actually register the user

		server := fmt.Sprintf(`http://%s:%d`, c.Config.Matrix.Server, c.Config.Matrix.Port)
		matrix.Prefix = "/_synapse/admin/v1"

		nonce, err := GetNonce(server)
		if err != nil {
			log.Println(err)
		}
		log.Println("nonce is ", nonce)

		//actually register the user
		mac, err := ConstructMac(&NewUser{
			Username: username,
			Password: password,
			Admin:    true,
		}, nonce, c.Config.Auth.SharedSecret)
		if err != nil {
			log.Println(err)
			panic(nil)
		}

		req := &gomatrix.ReqLegacyRegister{
			Username: username,
			Password: password,
			Type:     "org.matrix.login.shared_secret",
			Mac:      mac,
			Admin:    true,
			Nonce:    nonce,
		}

		log.Println("what is username?", username)

		resp, ui, err := matrix.LegacyRegister(req)

		log.Println("is resp nil?", resp == nil)
		log.Println("is ui nil?", ui == nil)
		log.Println("is err nil?", err == nil)

		if err != nil || resp == nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if resp != nil {
			log.Println(resp)
		}

		if ui != nil {
			log.Println(ui)
		}

		//store session
		token := RandomString(64)
		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			RefreshToken:      resp.RefreshToken,
		}

		matrix.SetCredentials(resp.UserID, resp.AccessToken)
		matrix.Prefix = "/_matrix/client/r0"

		serialized, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(token, resp.UserID, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		err = c.Store.Set(resp.UserID, serialized, 0).Err()
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		s.Values["access_token"] = token
		s.Values["device_id"] = resp.DeviceID

		s.AddFlash("Signed Up", "signed-up")
		s.Save(r, w)

		ds := fmt.Sprintf(`#commune:%s`, c.Config.Matrix.FederationServer)
		ra, err := matrix.ResolveAlias(ds)
		if err != nil {
			log.Println(err)
		}

		log.Println("resp is ", ra.RoomID)

		roomID := string(ra.RoomID)

		log.Println("WHAT")

		// JOIN main commune room first
		_, err = matrix.JoinRoom(roomID, "", nil)
		if err != nil {
			log.Println(err)
		}

		rooms := []string{}

		//grab room state
		state, err := matrix.RoomState(roomID)
		if err != nil {
			log.Println(err)
		}

		if state != nil {
			for _, x := range state {
				if x.Type == "m.space.child" {
					rooms = append(rooms, *x.StateKey)
					/*
						if val, ok := x.Content["default"]; ok {
							if val.(bool) {
								rooms = append(rooms, *x.StateKey)
							}
						}
					*/
					if streams, ok := x.Content["streams"]; ok {

						if chat, ok := streams.(map[string]interface{})["chat"]; ok {
							rooms = append(rooms, chat.(string))
						}
						if topics, ok := streams.(map[string]interface{})["topics"]; ok {
							rooms = append(rooms, topics.(string))
						}
					}
				}
			}
		}

		for _, room := range rooms {

			jr, err := matrix.JoinRoom(room, "", nil)
			if err != nil {
				log.Println(err)
			}
			if jr != nil {
				log.Println(jr)
			}
		}

		// create user profile room

		go func() {

			pl := gomatrix.Event{
				Type: "m.room.power_levels",
				Content: map[string]interface{}{
					"ban": 60,
					"events": map[string]interface{}{
						"m.room.name":         60,
						"m.room.power_levels": 100,
						"m.room.create":       10,
						"m.space.child":       10,
						"m.space.parent":      10,
					},
					"events_default": 10,
					"invite":         10,
					"kick":           60,
					"notifications": map[string]interface{}{
						"room": 20,
					},
					"redact":        10,
					"state_default": 10,
					"users": map[string]interface{}{
						resp.UserID: 100,
					},
					"users_default": 10,
				},
			}

			initState := []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "commune.room",
					Content: map[string]interface{}{
						"room_type": "profile",
					},
				}, gomatrix.Event{
					Type: "m.room.type",
					Content: map[string]interface{}{
						"type":  "m.space",
						"alias": username,
					},
				},
				pl,
			}

			creq := &gomatrix.ReqCreateRoom{
				RoomAliasName: fmt.Sprintf(`@%s`, username),
				RoomVersion:   "9",
				Preset:        "public_chat",
				Visibility:    "public",
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: initState,
			}

			crr, err := matrix.CreateRoom(creq)

			if err != nil || crr == nil {
				log.Println(err)
				log.Println(err)
				http.Error(w, err.Error(), 400)
				return
			}

			log.Println("Was Room created?", crr)
		}()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
}

//Log user in via API
func (c *Client) APILogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//user := LoggedInUser(r)

		s, err := GetSession(r, c)
		if err != nil || s == nil {
			log.Println(err)
			type Response struct {
				Error string `json:"error"`
			}

			ff := Response{
				Error: "Authentication error.",
			}
			js, err := json.Marshal(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		type payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var pay payload
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		username := pay.Username
		password := pay.Password

		if username == "" || password == "" {
			type Response struct {
				Error string `json:"error"`
			}

			ff := Response{
				Error: "Invalid username or password.",
			}
			js, err := json.Marshal(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		//port is only for my dev environment, this needs to go, or i'm just
		//confused
		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		matrix, err := gomatrix.NewClient(serverName, "", "")
		if err != nil {
			log.Println(err)
			c.Error(w, r)
			return
		}

		rl := &gomatrix.ReqLogin{
			Type:     "m.login.password",
			User:     username,
			Password: password,
		}

		resp, err := matrix.Login(rl)
		if err != nil {
			type Response struct {
				Error string `json:"error"`
			}

			ff := Response{
				Error: "Invalid username or password.",
			}
			js, err := json.Marshal(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		matrix.SetCredentials(resp.UserID, resp.AccessToken)

		// store user session
		token := RandomString(64)

		u := User{
			AccessToken:       token,
			MatrixAccessToken: resp.AccessToken,
			DeviceID:          resp.DeviceID,
			HomeServer:        resp.HomeServer,
			UserID:            resp.UserID,
			WellKnown:         serverName,
		}

		type Response struct {
			Authenticated bool        `json:"authenticated"`
			Identity      interface{} `json:"identity"`
		}

		acc, err := c.ConstructAccount(&u)

		if err != nil {
			log.Println(err)
		}

		us := acc.User
		us.SyncState = acc.SyncState
		us.Spaces = acc.Spaces
		us.TimelineState = acc.TimelineState
		us.Notifications = acc.Notifications
		us.UserStatus = acc.UserStatus
		us.AccountData = acc.AccountData

		ff := Response{
			Authenticated: true,
			Identity:      us,
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {

			tok, ok := s.Values["access_token"].(string)
			if ok {
				userid, err := c.Store.Get(tok).Result()
				if err != nil {
					log.Println(err)
				}

				user, err := c.Store.Get(userid).Result()
				if err != nil {
					log.Println(err)
				}

				var us User
				err = json.Unmarshal([]byte(user), &us)
				if err != nil {
					log.Println(err)
				}

				us.AltAccounts = append(us.AltAccounts, &u)

				serialized, err := json.Marshal(us)
				if err != nil {
					log.Println(err)
				}

				err = c.Store.Set(userid, serialized, 0).Err()
				if err != nil {
					log.Println(err)
				}

			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (c *Client) HandleDiscordAuthorized() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		code := query.Get("code")

		clientid := c.Config.Oauth2.Discord.ClientID
		clientsecret := c.Config.Oauth2.Discord.ClientSecret

		data := url.Values{
			"client_id":     {clientid},
			"client_secret": {clientsecret},
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {"http://localhost:8989/oauth2/discord/authorize"},
		}

		resp, err := http.PostForm("https://discord.com/api/oauth2/token", data)

		if err != nil {
			log.Fatal(err)
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		fmt.Println(res)

		w.Write([]byte("lol"))
	}
}

func (c *Client) HandleDiscordExchanged() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("exchanged"))
	}
}
