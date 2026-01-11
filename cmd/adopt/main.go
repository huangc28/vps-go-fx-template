package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets/*
var assets embed.FS

func main() {
	var (
		dir   string
		force bool
	)

	flag.StringVar(&dir, "dir", ".", "target repo directory to write guidance files into")
	flag.BoolVar(&force, "force", false, "overwrite existing files")
	flag.Parse()

	if err := run(dir, force); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(dir string, force bool) error {
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
		filepath.Join("codex", "skills", "adopt", "SKILL.md"),
		"assets/codex-skills-adopt-SKILL.md",
		force,
	); err != nil {
		return err
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
