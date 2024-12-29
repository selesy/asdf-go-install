// Package pkgsite includes functions needed to scrape tool information
// from pkg.go.dev.
package pkgsite

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly"
	"github.com/lmittmann/tint"

	"github.com/selesy/asdf-go-install/internal/config"
	"github.com/selesy/asdf-go-install/internal/gover"
)

const (
	packageSiteBaseURL         = "https://pkg.go.dev"
	packageSiteTabKey          = "tab"
	packageSiteVersionTabValue = "versions"
)

// Repository scrapes the URL of the Go package's Git repository from
// the pkg.go.dev web-site.
func Repository(cfg *config.Config, pkg string) (*url.URL, error) {
	var (
		u   *url.URL
		err error
	)

	u, err = url.Parse(packageSiteBaseURL)
	if err != nil {
		return nil, err
	}

	u.Path = "/" + pkg

	cfg.Log().Debug(
		"Scraping target repository URL",
		slog.String("url", u.String()),
		slog.String("goal", "repository"),
	)

	col := colly.NewCollector()

	var repo *url.URL

	col.OnError(func(r *colly.Response, e error) {
		cfg.Log().Error("Colly error", tint.Err(err))
		err = e
	})

	col.OnHTML("html body main aside div.UnitMeta div.UnitMeta-repo a", func(h *colly.HTMLElement) {
		repo, err = url.Parse(h.Attr("href"))
	})

	if err := col.Visit(u.String()); err != nil {
		return nil, err
	}

	col.Wait()

	return repo, err
}

var _ gover.Collector = Versions

// Versions scrapes the available versions of the Go package from the
// pkg.go.dev web-site.
func Versions(cfg *config.Config, pkg string) (*gover.Collection, error) {
	u, err := url.Parse(packageSiteBaseURL)
	if err != nil {
		return nil, err
	}

	u.Path = "/" + pkg
	u.RawQuery = "tab=versions"

	cfg.Log().Debug(
		"Scraping target",
		slog.String("url", u.String()),
		slog.String("goal", "versions"),
	)

	col := colly.NewCollector()

	col.OnError(func(r *colly.Response, err error) {
		cfg.Log().Error("Colly error", tint.Err(err))
		// TODO: this err should stop the plugin execution
	})

	var vers semver.Collection

	col.OnHTML("html body main article div.Versions div.Versions-list div.Version-tag a", func(h *colly.HTMLElement) {
		ver, err := gover.NewVersion(h.Text)
		if err != nil {
			cfg.Log().Warn(
				"Skipping invalid Go version",
				slog.String("package", pkg),
				slog.String("version", strings.TrimPrefix(h.Text, "v")),
				tint.Err(err),
			)

			return
		}

		vers = append(vers, ver)
	})

	if err := col.Visit(u.String()); err != nil {
		return nil, err
	}

	col.Wait()

	return gover.NewCollection(vers...), nil
}
