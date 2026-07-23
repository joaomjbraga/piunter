package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joaomjbraga/piunter/internal/config"
	"github.com/joaomjbraga/piunter/internal/core"
	"github.com/joaomjbraga/piunter/internal/modules"
	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const Version = "1.6.0"

var (
	all                bool
	analyze            bool
	dryRun             bool
	force              bool
	list               bool
	verbose            bool
	interactive        bool
	threshold          int
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:     "piunter",
		Version: Version,
		Short:   "CLI para limpeza e otimização de sistemas Linux",
		Long: `piunter - CLI para limpeza e otimização de sistemas Linux

Este projeto é voltado exclusivamente para ambientes Linux. Não há suporte para macOS, Windows ou outras plataformas.

Execute com módulos específicos ou use --all para executar todos.`,
		Run: func(cmd *cobra.Command, args []string) {
			runMain(args, cmd.Flags())
		},
	}

	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Executa todos os módulos")
	rootCmd.Flags().BoolVar(&analyze, "analyze", false, "Analisa sem limpar")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Simula a execução")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Pula confirmações")
	rootCmd.Flags().BoolVar(&list, "list", false, "Lista módulos disponíveis")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Mostra detalhes adicionais da execução")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Modo interativo para selecionar módulos")
	rootCmd.Flags().IntVar(&threshold, "threshold", 100, "Tamanho mínimo para arquivos grandes (MB)")
	rootCmd.AddCommand(completionCmd())
	rootCmd.Flags().Bool("cache", false, "Limpa cache do usuário")
	rootCmd.Flags().Bool("packages", false, "Remove pacotes órfãos")
	rootCmd.Flags().Bool("package-cache", false, "Limpa o cache de downloads dos gerenciadores de pacotes")
	rootCmd.Flags().Bool("temp-files", false, "Limpa arquivos temporários antigos em /tmp e /var/tmp")
	rootCmd.Flags().Bool("shell-history", false, "Limpa arquivos de histórico de shell do usuário")
	rootCmd.Flags().Bool("dev-cache", false, "Limpa caches de npm, pip e cargo do usuário")
	rootCmd.Flags().Bool("browser-cache", false, "Limpa caches de navegadores como Firefox e Chrome")
	rootCmd.Flags().Bool("editor-cache", false, "Limpa caches de editores como VS Code, IntelliJ e Vim")
	rootCmd.Flags().Bool("media-cache", false, "Limpa caches de mídia e aplicativos gráficos")
	rootCmd.Flags().Bool("game-cache", false, "Limpa caches de jogos e ferramentas de gaming")
	rootCmd.Flags().Bool("container-cache", false, "Limpa caches e artefatos de containers e VMs")
	rootCmd.Flags().Bool("build-cache", false, "Limpa caches de ferramentas de build e IA")
	rootCmd.Flags().Bool("ides-cache", false, "Limpa temporários de IDEs e resíduos de desinstalação")
	rootCmd.Flags().Bool("browser-plugins", false, "Limpa resíduos de extensões e plugins de navegadores")
	rootCmd.Flags().Bool("old-installers", false, "Remove instaladores antigos em Downloads e pastas semelhantes")
	rootCmd.Flags().Bool("swap-files", false, "Remove arquivos swap e temporários de editores")
	rootCmd.Flags().Bool("app-logs", false, "Limpa logs de aplicativos e shells")
	rootCmd.Flags().Bool("downloads-old", false, "Limpa arquivos antigos em Downloads e pastas semelhantes")
	rootCmd.Flags().Bool("docker", false, "Limpa Docker")
	rootCmd.Flags().Bool("logs", false, "Limpa logs do sistema")
	rootCmd.Flags().Bool("flatpak", false, "Limpa Flatpak")
	rootCmd.Flags().Bool("snap", false, "Limpa Snap")
	rootCmd.Flags().Bool("large-files", false, "Encontra arquivos grandes")
	rootCmd.Flags().Bool("appimage", false, "Remove AppImages")
	rootCmd.Flags().Bool("thumbs", false, "Limpa miniaturas")
	rootCmd.Flags().Bool("recent", false, "Lista arquivos recentes")
	rootCmd.Flags().Bool("trash", false, "Esvazia a lixeira")

	rootCmd.Flags().SortFlags = false
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func completionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh]",
		Short: "Gera autocompletar para Bash ou Zsh",
		Long:  "Gera scripts de autocompletar para Bash ou Zsh para o comando piunter.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			switch shell {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			default:
				return fmt.Errorf("shell suportado: bash ou zsh")
			}
		},
	}
	return cmd
}

var allModuleFlags = modules.GetModuleIDs()

