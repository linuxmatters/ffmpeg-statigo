package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ulikunitz/xz"
)

func writeLines(dst string, lines []string) {
	f, err := os.Create(dst)
	if err != nil {
		log.Panicln(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(f)

	for _, line := range lines {
		if _, err := f.WriteString(fmt.Sprintf("%v\n", line)); err != nil {
			log.Panicln(err)
		}
	}
}

func copyFile(dst string, src string) {

	srcF, err := os.Open(src)
	if err != nil {
		log.Panicln(err)
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
		log.Panicln(err)
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		log.Panicln(err)
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		log.Panicln(err)
	}
}

func run(prefix string, cmd *exec.Cmd) {
	wg := &sync.WaitGroup{}

	wg.Add(2)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicln(err)
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Panicln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Panicln(err)
	}

	scanner := bufio.NewScanner(outPipe)
	//scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			log.Println(prefix, scanner.Text())
		}

		wg.Done()
	}()

	scanner2 := bufio.NewScanner(errPipe)
	//scanner.Split(bufio.ScanLines)
	go func() {
		for scanner2.Scan() {
			log.Println(prefix, scanner2.Text())
		}

		wg.Done()
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Panicln(err)
	}

}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	log.Panicln(err)
	return false
}

func download(url string, path string) {
	const maxRetries = 3
	const retryDelay = 5 * time.Second

	// Extract filename from path for cleaner logging
	filename := filepath.Base(path)
	log.Printf("[DOWNLOAD] Starting: %s", filename)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := downloadWithResume(url, path, filename)
		if err == nil {
			log.Printf("[DOWNLOAD] Complete: %s", filename)
			return
		}

		if attempt < maxRetries {
			log.Printf("[DOWNLOAD] Attempt %d/%d failed: %v. Retrying in %v...", attempt, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Panicln("[DOWNLOAD] Failed after", maxRetries, "attempts:", filename, err)
		}
	}
}

func downloadWithResume(url string, path string, filename string) error {
	// Check if file exists and get its size for resume support
	var existingSize int64
	if fi, err := os.Stat(path); err == nil {
		existingSize = fi.Size()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add Range header for resume if file exists
	if existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
	}

	client := &http.Client{
		Timeout: 10 * time.Minute, // 10 minute timeout for large files
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle resume: 206 Partial Content means resume succeeded
	// 200 OK means server doesn't support resume, start from scratch
	var f *os.File
	var totalSize int64

	if resp.StatusCode == http.StatusPartialContent {
		// Resume: append to existing file
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		totalSize = existingSize + resp.ContentLength
		log.Printf("[DOWNLOAD] Resuming: %s from %.2f MB of %.2f MB",
			filename,
			float64(existingSize)/(1024*1024),
			float64(totalSize)/(1024*1024))
	} else if resp.StatusCode == http.StatusOK {
		// Start fresh
		f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		totalSize = resp.ContentLength
		existingSize = 0
		log.Printf("[DOWNLOAD] Size: %.2f MB", float64(totalSize)/(1024*1024))
	} else {
		return fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("[DOWNLOAD] Error closing file:", err)
		}
	}(f)

	// Simple progress reporting without animations
	// Report every 10 MB downloaded
	const reportInterval = 10 * 1024 * 1024
	var bytesDownloaded int64
	var lastReport int64

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := f.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("write failed: %w", writeErr)
			}
			bytesDownloaded += int64(n)

			// Report progress every 10MB
			if bytesDownloaded-lastReport >= reportInterval {
				currentMB := float64(existingSize+bytesDownloaded) / (1024 * 1024)
				totalMB := float64(totalSize) / (1024 * 1024)
				percent := float64(existingSize+bytesDownloaded) / float64(totalSize) * 100
				log.Printf("[DOWNLOAD] Progress: %s - %.2f/%.2f MB (%.0f%%)",
					filename, currentMB, totalMB, percent)
				lastReport = bytesDownloaded
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read failed: %w", err)
		}
	}

	return nil
}

