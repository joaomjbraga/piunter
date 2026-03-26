# Contributing to piunter

Thank you for your interest in contributing to piunter!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/joaomjbraga/piunter.git
cd piunter
```

2. Install dependencies:
```bash
npm install
```

3. Build the project:
```bash
npm run build
```

4. Run tests:
```bash
npm test
```

## Coding Standards

- Use TypeScript for all new code
- Run `npm run lint` before committing
- Run `npm run format` to format code
- Write tests for new features

## Project Structure

```
src/
├── cli.ts          # CLI entry point
├── core/           # Core logic (analyzer, cleaner)
├── modules/        # Cleaning modules
├── utils/          # Utilities
└── types/          # TypeScript types
```

## Submitting Changes

1. Create a feature branch:
```bash
git checkout -b feature/my-feature
```

2. Make your changes and commit:
```bash
git commit -m "feat: add new feature"
```

3. Push to your fork:
```bash
git push origin feature/my-feature
```

4. Open a Pull Request

## Commit Message Format

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance

## Questions?

Open an issue on GitHub for questions about contributing.
