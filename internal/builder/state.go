package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// BuildState tracks the build state for incremental builds
type BuildState struct {
	URL        string    `json:"url"`
	ConfigHash string    `json:"config_hash"`
	BuildTime  time.Time `json:"build_time"`

	lib       *Library
	buildRoot string
	statePath string
}

// NewBuildState creates a new build state tracker
func NewBuildState(lib *Library, buildRoot string) *BuildState {
	buildDir := filepath.Join(buildRoot, "build", lib.Name)
	statePath := filepath.Join(buildDir, ".build-state")

	state := &BuildState{
		lib:       lib,
		buildRoot: buildRoot,
		statePath: statePath,
	}

	// Try to load existing state
	state.Load()

	return state
}

// Load loads the build state from disk
func (s *BuildState) Load() error {
	data, err := os.ReadFile(s.statePath)
	if err != nil {
		return err // file doesn't exist or can't be read
	}

	return json.Unmarshal(data, s)
}

// Save saves the current build state to disk
func (s *BuildState) Save() error {
	s.URL = s.lib.URL
	s.ConfigHash = s.lib.ConfigHash()
	s.BuildTime = time.Now()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(s.statePath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.statePath, data, 0644)
}

// CanSkip checks if we can skip building this library
func (s *BuildState) CanSkip(installDir string) bool {
	// No state file = must build
	if s.URL == "" {
		return false
	}

	// URL changed = must rebuild
	if s.URL != s.lib.URL {
		return false
	}

	// Config changed = must rebuild
	if s.ConfigHash != s.lib.ConfigHash() {
		return false
	}

	// Check if outputs exist
	// For header-only libraries (LinkLibs == nil), we can skip if we built before
	if s.lib.LinkLibs == nil {
		return true
	}

	// For libraries with LinkLibs, check that all expected .a files exist
	for _, dir := range []string{"lib"} {
		libDir := filepath.Join(installDir, dir)
		allFound := true
		for _, libName := range s.lib.LinkLibs {
			libPath := filepath.Join(libDir, libName+".a")
			if !fileExists(libPath) {
				allFound = false
				break
			}
		}
		if allFound {
			return true
		}
	}

	return false
}
