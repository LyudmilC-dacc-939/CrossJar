package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const goCodeTemplate = `package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed {{.JarName}}
var jarData []byte

func main() {
	tmpDir, err := os.MkdirTemp("", "jar2exe")
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	jarPath := filepath.Join(tmpDir, "{{.JarName}}")
	if err := os.WriteFile(jarPath, jarData, 0644); err != nil {
		fmt.Println("Error writing JAR file:", err)
		os.Exit(1)
	}

	cmd := exec.Command("java", "-jar", jarPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running Java:", err)
		os.Exit(1)
	}
}
`

func main() {
	inputJar := flag.String("input", "", "Path to the input JAR file")
	outputExe := flag.String("output", "", "Output executable path")
	targetOS := flag.String("os", runtime.GOOS, "Target OS (windows, linux, darwin)")
	targetArch := flag.String("arch", runtime.GOARCH, "Target architecture (amd64, arm64, etc.)")
	flag.Parse()

	if *inputJar == "" || *outputExe == "" {
		fmt.Println("Usage: jar2exe -input <input.jar> -output <output> [-os windows|linux|darwin] [-arch amd64|arm64]")
		os.Exit(1)
	}

	validOS := map[string]bool{"windows": true, "linux": true, "darwin": true}
	if !validOS[strings.ToLower(*targetOS)] {
		fmt.Printf("Invalid OS: %s. Must be windows, linux, or darwin.\n", *targetOS)
		os.Exit(1)
	}

	// Auto-append .exe for Windows if missing
	if *targetOS == "windows" && filepath.Ext(*outputExe) != ".exe" {
		*outputExe += ".exe"
	}

	tmpDir, err := os.MkdirTemp("", "jar2exe-build")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	jarName := "app.jar"
	destJarPath := filepath.Join(tmpDir, jarName)
	if err := copyFile(*inputJar, destJarPath); err != nil {
		panic(err)
	}

	if err := generateMainGo(tmpDir, jarName); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(*outputExe), 0755); err != nil {
		panic(err)
	}

	// Prepare build command with cross-compilation settings
	cmd := exec.Command("go", "build", "-o", *outputExe, tmpDir)
	
	// Set up environment for cross-compilation
	env := os.Environ()
	env = append(env, "CGO_ENABLED=0") // Disable CGO for static binaries
	env = append(env, "GOOS="+*targetOS)
	env = append(env, "GOARCH="+*targetArch)
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	fmt.Printf("Successfully created %s for %s/%s\n", *outputExe, *targetOS, *targetArch)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func generateMainGo(tmpDir, jarName string) error {
	tmpl, err := template.New("gocode").Parse(goCodeTemplate)
	if err != nil {
		return err
	}

	mainGoPath := filepath.Join(tmpDir, "main.go")
	file, err := os.Create(mainGoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct{ JarName string }{JarName: jarName}
	return tmpl.Execute(file, data)
}
