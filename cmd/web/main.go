package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/FloodedRealms/borderland-keep-2.0/assets"
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

// The application struct holds the dependencies needed for our handlers,
// including a htmlRenderer type.
type application struct {
    logger *slog.Logger
    html   *htmlRenderer
	pages []page
}

type page interface {
	RegisterRoutes(router *http.ServeMux)
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
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    // Initialize a new htmlRenderer, parsing the base template and all partial
    // templates from assets/html into the shared template set.
    htmlRenderer, err := newHTMLRenderer(assets.HTMLFiles, "base.html", "printableBase.html", "partials/*.html")
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

	//pages
	calculatorPage := NewPageAdventureCalculator(htmlRenderer, logger)
	weatherPage := NewPageWeather(htmlRenderer, logger)

	app := &application{
        logger: logger,
        html:   htmlRenderer,
		pages: []page{calculatorPage, weatherPage},
    }

	go cleanupSessions(10*time.Minute, 1*time.Hour, []SessionDataUser{calculatorPage, weatherPage})
	//router
	router := http.NewServeMux()

	fs := http.FileServerFS(assets.StaticFiles)
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	for _, page := range(app.pages) {
		page.RegisterRoutes(router)
	}

	// "Static" pages

	router.HandleFunc("/{$}",  app.home)
	//router.Handle("/about",  renderStaticPage("about.html", staticRenderer))
	//router.Handle("/donate", renderStaticPage("donate.html", staticRenderer))


	// Tool Pages
	// calculatorPage.RegisterRoutes(router)
	// weatherPage.RegisterRoutes(router)

	//User Pages

	server := &http.Server{
		Addr:    ":9090",
		Handler: sessionMiddleware(router),
	}
    logger.Info("starting server", "port", 9090)
    err = server.ListenAndServe()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
	}

}
