import path from 'node:path';
import { parseArgs } from 'node:util';
import { startLoggedProcess } from './start.mjs';

try {
  const separator = process.argv.indexOf('--', 2);
  if (separator < 0 || !process.argv[separator + 1]) {
    throw new Error('Usage: node tools/dev-process/cli.mjs --logs <directory> [--cwd <directory>] -- <executable> [arguments...]');
  }
  const { values } = parseArgs({
    args: process.argv.slice(2, separator),
    options: { logs: { type: 'string' }, cwd: { type: 'string', default: process.cwd() } },
  });
  if (!values.logs) throw new Error('--logs is required; use a separate log directory for each service.');
  const stdoutPath = path.resolve(values.logs, 'stdout.log');
  const stderrPath = path.resolve(values.logs, 'stderr.log');
  const result = await startLoggedProcess({
    command: process.argv[separator + 1], args: process.argv.slice(separator + 2),
    cwd: values.cwd, stdoutPath, stderrPath,
  });
  console.log(JSON.stringify({ ...result, stdoutPath, stderrPath }));
} catch (error) {
  console.error(error.message);
  process.exitCode = 1;
}
