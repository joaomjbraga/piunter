package types

type PackageManager string

const (
	PackageManagerApt     PackageManager = "apt"
	PackageManagerPacman  PackageManager = "pacman"
	PackageManagerDnf     PackageManager = "dnf"
	PackageManagerUnknown PackageManager = "unknown"
)

type DistroInfo struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Version        string         `json:"version"`
	PackageManager PackageManager `json:"packageManager"`
}

type CommandResult struct {
	Success bool
	Stdout  string
	Stderr  string
	Code    int
}

type CleanableItem struct {
	Path        string
	Size        int64
	Type        string
	Description string
}

type CleaningResult struct {
	Module       string
	Success      bool
	SpaceFreed   int64
	ItemsRemoved int
	Errors       []string
}

type AnalysisResult struct {
	Module    string
	Items     []CleanableItem
	TotalSize int64
}
type Report struct {
	StartTime        string
	EndTime          string
	Modules          []CleaningResult
	TotalSpaceFreed  int64
	TotalItemsRemoved int
	Errors           []string
}
