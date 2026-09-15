package installer_test

import (
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

// An organization mirroring CLI releases into its own GitLab project installs
// from a generic package, which is addressed by version and file name and needs
// no release lookup. The fixture answers only requests carrying the mirror's own
// header, which is what proves the credential reaches the selected registry in
// the form that registry expects, and never as a bearer token meant for GitHub.
func TestMirrorInstallReadsGenericPackages(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("install.ps1 uses its own HTTP client; this fixture drives curl")
	}
	root := t.TempDir()
	const token = "glpat_fixture_TOKEN-0123456789"
	const project = "acme/platform/context-circuit"
	version := "2.0.0-mirror"
	archive, sums := packageFixture(t, root, version)
	packageName := filepath.Base(archive)

	// The project path is carried encoded, as one path segment.
	base := "/api/v4/projects/acme%2Fplatform%2Fcontext-circuit/packages/generic/context-circuit-cli/" + version + "/"
	var served, unauthorized int
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if bearer := r.Header.Get("Authorization"); bearer != "" {
			t.Errorf("mirror received a GitHub-style credential: %q", bearer)
		}
		if r.Header.Get("PRIVATE-TOKEN") != token {
			unauthorized++
			http.Error(w, `{"message":"404 Not Found"}`, http.StatusNotFound)
			return
		}
		path := r.URL.EscapedPath()
		if !strings.HasPrefix(path, base) {
			http.Error(w, `{"message":"404 Not Found"}`, http.StatusNotFound)
			return
		}
		switch strings.TrimPrefix(path, base) {
		case packageName:
			served++
			http.ServeFile(w, r, archive)
		case "SHA256SUMS":
			served++
			http.ServeFile(w, r, sums)
		default:
			http.Error(w, `{"message":"404 Not Found"}`, http.StatusNotFound)
		}
	}))
	defer server.Close()

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
		command.Env = append(os.Environ(), append([]string{"CURL_CA_BUNDLE=" + bundle}, environment...)...)
		return command.CombinedOutput()
	}
	mirror := server.URL + "/" + project

	data, err := run(nil, "--token", token, "--gitlab-url", mirror)
	if err != nil {
		t.Fatalf("mirror install: %v\n%s", err, data)
	}
	if served != 2 {
		t.Fatalf("expected the package and checksums to be served, got %d", served)
	}
	installed, err := exec.Command(filepath.Join(binDir, "context-circuit-cli"), "version").CombinedOutput()
	if err != nil || strings.TrimSpace(string(installed)) != version {
		t.Fatalf("installed command: %v %s", err, installed)
	}

	// The same selection and credential are accepted from the environment.
	if data, err := run([]string{
		"CONTEXT_CIRCUIT_TOKEN=" + token,
		"CONTEXT_CIRCUIT_GITLAB_URL=" + mirror,
	}); err != nil {
		t.Fatalf("environment mirror install: %v\n%s", err, data)
	}

	// A trailing slash and a .git suffix name the same project.
	if data, err := run(nil, "--token", token, "--gitlab-url", mirror+".git/"); err != nil {
		t.Fatalf("mirror install with a .git suffix: %v\n%s", err, data)
	}

	// A rejected credential reaches the fixture as a 404, which the installer
	// must report rather than proceed from.
	data, err = run([]string{"CONTEXT_CIRCUIT_TOKEN="}, "--token", "wrongtoken", "--gitlab-url", mirror)
	if err == nil {
		t.Fatalf("accepted a rejected token:\n%s", data)
	}
	if !strings.Contains(string(data), "confirm the version is published there and the token grants access") {
		t.Fatalf("unhelpful failure: %s", data)
	}
	if unauthorized == 0 {
		t.Fatal("the fixture never saw an unauthorized request")
	}

	for _, broken := range []struct{ url, reason string }{
		{"http://" + strings.TrimPrefix(server.URL, "https://") + "/" + project, "must begin with https://"},
		{server.URL, "needs a host and a project path"},
		{server.URL + "/acme/proj ect", "contains unexpected characters"},
	} {
		data, err := run(nil, "--token", token, "--gitlab-url", broken.url)
		if err == nil || !strings.Contains(string(data), broken.reason) {
			t.Fatalf("accepted %q: %v\n%s", broken.url, err, data)
		}
	}
}
