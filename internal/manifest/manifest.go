package manifest

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/validator/v10"

	"github.com/selesy/asdf-go-install/internal/config"
	"github.com/selesy/asdf-go-install/internal/plugin"
)

const (
	ManifestFilename  = "manifest.json"
	manifestVersionV1 = "v1"
)

var (
	_ json.Marshaler   = (*payload)(nil)
	_ json.Unmarshaler = (*payload)(nil)
)

type payload struct {
	PluginName    string              `json:"pluginName" validate:"required"`
	PackageName   string              `json:"packageName" validate:"required"`
	GitRepository *url.URL            `json:"gitRepository" validate:"required"`
	GitReference  *plumbing.Reference `json:"gitReference"`
}

// MarshalJSON implements json.Marshaler.
func (p *payload) MarshalJSON() ([]byte, error) {
	type gitReference struct {
		Name string `json:"name"`
		Hash string `json:"hash"`
	}

	var ref *gitReference

	if p.GitReference != nil {
		tkns := p.GitReference.Strings()

		ref = &gitReference{
			Name: tkns[0],
			Hash: tkns[1],
		}
	}

	type alias payload

	return json.Marshal(&struct {
		GitRepository string        `json:"gitRepository" validate:"required"`
		GitReference  *gitReference `json:"gitReference,omitempty"`
		*alias
	}{
		GitRepository: p.GitRepository.String(),
		GitReference:  ref,
		alias:         (*alias)(p),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *payload) UnmarshalJSON(data []byte) error {
	type alias payload
	clone := &struct {
		GitRepository string `json:"gitRepository" validate:"required"`
		GitReference  struct {
			Name string
			Hash string
		} `json:"gitReference"`
		*alias
	}{
		alias: (*alias)(p),
	}

	if err := json.Unmarshal(data, &clone); err != nil {
		return err
	}

	var err error

	p.GitRepository, err = url.Parse(clone.GitRepository)
	if err != nil {
		return err
	}

	p.GitReference = plumbing.NewReferenceFromStrings(clone.GitReference.Name, clone.GitReference.Hash)

	return nil
}

var _ json.Marshaler = (*manifest)(nil)

type manifest struct {
	ManifestVersion *semver.Version `json:"manifestVersion" validate:"required"`
	Payload         *payload        `json:"manifestPayload" validate:"required"`
}

// MarshalJSON implements json.Marshaler.
func (m *manifest) MarshalJSON() ([]byte, error) {
	type alias manifest

	return json.Marshal(&struct {
		ManifestVersion string `json:"manifestVersion"`
		*alias
	}{
		ManifestVersion: m.ManifestVersion.Original(),
		alias:           (*alias)(m),
	})
}

// Manifest information needed to allow the Plugin to manage
// installations of the desired tool.
type Manifest struct {
	manifest *manifest
}

// New creates an immutable instance of a Manifest.
func New(name string, pkg string, repo *url.URL) *Manifest {
	vers := semver.MustParse(manifestVersionV1)

	return &Manifest{
		manifest: &manifest{
			ManifestVersion: vers,
			Payload: &payload{
				PluginName:    name,
				PackageName:   pkg,
				GitRepository: repo,
			},
		},
	}
}

// Read opens the manifest file in the plugin's top-level directory and
// decodes the JSON into a Manifest.
func Read(cfg *config.Config, pluginName string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(cfg.Env().DataDir(), "plugins", pluginName, ManifestFilename))
	if err != nil {
		return nil, err
	}

	var man manifest

	if err := json.Unmarshal(data, &man); err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(man); err != nil {
		return nil, err
	}

	return &Manifest{
		manifest: &man,
	}, nil
}

// GitReferenece returns the plugin's Git reference or nil if no Git
// reference is defined.
func (m *Manifest) GitReference() *plumbing.Reference {
	return m.manifest.Payload.GitReference
}

// GitRepository returns the URL of the plugin's Git repository.
func (m *Manifest) GitRepository() *url.URL {
	return m.manifest.Payload.GitRepository
}

// ManifestVersion returns the version of the Manifest.
func (m *Manifest) ManifestVersion() *semver.Version {
	return m.manifest.ManifestVersion
}

// PluginName returns the plugin's name.
func (m *Manifest) PluginName() string {
	return m.manifest.Payload.PluginName
}

// PluginPackage returns the plugin's package.
func (m *Manifest) PluginPackage() string {
	return m.manifest.Payload.PackageName
}

// WithGitReference creates a clone of the Manifest that includes the
// provided Git reference.
func (m *Manifest) WithGitReference(ref *plumbing.Reference) *Manifest {
	return &Manifest{
		manifest: &manifest{
			ManifestVersion: m.manifest.ManifestVersion,
			Payload: &payload{
				PluginName:    m.manifest.Payload.PluginName,
				PackageName:   m.manifest.Payload.PackageName,
				GitRepository: m.manifest.Payload.GitRepository,
				GitReference:  ref,
			},
		},
	}
}

// Write encodes the Manifest to JSON and creates the relevant file in
// the plugin's top-level directory.
func (m *Manifest) Write(cfg *config.Config, pluginName string) error {
	data, err := json.Marshal(m.manifest)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(plugin.Path(cfg, pluginName), ManifestFilename), data, 0o644)
}
