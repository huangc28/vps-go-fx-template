package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed assets
var assets embed.FS

func main() {
	var (
		dir      string
		force    bool
		scaffold bool
	)

	flag.StringVar(&dir, "dir", ".", "target repo directory to write guidance files into")
	flag.BoolVar(&force, "force", false, "overwrite existing files")
	flag.BoolVar(&scaffold, "scaffold", false, "write template code scaffold into the target repo (non-destructive unless --force)")
	flag.Parse()

	if err := run(dir, force, scaffold); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(dir string, force bool, scaffold bool) error {
	targetDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve --dir: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(targetDir, "architecture"), 0o755); err != nil {
		return fmt.Errorf("create architecture dir: %w", err)
	}

	if err := writeAsset(targetDir, "AGENTS.md", "assets/AGENTS.md", force); err != nil {
		return err
	}

	if err := writeAsset(
		targetDir,
		filepath.Join("architecture", "vps-go-fx-template.md"),
		"assets/architecture-vps-go-fx-template.md",
		force,
	); err != nil {
		return err
	}

	if err := writeAsset(
		targetDir,
		filepath.Join("architecture", "config-go.md"),
		"assets/architecture-config-go.md",
		force,
	); err != nil {
		return err
	}

	if err := writeAsset(
		targetDir,
		filepath.Join("architecture", "db-go.md"),
		"assets/architecture-db-go.md",
		force,
	); err != nil {
		return err
	}

	if err := writeAsset(
		targetDir,
		filepath.Join("architecture", "cache-redis-go.md"),
		"assets/architecture-cache-redis-go.md",
		force,
	); err != nil {
		return err
	}

	if err := writeAsset(
		targetDir,
		filepath.Join("codex", "skills", "adopt", "SKILL.md"),
		"assets/codex-skills-adopt-SKILL.md",
		force,
	); err != nil {
		return err
	}

	if scaffold {
		modulePath, err := readModulePath(filepath.Join(targetDir, "go.mod"))
		if err != nil {
			return err
		}

		if err := writeTemplateTree(targetDir, modulePath, force); err != nil {
			return err
		}
	}

	return nil
}

func writeAsset(targetDir, relTargetPath, assetPath string, force bool) error {
	b, err := fs.ReadFile(assets, assetPath)
	if err != nil {
		return fmt.Errorf("read embedded asset %q: %w", assetPath, err)
	}

	dst := filepath.Join(targetDir, relTargetPath)
	if !force {
		if _, statErr := os.Stat(dst); statErr == nil {
			return fmt.Errorf("refusing to overwrite existing file %q (use --force)", dst)
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return fmt.Errorf("stat %q: %w", dst, statErr)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create parent dir for %q: %w", dst, err)
	}

	if err := os.WriteFile(dst, b, 0o644); err != nil {
		return fmt.Errorf("write %q: %w", dst, err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "wrote %s\n", dst)
	return nil
}

func readModulePath(goModPath string) (string, error) {
	b, err := os.ReadFile(goModPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("missing %q; create a Go module first (e.g. `go mod init <module>`) before running with --scaffold", goModPath)
		}
		return "", fmt.Errorf("read %q: %w", goModPath, err)
	}

	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if modulePath == "" {
				return "", fmt.Errorf("invalid module line in %q", goModPath)
			}
			return modulePath, nil
		}
	}

	return "", fmt.Errorf("could not find module path in %q", goModPath)
}

func writeTemplateTree(targetDir, modulePath string, force bool) error {
	const templateRoot = "assets/template"

	return fs.WalkDir(assets, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, ok := strings.CutPrefix(path, templateRoot+string(filepath.Separator))
		if !ok {
			rel, ok = strings.CutPrefix(path, templateRoot+"/")
			if !ok {
				return fmt.Errorf("unexpected template path %q", path)
			}
		}

		b, err := fs.ReadFile(assets, path)
		if err != nil {
			return fmt.Errorf("read embedded template %q: %w", path, err)
		}

		b = []byte(strings.ReplaceAll(string(b), "{{MODULE}}", modulePath))

		if strings.HasSuffix(rel, ".tmpl") {
			rel = strings.TrimSuffix(rel, ".tmpl")
		}

		dst := filepath.Join(targetDir, rel)
		if !force {
			if _, statErr := os.Stat(dst); statErr == nil {
				return fmt.Errorf("refusing to overwrite existing file %q (use --force)", dst)
			} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
				return fmt.Errorf("stat %q: %w", dst, statErr)
			}
		}

		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("create parent dir for %q: %w", dst, err)
		}

		if err := os.WriteFile(dst, b, 0o644); err != nil {
			return fmt.Errorf("write %q: %w", dst, err)
		}

		_, _ = fmt.Fprintf(os.Stdout, "wrote %s\n", dst)
		return nil
	})
}
