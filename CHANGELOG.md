# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

### Features
- TypeScript with strict mode
- Modular architecture (Clean Architecture)
- CLI with multiple flags
- Config file support (~/.piunter.json)
- Comprehensive error handling
- Command availability detection
