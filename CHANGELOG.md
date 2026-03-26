# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-03-26

### Added

- **New Cleaning Modules:**
  - Snap support
  - AppImage detection and cleanup
  - System thumbnails cleaner
  - Recent files cleaner

- **Shell Completion:** Bash and Zsh auto-completion support
- **Plugin System:** Extensible architecture for custom plugins
- **Progress Bars:** Visual progress indicators for operations
- **Wiki Documentation:** Comprehensive guides in `/wiki` folder

### Fixed

- Removed unused dependencies (commander)
- ESLint/TypeScript linting configuration
- Error handling and fallbacks for all modules
- Code quality improvements across all modules

### Changed

- Improved CLI help output
- Better error messages with sudo suggestions
- Enhanced module availability detection

## [1.0.0] - 2026-03-26

### Added

- Initial release
- Detection of Linux distribution (Debian, Arch, Fedora)
- Package manager support (APT, Pacman, DNF)
- Cleaning modules:
  - User cache (~/.cache)
  - NPM, Yarn, PNPM
  - Flatpak
  - Docker
  - System logs (journalctl)
  - Package manager cache
  - Large files detection
- Interactive mode with checkbox selection
- Dry-run mode for simulation
- Confirmation prompts for destructive operations
- Detailed reports with space freed
- Automatic sudo elevation for privileged operations
- ASCII art banner
- Config file support (~/.piunter.json)
- Comprehensive error handling
- Command availability detection
- TypeScript with strict mode
- Modular architecture (Clean Architecture)
- Test suite with Vitest
