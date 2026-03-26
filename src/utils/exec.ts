import { execFileSync } from 'child_process';
import type { CommandResult } from '../types/index.js';

export async function exec(command: string, args: string[] = [], options: { sudo?: boolean } = {}): Promise<CommandResult> {
  const actualArgs = options.sudo ? ['-S', command, ...args] : args;
  const actualCommand = options.sudo ? 'sudo' : command;

  try {
    const result = execFileSync(actualCommand, actualArgs, {
      encoding: 'utf-8',
      timeout: 300000,
      maxBuffer: 50 * 1024 * 1024,
    });

    return {
      success: true,
      stdout: result,
      stderr: '',
      code: 0,
    };
  } catch (error: unknown) {
    const err = error as { status?: number; message?: string };
    return {
      success: false,
      stdout: '',
      stderr: err.message || '',
      code: err.status || 1,
    };
  }
}

export async function execWithOutput(
  command: string,
  args: string[] = [],
  options: { sudo?: boolean } = {}
): Promise<string> {
  const result = await exec(command, args, options);
  return result.stdout;
}

export function isCommandAvailable(command: string): boolean {
  try {
    execFileSync('which', [command], { encoding: 'utf-8', timeout: 5000 });
    return true;
  } catch {
    return false;
  }
}

export async function getCommandPath(command: string): Promise<string | null> {
  try {
    const result = execFileSync('which', [command], { encoding: 'utf-8', timeout: 5000 });
    return result.trim();
  } catch {
    return null;
  }
}
