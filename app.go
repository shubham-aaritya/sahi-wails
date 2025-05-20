package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) ProcessFiles(nsePath, bsePath, cdslPath string) (string, error) {
	nseUCCs, err := loadNSEUCCs(nsePath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error loading NSE file: %v", err)
		return "", err
	}

	bseUCCs, err := loadBSEUCCs(bsePath)
	if err != nil {
		return "", err
	}

	records, err := processCDSL(cdslPath, nseUCCs, bseUCCs)
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(os.TempDir(), "boid_cm_filtered_output.csv")
	if err := writeCSV(records, outputPath); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (a *App) SaveFile(path string, data []byte) error {
	fmt.Println("Saving file to:", path)
	return os.WriteFile(path, data, 0644)
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
