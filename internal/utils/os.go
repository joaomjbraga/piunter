package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joaomjbraga/piunter/pkg/types"
)

func readOutput(stdout io.ReadCloser, stderr io.ReadCloser) (string, string, error) {
	var outBuf, errBuf bytes.Buffer
	var outErr, errErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, outErr = io.Copy(&outBuf, stdout)
	}()
	go func() {
		defer wg.Done()
		_, errErr = io.Copy(&errBuf, stderr)
	}()
	wg.Wait()

	var resultErr error
	if outErr != nil {
		resultErr = outErr
	} else {
		resultErr = errErr
	}

	return outBuf.String(), errBuf.String(), resultErr
}

func Exec(command string, args ...string) types.CommandResult {
	cmd := exec.Command(command, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return types.CommandResult{
			Success: false,
			Stderr:  err.Error(),
			Code:    1,
		}
	}

	output, stderrOutput, _ := readOutput(stdout, stderr)
	cmd.Wait()

	return types.CommandResult{
		Success: cmd.ProcessState.ExitCode() == 0,
		Stdout:  output,
		Stderr:  stderrOutput,
		Code:    cmd.ProcessState.ExitCode(),
	}
}

func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

func HasSudoPassword() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}

func RequestSudo() bool {
	cmd := exec.Command("sudo", "true")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err == nil
}

func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		Debug(fmt.Sprintf("falha ao obter home directory: %s", err))
		if IsRoot() {
			return "/root"
		}
		return "/tmp"
	}
	return home
}

func GetCacheDir() string {
	home := GetHomeDir()
	return fmt.Sprintf("%s/.cache", home)
}

func GetDistroInfo() types.DistroInfo {
	distro := types.DistroInfo{
		ID:      "unknown",
		Name:    "Linux",
		Version: "",
	}

	osRelease, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return distro
	}

	lines := strings.Split(string(osRelease), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distro.ID = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		}
		if strings.HasPrefix(line, "NAME=") {
			distro.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			distro.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}

	distro.PackageManager = detectPackageManager(distro.ID)
	return distro
}

func detectPackageManager(distroID string) types.PackageManager {
	if IsCommandAvailable("apt") {
		return types.PackageManagerApt
	}
	if IsCommandAvailable("pacman") {
		return types.PackageManagerPacman
	}
	if IsCommandAvailable("dnf") {
		return types.PackageManagerDnf
	}
	return types.PackageManagerUnknown
}

func GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func GetDirSizeSafe(path string) int64 {
	size, err := GetDirSize(path)
	if err != nil {
		Debug(fmt.Sprintf("falha ao medir diretório %s: %s", path, err))
	}
	return size
}

func RemovePath(path string, recursive bool) error {
	if recursive {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}