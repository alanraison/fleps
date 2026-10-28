// Package logos finds team badge image URLs on football-logos.cc.
//
// Badges live at /{country}/{club}/ on the site, with image filenames that
// include a content hash. The hashes are taken from the catalog published by
// the football-logos npm package and served through jsDelivr.
package logos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"unicode"

	"github.com/alanraison/fleps/pkg/model"
)

const (
	defaultCatalogURL = "https://cdn.jsdelivr.net/gh/julianYaman/football-logos@main/catalog/v1/countries/%s.json"
	assetURL          = "https://assets.football-logos.cc/logos/%s/512x512/%s.%s.png"
	country           = "england"
)

type catalog struct {
	Country struct {
		Clubs map[string]struct {
			Slug string `json:"slug"`
			Name string `json:"name"`
			Hash string `json:"hash"`
		} `json:"clubs"`
	} `json:"country"`
}

// Resolver looks up badge URLs for teams.
type Resolver struct {
	catalogURL string
	client     *http.Client

	mu   sync.Mutex
	urls map[string]string
}

// NewResolver returns a Resolver using the live catalog.
func NewResolver() *Resolver {
	return &Resolver{catalogURL: fmt.Sprintf(defaultCatalogURL, country), client: http.DefaultClient}
}

// BadgeURL returns the badge URL for the team, or "" if none could be matched.
// A non-empty slug is used directly; otherwise the team's name is matched.
func (r *Resolver) BadgeURL(ctx context.Context, team model.Team, slug string) (string, error) {
	urls, err := r.load(ctx)
	if err != nil {
		return "", err
	}
	return match(urls, team, slug), nil
}

func (r *Resolver) load(ctx context.Context) (map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.urls != nil {
		return r.urls, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.catalogURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating logo catalog request: %w", err)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching logo catalog: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching logo catalog: unexpected status %s", resp.Status)
	}
	var c catalog
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, fmt.Errorf("decoding logo catalog: %w", err)
	}
	urls := make(map[string]string, len(c.Country.Clubs))
	for key, club := range c.Country.Clubs {
		if club.Slug == "" {
			club.Slug = key
		}
		if club.Hash != "" {
			urls[club.Slug] = fmt.Sprintf(assetURL, country, club.Slug, club.Hash)
		}
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("logo catalog has no clubs")
	}
	r.urls = urls
	return urls, nil
}

// match returns the logo for the given slug if it is in the catalog, otherwise
// finds a logo by the slug of the team's full or short name, falling back
// to the longest logo slug that is a prefix of the full name's slug (so that
// "Brighton & Hove Albion" finds "brighton").
func match(urls map[string]string, team model.Team, slug string) string {
	if u, ok := urls[slug]; ok {
		return u
	}
	full := slugify(team.FullName)
	for _, s := range []string{full, slugify(team.ShortName)} {
		if u, ok := urls[s]; ok {
			return u
		}
	}
	best := ""
	for s := range urls {
		if strings.HasPrefix(full, s+"-") && len(s) > len(best) {
			best = s
		}
	}
	return urls[best]
}

func slugify(name string) string {
	var b strings.Builder
	pendingDash := false
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if pendingDash && b.Len() > 0 {
				b.WriteByte('-')
			}
			pendingDash = false
			b.WriteRune(r)
		} else if r != '\'' {
			pendingDash = true
		}
	}
	return b.String()
}
