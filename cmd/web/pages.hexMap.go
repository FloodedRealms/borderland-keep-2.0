package main

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
)

const hexKey contextKey = "hexmap"

type HexMap struct {

}

func NewHexMap() *HexMap {
	return &HexMap{}
}

type PageHexEditor struct {
	renderer *htmlRenderer
	logger *slog.Logger
	templateName string
	formTemplateName string
	sessionData map[string]*HexMap
	mu *sync.RWMutex
}

func (p PageHexEditor) RegisterRoutes(router *http.ServeMux) {
	router.Handle("/tools/hexhistory", p.addHexToContext(p.index()))
}

func (p PageHexEditor) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r* http.Request) {

		h := p.getHexFromContext(r)
		//TODO: Create Hex Map Pages
		err := p.renderer.render(w, http.StatusOK, h, "base")
		if err != nil {
			p.logger.Info("Error rendering initial load of hex map:")
			p.logger.Error(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

	}
}


func (p PageHexEditor) addHexToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		session := getSessionFromCtx(r)
		if p.sessionData[session.Id] == nil {
			p.sessionData[session.Id] = NewHexMap()
		}
		a := p.sessionData[session.Id]
		ctx := context.WithValue(r.Context(), hexKey, a)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (p PageHexEditor) getHexFromContext(r *http.Request) *HexMap {
	p.mu.RLock()
	a, _ := r.Context().Value(hexKey).(*HexMap)
	p.mu.Unlock()
	return a
}
