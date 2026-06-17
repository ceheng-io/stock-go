package stock

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGoModulePathsUseCanonicalGitHubSlug(t *testing.T) {
	assertModulePath(t, ".", "github.com/ceheng-io/stock-go")
	assertModulePath(t, "apps/web", "github.com/ceheng-io/stock-go/apps/web")
}

func assertModulePath(t *testing.T, dir string, expected string) {
	t.Helper()
	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -m in %s failed: %v\n%s", dir, err, output)
	}
	if got := strings.TrimSpace(string(output)); got != expected {
		t.Fatalf("module path in %s = %q, want %q", dir, got, expected)
	}
}

func TestGoPackageEnumerationExcludesFrontendDependencies(t *testing.T) {
	output, err := exec.Command("go", "list", "./...").CombinedOutput()
	if err != nil {
		t.Fatalf("go list ./... failed: %v\n%s", err, output)
	}
	for _, pkg := range strings.Fields(string(output)) {
		if strings.Contains(pkg, "/node_modules/") {
			t.Fatalf("go list ./... included frontend dependency package %q", pkg)
		}
	}
}

func TestFrontendDirectoryHasGoModuleBoundary(t *testing.T) {
	if _, err := os.Stat("apps/web/go.mod"); err != nil {
		t.Fatalf("apps/web/go.mod is missing: %v", err)
	}
}

func TestRecommendedFixtureDirectoryExists(t *testing.T) {
	info, err := os.Stat("testdata")
	if err != nil {
		t.Fatalf("testdata directory is missing: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("testdata exists but is not a directory")
	}
	if _, err := os.Stat("testdata/README.md"); err != nil {
		t.Fatalf("testdata README is missing: %v", err)
	}
}
