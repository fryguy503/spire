import { expect, Page, test } from '@playwright/test';

type Session = {
  env?: 'desktop' | 'production';
  authEnabled?: boolean;
  user?: { id: number; is_admin: boolean } | null;
};

async function mockServer(page: Page, {
  env = 'desktop',
  authEnabled = true,
  user = { id: 2, is_admin: false },
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
  await page.goto('/admin/configuration/motd');
  await expect(page.getByRole('textbox')).toHaveValue('Welcome to the test server');
});
