import fs from 'node:fs';
import fsp from 'node:fs/promises';
import path from 'node:path';
import { spawn, execFile } from 'node:child_process';
import { createHash, randomUUID } from 'node:crypto';
import { fileURLToPath } from 'node:url';
import { promisify } from 'node:util';

const execFileAsync = promisify(execFile);
const windows = process.platform === 'win32';
const builds = new Map();

async function windowsStarter(directory, env) {
  const source = fileURLToPath(new URL('./start-hidden-windows.go', import.meta.url));
  const digest = createHash('sha256').update(await fsp.readFile(source)).digest('hex').slice(0, 16);
  const target = path.join(directory, `start-hidden-${process.arch}-${digest}.exe`);
  if (fs.existsSync(target)) return target;
  if (builds.has(target)) return builds.get(target);

  const build = (async () => {
    // A unique output keeps concurrent launches from replacing a running helper.
    const temporary = `${target}.${randomUUID()}.exe`;
    try {
      await execFileAsync('go', ['build', '-o', temporary, source], {
        env: { ...env, GOOS: 'windows', GOARCH: { x64: 'amd64', arm64: 'arm64', ia32: '386' }[process.arch], CGO_ENABLED: '0' },
        windowsHide: true,
        timeout: 120000,
      });
      try {
        if (!fs.existsSync(target)) await fsp.rename(temporary, target);
      } catch (error) {
        if (!fs.existsSync(target)) throw error;
      }
      return target;
    } finally {
      await fsp.rm(temporary, { force: true });
    }
  })();
  builds.set(target, build);
  try { return await build; }
  finally { builds.delete(target); }
}

// All paths are resolved from the caller's working directory. This starts a
// native executable directly; callers must explicitly select a shell if needed.
export async function startLoggedProcess({ command, args = [], cwd = process.cwd(), stdoutPath, stderrPath, env = process.env }) {
  cwd = path.resolve(cwd);
  stdoutPath = path.resolve(stdoutPath);
  stderrPath = path.resolve(stderrPath);
  await fsp.mkdir(path.dirname(stdoutPath), { recursive: true });
  await fsp.mkdir(path.dirname(stderrPath), { recursive: true });
  const starter = windows ? await windowsStarter(path.dirname(stdoutPath), env) : null;
  const pidPath = windows ? path.join(path.dirname(stdoutPath), `start-${randomUUID()}.pid`) : null;
  let stdoutFd;
  let stderrFd;
  try {
    stdoutFd = fs.openSync(stdoutPath, 'a');
    stderrFd = fs.openSync(stderrPath, 'a');
    // Windows can open a directory with append flags. Reject it before launch
    // instead of starting a worker whose output cannot reach the log.
    if (!fs.fstatSync(stdoutFd).isFile() || !fs.fstatSync(stderrFd).isFile()) {
      throw new Error('Log destinations must be regular files.');
    }
    const child = spawn(windows ? starter : command, windows ? [pidPath, command, ...args] : args, {
      cwd, env, detached: true, windowsHide: true, shell: false,
      stdio: ['ignore', stdoutFd, stderrFd],
    });

    if (windows) {
      // Node's detached mode has no console. The native helper gives the actual
      // worker a windowless console that its children inherit, then exits.
      const code = await new Promise((resolve, reject) => {
        child.once('error', reject);
        child.once('close', resolve);
      });
      if (code !== 0) throw new Error(`Failed to start ${command}; see ${stderrPath}.`);
      const pid = Number(await fsp.readFile(pidPath, 'utf8'));
      if (!Number.isSafeInteger(pid) || pid <= 0) throw new Error(`Invalid background PID for ${command}.`);
      return { pid };
    }

    await new Promise((resolve, reject) => {
      child.once('error', reject);
      child.once('spawn', resolve);
    });
    child.unref();
    return { pid: child.pid };
  } finally {
    if (stdoutFd !== undefined) fs.closeSync(stdoutFd);
    if (stderrFd !== undefined) fs.closeSync(stderrFd);
    if (pidPath) await fsp.rm(pidPath, { force: true });
  }
}
