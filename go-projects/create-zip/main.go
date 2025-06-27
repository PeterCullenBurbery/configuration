package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/PeterCullenBurbery/go-functions"
	yekazip "github.com/yeka/zip"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	yamlPath := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\install.yaml`
	zipSource := `C:\Users\Administrator\Desktop\GitHub-repositories\configuration\host\nirsoft_package_enc_1.30.19.zip`

	// Load YAML
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML: %v", err)
	}
	var root map[string]interface{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		log.Fatalf("❌ Failed to parse YAML: %v", err)
	}

	// Parse relevant paths
	install := gofunctions.GetCaseInsensitiveMap(root, "install")
	downloads := gofunctions.GetCaseInsensitiveMap(install, "downloads")
	globalDownloadDir := strings.TrimSpace(gofunctions.GetCaseInsensitiveString(downloads, "global download directory"))
	perAppDownloads := gofunctions.GetCaseInsensitiveMap(downloads, "per app download directories")
	subDownload := strings.TrimSpace(gofunctions.GetCaseInsensitiveString(perAppDownloads, "Nirsoft"))
	nirsoftBase := filepath.Join(globalDownloadDir, subDownload)

	// Defender Exclusion
	if err := addDefenderExclusion(nirsoftBase); err != nil {
		log.Printf("⚠️ Failed to add Defender exclusion: %v", err)
	} else {
		fmt.Println("✅ Added Defender exclusion for:", nirsoftBase)
	}

	// T1: Copy ZIP
	rawT1, _ := gofunctions.DateTimeStamp()
	T1 := gofunctions.SafeTimeStamp(rawT1, 1)
	copyDir := filepath.Join(nirsoftBase, T1)
	_ = os.MkdirAll(copyDir, os.ModePerm)
	zipCopyPath := filepath.Join(copyDir, filepath.Base(zipSource))
	_ = copyFile(zipSource, zipCopyPath)
	fmt.Println("✅ Original ZIP copied to:", zipCopyPath)

	// T2: Extract ZIP
	rawT2, _ := gofunctions.DateTimeStamp()
	T2 := gofunctions.SafeTimeStamp(rawT2, 1)
	extractDir := filepath.Join(copyDir, T2)
	if err := unzip(zipCopyPath, extractDir); err != nil {
		log.Fatalf("❌ Failed to extract ZIP: %v", err)
	}
	fmt.Println("✅ Original ZIP extracted to:", extractDir)

	// T3: Create protected ZIP
	rawT3, _ := gofunctions.DateTimeStamp()
	T3 := gofunctions.SafeTimeStamp(rawT3, 1)
	secureDir := filepath.Join(nirsoftBase, T3)
	_ = os.MkdirAll(secureDir, os.ModePerm)
	secureZipPath := filepath.Join(secureDir, "nirsoft_package_enc_1.30.19.zip")
	if err := zipWithPassword(extractDir, secureZipPath, "nirsoft9876$"); err != nil {
		log.Fatalf("❌ Failed to create password-protected ZIP: %v", err)
	}
	fmt.Println("✅ Password-protected ZIP created:", secureZipPath)

	// T4: Extract protected ZIP
	rawT4, _ := gofunctions.DateTimeStamp()
	T4 := gofunctions.SafeTimeStamp(rawT4, 1)
	finalExtract := filepath.Join(secureDir, T4)
	if err := unzipWithPassword(secureZipPath, finalExtract, "nirsoft9876$"); err != nil {
		log.Fatalf("❌ Failed to extract protected ZIP: %v", err)
	}
	fmt.Println("✅ Password-protected ZIP extracted to:", finalExtract)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func unzip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		fullPath := filepath.Join(destDir, file.Name)
		if !strings.HasPrefix(fullPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fullPath)
		}
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(fullPath, os.ModePerm)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		srcFile, _ := file.Open()
		defer srcFile.Close()
		destFile, _ := os.Create(fullPath)
		defer destFile.Close()
		_, _ = io.Copy(destFile, srcFile)
	}
	return nil
}

func zipWithPassword(sourceDir, outputZip, password string) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := yekazip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(sourceDir, path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, _ := yekazip.FileInfoHeader(info)
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.Encrypt(header.Name, password, yekazip.AES256Encryption)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		return err
	})
}

func unzipWithPassword(zipPath, destDir, password string) error {
	r, err := yekazip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		f.SetPassword(password)
		outPath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(outPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", outPath)
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(outPath, os.ModePerm)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		dstFile, _ := os.Create(outPath)
		defer dstFile.Close()
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func addDefenderExclusion(path string) error {
	checkCmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass",
		"-Command", fmt.Sprintf(`(Get-MpPreference).ExclusionPath -contains "%s"`, path))
	var out bytes.Buffer
	checkCmd.Stdout = &out
	checkCmd.Stderr = &out
	if err := checkCmd.Run(); err == nil && strings.Contains(out.String(), "True") {
		log.Printf("ℹ️ Defender exclusion already exists:\n↳ %s", path)
		return nil
	}

	psCommand := fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, path)
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
