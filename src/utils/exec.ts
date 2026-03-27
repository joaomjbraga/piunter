import { execFileSync, spawn } from 'child_process';
import chalk from 'chalk';
import type { CommandResult } from '../types/index.js';

let sudoPassword: string | null = null;

export function hasSudoPassword(): boolean {
  return sudoPassword !== null;
}

export async function requestSudo(): Promise<boolean> {
  if (sudoPassword) return true;

  return new Promise((resolve) => {
    process.stdout.write(chalk.yellow('  Senha sudo: '));

    const child = spawn('bash', ['-c', 'read -s -p "" pass && echo "$pass"'], {
      stdio: ['inherit', 'pipe', 'inherit'],
      windowsHide: true,
    });

    let password = '';
    const timeoutId = setTimeout(() => {
      child.kill();
      console.log(chalk.red('  Timeout. Operacoes que requerem sudo serao puladas.'));
      resolve(false);
    }, 30000);

    child.stdout?.on('data', (data: Buffer) => {
      password += data.toString();
    });

    child.on('close', () => {
      clearTimeout(timeoutId);
      password = password.trim();

      if (!password) {
        console.log(chalk.red('  Senha vazia. Operacoes que requerem sudo serao puladas.'));
        resolve(false);
        return;
      }

      try {
        execFileSync('sudo', ['-S', 'true'], {
          input: password + '\n',
          timeout: 10000,
        });
        sudoPassword = password;
        console.log(chalk.green('  Sudo confirmado.'));
        resolve(true);
      } catch {
        console.log(chalk.red('  Senha incorreta.'));
        resolve(false);
      }
    });

    child.on('error', () => {
      clearTimeout(timeoutId);
      console.log(chalk.red('  Erro ao pedir senha.'));
      resolve(false);
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
