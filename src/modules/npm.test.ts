import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { NpmModule, YarnModule, PnpmModule } from './npm.js';
import * as fs from 'fs';
import * as path from 'path';
import { exec, isCommandAvailable } from '../utils/exec.js';
import { getHomeDir } from '../utils/os.js';
import { getDirSize } from '../utils/fs.js';

const HOME_DIR = '/home/testuser';

vi.mock('fs');
vi.mock('../utils/exec.js', () => ({
  exec: vi.fn(),
  isCommandAvailable: vi.fn().mockReturnValue(true),
}));

vi.mock('../utils/os.js', () => ({
  getHomeDir: vi.fn(() => HOME_DIR),
}));

vi.mock('../utils/fs.js', () => ({
  getDirSize: vi.fn((p: string) => {
    if (p.includes('.npm')) return 1000000;
    if (p.includes('.yarn')) return 500000;
    if (p.includes('.pnpm')) return 300000;
    return 0;
  }),
}));

describe('NpmModule', () => {
  let npmModule: NpmModule;

  beforeEach(() => {
    vi.clearAllMocks();
    npmModule = new NpmModule();
  });

  it('should have correct id', () => {
    expect(npmModule.id).toBe('npm');
  });

  it('should have correct name', () => {
    expect(npmModule.name).toBe('NPM');
  });

  it('should be available when npm is installed', () => {
    expect(npmModule.isAvailable()).toBe(true);
  });

  it('should return correct cache paths', () => {
    const paths = npmModule.getCachePaths();
    expect(paths).toContain(path.join(HOME_DIR, '.npm'));
  });

  it('should return correct clean command', () => {
    const cmd = npmModule.getCleanCommand();
    expect(cmd).toEqual(['cache', 'clean', '--force']);
  });

  it('should analyze npm cache', async () => {
    vi.mocked(fs.existsSync).mockReturnValue(true);

    const result = await npmModule.analyze();

    expect(result.module).toBe('npm');
    expect(result.items).toHaveLength(1);
    expect(result.items[0].path).toContain('.npm');
    expect(result.totalSize).toBeGreaterThan(0);
  });
});

describe('YarnModule', () => {
  let yarnModule: YarnModule;

  beforeEach(() => {
    vi.clearAllMocks();
    yarnModule = new YarnModule();
  });

  it('should have correct id', () => {
    expect(yarnModule.id).toBe('yarn');
  });

  it('should have correct name', () => {
    expect(yarnModule.name).toBe('Yarn');
  });

  it('should return correct cache paths', () => {
    const paths = yarnModule.getCachePaths();
    expect(paths).toContain(path.join(HOME_DIR, '.yarn', 'cache'));
  });

  it('should return correct clean command', () => {
    const cmd = yarnModule.getCleanCommand();
    expect(cmd).toEqual(['cache', 'clean']);
  });
});

describe('PnpmModule', () => {
  let pnpmModule: PnpmModule;

  beforeEach(() => {
    vi.clearAllMocks();
    pnpmModule = new PnpmModule();
  });

  it('should have correct id', () => {
    expect(pnpmModule.id).toBe('pnpm');
  });

  it('should have correct name', () => {
    expect(pnpmModule.name).toBe('PNPM');
  });

  it('should return correct clean command', () => {
    const cmd = pnpmModule.getCleanCommand();
    expect(cmd).toEqual(['store', 'prune']);
  });
});
