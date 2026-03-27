import { describe, it, expect, vi } from 'vitest';
import { parseFlags, getModulesFromFlags } from './cli/index.js';

describe('CLI Flags', () => {
  describe('parseFlags', () => {
    it('should parse --all flag', () => {
      const flags = parseFlags(['--all']);
      expect(flags.all).toBe(true);
    });

    it('should parse -a as alias for --all', () => {
      const flags = parseFlags(['-a']);
      expect(flags.all).toBe(true);
    });

    it('should parse module flags', () => {
      const flags = parseFlags(['--cache', '--npm', '--docker']);
      expect(flags.cache).toBe(true);
      expect(flags.npm).toBe(true);
      expect(flags.docker).toBe(true);
    });

    it('should parse --dry-run and -n', () => {
      const flags = parseFlags(['--dry-run']);
      expect(flags.dryRun).toBe(true);

      const flags2 = parseFlags(['-n']);
      expect(flags2.dryRun).toBe(true);
    });

    it('should parse --force and -f', () => {
      const flags = parseFlags(['--force']);
      expect(flags.force).toBe(true);

      const flags2 = parseFlags(['-f']);
      expect(flags2.force).toBe(true);
    });

    it('should parse --interactive and -i', () => {
      const flags = parseFlags(['--interactive']);
      expect(flags.interactive).toBe(true);

      const flags2 = parseFlags(['-i']);
      expect(flags2.interactive).toBe(true);
    });

    it('should parse --threshold value', () => {
      const flags = parseFlags(['--threshold=500']);
      expect(flags.largeFilesThreshold).toBe(500);
    });

    it('should use default threshold when not provided', () => {
      const flags = parseFlags([]);
      expect(flags.largeFilesThreshold).toBe(100);
    });

    it('should clamp threshold to max', () => {
      const flags = parseFlags(['--threshold=50000']);
      expect(flags.largeFilesThreshold).toBe(10000);
    });

    it('should handle empty threshold value', () => {
      const flags = parseFlags(['--threshold=']);
      expect(flags.largeFilesThreshold).toBe(100);
    });

    it('should parse --analyze flag', () => {
      const flags = parseFlags(['--analyze']);
      expect(flags.analyze).toBe(true);
    });

    it('should parse all module flags', () => {
      const flags = parseFlags([
        '--cache',
        '--npm',
        '--yarn',
        '--pnpm',
        '--flatpak',
        '--snap',
        '--docker',
        '--logs',
        '--packages',
        '--large-files',
        '--appimage',
        '--thumbs',
        '--recent',
      ]);
      expect(flags.cache).toBe(true);
      expect(flags.npm).toBe(true);
      expect(flags.yarn).toBe(true);
      expect(flags.pnpm).toBe(true);
      expect(flags.flatpak).toBe(true);
      expect(flags.snap).toBe(true);
      expect(flags.docker).toBe(true);
      expect(flags.logs).toBe(true);
      expect(flags.packages).toBe(true);
      expect(flags.largeFiles).toBe(true);
      expect(flags.appimage).toBe(true);
      expect(flags.thumbs).toBe(true);
      expect(flags.recent).toBe(true);
    });
  });

  describe('getModulesFromFlags', () => {
    it('should return specific modules from flags', () => {
      const flags = {
        all: false,
        cache: true,
        npm: true,
        docker: false,
        yarn: false,
        pnpm: false,
        flatpak: false,
        snap: false,
        logs: false,
        packages: false,
        largeFiles: false,
        appimage: false,
        thumbs: false,
        recent: false,
      } as any;

      const modules = getModulesFromFlags(flags);
      expect(modules).toEqual(['cache', 'npm']);
    });

    it('should return empty array when no modules specified', () => {
      const flags = {
        all: false,
        cache: false,
        npm: false,
        yarn: false,
        pnpm: false,
        flatpak: false,
        snap: false,
        docker: false,
        logs: false,
        packages: false,
        largeFiles: false,
        appimage: false,
        thumbs: false,
        recent: false,
      } as any;

      const modules = getModulesFromFlags(flags);
      expect(modules).toEqual([]);
    });

    it('should return large-files module id', () => {
      const flags = {
        all: false,
        cache: false,
        npm: false,
        yarn: false,
        pnpm: false,
        flatpak: false,
        snap: false,
        docker: false,
        logs: false,
        packages: false,
        largeFiles: true,
        appimage: false,
        thumbs: false,
        recent: false,
      } as any;

      const modules = getModulesFromFlags(flags);
      expect(modules).toEqual(['large-files']);
    });
  });
});
