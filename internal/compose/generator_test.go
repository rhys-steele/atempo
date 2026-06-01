package compose

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateDockerCompose(t *testing.T) {
	dir := t.TempDir()

	cfg := &Config{
		Name:      "test-project",
		Framework: "laravel",
		Language:  "php",
		Services: map[string]Service{
			"webserver": {
				Type:  "image",
				Image: "nginx:alpine",
				Ports: []string{"8000:80"},
			},
			"mysql": {
				Type:  "image",
				Image: "mysql:8.0",
				Ports: []string{"3306:3306"},
				Environment: map[string]string{
					"MYSQL_DATABASE": "app",
					"MYSQL_ROOT_PASSWORD": "secret",
				},
			},
		},
		Volumes: map[string]Volume{
			"mysql_data": {Driver: "local"},
		},
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "atempo.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	assignments, err := GenerateDockerCompose(dir)
	if err != nil {
		t.Fatalf("GenerateDockerCompose: %v", err)
	}

	// Two services → at least two port assignments
	if len(assignments) < 2 {
		t.Errorf("expected ≥2 port assignments, got %d", len(assignments))
	}

	// docker-compose.yml must exist
	composePath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Fatalf("docker-compose.yml not created: %v", err)
	}

	// atempo.json must be updated with resolved ports (no "{{" template vars)
	updatedCfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("reload atempo.json: %v", err)
	}
	for svcName, svc := range updatedCfg.Services {
		for _, p := range svc.Ports {
			if len(p) == 0 {
				t.Errorf("service %s: empty port mapping", svcName)
			}
		}
	}

	// All assigned external ports must be distinct
	seen := map[int]string{}
	for _, a := range assignments {
		if prev, ok := seen[a.External]; ok {
			t.Errorf("port %d assigned to both %s and %s", a.External, prev, a.Service)
		}
		seen[a.External] = a.Service
	}
}

func TestParsePortMapping(t *testing.T) {
	cases := []struct {
		input    string
		wantHost int
		wantCont int
		wantErr  bool
	}{
		{"8000:80", 8000, 80, false},
		{"3306:3306", 3306, 3306, false},
		{"80", 80, 80, false},
		{"bad:port", 0, 0, true},
		{"", 0, 0, true},
	}
	for _, tc := range cases {
		h, c, err := parsePortMapping(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parsePortMapping(%q): expected error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parsePortMapping(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if h != tc.wantHost || c != tc.wantCont {
			t.Errorf("parsePortMapping(%q): got (%d,%d), want (%d,%d)",
				tc.input, h, c, tc.wantHost, tc.wantCont)
		}
	}
}
