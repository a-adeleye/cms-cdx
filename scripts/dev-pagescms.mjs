import { spawn } from 'node:child_process';
import path from 'node:path';

const workspaceRoot = process.cwd();
const child = spawn(process.execPath, [path.join(workspaceRoot, 'node_modules/next/dist/bin/next'), 'dev', '--hostname', '127.0.0.1', '--port', '3000'], {
  cwd: path.join(workspaceRoot, 'apps/pagescms'),
  stdio: 'inherit',
  env: {
    ...process.env,
    NODE_ENV: 'development',
    PAGESCMS_LOCAL_MODE: 'true',
    NEXT_PUBLIC_PAGESCMS_LOCAL_MODE: 'true',
    PAGESCMS_LOCAL_ROOT: path.resolve(workspaceRoot),
    PAGESCMS_LOCAL_OWNER: 'local',
    PAGESCMS_LOCAL_REPO: path.basename(workspaceRoot),
    PAGESCMS_LOCAL_BRANCH: 'working-tree',
    DATABASE_URL: process.env.DATABASE_URL || 'postgresql://pagescms:pagescms@127.0.0.1:5434/pagescms',
    BETTER_AUTH_SECRET: process.env.BETTER_AUTH_SECRET || 'local-development-only-secret-change-me',
    CRYPTO_KEY: process.env.CRYPTO_KEY || 'local-development-only-crypto-key',
    AI_PROVIDER: process.env.AI_PROVIDER || 'mock',
    AI_REQUEST_TIMEOUT_MS: process.env.AI_REQUEST_TIMEOUT_MS || '30000',
  },
});

child.on('exit', (code) => process.exit(code ?? 1));
