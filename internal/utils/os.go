package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joaomjbraga/piunter/pkg/types"
)

func readOutput(stdout io.ReadCloser, stderr io.ReadCloser) (string, string, error) {
	var outBuf, errBuf bytes.Buffer
	outCh := make(chan error, 1)

	go func() {
		_, err := io.Copy(&outBuf, stdout)
		outCh <- err
	}()

	_, copyErr := io.Copy(&errBuf, stderr)
	<-outCh

	return outBuf.String(), errBuf.String(), copyErr
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

	output, _, _ := readOutput(stdout, stderr)
	cmd.Wait()

	return types.CommandResult{
		Success: cmd.ProcessState.ExitCode() == 0,
		Stdout:  output,
		Stderr:  "",
		Code:    cmd.ProcessState.ExitCode(),
	}
}

func ExecWithEnv(env []string, command string, args ...string) types.CommandResult {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), env...)
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
	fmt.Print("Senha de sudo: ")
	cmd := exec.Command("sudo", "true")
	err := cmd.Run()
	return err == nil
}

func GetHomeDir() string {
	home, _ := os.UserHomeDir()
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

func GetDirSizeAsync(path string) int64 {
	size, _ := GetDirSize(path)
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

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}