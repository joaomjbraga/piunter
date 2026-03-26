import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getDistroInfo, getHomeDir, getCacheDir, getConfigDir } from './os.js';
import { logger } from './logger.js';
import { exec, isCommandAvailable } from './exec.js';

describe('OS Utils', () => {
  it('should get distro info', () => {
    const distro = getDistroInfo();
    expect(distro).toBeDefined();
    expect(distro.id).toBeDefined();
    expect(distro.name).toBeDefined();
    expect(distro.packageManager).toBeDefined();
  });

  it('should get home directory', () => {
    const home = getHomeDir();
    expect(home).toBeDefined();
    expect(home.length).toBeGreaterThan(0);
    expect(home).toContain('/');
  });

  it('should get cache directory', () => {
    const cache = getCacheDir();
    expect(cache).toBeDefined();
    expect(cache).toContain('.cache');
  });

  it('should get config directory', () => {
    const config = getConfigDir();
    expect(config).toBeDefined();
    expect(config).toContain('.config');
  });
});

describe('Logger', () => {
  it('should format bytes correctly', () => {
    expect(logger.formatBytes(0)).toBe('0 B');
    expect(logger.formatBytes(1024)).toBe('1 KB');
    expect(logger.formatBytes(1024 * 1024)).toBe('1 MB');
    expect(logger.formatBytes(1024 * 1024 * 1024)).toBe('1 GB');
    expect(logger.formatBytes(1024 * 1024 * 1024 * 1024)).toBe('1 TB');
  });

  it('should format bytes with decimal precision', () => {
    expect(logger.formatBytes(1536)).toBe('1.5 KB');
    expect(logger.formatBytes(1572864)).toBe('1.5 MB');
  });
});

describe('Exec Utils', () => {
  it('should check if commands are available', () => {
    const nodeAvailable = isCommandAvailable('node');
    expect(nodeAvailable).toBe(true);
  });

  it('should return false for non-existent commands', () => {
    const fakeCommandAvailable = isCommandAvailable('thiscommanddoesnotexist12345');
    expect(fakeCommandAvailable).toBe(false);
  });
});
