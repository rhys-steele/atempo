package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestExpandTokens(t *testing.T) {
	cases := []struct {
		in, project, version, want string
	}{
		{"Hello {{project}}", "my-app", "12", "Hello my-app"},
		{"v{{version}}", "x", "11.3", "v11.3"},
		{"major: {{majorVersion}}", "x", "11.3", "major: 11"},
		{"major: {{majorVersion}}", "x", "", "major: latest"},
		{"{{project}} v{{majorVersion}}", "api", "3.2", "api v3"},
	}

	for _, tc := range cases {
		got := expandTokens(tc.in, tc.project, tc.version)
		if got != tc.want {
			t.Errorf("expandTokens(%q, %q, %q) = %q, want %q",
				tc.in, tc.project, tc.version, got, tc.want)
		}
	}
}

func TestMajorVersion(t *testing.T) {
	cases := []struct{ in, want string }{
		{"12.3.1", "12"},
		{"3", "3"},
		{"", "latest"},
	}
	for _, tc := range cases {
		if got := majorVersion(tc.in); got != tc.want {
			t.Errorf("majorVersion(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseFrameworkArgInScaffold(t *testing.T) {
	base := []string{"composer", "create-project", "laravel/laravel", "src"}

	cp := func() []string {
		c := make([]string, len(base))
		copy(c, base)
		return c
	}

	got := applyVersionConstraints(cp(), "laravel", "12")
	if got[2] != "laravel/laravel:^12.0" {
		t.Errorf("expected version constraint, got %q", got[2])
	}

	got2 := applyVersionConstraints(cp(), "laravel", "")
	if got2[2] != "laravel/laravel" {
		t.Errorf("empty version should not add constraint, got %q", got2[2])
	}
}

func TestLoadMetadataUnknownFramework(t *testing.T) {
	fsys := fstest.MapFS{}
	_, err := loadMetadata(fsys, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown framework")
	}
}

func TestLoadMetadataValid(t *testing.T) {
	fsys := fstest.MapFS{
		"frameworks/bun/atempo.json": &fstest.MapFile{
			Data: []byte(`{"name":"{{project}}","framework":"bun","language":"typescript","installer":{"type":"template","command":[]}}`),
		},
	}
	meta, err := loadMetadata(fsys, "bun")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Framework != "bun" {
		t.Errorf("got framework %q, want bun", meta.Framework)
	}
}

func TestCopyTemplateFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"frameworks/fastapi/atempo.json": &fstest.MapFile{
			Data: []byte(`{"name":"{{project}}"}`),
		},
		"frameworks/fastapi/README.md": &fstest.MapFile{
			Data: []byte("# {{project}}\n"),
		},
		"frameworks/fastapi/src/main.py": &fstest.MapFile{
			Data: []byte(`print("{{project}}")`),
		},
	}

	dst := t.TempDir()
	if err := copyTemplateFiles(fsys, "fastapi", dst, "my-api", "0.1"); err != nil {
		t.Fatal(err)
	}

	// atempo.json is skipped by copyTemplateFiles (handled by registerAndGenerate).
	if _, err := os.Stat(filepath.Join(dst, "atempo.json")); err == nil {
		t.Error("atempo.json should not be copied by copyTemplateFiles")
	}

	data, err := os.ReadFile(filepath.Join(dst, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# my-api\n" {
		t.Errorf("token not expanded in README.md: %q", string(data))
	}

	src, err := os.ReadFile(filepath.Join(dst, "src", "main.py"))
	if err != nil {
		t.Fatal(err)
	}
	if string(src) != `print("my-api")` {
		t.Errorf("token not expanded in src/main.py: %q", string(src))
	}
}
