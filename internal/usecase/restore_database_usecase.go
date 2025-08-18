package usecase

import (
	"context"
	"ez-snapshot/internal/repository/backup"
	"ez-snapshot/internal/repository/storage"
	"fmt"
	"os"
	"path/filepath"
)

type RestoreDatabaseUseCase struct {
	backup  backup.Repository
	storage storage.Repository
}

func NewRestoreDatabaseUseCase(
	backup backup.Repository,
	storage storage.Repository,
) *RestoreDatabaseUseCase {
	return &RestoreDatabaseUseCase{
		backup:  backup,
		storage: storage,
	}
}

func (uc *RestoreDatabaseUseCase) Execute(ctx context.Context, key string) error {

	fmt.Println("Backup existing database...")
	// Step 1: Dump database
	dumpPath, err := uc.backup.Dump(ctx)
	if err != nil {
		return fmt.Errorf("❌ dump failed: %w", err)
	}
	fmt.Println("✅ Backup created")

	// Step 2: Rename file to "backup_xxxxx.tar.gz"
	current, err := os.Stat(dumpPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(dumpPath)
	newName := fmt.Sprintf("backup_%s", current.Name())
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(dumpPath, newPath); err != nil {
		return fmt.Errorf("❌ failed to rename dump: %w", err)
	}

	fmt.Println("Upload to file storage ...")

	// Step 3: Upload to storage
	f, err := os.Open(newPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := uc.storage.Upload(ctx, filepath.Base(newPath), f); err != nil {
		return fmt.Errorf("❌ backup upload failed: %w", err)
	}

	fmt.Println("✅Backup has been complete")

	// Step 4: Drop all tables
	fmt.Println("Begin downloading snapshot file ...")
	b, err := uc.storage.Download(ctx, key)
	if err != nil {
		return fmt.Errorf("❌ can't download snapshot: %w", err)
	}
	defer b.Close()

	fmt.Println("✅Snapshot has been downloaded")

	fmt.Println("Dropping all tables ...")

	// Step 5: Drop all tables
	if err := uc.backup.DropAllTables(ctx); err != nil {
		return fmt.Errorf("❌ drop all tables failed: %w", err)
	}

	fmt.Println("✅ Table has been dropped")

	fmt.Println("Begin restore process ...")
	if err := uc.backup.Restore(ctx, b); err != nil {
		return fmt.Errorf("❌ restore failed: %w", err)
	}
	fmt.Println("✅ Restore has been complete")

	return nil
}