func untar(src string, dest string, prefix string) {
	filename := filepath.Base(src)
	log.Printf("[EXTRACT] Starting: %s", filename)

	os.RemoveAll(dest)

	if err := os.MkdirAll(dest, 0755); err != nil {
		log.Panicln(err)
	}

	gzipStream, err := os.Open(src)
	if err != nil {
		log.Panicln(err)
	}

	defer gzipStream.Close()

	var uncompressedStream io.Reader

	if strings.HasSuffix(src, ".xz") {
		uncompressedStream, err = xz.NewReader(gzipStream)
		if err != nil {
			log.Fatal("ExtractTarXz: NewReader failed", err)
		}
	} else if strings.HasSuffix(src, ".bz2") {
		uncompressedStream = bzip2.NewReader(gzipStream)
	} else {
		uncompressedStream, err = gzip.NewReader(gzipStream)
		if err != nil {
			log.Fatal("ExtractTarGz: NewReader failed", err)
		}
	}

	tarReader := tar.NewReader(uncompressedStream)

	var fileCount int

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		header.Name = strings.TrimPrefix(header.Name, prefix)

		if header.Name == "" {
			continue
		}

		path := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()
			fileCount++

		case tar.TypeXGlobalHeader:
			log.Println("Ignoring TypeXGlobalHeader")
			// Ignore?
		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}

	}

	log.Printf("[EXTRACT] Complete: %s (%d files)", filename, fileCount)
}

func unzip(src string, dest string) {
	filename := filepath.Base(src)
	log.Printf("[EXTRACT] Starting: %s", filename)

	os.RemoveAll(dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Panicln(err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		log.Panicln(err)
	}

	prefixes := make(map[string]struct{})

	for _, f := range r.File {
		parts := strings.SplitN(f.Name, string(os.PathSeparator), 2)
		prefixes[parts[0]] = struct{}{}
	}

	if len(prefixes) != 1 {
		log.Panicln("Unexpected prefixes", prefixes)
	}

	var prefix string
	for s, _ := range prefixes {
		prefix = fmt.Sprintf("%v%v", s, string(os.PathSeparator))
	}

	var fileCount int
	for _, f := range r.File {
		f.Name = strings.TrimPrefix(f.Name, prefix)

		if f.Name == "" {
			continue
		}

		if err := extractAndWriteFile(f, dest); err != nil {
			log.Panicln(err)
		}
		if !f.FileInfo().IsDir() {
			fileCount++
		}
	}

	log.Printf("[EXTRACT] Complete: %s (%d files)", filename, fileCount)
}

func extractAndWriteFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	path := filepath.Join(dest, f.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(path, f.Mode()); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Panicln(err)
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmd(name string, dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(name)
	cmd.Dir = dir
	cmd.Env = os.Environ()

	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("CFLAGS=-I%v", incDir),
		fmt.Sprintf("CPPFLAGS=-I%v", incDir),
		fmt.Sprintf("CXXFLAGS=-I%v", incDir),
		fmt.Sprintf("LDFLAGS=-L%v", libDir),
		fmt.Sprintf("PKG_CONFIG_PATH=%v/pkgconfig", libDir),
	)

	cmd.Args = append(cmd.Args, args...)

	return cmd
}

func modify(path string, mod func([]byte) []byte) {
	s, err := os.Stat(path)
	if err != nil {
		log.Panicln(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}

	content = mod(content)

	if err := os.WriteFile(path, content, s.Mode()); err != nil {
		log.Panicln(err)
	}
}

// touchAutomakeFiles updates timestamps on automake-generated files to prevent
// Make from trying to regenerate them. This works around automake version mismatches
// where the tarball was configured with a different automake version than installed.
func touchAutomakeFiles(srcPath string) {
	now := time.Now()

	// Touch all automake-generated files to be newer than their sources
	// This prevents Make from invoking automake to regenerate them
	files := []string{
		"aclocal.m4",
		"Makefile.in",
		"config.h.in",
		"configure",
	}

	for _, file := range files {
		fullPath := filepath.Join(srcPath, file)
		if _, err := os.Stat(fullPath); err == nil {
			// File exists, update its timestamp
			if err := os.Chtimes(fullPath, now, now); err != nil {
				log.Printf("Warning: failed to touch %s: %v", file, err)
			}
		}
	}

	// Also touch any Makefile.in files in subdirectories
	filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "Makefile.in" {
			os.Chtimes(path, now, now)
		}
		return nil
	})
}
