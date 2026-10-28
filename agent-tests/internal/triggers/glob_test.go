package triggers

import "testing"

func TestCompileGlob(t *testing.T) {
	cases := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"README.md", "README.md", true},
		{"README.md", "template/intent/README.md", false},
		{"**/*_test.go", "internal/cli/cli_test.go", true},
		{"**/*_test.go", "cli_test.go", true},
		{"**/*_test.go", "internal/cli/cli.go", false},
		{"scripts/*.sh", "scripts/build.sh", true},
		{"scripts/*.sh", "scripts/release-manifest.txt", false},
		{"internal/cow/**", "internal/cow/copy.go", true},
		{"internal/cow/**", "internal/cli/cli.go", false},
		{"product/launcher/*", "product/launcher/cli.sh", true},
		{"product/launcher/*", "product/launcher/extra/cli.sh", false},
		{"product/skills/cc-cli/scripts/install.*", "product/skills/cc-cli/scripts/install.sh", true},
	}
	for _, tc := range cases {
		re, err := compileGlob(tc.pattern)
		if err != nil {
			t.Fatal(err)
		}
		if re.MatchString(tc.path) != tc.want {
			t.Errorf("%s on %s = %v, want %v", tc.pattern, tc.path, !tc.want, tc.want)
		}
	}
	if _, err := compileGlob("file[abc].go"); err == nil {
		t.Fatal("character class accepted")
	}
}
