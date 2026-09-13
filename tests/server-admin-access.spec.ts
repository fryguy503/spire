import { expect, Page, test } from '@playwright/test';

type Session = {
  env?: 'desktop' | 'production';
  authEnabled?: boolean;
  user?: { id: number; is_admin: boolean } | null;
  permissions?: { connection_id: number; read_all: boolean; write_all: boolean; read: string[]; write: string[] };
};

async function mockServer(page: Page, {
  env = 'desktop',
  authEnabled = true,
  user = { id: 2, is_admin: false },
  permissions = { connection_id: 1, read_all: true, write_all: true, read: [], write: [] },
}: Session = {}) {
  if (user) {
    await page.addInitScript(() => {
      localStorage.setItem(`spire-web-access-token-${location.host}`, 'test-session');
    });
  }

  // Broad handlers first: Playwright checks routes in reverse registration order.
  await page.route('**/api/v1/**', route => route.fulfill({ json: [] }));
  await page.routeWebSocket('**/api/v1/websocket*', () => {});
  await page.route('https://api.github.com/**', route => route.abort());
  await page.route('**/api/v1/app/env', route => route.fulfill({ json: { data: {
    env,
    version: '5.6.1',
    os: 'linux',
    is_spire_initialized: true,
    features: { github_auth_enabled: env === 'production' && authEnabled },
    settings: env === 'production' ? [] : [{ setting: 'AUTH_ENABLED', value: String(authEnabled) }],
  } } }));
  await page.route('**/api/v1/me', route => route.fulfill({
    json: user ? { ...user, user_name: 'Developer', provider: 'local' } : { error: 'User context not found' },
  }));
  await page.route('**/api/v1/connections', route => route.fulfill({ json: { data: [] } }));
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ json: permissions }));
  await page.route('**/api/v1/app/changelog', route => route.fulfill({ json: { data: '' } }));
  await page.route('**/api/v1/eqemuserver/server-stats', route => route.fulfill({ json: {
    server_name: '', zone_count: 0, players_online: 0, uptime: '', main_process_stats: [],
  } }));
  await page.route('**/api/v1/eqemuserver/get-lock-status', route => route.fulfill({ json: { locked: false } }));
  await page.route('**/api/v1/variables*', route => route.fulfill({
    json: [{ id: 1, varname: 'MOTD', value: 'Welcome to the test server' }],
  }));
}

test('a non-admin connection user can open Server Admin from Home', async ({ page }) => {
  await mockServer(page);
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Version (desktop) 5.6.1' })).toBeVisible();
  await page.getByRole('link', { name: 'Server Admin', exact: false }).click();
  await expect(page).toHaveURL(/\/admin\/?$/);
  await expect(page.getByText('Server Processes', { exact: true })).toBeVisible();
});

test('a non-admin connection user can open a nested admin route directly', async ({ page }) => {
  await mockServer(page);
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  await expect(page).toHaveURL(/\/admin\/configuration\/motd$/);
});

test('navigation search includes admin pages for a connection user', async ({ page }) => {
  await mockServer(page);
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Version (desktop) 5.6.1' })).toBeVisible();
  await page.getByRole('link', { name: 'Nav Search', exact: false }).click();
  await page.locator('input[placeholder="Where would you like to go?"]').fill('MOTD');
  await page.getByText('[Admin] [Configuration] MOTD', { exact: true }).click();
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
});

test('a denied write still displays the API permission error', async ({ page }) => {
  await mockServer(page);
  await page.route('**/api/v1/variable/1', route => route.fulfill({
    status: 403,
    json: { error: 'You do not have permission to write to this resource [/api/v1/variable/1]' },
  }));
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  await page.getByRole('textbox').fill('A denied update');
  const deniedWrite = page.waitForResponse(response => response.url().endsWith('/variable/1') && response.status() === 403);
  await page.getByRole('button', { name: /Save$/ }).click();
  await deniedWrite;
  await expect(page.getByText('You do not have permission to write to this resource', { exact: false })).toBeVisible();
  await expect(page.getByText('Message of the day updated!', { exact: true })).toBeHidden();
});

for (const env of ['desktop', 'production'] as const) {
  test(`unauthenticated users cannot load admin pages when ${env} auth is enabled`, async ({ page }) => {
    await mockServer(page, { env, user: null });
    const adminRequests: string[] = [];
    page.on('request', request => {
      if (request.url().includes('/api/v1/admin/') || request.url().includes('/api/v1/variables')) {
        adminRequests.push(request.url());
      }
    });
    await page.goto('/admin/configuration/motd');
    if (env === 'desktop') {
      await expect(page).toHaveURL(/\/login$/);
      await expect(page.getByRole('button', { name: /Spire Login$/ })).toBeVisible();
    } else {
      await expect(page).toHaveURL(/\/$/);
      await expect(page.getByText('Spire Changelog', { exact: true })).toBeVisible();
    }
    expect(adminRequests).toEqual([]);
  });
}

