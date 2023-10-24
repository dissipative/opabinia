package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dissipative/opabinia/internal/infra/config"
	"github.com/dissipative/opabinia/internal/infra/fs"
)

func TestCommander_DoCompile(t *testing.T) {
	c := &App{config: config.NewConfig(), logger: nil} // Replace with your actual Config type and logger if applicable

	projectName := "CompileTestProject"
	defer os.RemoveAll(projectName)

	os.Args = []string{Name, "-i", projectName}
	err := c.DoInit()
	if err != nil {
		t.Fatalf("DoInit() failed: %v", err)
	}

	// go to project dir to compile
	err = os.Chdir(projectName)
	if err != nil {
		t.Fatalf("Change dir failed: %v", err)
	}

	err = c.Compile()
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	// return back
	err = os.Chdir("../")
	if err != nil {
		t.Fatalf("Change dir failed: %v", err)
	}

	for _, file := range projectFiles {
		filename := filepath.Join(projectName, fs.Dist, file.dir, file.name)

		if strings.HasSuffix(file.name, ".md") {
			filename = strings.ReplaceAll(filename, ".md", ".html")
		}
	}
}