func resolveModuleIds(flags *pflag.FlagSet, cfg config.Config) []string {
	if cfg.All || all {
		return allModuleFlags
	}

	if len(cfg.Modules) > 0 {
		return cfg.Modules
	}

	var ids []string

	for _, flag := range allModuleFlags {
		enabled, err := flags.GetBool(flag)
		if err == nil && enabled {
			ids = append(ids, flag)
		}
	}

	return ids
}

func applyConfigOverrides(allFlag, analyzeFlag, dryRunFlag, forceFlag, verboseFlag *bool, thresholdFlag *int, cfg config.Config) {
	if cfg.All {
		*allFlag = true
	}
	if cfg.Analyze {
		*analyzeFlag = true
	}
	if cfg.DryRun {
		*dryRunFlag = true
	}
	if cfg.Force {
		*forceFlag = true
	}
	if cfg.Verbose {
		*verboseFlag = true
	}
	if cfg.ThresholdMB > 0 {
		*thresholdFlag = cfg.ThresholdMB
	}
}

func requiresSudo(ids []string) bool {
	sudoModules := map[string]bool{
		"packages": true,
		"logs":     true,
		"flatpak":  true,
		"docker":   true,
		"snap":     true,
	}
	for _, id := range ids {
		if sudoModules[id] {
			return true
		}
	}
	return false
}

func formatModuleList(moduleIDs []string) string {
	if len(moduleIDs) == 0 {
		return "nenhum módulo selecionado"
	}

	return strings.Join(moduleIDs, ", ")
}

func runMain(args []string, flags *pflag.FlagSet) {
	cfg := config.Load()
	distro := utils.GetDistroInfo()

	if list {
		printList()
		return
	}

	if interactive {
		moduleIds := runInteractiveSelection()
		if len(moduleIds) == 0 {
			fmt.Println("  \033[90mNenhum módulo selecionado.\033[0m")
			return
		}
		runWizardConfig()
		wizardMode := runWizard()
		if wizardMode.analyze {
			printHeader(distro, cfg)
			runAnalyze(moduleIds)
			return
		}
		printHeader(distro, cfg)
		runClean(moduleIds)
		return
	}

	moduleIds := resolveModuleIds(flags, cfg)
	applyConfigOverrides(&all, &analyze, &dryRun, &force, &verbose, &threshold, cfg)

	if len(moduleIds) == 0 {
		printHeader(distro, cfg)
		printList()
		return
	}

	if analyze {
		printHeader(distro, cfg)
		runAnalyze(moduleIds)
		return
	}

	printHeader(distro, cfg)

	if dryRun {
		fmt.Printf("  \033[33mModo dry-run ativo\033[0m\n\n")
	}

	fmt.Printf("  \033[90mMódulos selecionados:\033[0m %s\n", formatModuleList(moduleIds))
	if requiresSudo(moduleIds) {
		fmt.Printf("  \033[33mRequer privilégios administrativos\033[0m\n")
	}
	if dryRun {
		fmt.Printf("  \033[33mNenhuma alteração será aplicada\033[0m\n")
	}
	fmt.Println()

	if requiresSudo(moduleIds) && !utils.IsRoot() && !utils.HasSudoPassword() {
		fmt.Println()
		fmt.Println("  \033[33mAlguns módulos requerem privilégios de administrador.\033[0m")
		if !utils.RequestSudo() {
			fmt.Println("  \033[90mMódulos que requerem sudo serão pulados.\033[0m")
			var filtered []string
			for _, id := range moduleIds {
				if !requiresSudo([]string{id}) {
					filtered = append(filtered, id)
				}
			}
			moduleIds = filtered
		}
		fmt.Println()
	}

	runClean(moduleIds)
}

func printHeader(distro types.DistroInfo, cfg config.Config) {
	fmt.Println()
	fmt.Printf("  \033[36;1mpiunter\033[0m \033[90m· CLI para Linux\033[0m\n")

	if !cfg.SkipUpdateCheck {
		if latest, err := utils.CheckForUpdate(Version); err == nil && latest != "" {
			fmt.Printf("  \033[33m!\033[0m Nova versão: \033[36m%s\033[0m\n", latest)
			fmt.Printf("  \033[90m  https://github.com/joaomjbraga/piunter/releases\033[0m\n")
		}
	}

	fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 30))
	fmt.Println()
	fmt.Printf("  \033[90mSistema: %s\033[0m\n", distro.Name)
	fmt.Println()
}

func printList() {
	fmt.Println()
	fmt.Println("  \033[1mMódulos disponíveis:\033[0m")
	fmt.Println()

	for _, m := range modules.GetAllModuleInfos() {
		status := "\033[32m*\033[0m"
		if !m.Available {
			status = "\033[31m-\033[0m"
		}
		fmt.Printf("  %s \033[37m%-15s\033[0m %s\n", status, m.Name, m.Description)
	}
	fmt.Println()
}

