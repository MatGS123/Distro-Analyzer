package collect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"distroanalyzer/profile"
)

// GitHubCollector recolecta datos de perfiles públicos de GitHub.
type GitHubCollector struct {
	client *http.Client
	token  string // Token opcional para aumentar rate limits
}

// NewGitHubCollector crea un nuevo collector para GitHub.
// El token es opcional pero recomendado para evitar rate limits (60 req/h sin auth).
func NewGitHubCollector(token string) *GitHubCollector {
	return &GitHubCollector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

// Collect obtiene bio, repos y README del usuario de GitHub.
func (g *GitHubCollector) Collect(username string) (*profile.RawData, error) {
	userURL := fmt.Sprintf("https://api.github.com/users/%s", username)

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	g.setHeaders(req)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user githubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	repos, err := g.fetchRepos(username)
	if err != nil {
		repos = []string{}
	}

	readme := g.fetchReadme(username, repos)

	return &profile.RawData{
		Bio:          user.Bio,
		Repositories: repos,
		Website:      user.Blog,
		Location:     user.Location,
		Email:        user.Email,
		ReadmeText:   readme,
	}, nil
}

func (g *GitHubCollector) fetchRepos(username string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=10", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	g.setHeaders(req)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []githubRepo
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	names := make([]string, len(repos))
	for i, r := range repos {
		names[i] = r.Name
	}

	return names, nil
}

func (g *GitHubCollector) fetchReadme(username string, repos []string) *string {
	if len(repos) == 0 {
		return nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", username, repos[0])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	text := string(body)
	return &text
}

// setHeaders configura headers comunes para requests a GitHub API.
func (g *GitHubCollector) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}
}

// githubUser representa la respuesta de la API de GitHub para un usuario.
type githubUser struct {
	Bio      string `json:"bio"`
	Blog     string `json:"blog"`
	Location string `json:"location"`
	Email    string `json:"email"`
}

// githubRepo representa un repositorio en la respuesta de GitHub.
type githubRepo struct {
	Name string `json:"name"`
}
