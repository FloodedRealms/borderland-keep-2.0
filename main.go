package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id string
	LastSeen time.Time
}

var (
	sessions = make(map[string]*Session)
	sessionMu sync.RWMutex
)

type contextKey string
const sessionKey contextKey = "session"

func generateSessionID() string {
	return uuid.NewString()
}

type SessionDataUser interface {
	CleanupSessionData(string)
}

func cleanupSessions(interval, maxAge time.Duration, consumers []SessionDataUser) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        cutoff := time.Now().Add(-maxAge)
        sessionMu.Lock()
        for id, s := range sessions {
            if s.LastSeen.Before(cutoff) {
                delete(sessions, id)
				for _, c := range consumers {
					c.CleanupSessionData(id)
				}
            }
        }
        sessionMu.Unlock()
    }
}

func getSession(r *http.Request) (string, *Session) {
    cookie, err := r.Cookie("session_id")

    if err == nil {
        sessionMu.RLock()
        session, exists := sessions[cookie.Value]
        sessionMu.RUnlock()
        if exists {
            session.LastSeen = time.Now()
            return cookie.Value, session
        }
    }

    // No valid session — create one
    id := generateSessionID()
    session := &Session{LastSeen: time.Now(), Id: id}
    sessionMu.Lock()
    sessions[id] = session
    sessionMu.Unlock()
    return id, session
}

func sessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sessionID, session := getSession(r)

        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    sessionID,
            HttpOnly: true,
            Path:     "/",
        })

        // Attach session to the request context
        ctx := context.WithValue(r.Context(), sessionKey, session)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func getSessionFromCtx(r *http.Request) *Session {
    session, _ := r.Context().Value(sessionKey).(*Session)
    return session
}

func main() {


	sessionMu = sync.RWMutex{}
	//pages
	calculatorPage := NewPageAdventureCalculator()


	go cleanupSessions(10*time.Minute, 1*time.Hour, []SessionDataUser{calculatorPage})
	staticRenderer := *NewRenderer()
	//router
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// "Static" pages

	router.Handle("/", renderStaticPage("index.html", staticRenderer))
	router.Handle("/index",  renderStaticPage("index.html", staticRenderer))
	router.Handle("/dicegoblin",  renderStaticPage("dicegoblin.html", staticRenderer))
	router.Handle("/legal",  renderStaticPage("legal.html", staticRenderer))
	// Tool Pages
	calculatorPage.RegisterRoutes(router)

	//User Pages

	server := &http.Server{
		Addr:    ":9090",
		Handler: sessionMiddleware(router),
	}
	log.Print("Listening on 9090")
	for {
		server.ListenAndServe()
		log.Print("Server crash... attempting restart")
	}
}

func renderStaticPage(pageName string, pr Renderer) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		renderedPage, err := pr.RenderPage(pageName, nil)
		if err != nil {
			log.Printf("Error rendering Static Page %s: %v\n", pageName, err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}
