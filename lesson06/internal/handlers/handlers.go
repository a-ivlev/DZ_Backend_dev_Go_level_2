package handlers

import (
	"context"
	"fmt"
	"lesson06/internal/app/sessions"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

var CreateFormTmpl = []byte(`
<html>
	<body>
		<form action="/verify" method="post">
			Login: <input type="text" name="login">
			Password: <input type="password" name="password">
			RePassword: <input type="password" name="passwordRepeat">
			<input type="submit" value="Verify Account">
		</form>
	</body>
</html>
`)

var VerifyFormTmpl = []byte(`
<html>
	<body>
		<form action="/create" method="post">
			Login: <input type="text" name="login">
			Code: <input type="code" name="code">
			<input type="submit" value="Verify">
		</form>
	</body>
</html>
`)

var loginFormTmpl = []byte(`
<html>
	<body>
		<form action="/login" method="post">
			Login: <input type="text" name="login">
			Password: <input type="password" name="password">
			<input type="submit" value="Login">
		</form>
		<form action="/registry" method="post">
			<input type="submit" value="Create an Account">
		</form>
	</body>
</html>
`)

const (
	loginValue          = "login"
	passwordValue       = "password"
	passwordValueRepeat = "passwordRepeat"
	codeValue           = "code"
)

var welcome = "Welcome, %s <br />\nSession User-Agent: %s <br />\n<a href=\"/logout\">logout</a>"

var badRepeatPassword = "The passwords you entered do not match, %s %s <br />\n<a href=\"/registry\">Create an Account</a>"

var badCode = "The code you entered do not match, %s <br />\n<a href=\"/registry\">Create an Account</a>"
var badPass = "The password you entered do not match, %s <br />\n<a href=\"/registry\">Create an Account</a>"

type Router struct {
	*http.ServeMux
	sc  *sessions.SessionCache
	ttl time.Duration
}

func NewRouter(sessionCache *sessions.SessionCache, ttl time.Duration) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		sc:       sessionCache,
		ttl:      ttl,
	}

	r.Handle("/", http.HandlerFunc(r.RootHandler))
	r.Handle("/registry", http.HandlerFunc(r.RegistryUserHandler))
	r.Handle("/create", http.HandlerFunc(r.CreateHandler))
	r.Handle("/verify", http.HandlerFunc(r.VerifyUserHandler))
	r.Handle("/login", http.HandlerFunc(r.LoginHandler))
	r.Handle("/logout", http.HandlerFunc(r.LogoutHandler))
	return r
}

func (rt *Router) checkSession(r *http.Request) (*sessions.Session, error) {
	cookieSessionID, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess, err := rt.sc.Check(sessions.SessionID{ID: cookieSessionID.Value})
	if err != nil {
		return nil, fmt.Errorf("check sessions value %q: %w", cookieSessionID.Value,
			err)
	}
	return sess, nil
}

func (rt *Router) RootHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := rt.checkSession(r)
	if err != nil {
		err = fmt.Errorf("check sessions: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sess == nil {
		_, _ = w.Write(loginFormTmpl)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintln(w, fmt.Sprintf(welcome, sess.Login, sess.Useragent))
}

func (rt *Router) RegistryUserHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(CreateFormTmpl)
}

func (rt *Router) VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue(loginValue)

	inputPass := r.FormValue(passwordValue)
	inputPassRepeat := r.FormValue(passwordValueRepeat)

	if inputPass != inputPassRepeat {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintln(w, fmt.Sprintf(badRepeatPassword, inputPass, inputPassRepeat))
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "inputLogin", inputLogin))

	rawCode, _ := uuid.NewUUID()
	code := strings.Split((rawCode).String(), "-")

	fmt.Println("Отправка кода смской: ", code[3])

	err := rt.sc.SetCache(inputLogin+"code", []byte(code[3]))
	if err != nil {
		err = fmt.Errorf("create code Redis DB: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.sc.SetCache(inputLogin+"pass", []byte(inputPass))
	if err != nil {
		err = fmt.Errorf("error writing user to cache: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(VerifyFormTmpl)
}

func (rt *Router) CreateHandler(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue(loginValue)

	inputCode := r.FormValue(codeValue)
	code, err := rt.sc.GetRecordCache(inputLogin + "code")
	if string(code) != inputCode {
		log.Println("ERROR ENTER BAD CODE")
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintln(w, fmt.Sprintf(badCode, inputCode))
		return
	}

	//pass, err := rt.sc.GetRecordCache(inputLogin + "pass")
	//log.Printf("pass = %s err = %v", string(pass), err)

	sess, err := rt.sc.Create(sessions.Session{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		err = fmt.Errorf("create sessions: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   sess.ID,
		Expires: time.Now().Add(rt.ttl),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

const cookieName = "session_id"

func (rt *Router) LoginHandler(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue(loginValue)
	inputPass := r.FormValue(passwordValue)

	pass, err := rt.sc.GetRecordCache(inputLogin + "pass")
	if string(pass) != inputPass {
		log.Println("ERROR ENTER BAD PASSWORD")
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintln(w, fmt.Sprintf(badPass, inputPass))
		return
	}

	sess, err := rt.sc.Create(sessions.Session{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		err = fmt.Errorf("create sessions: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   sess.ID,
		Expires: time.Now().Add(rt.ttl),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (rt *Router) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		err = fmt.Errorf("read cookie %q: %w", cookieName, err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.sc.Delete(sessions.SessionID{ID: session.Value})
	if err != nil {
		err = fmt.Errorf("delete sessions value %q: %w", session.Value, err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}
