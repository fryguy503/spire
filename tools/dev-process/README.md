# Run development processes in the background

This launcher keeps Windows console processes and their children windowless.
It preserves arguments, the working directory, environment, logs, and child
IPC. On Linux and macOS it starts a detached process group.

Use Node.js 22 or newer. Windows also needs Go from Spire's development setup.
The Windows helper is compiled once into the log directory and rebuilt when
its source changes. It exits after starting the worker, so there is no resident
supervisor or polling loop.

From the repository root, start the frontend after installing its dependencies:

```powershell
$env:NODE_OPTIONS = '--openssl-legacy-provider'
npm run dev:background -- --cwd frontend --logs tmp/dev-frontend -- node --max_old_space_size=4096 --stack-size=10000 node_modules/@vue/cli-service/bin/vue-cli-service.js serve --host=127.0.0.1
```

To start an already built backend with your existing database configuration:

```powershell
npm run dev:background -- --logs tmp/dev-backend -- .\spire.exe http:serve --port=3010
```

The final JSON line contains the worker's `pid`, `stdoutPath`, and `stderrPath`.
Each service needs its own log directory. Relative working and log directories
are resolved from the directory where you run the launcher. Output is appended
to the logs. A returned PID confirms process creation; use the service's health
endpoint or logs to check readiness.

On Windows, stop that process and its children with its reported PID:

```powershell
taskkill /PID <pid> /T /F
```

On Linux or macOS, stop the reported process group with `kill -TERM -- -<pid>`.
Use the PID from the current launch, since operating systems can reuse PIDs.

Commands are launched as native executables without an implicit shell. For
Node programs, invoke `node` with the entry script, as above. To execute a shell
script, pass the shell executable explicitly and use its argument syntax.

## Use from another launcher

Import `startLoggedProcess` from `tools/dev-process/start.mjs` and await it with
`command`, `args`, `cwd`, `stdoutPath`, and `stderrPath`. An optional `env` replaces
the inherited environment; spread `process.env` when adding overrides. The
result is `{ pid }`. Record that PID immediately so a later startup failure can
clean up already started services.

## Validate the launcher

```sh
npm run test:dev-process
```

The tests launch real parent and child processes, then let their original
launcher exit. They check argument quoting, environment and working-directory
preservation, IPC, logs, shell grandchildren, startup failures, and complete
process-tree shutdown. Concurrent CLI launches exercise the helper cache and
log capture. Windows additionally checks that both generations have
windowless consoles and that a failed PID handoff starts no worker. All test
processes and temporary files are cleaned up.
