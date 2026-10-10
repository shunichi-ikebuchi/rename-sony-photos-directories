package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
)

// cmd/rspd は main package なので、この file が無いと overlay の coverage gate が
// 「最初の package (cmd/rspd) が 0.0%」で push を止める (2026-10-10 実測)。
// cobra の木と、引数から実処理へ渡す部分をここで実行する。

// グローバルな flag 変数は同じ process のテスト間で持ち越すので、各テストの前に戻す。
func resetFlags() {
	configPath = ""
	dryRun = false
	targetPath = ""
}

func TestLoadConfigurationFromExplicitPath(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := config.Save(config.Default(), path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	configPath = path
	cfg, err := loadConfiguration()
	if err != nil {
		t.Fatalf("loadConfiguration: %v", err)
	}
	if cfg == nil || cfg.TargetPath == "" {
		t.Fatalf("cfg = %+v, want non-empty TargetPath", cfg)
	}
}

func TestLoadConfigurationMissingExplicitPath(t *testing.T) {
	resetFlags()
	configPath = filepath.Join(t.TempDir(), "nope.yaml")
	if _, err := loadConfiguration(); err == nil {
		t.Fatal("loadConfiguration: want error for a missing explicit path")
	}
}

func TestLoadConfigurationFallsBackToDefault(t *testing.T) {
	resetFlags()
	// configPath が空なら LoadOrDefault を通る (config.yaml が無ければ既定値)。
	if _, err := loadConfiguration(); err != nil {
		t.Fatalf("loadConfiguration: %v", err)
	}
}

func TestHandleCreateConfig(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	t.Chdir(dir)

	if err := handleCreateConfig(); err != nil {
		t.Fatalf("handleCreateConfig: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err != nil {
		t.Fatalf("config.yaml が出来ていない: %v", err)
	}
}

func TestRunRenameOnlyDryRunDoesNotTouchFilesystem(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	before, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	cfg := &config.Config{TargetPath: dir}
	targetPath = dir
	dryRun = true
	if err := runRenameOnly(cfg); err != nil {
		t.Fatalf("runRenameOnly (dry-run): %v", err)
	}

	after, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(before) != len(after) {
		t.Fatalf("dry-run が dir を変えた: %d -> %d entries", len(before), len(after))
	}
}

func TestRunRenameOnlyUsesTargetPathFlagOverConfig(t *testing.T) {
	resetFlags()
	flagDir := t.TempDir()
	cfgDir := t.TempDir()

	cfg := &config.Config{TargetPath: cfgDir}
	targetPath = flagDir
	dryRun = true
	// flag が優先される。dry-run なので実際には何もしない。
	if err := runRenameOnly(cfg); err != nil {
		t.Fatalf("runRenameOnly: %v", err)
	}
}

func TestRunRenameOnlyDefaultsToCurrentDirectory(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	t.Chdir(dir)

	dryRun = true
	if err := runRenameOnly(&config.Config{}); err != nil {
		t.Fatalf("runRenameOnly (path 既定): %v", err)
	}
}

func TestRunRenameOnlyMissingPathIsError(t *testing.T) {
	resetFlags()
	cfg := &config.Config{TargetPath: filepath.Join(t.TempDir(), "does-not-exist")}
	if err := runRenameOnly(cfg); err == nil {
		t.Fatal("runRenameOnly: 存在しない path では error を返すはず")
	}
}

func TestRootCmdDryRun(t *testing.T) {
	resetFlags()
	dir := t.TempDir()

	rootCmd.SetArgs([]string{"--dry-run", "--path", dir})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute: %v", err)
	}
	_ = dir
}

func TestCreateConfigSubcommand(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	t.Chdir(dir)

	rootCmd.SetArgs([]string{"create-config"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("create-config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err != nil {
		t.Fatalf("config.yaml が出来ていない: %v", err)
	}
}

// 存在しない mount を指す config で動かす。実処理は途中で error を返すが、
// 「引数を受けて workflow に渡す」までの statement は通る。
func TestWorkflowSubcommandPropagatesError(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	cfg := &config.Config{
		TargetPath:      filepath.Join(dir, "no-such-source"),
		BackupPath:      filepath.Join(dir, "no-such-backup"),
		DestinationPath: filepath.Join(dir, "no-such-dest"),
		TmpDir:          filepath.Join(dir, "tmp"),
	}
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rootCmd.SetArgs([]string{"--config", path, "workflow"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("workflow: source が無いので error を返すはず")
	}
}

func TestBackupCleanupSubcommandPropagatesError(t *testing.T) {
	resetFlags()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	cfg := &config.Config{
		BackupPath:      filepath.Join(dir, "no-such-backup"),
		DestinationPath: filepath.Join(dir, "no-such-dest"),
		TmpDir:          filepath.Join(dir, "tmp"),
	}
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rootCmd.SetArgs([]string{"--config", path, "backup-cleanup", "--dry-run"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("backup-cleanup: 対象が無いので error を返すはず")
	}
}

func TestRunReturnsZeroOnSuccess(t *testing.T) {
	resetFlags()
	dir := t.TempDir()

	rootCmd.SetArgs([]string{"--dry-run", "--path", dir})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	if got := run(); got != 0 {
		t.Fatalf("run() = %d, want 0", got)
	}
}

func TestRunReturnsOneOnError(t *testing.T) {
	resetFlags()
	configPath = filepath.Join(t.TempDir(), "nope.yaml")

	rootCmd.SetArgs(nil)
	if got := run(); got != 1 {
		t.Fatalf("run() = %d, want 1 (存在しない --config)", got)
	}
}

func TestCommandTreeIsWired(t *testing.T) {
	resetFlags()
	want := map[string]bool{"workflow": false, "backup-cleanup": false, "create-config": false}
	for _, c := range rootCmd.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("subcommand %q が root に付いていない", name)
		}
	}

	if rootCmd.PersistentFlags().Lookup("config") == nil {
		t.Error("--config が無い")
	}
	if rootCmd.PersistentFlags().Lookup("dry-run") == nil {
		t.Error("--dry-run が無い")
	}
	if rootCmd.Flags().Lookup("path") == nil {
		t.Error("--path が無い")
	}
	if rootCmd.Use != "rspd" {
		t.Errorf("Use = %q, want %q (改名の回帰)", rootCmd.Use, "rspd")
	}
}