test('production still requires an instance administrator', async ({ page }) => {
  await mockServer(page, { env: 'production' });
  await page.goto('/admin/configuration/motd');
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByText('Spire Changelog', { exact: true })).toBeVisible();
});

test('instance administrators retain access', async ({ page }) => {
  await mockServer(page, { env: 'production', user: { id: 1, is_admin: true } });
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
});

test('local installations with authentication disabled retain access', async ({ page }) => {
  await mockServer(page, { authEnabled: false, user: null });
  let permissionRequests = 0;
  page.on('request', request => {
    if (request.url().endsWith('/permissions/me')) permissionRequests++;
  });
  await page.goto('/');
  await page.getByRole('link', { name: 'Server Admin', exact: false }).click();
  await expect(page.getByText('Server Processes', { exact: true })).toBeVisible();
  await expect(page.getByRole('link', { name: /Reloading \(Global\)/ })).toBeVisible();
  await expect(page.getByRole('link', { name: /Server Update$/ })).toBeVisible();
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  expect(permissionRequests).toBe(0);
});

const scopedPermissions = (read: string[] = [], write: string[] = []) => ({
  connection_id: 1, read_all: false, write_all: false, read, write,
});

test('mixed-case admin URLs cannot bypass the permission guard', async ({ page }) => {
  await mockServer(page, { permissions: scopedPermissions() });
  let variableRequests = 0;
  page.on('request', request => {
    if (request.url().includes('/api/v1/variables')) variableRequests++;
  });
  await page.goto('/ADMIN/configuration/motd');
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('textbox')).toHaveCount(0);
  expect(variableRequests).toBe(0);
});

test('permitted mixed-case admin URLs retain the admin menu', async ({ page }) => {
  await mockServer(page, { permissions: scopedPermissions(['variables']) });
  await page.goto('/ADMIN/CONFIGURATION/MOTD');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  await expect(page.getByRole('link', { name: /Editing Tools Home$/ })).toBeVisible();
  await expect(page.locator('#sidebar').getByText('Configuration', { exact: true })).toBeVisible();
});

for (const canReload of [false, true]) {
  test(`a rule save ${canReload ? 'with' : 'without'} reload permission respects the reload grant`, async ({ page }) => {
    await mockServer(page, { permissions: scopedPermissions(['rule_values'], ['rule_value', ...(canReload ? ['eqemuserver/reload'] : [])]) });
    const rule = { ruleset_id: 1, rule_name: 'Character:MaxLevel', rule_value: '65', notes: 'Maximum level' };
    await page.route('**/api/v1/rule_values*', route => route.fulfill({ json: [rule] }));
    await page.route('**/api/v1/rule_value/1*', async route => {
      rule.rule_value = '70';
      await route.fulfill({ json: rule });
    });
    let reloadRequests = 0;
    await page.route('**/api/v1/eqemuserver/reload/rules', route => {
      reloadRequests++;
      if (canReload) return route.fulfill({ json: { message: 'Reloaded' } });
      return route.fulfill({ status: 403, json: { error: 'Reload denied' } });
    });
    await page.goto('/admin/configuration/server-rules');
    const value = page.locator('tbody tr').filter({ hasText: 'Character:MaxLevel' }).getByRole('textbox');
    await value.fill('70');
    await value.press('Tab');
    if (canReload) {
      await expect(page.getByText('Server rules reloaded in-game!', { exact: true })).toBeVisible();
    } else {
      await expect(page.getByText('Updated rule (1) [Character:MaxLevel] to value (70)!', { exact: true })).toBeVisible();
    }
    await expect(page.getByText('Reload denied', { exact: true })).toHaveCount(0);
    expect(reloadRequests).toBe(canReload ? 1 : 0);
  });
}

