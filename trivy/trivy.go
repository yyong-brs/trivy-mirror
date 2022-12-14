package trivy

import (
	"context"
	"fmt"
	"github.com/aquasecurity/trivy-db/pkg/db"
	"github.com/aquasecurity/trivy-db/pkg/metadata"
	dbFile "github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/utils"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"sync"
)

type dbWorker struct {
	dbClient dbFile.Operation
}

func NewDBWorker(dbClient dbFile.Operation) dbWorker {
	return dbWorker{dbClient: dbClient}
}

func (w dbWorker) Update(ctx context.Context, cacheDir string,
	dbUpdateWg, requestWg *sync.WaitGroup) error {

	fmt.Println("Updating DB...")
	if err := w.hotUpdate(ctx, cacheDir, dbUpdateWg, requestWg); err != nil {
		return xerrors.Errorf("failed DB hot update: %w", err)
	}
	return nil
}

func (w dbWorker) hotUpdate(ctx context.Context, cacheDir string, dbUpdateWg, requestWg *sync.WaitGroup) error {
	tmpDir, err := os.MkdirTemp("", "db")
	if err != nil {
		return xerrors.Errorf("failed to create a temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err = w.dbClient.Download(ctx, tmpDir); err != nil {
		return xerrors.Errorf("failed to download vulnerability DB: %w", err)
	}

	fmt.Println("Suspending all requests during DB update")
	dbUpdateWg.Add(1)
	defer dbUpdateWg.Done()

	fmt.Println("Waiting for all requests to be processed before DB update...")
	requestWg.Wait()

	if err = db.Close(); err != nil {
		return xerrors.Errorf("failed to close DB: %w", err)
	}
	dbSrcPath := db.Path(tmpDir)
	dbDestPath := filepath.Join(cacheDir, "trivy.db")
	fmt.Printf("Copy db from %s to %s\n", dbSrcPath, dbDestPath)
	// Copy trivy.db
	if _, err = utils.CopyFile(dbSrcPath, dbDestPath); err != nil {
		return xerrors.Errorf("failed to copy the database file: %w", err)
	}

	metadataSrcPath := metadata.Path(tmpDir)
	metadataDestPath := filepath.Join(cacheDir, "metadata.json")
	fmt.Printf("Copy metadata from %s to %s\n", metadataSrcPath, metadataDestPath)
	// Copy metadata.json
	if _, err = utils.CopyFile(metadataSrcPath, metadataDestPath); err != nil {
		return xerrors.Errorf("failed to copy the metadata file: %w", err)
	}

	fmt.Println("Reopening DB...")
	if err = db.Init(cacheDir); err != nil {
		return xerrors.Errorf("failed to open DB: %w", err)
	}

	return nil
}
