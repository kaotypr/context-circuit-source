package installer_test

import (
	"encoding/pem"
	"fmt"
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

// A private repository serves release assets only through the API, so the
// installer resolves an asset id and presents a token. The fixture answers only
// authenticated requests, which is what proves the credential is sent.
func TestTokenInstallResolvesPrivateReleaseAssets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("install.ps1 uses its own HTTP client; this fixture drives curl")
	}
	root := t.TempDir()
	const token = "ghp_fixture_TOKEN-0123456789"
	version := "2.0.0-token"
	archive, sums := packageFixture(t, root, version)
	packageName := filepath.Base(archive)

	var served, unauthorized int
	mux := http.NewServeMux()
	var base string
	body := func(w http.ResponseWriter, r *http.Request, path string) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			unauthorized++
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		if r.Header.Get("Accept") != "application/octet-stream" {
			http.Error(w, "wrong accept", http.StatusNotAcceptable)
			return
		}
		served++
		http.ServeFile(w, r, path)
	}
	mux.HandleFunc("/repos/kaotypr/context-circuit-source/releases/assets/1", func(w http.ResponseWriter, r *http.Request) { body(w, r, archive) })
	mux.HandleFunc("/repos/kaotypr/context-circuit-source/releases/assets/2", func(w http.ResponseWriter, r *http.Request) { body(w, r, sums) })
	mux.HandleFunc("/repos/kaotypr/context-circuit-source/releases/tags/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			unauthorized++
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		if !strings.HasSuffix(r.URL.Path, "/cli-v"+version) {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// Shaped like the GitHub payload, including the nested uploader object
		// and a browser_download_url the installer must not mistake for the
		// asset API URL.
		asset := func(id int, name string) string {
			return fmt.Sprintf(`{"url":"%s/repos/kaotypr/context-circuit-source/releases/assets/%d","id":%d,"node_id":"RA_%d","name":"%s","label":null,`+
				`"uploader":{"login":"kaotypr","id":99,"node_id":"U_1","type":"User"},`+
				`"content_type":"application/octet-stream","state":"uploaded","size":1,"download_count":0,`+
				`"browser_download_url":"https://github.com/kaotypr/context-circuit-source/releases/download/cli-v%s/%s"}`,
				base, id, id, id, name, version, name)
		}
		fmt.Fprintf(w, `{"url":"%s/repos/kaotypr/context-circuit-source/releases/7","id":7,"node_id":"RE_7",`+
			`"tag_name":"cli-v%s","name":"Context Circuit CLI %s","draft":false,"prerelease":true,"assets":[%s,%s]}`,
			base, version, version, asset(1, packageName), asset(2, "SHA256SUMS"))
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
	command := exec.Command(filepath.Join(binDir, "context-circuit-cli"), "version")
	installed, err := command.CombinedOutput()
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
