# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.1] - 2026-03-27

### Fixed

- **Plugin Config Persistence:** `saveConfig()` now properly saves plugin configuration to disk (was empty before)
- **Docker itemsRemoved Calculation:** Fixed potential negative values when new Docker items are created during cleanup
- **Unused Variable:** Removed unused `totalSize` variable in `packages.ts`
- **Unused Module:** Removed unused `DiskUsageModule` class
- **Duplicate Tests:** Removed duplicate `getDistroInfo` tests in `logger.test.ts`

### Changed

- **Usage Model:** Changed from requiring global install to using `npx` for immediate execution without installation
- **Refactored npm/yarn/pnpm Modules:** Created `PackageCacheModule` base class to eliminate code duplication
- **Centralized parseSize():** Moved to `utils/fs.ts` to avoid duplication in `flatpak.ts` and `disk.ts`
- **Fixed Race Condition:** `LargeFilesModule` now uses local threshold variable instead of instance state
- **Security Enhancement:** Sudo password is now cleared from memory on process exit
- **Version Sync:** CLI version updated to match package.json (1.0.1)
- **Cleaner API:** Removed unused `_force` parameter from all module `clean()` methods
- Updated documentation to reflect `npx piunter` as the primary usage method
- Added ASCII art banner to README

## [1.2.0] - 2026-03-26

### Fixed

- **Command Injection Bug:** Pacman orphan removal now executes packages safely in loop instead of passing subshell string
- **Memory Leak:** Fixed stdin event listener not being cleaned up in sudo password prompt
- **Race Condition:** Added 30s timeout to sudo password prompt
- **Inaccurate Space Calculation:** All modules now calculate real space freed using before/after analysis
- **apt-orphans Detection:** APT module now detects orphan packages during analysis
- **promptYesNo Recursion:** Fixed potential stack overflow with recursive calls
- **Sudo Password Input:** Fixed stdin interference with inquirer, now properly accepts password input

### Changed

- Single-key confirmation (y/n) without needing Enter key
- Refactored duplicate code into shared utils (`src/utils/fs.ts`)
- Improved error handling and logging

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
