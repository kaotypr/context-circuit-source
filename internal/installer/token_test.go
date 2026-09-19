package installer_test

import (
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	assets "github.com/kaotypr/context-circuit-source"
)

// Field order matches the GitHub payload: an asset's own URL precedes the
// nested uploader object, so both survive a brace-delimited read.
type releaseAsset struct {
	URL                string         `json:"url"`
	ID                 int            `json:"id"`
	NodeID             string         `json:"node_id"`
	Name               string         `json:"name"`
	Label              *string        `json:"label"`
	Uploader           map[string]any `json:"uploader"`
	ContentType        string         `json:"content_type"`
	State              string         `json:"state"`
	Size               int            `json:"size"`
	DownloadCount      int            `json:"download_count"`
	BrowserDownloadURL string         `json:"browser_download_url"`
}

type release struct {
	URL        string         `json:"url"`
	ID         int            `json:"id"`
	NodeID     string         `json:"node_id"`
	TagName    string         `json:"tag_name"`
	Name       string         `json:"name"`
	Body       string         `json:"body"`
	Draft      bool           `json:"draft"`
	Prerelease bool           `json:"prerelease"`
	Assets     []releaseAsset `json:"assets"`
}

// A private repository serves release assets only through the API, so the
// installer resolves an asset id and presents a token. The fixture answers only
// authenticated requests, which is what proves the credential is sent. The API
// pretty-prints its payload, so the asset reader is exercised against that shape
// as well as a compact one.
func TestTokenInstallResolvesPrivateReleaseAssets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("install.ps1 uses its own HTTP client; this fixture drives curl")
	}
	for _, shape := range []struct {
		name   string
		indent bool
	}{{"pretty", true}, {"compact", false}} {
		t.Run(shape.name, func(t *testing.T) {
			runTokenInstall(t, shape.indent)
		})
	}
}

func runTokenInstall(t *testing.T, indent bool) {
	t.Helper()
	root := t.TempDir()
	const token = "ghp_fixture_TOKEN-0123456789"
	version := "2.0.0-token"
	archive, sums := packageFixture(t, root, version)
	packageName := filepath.Base(archive)

	var served, unauthorized int
	authorized := func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("Authorization") != "Bearer "+token {
			unauthorized++
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return false
		}
		return true
	}
	mux := http.NewServeMux()
	var base string
	serve := func(path string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !authorized(w, r) {
				return
			}
			if r.Header.Get("Accept") != "application/octet-stream" {
				http.Error(w, "wrong accept", http.StatusNotAcceptable)
				return
			}
			served++
			http.ServeFile(w, r, path)
		}
	}
	const prefix = "/repos/kaotypr/context-circuit-source/releases"
	mux.HandleFunc(prefix+"/assets/564138120", serve(archive))
	mux.HandleFunc(prefix+"/assets/564138121", serve(sums))
	mux.HandleFunc(prefix+"/tags/", func(w http.ResponseWriter, r *http.Request) {
		if !authorized(w, r) {
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/cli-v"+version) {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		asset := func(id int, name string) releaseAsset {
			return releaseAsset{
				URL:                base + prefix + "/assets/" + itoa(id),
				ID:                 id,
				NodeID:             "RA_" + itoa(id),
				Name:               name,
				Uploader:           map[string]any{"login": "kaotypr", "id": 99, "node_id": "U_1", "type": "User"},
				ContentType:        "application/octet-stream",
				State:              "uploaded",
				Size:               1,
				BrowserDownloadURL: "https://github.com/kaotypr/context-circuit-source/releases/download/cli-v" + version + "/" + name,
			}
		}
		payload := release{
			URL: base + prefix + "/7", ID: 7, NodeID: "RE_7",
			TagName: "cli-v" + version, Name: "Context Circuit CLI " + version,
			// Release notes carry braces and quotes, which must not be mistaken
			// for asset structure.
			Body:       "Run `context-circuit-cli help`.\nSeeds write {} for an empty container.\n",
			Prerelease: true,
			Assets:     []releaseAsset{asset(564138120, packageName), asset(564138121, "SHA256SUMS")},
		}
		var body []byte
		var err error
		if indent {
			body, err = json.MarshalIndent(payload, "", "  ")
		} else {
			body, err = json.Marshal(payload)
		}
		if err != nil {
			t.Error(err)
			http.Error(w, "fixture", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})
	server := httptest.NewTLSServer(mux)
	defer server.Close()
	base = server.URL

	// curl must trust the fixture without relaxing certificate verification.
	bundle := filepath.Join(root, "fixture-ca.pem")
	certificate := server.Certificate()
	if err := os.WriteFile(bundle, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw}), 0600); err != nil {
		t.Fatal(err)
	}

	files, err := assets.Files()
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(root, "install.sh")
	if err := os.WriteFile(script, files[".agents/skills/cc-cli/scripts/install.sh"], 0600); err != nil {
		t.Fatal(err)
	}
	binDir := filepath.Join(root, "bin")
	run := func(environment []string, arguments ...string) ([]byte, error) {
		command := exec.Command("sh", append([]string{script, "--version", version, "--bin-dir", binDir}, arguments...)...)
		command.Env = append(os.Environ(), append([]string{
			"CURL_CA_BUNDLE=" + bundle,
			"CONTEXT_CIRCUIT_API=" + server.URL,
		}, environment...)...)
		return command.CombinedOutput()
	}

	data, err := run(nil, "--token", token)
	if err != nil {
		t.Fatalf("token install: %v\n%s", err, data)
	}
	if served != 2 {
		t.Fatalf("expected the package and checksums to be served, got %d", served)
	}
	installed, err := exec.Command(filepath.Join(binDir, "context-circuit-cli"), "version").CombinedOutput()
	if err != nil || strings.TrimSpace(string(installed)) != version {
		t.Fatalf("installed command: %v %s", err, installed)
	}

	// The same credential is accepted from the environment.
	if data, err := run([]string{"CONTEXT_CIRCUIT_TOKEN=" + token}); err != nil {
		t.Fatalf("environment token install: %v\n%s", err, data)
	}

	// An unauthenticated request reaches the fixture as a 404, which the
	// installer must report rather than proceed from.
	data, err = run([]string{"CONTEXT_CIRCUIT_TOKEN="}, "--token", "wrongtoken")
	if err == nil {
		t.Fatalf("accepted a rejected token:\n%s", data)
	}
	if !strings.Contains(string(data), "confirm it exists and the token grants access") {
		t.Fatalf("unhelpful failure: %s", data)
	}
	if unauthorized == 0 {
		t.Fatal("the fixture never saw an unauthorized request")
	}

	// A token the curl configuration line cannot carry is refused outright.
	data, err = run(nil, "--token", `a"b`)
	if err == nil || !strings.Contains(string(data), "token contains unexpected characters") {
		t.Fatalf("accepted an unquotable token: %v\n%s", err, data)
	}
}

func itoa(value int) string {
	digits := ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	if digits == "" {
		return "0"
	}
	return digits
}