for (const [name, permissions] of [
  ['no grants', scopedPermissions()],
  ['content-only grants', scopedPermissions(['items'], ['items'])],
  ['write-only grants', scopedPermissions([], ['variables'])],
] as const) {
  test(`${name} hide Server Admin in navigation and search and deny direct entry`, async ({ page }) => {
    await mockServer(page, { permissions });
    const requests: string[] = [];
    page.on('request', request => {
      if (/\/api\/v1\/(admin\/|eqemuserver\/|variables)/.test(request.url())) requests.push(request.url());
    });
    await page.goto('/admin/configuration/motd');
    await expect(page).toHaveURL(/\/$/);
    await expect(page.getByRole('heading', { name: 'Version (desktop) 5.6.1' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
    await expect.poll(() => page.locator('ninja-keys').evaluate((element: any) =>
      element.data.filter((entry: any) => entry.title.startsWith('[Admin]')).length
    )).toBe(0);
    expect(requests).toEqual([]);
  });
}

test('MOTD-only access shows only MOTD, hides empty groups, and skips unrelated header requests', async ({ page }) => {
  await mockServer(page, { permissions: scopedPermissions(['variable', 'variables']) });
  const adminRequests: string[] = [];
  page.on('request', request => {
    if (/\/api\/v1\/(admin\/|eqemuserver\/)/.test(request.url())) adminRequests.push(request.url());
  });
  await page.goto('/');
  await page.getByRole('link', { name: 'Server Admin', exact: false }).click();
  await expect(page).toHaveURL(/\/admin\/configuration\/motd$/);
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  const sidebar = page.locator('#sidebar');
  await expect(sidebar.getByText('Configuration', { exact: true })).toBeVisible();
  await expect(sidebar.getByRole('link', { name: /MOTD$/ })).toBeVisible();
  for (const label of ['Players Online', 'Player Operations', 'Database', 'Logs', 'Zone Servers', 'Server Update', 'Server Config']) {
    await expect(sidebar.getByText(label, { exact: true })).toHaveCount(0);
  }
  const adminLinks = await sidebar.locator('a[href^="/admin"]').evaluateAll(links => links.map(link => link.getAttribute('href')));
  expect(adminLinks.length).toBeGreaterThan(0);
  expect(adminLinks.every(href => href === '/admin/configuration/motd')).toBe(true);
  await expect(page.locator('.admin-header-window')).toHaveCount(0);
  const search = await page.locator('ninja-keys').evaluate((element: any) =>
    element.data.filter((entry: any) => entry.title.startsWith('[Admin]')).map((entry: any) => entry.title)
  );
  expect(search.sort()).toEqual(['[Admin] Server Admin', '[Admin] [Configuration] MOTD']);
  expect(adminRequests).toEqual([]);
});

test('a forbidden nested route redirects to the first permitted tool before fetching its data', async ({ page }) => {
  await mockServer(page, { permissions: scopedPermissions(['variables']) });
  let forbiddenRequests = 0;
  page.on('request', request => {
    if (request.url().includes('/admin/serverconfig')) forbiddenRequests++;
  });
  await page.goto('/admin/configuration/server?s=Database');
  await expect(page).toHaveURL(/\/admin\/configuration\/motd$/);
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  expect(forbiddenRequests).toBe(0);
});

test('process-stats-only access renders that widget without unrelated monitoring or process controls', async ({ page }) => {
  await mockServer(page, { permissions: scopedPermissions(['eqemuserver/server-stats']) });
  const requests: string[] = [];
  page.on('request', request => {
    if (/\/api\/v1\/(admin\/|eqemuserver\/)/.test(request.url())) requests.push(new URL(request.url()).pathname);
  });
  await page.goto('/admin');
  await expect(page.getByText('Server Processes', { exact: true })).toBeVisible();
  await expect(page.locator('.admin-header-window')).toBeVisible();
  await expect(page.getByRole('button', { name: /Start Server|Restart|Stop Server/ })).toHaveCount(0);
  await expect(page.locator('[data-testid="admin-host-metrics"]')).toHaveCount(0);
  expect(requests.length).toBeGreaterThan(0);
  expect(requests.every(path => path === '/api/v1/eqemuserver/server-stats')).toBe(true);
});

test('read-only ALL grants hide command-only tools and server mutation controls', async ({ page }) => {
  await mockServer(page, { permissions: { ...scopedPermissions(), read_all: true } });
  await page.goto('/admin');
  await expect(page.getByText('Server Processes', { exact: true })).toBeVisible();
  const adminSearch = await page.locator('ninja-keys').evaluate((element: any) =>
    element.data.filter((entry: any) => entry.title.startsWith('[Admin]')).map((entry: any) => entry.title)
  );
  expect(adminSearch).toContain('[Admin] [Configuration] MOTD');
  expect(adminSearch).not.toContain('[Admin] [Database] Database Backups');
  expect(adminSearch).not.toContain('[Admin] Reloading (Global)');
  await expect(page.locator('.admin-header-window a').filter({ hasText: 'Unlocked' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: /Start Server|Restart|Stop Server/ })).toHaveCount(0);
});

test('instance-admin status alone does not grant access to connection tools', async ({ page }) => {
  await mockServer(page, { user: { id: 1, is_admin: true }, permissions: scopedPermissions() });
  await page.goto('/admin');
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
});

test('permission lookup failures fail closed', async ({ page }) => {
  await mockServer(page);
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ status: 503, json: { error: 'Unavailable' } }));
  await page.goto('/admin');
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
});

