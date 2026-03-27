export class PiunterError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly module?: string,
    public readonly details?: Record<string, unknown>
  ) {
    super(message);
    this.name = 'PiunterError';
    Error.captureStackTrace(this, PiunterError);
  }
}

export class ModuleNotAvailableError extends PiunterError {
  constructor(moduleName: string) {
    super(`Módulo ${moduleName} não está disponível`, 'MODULE_NOT_AVAILABLE', moduleName);
    this.name = 'ModuleNotAvailableError';
  }
}

export class SudoRequiredError extends PiunterError {
  constructor(operation: string) {
    super(`Operação requer privilégios sudo: ${operation}`, 'SUDO_REQUIRED');
    this.name = 'SudoRequiredError';
  }
}

export class CommandExecutionError extends PiunterError {
  constructor(
    command: string,
    message: string,
    public readonly exitCode?: number
  ) {
    super(`Erro ao executar ${command}: ${message}`, 'COMMAND_EXECUTION_ERROR', undefined, {
      command,
      exitCode,
    });
    this.name = 'CommandExecutionError';
  }
}

export function formatError(error: unknown): string {
  if (error instanceof PiunterError) {
    return `[${error.code}] ${error.message}`;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return String(error);
}

export function isRetryableError(error: unknown): boolean {
  if (error instanceof CommandExecutionError) {
    return error.exitCode === undefined || error.exitCode === 124;
  }
  return false;
}
