import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PackagesModule } from './packages.js';
import { exec } from '../utils/exec.js';

vi.mock('../utils/exec.js', () => ({
  exec: vi.fn(),
  isCommandAvailable: vi.fn().mockReturnValue(true),
}));

vi.mock('../utils/os.js', () => ({
  getDistroInfo: vi.fn(() => ({
    id: 'ubuntu',
    name: 'Ubuntu',
    version: '22.04',
    packageManager: 'apt',
  })),
}));

vi.mock('../utils/fs.js', () => ({
  parseSize: vi.fn((s: string) => {
    if (s.includes('MB')) return parseFloat(s) * 1024 * 1024;
    if (s.includes('GB')) return parseFloat(s) * 1024 * 1024 * 1024;
    return parseFloat(s);
  }),
}));

describe('PackagesModule', () => {
  let packagesModule: PackagesModule;

  beforeEach(() => {
    vi.clearAllMocks();
    packagesModule = new PackagesModule();
  });

  it('should have correct id', () => {
    expect(packagesModule.id).toBe('packages');
  });

  it('should have correct name', () => {
    expect(packagesModule.name).toBe('Gerenciador de Pacotes');
  });

  it('should be available on known package managers', () => {
    expect(packagesModule.isAvailable()).toBe(true);
  });

  it('should analyze apt cache', async () => {
    vi.mocked(exec)
      .mockResolvedValueOnce({
        success: true,
        stdout: '1000000\t/var/cache/apt/archives',
        stderr: '',
        code: 0,
      })
      .mockResolvedValueOnce({
        success: true,
        stdout: '3 packages',
        stderr: '',
        code: 0,
      });

    const result = await packagesModule.analyze();

    expect(result.module).toBe('packages');
    expect(result.items.length).toBeGreaterThanOrEqual(1);
  });

  it('should calculate space freed correctly', async () => {
    vi.mocked(exec).mockResolvedValue({
      success: true,
      stdout: '',
      stderr: '',
      code: 0,
    });

    const result = await packagesModule.clean(false);

    expect(result.module).toBe('packages');
    expect(result.spaceFreed).toBeGreaterThanOrEqual(0);
  });

  it('should return error on unsupported package manager', async () => {
    const customModule = Object.create(PackagesModule.prototype);
    customModule.packageManager = 'unknown';

    const result = await customModule.clean(false);
    expect(result.errors).toContain('Gerenciador de pacotes não suportado');
  });
});