test('malformed permission responses fail closed', async ({ page }) => {
  await mockServer(page);
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ json: { read_all: 'false', read: '*' } }));
  await page.goto('/admin');
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
});

for (const tool of [
  { path: '/admin/log-settings', resource: 'logsys_categories', title: 'Log Settings' },
  { path: '/admin/player-event-logs/settings', resource: 'player_event_log_settings', title: 'Log Settings' },
  { path: '/admin/player-event-logs/explorer', resource: 'player_event_logs', title: 'Player Event Log Explorer' },
  { path: '/admin/zones', resource: 'eqemuserver/zoneserver-list', title: 'Zone Servers' },
]) {
  test(`${tool.resource} access does not expose optional tools or fetch their data`, async ({ page }) => {
    await mockServer(page, { permissions: scopedPermissions([tool.resource]) });
    if (tool.path === '/admin/zones') {
      await page.route('**/api/v1/eqemuserver/zoneserver-list', route => route.fulfill({ json: [{
        id: 1, zone_id: 1, zone_name: 'qeynos', number_players: 0, clients: [],
        zone_os_pid: 1234, zone_server_address: '127.0.0.1', client_port: 7000, cpu: 0,
      }] }));
    }
    const forbidden: string[] = [];
    page.on('request', request => {
      if (/\/api\/v1\/(discord_webhooks|guilds|eqemuserver\/player-event-logs\/etl-settings)/.test(request.url())) {
        forbidden.push(request.url());
      }
    });
    await page.goto(tool.path);
    await expect(page).toHaveURL(new RegExp(`${tool.path}$`));
    await expect(page.locator('.main-content').getByText(tool.title, { exact: true }).first()).toBeVisible();
    await expect(page.locator('.main-content').getByRole('link', { name: /Discord Webhook/ })).toHaveCount(0);
    if (tool.path === '/admin/zones') {
      await expect(page.getByTitle('Logs', { exact: true })).toHaveCount(0);
      await expect(page.getByTitle('Kill Zone', { exact: true })).toHaveCount(0);
    }
    expect(forbidden).toEqual([]);
  });
}

test('navigation refreshes permissions after grants change on the active connection', async ({ page }) => {
  await mockServer(page);
  let permissions = scopedPermissions(['variables']);
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ json: permissions }));
  await page.goto('/admin');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  permissions = scopedPermissions(['items']);
  await page.getByRole('link', { name: 'Editing Tools Home', exact: false }).click();
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
  await expect.poll(() => page.locator('ninja-keys').evaluate((element: any) =>
    element.data.filter((entry: any) => entry.title.startsWith('[Admin]')).length
  )).toBe(0);
});

test('revoking the current tool does not remount it during navigation', async ({ page }) => {
  await mockServer(page);
  let permissions = scopedPermissions(['variables']);
  let variableRequests = 0;
  page.on('request', request => {
    if (request.url().includes('/api/v1/variables')) variableRequests++;
  });
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ json: permissions }));
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
  const before = variableRequests;
  permissions = scopedPermissions();
  await page.getByRole('link', { name: /Editing Tools Home$/ }).click();
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole('link', { name: /Server Admin$/ })).toHaveCount(0);
  expect(variableRequests).toBe(before);
});

test('switching connections refreshes the menu without navigating away', async ({ page }) => {
  await mockServer(page);
  let active = 1;
  const connection = (id: number) => ({
    id, active: active === id ? 1 : 0, server_database_connection_id: id,
    database_connection: {
      id, name: `QA connection ${id}`, created_by: 1, db_host: 'localhost', db_port: 3306,
      db_name: 'qa', db_username: 'qa', content_db_username: '', user_server_database_connections: [],
    },
  });
  await page.route('**/api/v1/connections', route => route.fulfill({ json: { data: [connection(1), connection(2)] } }));
  await page.route('**/api/v1/connection-check/*', route => route.fulfill({ json: { data: { message: 'Online' } } }));
  await page.route('**/api/v1/permissions/me', route => route.fulfill({ json: {
    ...scopedPermissions(active === 1 ? ['variables'] : ['items']), connection_id: active,
  } }));
  await page.route('**/api/v1/connection/*/set-active', route => {
    active = Number(new URL(route.request().url()).pathname.split('/')[4]);
    return route.fulfill({ json: { data: true } });
  });
  await page.goto('/connections');
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toBeVisible();
  await page.getByText('Set Active', { exact: true }).click();
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toHaveCount(0);
  await expect(page).toHaveURL(/\/connections$/);
  await page.getByText('Set Active', { exact: true }).click();
  await expect(page.getByRole('link', { name: 'Server Admin', exact: false })).toBeVisible();
});