type wizardConfig struct {
	analyze bool
}

type moduleSelectionInfo struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

func getModuleCategory(moduleID string) string {
	categories := map[string]string{
		"cache":               "Cache e temporários",
		"package-cache":       "Cache e temporários",
		"temp-files":          "Cache e temporários",
		"shell-history":       "Cache e temporários",
		"dev-cache":           "Cache e temporários",
		"browser-cache":       "Cache e temporários",
		"editor-cache":        "Cache e temporários",
		"media-cache":         "Cache e temporários",
		"game-cache":          "Cache e temporários",
		"container-cache":     "Cache e temporários",
		"build-cache":         "Cache e temporários",
		"ides-cache":          "Cache e temporários",
		"browser-plugins":     "Cache e temporários",
		"old-installers":      "Cache e temporários",
		"swap-files":          "Cache e temporários",
		"app-logs":            "Logs e arquivos",
		"downloads-old":       "Logs e arquivos",
		"docker":              "Containers e ambientes",
		"logs":                "Logs e arquivos",
		"flatpak":             "Pacotes e ambientes",
		"snap":                "Pacotes e ambientes",
		"packages":            "Pacotes e ambientes",
		"large-files":         "Arquivos e disco",
		"appimage":            "Pacotes e ambientes",
		"thumbs":              "Cache e temporários",
		"recent":              "Arquivos e disco",
		"trash":               "Arquivos e disco",
	}
	if category, ok := categories[moduleID]; ok {
		return category
	}
	return "Outros"
}

func normalizeSelectionToken(token string) string {
	value := strings.ToLower(strings.TrimSpace(token))
	value = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == ' ', r == '-', r == '_':
			return '-'
		default:
			return -1
		}
	}, value)
	return strings.Trim(value, "-")
}

func getInteractiveSuggestions(input string, infos []moduleSelectionInfo) []string {
	query := normalizeSelectionToken(input)
	if query == "" {
		return nil
	}

	var matches []string
	for _, info := range infos {
		if !info.Available {
			continue
		}
		if strings.Contains(normalizeSelectionToken(info.ID), query) || strings.Contains(normalizeSelectionToken(info.Name), query) {
			matches = append(matches, info.ID)
		}
	}
	return matches
}

func parseInteractiveSelection(input string, infos []moduleSelectionInfo) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	parts := strings.Fields(input)
	selected := make([]string, 0, len(parts))
	seen := make(map[string]bool)

	for _, part := range parts {
		if strings.EqualFold(part, "all") || strings.EqualFold(part, "todos") || strings.EqualFold(part, "tudo") {
			for _, info := range infos {
				if info.Available && !seen[info.ID] {
					selected = append(selected, info.ID)
					seen[info.ID] = true
				}
			}
			continue
		}

		if index, err := strconv.Atoi(part); err == nil {
			if index < 1 || index > len(infos) {
				continue
			}
			info := infos[index-1]
			if info.Available && !seen[info.ID] {
				selected = append(selected, info.ID)
				seen[info.ID] = true
			}
			continue
		}

		normalized := normalizeSelectionToken(part)
		for _, info := range infos {
			if !info.Available || seen[info.ID] {
				continue
			}
			if normalizeSelectionToken(info.ID) == normalized || normalizeSelectionToken(info.Name) == normalized {
				selected = append(selected, info.ID)
				seen[info.ID] = true
				break
			}
		}
	}

	return selected
}

func runWizard() wizardConfig {
	fmt.Println()
	fmt.Println("  \033[1mAssistente rápido\033[0m")
	fmt.Println("  1. Analisar apenas")
	fmt.Println("  2. Limpar (modo padrão)")
	fmt.Println("  3. Simular execução (dry-run)")
	fmt.Println("  4. Modo interativo completo")
	fmt.Println()
	fmt.Print("  Escolha uma opção: ")

	var input string
	fmt.Scanln(&input)
	response := strings.ToLower(strings.TrimSpace(input))

	switch response {
	case "1", "analisar", "a":
		return wizardConfig{analyze: true}
	case "3", "dry", "dry-run", "simular", "s":
		dryRun = true
		return wizardConfig{}
	case "4", "interativo", "i":
		interactive = true
		return wizardConfig{}
	default:
		return wizardConfig{}
	}
}

func runWizardConfig() {
	fmt.Println()
	fmt.Println("  \033[1mConfiguração rápida\033[0m")
	fmt.Print("  Usar todos os módulos? (s/N): ")
	var allInput string
	fmt.Scanln(&allInput)
	if strings.ToLower(strings.TrimSpace(allInput)) == "s" {
		all = true
	}

	fmt.Print("  Ativar força para pular confirmações? (s/N): ")
	var forceInput string
	fmt.Scanln(&forceInput)
	if strings.ToLower(strings.TrimSpace(forceInput)) == "s" {
		force = true
	}

	fmt.Print("  Mostrar saída detalhada? (s/N): ")
	var verboseInput string
	fmt.Scanln(&verboseInput)
	if strings.ToLower(strings.TrimSpace(verboseInput)) == "s" {
		verbose = true
	}
}

