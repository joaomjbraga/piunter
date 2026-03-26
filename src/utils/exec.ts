import { execFileSync } from 'child_process';
import chalk from 'chalk';
import type { CommandResult } from '../types/index.js';

let sudoPassword: string | null = null;

export function hasSudoPassword(): boolean {
  return sudoPassword !== null;
}

export async function requestSudo(): Promise<boolean> {
  if (sudoPassword) return true;

  return new Promise((resolve) => {
    const password: string[] = [];
    let cursorPos = 0;

    process.stdout.write(chalk.yellow('  Senha sudo: '));
    process.stdin.setRawMode(true);
    process.stdin.resume();

    const cleanup = () => {
      process.stdin.setRawMode(false);
      process.stdin.pause();
    };

    const finish = () => {
      cleanup();
      process.stdout.write('\n');
      
      const pass = password.join('');
      if (!pass) {
        console.log(chalk.red('  Senha vazia. Operacoes que requerem sudo serao puladas.'));
        resolve(false);
        return;
      }

      try {
        execFileSync('sudo', ['-S', 'true'], {
          input: pass + '\n',
          timeout: 10000,
        });
        sudoPassword = pass;
        console.log(chalk.green('  Sudo confirmado.'));
        resolve(true);
      } catch {
        console.log(chalk.red('  Senha incorreta.'));
        resolve(false);
      }
    };

    process.stdin.on('data', (chunk) => {
      const char = chunk.toString();
      
      if (char === '\x03') {
        cleanup();
        process.stdout.write('^C\n');
        process.exit(0);
      }

      if (char === '\r' || char === '\n') {
        finish();
        return;
      }

      if (char === '\x7f') {
        if (cursorPos > 0) {
          password.splice(--cursorPos, 1);
          process.stdout.write('\b \b');
        }
        return;
      }

      if (char === '\x1b[D') {
        if (cursorPos > 0) {
          cursorPos--;
          process.stdout.write('\x1b[D');
        }
        return;
      }

      if (char === '\x1b[C') {
        if (cursorPos < password.length) {
          cursorPos++;
          process.stdout.write('\x1b[C');
        }
        return;
      }

      password.splice(cursorPos, 0, char);
      cursorPos++;
      process.stdout.write('*');
    });
  });
}

export async function exec(
  command: string,
  args: string[] = [],
  options: { sudo?: boolean } = {}
): Promise<CommandResult> {
  let actualArgs = args;
  let actualCommand = command;

  if (options.sudo) {
    if (!sudoPassword) {
      return {
        success: false,
        stdout: '',
        stderr: 'Sudo requerido mas senha nao foi fornecida',
        code: 1,
      };
    }
    actualArgs = ['-S', command, ...args];
    actualCommand = 'sudo';
  }

  try {
    const result = execFileSync(actualCommand, actualArgs, {
      encoding: 'utf-8',
      timeout: 300000,
      maxBuffer: 50 * 1024 * 1024,
      input: options.sudo ? sudoPassword + '\n' : undefined,
      stdio: options.sudo ? 'pipe' : 'pipe',
    });

    return {
      success: true,
      stdout: result,
      stderr: '',
      code: 0,
    };
  } catch (error: unknown) {
    const err = error as { status?: number; message?: string; signal?: string };
    
    if (err.signal === 'SIGTERM' || err.message?.includes('sudo')) {
      return {
        success: false,
        stdout: '',
        stderr: 'Senha sudo incorreta ou expirada',
        code: 1,
      };
    }

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
