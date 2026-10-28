package logos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/alanraison/fleps/pkg/model"
)

const catalogJSON = `{"country":{"clubs":{
"arsenal":{"slug":"arsenal","name":"Arsenal","hash":"aaaa1111"},
"brighton":{"slug":"brighton","name":"Brighton","hash":"bbbb2222"},
"manchester-city":{"slug":"manchester-city","name":"Manchester City","hash":"dddd4444"},
"wolves":{"slug":"wolves","name":"Wolverhampton Wanderers","hash":"eeee5555"}}}}`

func TestBadgeURL(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Write([]byte(catalogJSON))
	}))
	defer srv.Close()
	r := &Resolver{catalogURL: srv.URL, client: srv.Client()}

	const base = "https://assets.football-logos.cc/logos/england/512x512/"
	tests := []struct {
		team model.Team
		slug string
		want string
	}{
		{model.Team{FullName: "Arsenal", ShortName: "Arsenal"}, "", base + "arsenal.aaaa1111.png"},
		{model.Team{FullName: "Brighton & Hove Albion", ShortName: "Brighton"}, "", base + "brighton.bbbb2222.png"},
		{model.Team{FullName: "Manchester City", ShortName: "Man City"}, "", base + "manchester-city.dddd4444.png"},
		{model.Team{FullName: "Wolverhampton Wanderers", ShortName: "Wolves"}, "", base + "wolves.eeee5555.png"},
		{model.Team{FullName: "Anything"}, "wolves", base + "wolves.eeee5555.png"},
		{model.Team{FullName: "Nowhere FC", ShortName: "Nowhere"}, "", ""},
	}
	for _, tt := range tests {
		got, err := r.BadgeURL(context.Background(), tt.team, tt.slug)
		if err != nil || got != tt.want {
			t.Errorf("%s: got %q, %v; want %q", tt.team.FullName, got, err, tt.want)
		}
	}
	if hits.Load() != 1 {
		t.Errorf("expected catalog to be fetched once, got %d", hits.Load())
	}
}

func TestBadgeURLErrors(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	r := &Resolver{catalogURL: srv.URL, client: srv.Client()}
	if _, err := r.BadgeURL(context.Background(), model.Team{FullName: "Arsenal"}, ""); err == nil {
		t.Error("expected error")
	}
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"Brighton & Hove Albion": "brighton-hove-albion", "B'mth": "bmth", "  Man  City ": "man-city"} {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q)=%q want %q", in, got, want)
		}
	}
}
