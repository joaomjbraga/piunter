package modules

var AllModules = []Module{}

func init() {
	AllModules = []Module{
		NewPackagesModule(),
		NewCacheModule(),
		NewNpmModule(),
		NewYarnModule(),
		NewPnpmModule(),
		NewFlatpakModule(),
		NewSnapModule(),
		NewDockerModule(),
		NewLogsModule(),
		NewLargeFilesModule(),
		NewAppImageModule(),
		NewThumbsModule(),
		NewRecentFilesModule(),
		NewTrashModule(),
	}
}

func GetModule(id string) Module {
	for _, m := range AllModules {
		if m.ID() == id {
			return m
		}
	}
	return nil
}

func GetModulesByIds(ids []string) []Module {
	var result []Module
	for _, id := range ids {
		m := GetModule(id)
		if m != nil {
			result = append(result, m)
		}
	}
	return result
}

func GetAvailableModules() []Module {
	var result []Module
	for _, m := range AllModules {
		if m.IsAvailable() {
			result = append(result, m)
		}
	}
	return result
}

func GetAllModuleInfos() []struct {
	ID          string
	Name        string
	Description string
	Available   bool
} {
	var result []struct {
		ID          string
		Name        string
		Description string
		Available   bool
	}
	for _, m := range AllModules {
		result = append(result, struct {
			ID          string
			Name        string
			Description string
			Available   bool
		}{
			ID:          m.ID(),
			Name:        m.Name(),
			Description: m.Description(),
			Available:   m.IsAvailable(),
		})
	}
	return result
}