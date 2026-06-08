package repositories

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"atempo/internal/registry"
)

// newTestRegistry returns a FileRegistryRepository backed by a temp file.
func newTestRegistry(t *testing.T) *FileRegistryRepository {
	t.Helper()
	dir := t.TempDir()
	cache := NewMemoryCacheRepository()
	r := &FileRegistryRepository{
		path:  filepath.Join(dir, "registry.json"),
		cache: cache,
	}
	return r
}

func TestRegistryEmpty(t *testing.T) {
	r := newTestRegistry(t)
	projects, err := r.ListProjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected empty registry, got %d projects", len(projects))
	}
}

func TestRegistryAddAndFind(t *testing.T) {
	r := newTestRegistry(t)

	proj := &registry.Project{
		Name:      "my-app",
		Path:      "/tmp/my-app",
		Framework: "laravel",
		Version:   "12",
	}
	if err := r.AddProject(proj); err != nil {
		t.Fatal(err)
	}

	found, err := r.FindProject("my-app")
	if err != nil {
		t.Fatal(err)
	}
	if found.Path != "/tmp/my-app" {
		t.Errorf("got path %q, want %q", found.Path, "/tmp/my-app")
	}
	if found.Framework != "laravel" {
		t.Errorf("got framework %q, want %q", found.Framework, "laravel")
	}
}

func TestRegistryAddUpdatesExisting(t *testing.T) {
	r := newTestRegistry(t)

	r.AddProject(&registry.Project{Name: "app", Path: "/old", Framework: "bun"})
	r.AddProject(&registry.Project{Name: "app", Path: "/new", Framework: "bun"})

	projects, _ := r.ListProjects()
	if len(projects) != 1 {
		t.Fatalf("expected 1 project after update, got %d", len(projects))
	}
	if projects[0].Path != "/new" {
		t.Errorf("expected updated path /new, got %q", projects[0].Path)
	}
}

func TestRegistryRemove(t *testing.T) {
	r := newTestRegistry(t)
	r.AddProject(&registry.Project{Name: "app", Path: "/tmp/app", Framework: "fastapi"})

	if err := r.RemoveProject("app"); err != nil {
		t.Fatal(err)
	}

	_, err := r.FindProject("app")
	if err == nil {
		t.Fatal("expected error finding removed project")
	}
}

func TestRegistryRemoveNotFound(t *testing.T) {
	r := newTestRegistry(t)
	if err := r.RemoveProject("nonexistent"); err == nil {
		t.Fatal("expected error removing nonexistent project")
	}
}

func TestRegistryPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "registry.json")

	r1 := &FileRegistryRepository{path: path, cache: NewMemoryCacheRepository()}
	r1.AddProject(&registry.Project{Name: "app", Path: "/tmp/app", Framework: "python"})

	// New instance, same file — no cache warm.
	r2 := &FileRegistryRepository{path: path, cache: NewMemoryCacheRepository()}
	found, err := r2.FindProject("app")
	if err != nil {
		t.Fatalf("project not persisted: %v", err)
	}
	if found.Framework != "python" {
		t.Errorf("got framework %q, want python", found.Framework)
	}
}

func TestRegistryCacheServedOnSecondLoad(t *testing.T) {
	r := newTestRegistry(t)
	r.AddProject(&registry.Project{Name: "app", Path: "/tmp/app", Framework: "bun"})

	// Delete the backing file — subsequent reads must come from cache.
	os.Remove(r.path)

	_, err := r.FindProject("app")
	if err != nil {
		t.Fatal("expected cache hit after file deleted, got:", err)
	}
}

func TestRegistryCreatedAtSet(t *testing.T) {
	r := newTestRegistry(t)
	before := time.Now()
	r.AddProject(&registry.Project{Name: "app", Path: "/tmp/app", Framework: "bun"})

	found, _ := r.FindProject("app")
	if found.CreatedAt.Before(before) {
		t.Error("CreatedAt should be set on add")
	}
}