func confirmInteractiveSelection(selected []string, analyzeMode bool) bool {
	if len(selected) == 0 {
		return false
	}

	fmt.Println()
	fmt.Println("  \033[1mRevisão final\033[0m")
	fmt.Printf("  Ação: %s\n", map[bool]string{true: "analisar", false: "limpar"}[analyzeMode])
	fmt.Printf("  Módulos selecionados: %s\n", strings.Join(selected, ", "))
	if dryRun {
		fmt.Println("  Modo: dry-run")
	}
	if force {
		fmt.Println("  Força: ativada")
	}
	if verbose {
		fmt.Println("  Verbose: ativado")
	}
	fmt.Println("  1. Confirmar e executar")
	fmt.Println("  2. Voltar e editar a seleção")
	fmt.Println("  3. Cancelar")
	fmt.Println()
	fmt.Print("  Escolha uma opção: ")

	var input string
	fmt.Scanln(&input)
	response := strings.ToLower(strings.TrimSpace(input))

	switch response {
	case "1", "confirmar", "c", "executar", "e":
		return true
	case "2", "editar", "voltar", "v":
		return false
	default:
		return false
	}
}

func runInteractiveSelection() []string {
	fmt.Println()
	fmt.Println("  \033[1mModo interativo\033[0m")
	fmt.Println("  Selecione os módulos por número, ID ou nome.")
	fmt.Println("  Exemplos: 1 3 5 | package-cache | cache")
	fmt.Println()

	infos := make([]moduleSelectionInfo, 0, len(modules.GetAllModuleInfos()))
	for _, info := range modules.GetAllModuleInfos() {
		infos = append(infos, moduleSelectionInfo{
			ID:          info.ID,
			Name:        info.Name,
			Description: info.Description,
			Available:   info.Available,
		})
	}

	currentCategory := ""
	for i, info := range infos {
		category := getModuleCategory(info.ID)
		if currentCategory != category {
			currentCategory = category
			fmt.Println()
			fmt.Printf("  \033[36m%s\033[0m\n", category)
		}
		status := "*"
		if !info.Available {
			status = "-"
		}
		fmt.Printf("  %d. %s - %s\n", i+1, info.Name, info.Description)
		_ = status
	}
	fmt.Println()
	fmt.Print("  Escolha os módulos: ")

	var input string
	fmt.Scanln(&input)
	selected := parseInteractiveSelection(input, infos)
	if len(selected) == 0 {
		suggestions := getInteractiveSuggestions(input, infos)
		if len(suggestions) > 0 {
			fmt.Println()
			fmt.Printf("  \033[90mSugestões: %s\033[0m\n", strings.Join(suggestions, ", "))
		}
	} else {
		fmt.Println()
		fmt.Printf("  \033[32mMódulos selecionados:\033[0m %s\n", strings.Join(selected, ", "))
	}

	if !confirmInteractiveSelection(selected, false) {
		return nil
	}

	return selected
}

func runAnalyze(moduleIds []string) {
	analyzer := core.NewAnalyzer(moduleIds, threshold)
	results, err := analyzer.Analyze()
	if err != nil {
		utils.Error(fmt.Sprintf("Erro ao analisar: %s", err.Error()))
		return
	}
	analyzer.PrintAnalysis(results)
}

func runClean(moduleIds []string) {
	fmt.Println()
	fmt.Printf("  \033[33mResumo da operação\033[0m\n")
	fmt.Printf("  %s\n", core.BuildConfirmationSummary(moduleIds, dryRun))

	if !force && !dryRun {
		fmt.Println()
		fmt.Println("  \033[31mAtenção:\033[0m esta ação pode remover arquivos permanentemente.")
		fmt.Print("  Confirmar limpeza? (y/s/N) ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if response != "y" && response != "s" {
			fmt.Println("  \033[90mOperação cancelada.\033[0m")
			return
		}
		fmt.Println()
		fmt.Println("  \033[32mIniciando limpeza...\033[0m")
	} else {
		fmt.Println()
		fmt.Printf("  \033[32m%s\033[0m\n", map[bool]string{true: "Execução preparada", false: "Execução iniciada"}[dryRun])
	}

	cleaner := core.NewCleanerWithOptions(moduleIds, dryRun, verbose)
	report, err := cleaner.Clean()
	if err != nil {
		utils.Error(fmt.Sprintf("Erro ao limpar: %s", err.Error()))
		return
	}
	cleaner.PrintReport(report)
}