package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joaomjbraga/piunter/internal/core"
	"github.com/joaomjbraga/piunter/internal/modules"
	"github.com/joaomjbraga/piunter/internal/utils"
	"github.com/joaomjbraga/piunter/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const VERSION = "1.3.0"

var (
	all                bool
	analyze            bool
	dryRun             bool
	force              bool
	interactive        bool
	list               bool
	threshold          int
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "piunter",
		Short: "CLI para limpeza e otimização de sistemas Linux",
		Long: `piunter - CLI para limpeza e otimização de sistemas Linux

Execute com módulos específicos ou use --all para executar todos.`,
		Run: func(cmd *cobra.Command, args []string) {
			runMain(args, cmd.Flags())
		},
	}

	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Executa todos os módulos")
	rootCmd.Flags().BoolVar(&analyze, "analyze", false, "Analisa sem limpar")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Simula a execução")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Pula confirmações")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Modo interativo")
	rootCmd.Flags().BoolVar(&list, "list", false, "Lista módulos disponíveis")
	rootCmd.Flags().IntVar(&threshold, "threshold", 100, "Tamanho mínimo para arquivos grandes (MB)")

	rootCmd.Flags().Bool("cache", false, "Limpa cache do usuário")
	rootCmd.Flags().Bool("npm", false, "Limpa cache do NPM")
	rootCmd.Flags().Bool("yarn", false, "Limpa cache do Yarn")
	rootCmd.Flags().Bool("pnpm", false, "Limpa cache do PNPM")
	rootCmd.Flags().Bool("nvm", false, "Limpa cache do NVM (Node Version Manager)")
	rootCmd.Flags().Bool("sdkman", false, "Limpa cache do SDKMAN")
	rootCmd.Flags().Bool("packages", false, "Remove pacotes órfãos")
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

var allModuleFlags = []string{
	"cache", "npm", "yarn", "pnpm", "nvm", "sdkman", "packages", "docker", "logs",
	"flatpak", "snap", "large-files", "appimage", "thumbs", "recent", "trash",
}

func getModuleIdsFromFlags(flags *pflag.FlagSet) []string {
	if all {
		return allModuleFlags
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

func runMain(args []string, flags *pflag.FlagSet) {
	distro := utils.GetDistroInfo()

	if list {
		printList()
		return
	}

	moduleIds := getModuleIdsFromFlags(flags)

	if analyze {
		printHeader(distro)
		runAnalyze(moduleIds)
		return
	}

	if interactive || len(moduleIds) == 0 {
		printHeader(distro)
		runInteractive(distro)
		return
	}

	printHeader(distro)

	if dryRun {
		fmt.Printf("  \033[33mModo dry-run ativo\033[0m\n\n")
	}

	if requiresSudo(moduleIds) && !utils.IsRoot() && !utils.HasSudoPassword() {
		fmt.Println()
		fmt.Println("  \033[33mAlguns módulos requerem privilégios de administrador.\033[0m")
		if !utils.RequestSudo() {
			fmt.Println("  \033[90mMódulos que requerem sudo serão pulados.\033[0m")
		}
		fmt.Println()
	}

	runClean(moduleIds)
}

func printHeader(distro types.DistroInfo) {
	fmt.Println()
	fmt.Printf("  \033[36;1mpiunter\033[0m \033[90m· CLI para Linux\033[0m\n")
	fmt.Printf("  \033[90m%s\033[0m\n", strings.Repeat("─", 30))
	fmt.Println()
	fmt.Printf("  \033[90mSistema: %s\033[0m\n", distro.Name)
	fmt.Printf("  \033[90mGerenciador: %s\033[0m\n", distro.PackageManager)
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
	if !force && !dryRun {
		fmt.Print("\n  Confirmar limpeza? (y/s/N) ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if response != "y" && response != "s" {
			fmt.Println("  \033[90mOperação cancelada.\033[0m")
			return
		}
	}

	cleaner := core.NewCleaner(moduleIds, dryRun)
	report, err := cleaner.Clean()
	if err != nil {
		utils.Error(fmt.Sprintf("Erro ao limpar: %s", err.Error()))
		return
	}
	cleaner.PrintReport(report)
}

func runInteractive(distro types.DistroInfo) {
	fmt.Println("  \033[1mSelecione os módulos:\033[0m")
	fmt.Println()

	available := modules.GetAvailableModules()
	for i, m := range available {
		fmt.Printf("  %d) %s - %s\n", i+1, m.Name(), m.Description())
	}
	fmt.Println()

	fmt.Print("  Digite os números separados por vírgula (ex: 1,2,3): ")
	var input string
	fmt.Scanln(&input)

	ids := parseModuleSelection(input, available)
	if len(ids) == 0 {
		fmt.Println("  \033[33mNenhum módulo selecionado.\033[0m")
		return
	}

	if !force {
		fmt.Print("\n  Confirmar limpeza? (y/s/N) ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if response != "y" && response != "s" {
			fmt.Println("  \033[90mOperação cancelada.\033[0m")
			return
		}
	}

	runClean(ids)
}

func parseModuleSelection(input string, available []modules.Module) []string {
	var ids []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var idx int
		if _, err := fmt.Sscanf(part, "%d", &idx); err == nil {
			if idx > 0 && idx <= len(available) {
				ids = append(ids, available[idx-1].ID())
			}
		}
	}
	return ids
}