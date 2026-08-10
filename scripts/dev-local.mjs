import { spawn } from 'node:child_process';
import path from 'node:path';

const workspaceRoot = process.cwd();
const children = [
  spawn(process.execPath, [path.join(workspaceRoot, 'scripts/dev-pagescms.mjs')], { cwd: workspaceRoot, stdio: 'inherit' }),
  spawn(process.execPath, [path.join(workspaceRoot, 'node_modules/astro/astro.js'), 'dev', '--host', '127.0.0.1', '--port', '4321'], { cwd: path.join(workspaceRoot, 'apps/builder'), stdio: 'inherit' }),
];

const stop = () => children.forEach((child) => child.kill('SIGTERM'));
process.on('SIGINT', stop);
process.on('SIGTERM', stop);
for (const child of children) child.on('exit', (code) => { if (code) { stop(); process.exitCode = code; } });
