import { expect, test, Page } from '@playwright/test';
import inventory from '../internal/http/staticmaps/race-inventory-map.json';

async function mockViewer(page: Page, options: { inventoryFails?: boolean; previewsFail?: boolean } = {}) {
  // Narrow overrides follow the fallback because route handlers resolve LIFO.
  await page.route('**/api/v1/**', route => route.fulfill({ json: [] }));
  await page.route('**/api/v1/app/env**', route => route.fulfill({ json: { data: {
    is_spire_initialized: true, env: 'local', version: 'test', features: {}, settings: [], os: 'darwin',
  } } }));
  await page.route('**/api/v1/static-map/race-inventory-map.json', route => options.inventoryFails
    ? route.fulfill({ status: 503, json: { error: 'Unavailable' } })
    : route.fulfill({ json: inventory }));
  await page.route('**/api/v1/static-map/npc-models-map.json', route => options.previewsFail
    ? route.fulfill({ status: 503, json: { error: 'Unavailable' } })
    : route.fulfill({ json: [{ contents: [
      { name: 'CTN_1_0_0_0.png' }, { name: 'CTN_1_1_0_0.png' },
      { name: 'CTN_13_2_0_0.png' }, { name: 'CTN_13_2_1_0.png' },
    ] }] }));
}

test('model code search reveals archive-specific local and imported zones', async ({ page }) => {
  await mockViewer(page);
  await page.goto('/race-viewer');
  await page.getByLabel('Search models', { exact: true }).fill('AVI');
  await page.getByRole('button', { name: 'Aviak, race 13. View model locations', exact: true }).click();
  const inspector = page.getByTestId('race-inspector');
  await expect(inspector.getByRole('heading', { name: 'Aviak', exact: true })).toBeVisible();
  await expect(inspector.getByRole('button', { name: 'Copy neutral model code AVI' })).toBeVisible();
  const source = inspector.locator('.model-source').filter({ hasText: 'southkarana_chr.s3d' });
  await expect(source).toContainText('AVI Neutral');
  await expect(source.locator('li').filter({ hasText: 'southkarana' })).toContainText('Local');
  await expect(source.locator('li').filter({ hasText: 'freeportwest' })).toContainText('Imported');
  await expect(page).toHaveURL(/race=13/);
});

test('zone filtering includes globals optionally and composes with text search', async ({ page }) => {
  await mockViewer(page);
  await page.goto('/race-viewer?zoneSearch=14&race=13');
  const inspector = page.getByTestId('race-inspector');
  await expect(inspector).toContainText('southkarana_chr.s3d');
  await expect(inspector).not.toContainText('freportw_chr.s3d');
  await expect(page.getByRole('button', { name: 'Human, race 1. View model locations', exact: true })).toBeVisible();
  await page.getByLabel('Include global models').uncheck();
  await expect(page.getByRole('button', { name: 'Human, race 1. View model locations', exact: true })).toHaveCount(0);
  await page.getByLabel('Search models', { exact: true }).fill('AVI');
  await expect(page.getByRole('button', { name: 'Aviak, race 13. View model locations', exact: true })).toBeVisible();
  await page.getByLabel('Search models', { exact: true }).fill('no-such-model');
  await expect(page.getByRole('heading', { name: 'No matching races' })).toBeVisible();
  await expect(inspector).toHaveCount(0);
});

test('zone links and browser history restore the same integrated selection', async ({ page }) => {
  await mockViewer(page);
  await page.goto('/race-viewer?raceSearch=AVI&race=13');
  await page.getByRole('button', { name: 'Find models in South Karana (southkarana)', exact: true }).click();
  await expect(page.getByLabel('Available in zone')).toHaveValue('14');
  await expect(page.getByLabel('Search models', { exact: true })).toHaveValue('');
  await expect(page.getByTestId('race-inspector').getByRole('heading', { name: 'Aviak', exact: true })).toBeVisible();
  await page.goBack();
  await expect(page.getByLabel('Available in zone')).toHaveValue('0');
  await expect(page.getByLabel('Search models', { exact: true })).toHaveValue('AVI');
  await expect(page.getByTestId('race-inspector')).toContainText('freportw_chr.s3d');
});

test('races beyond the preview range and archives without zones remain discoverable', async ({ page }) => {
  await mockViewer(page, { previewsFail: true });
  await page.goto('/race-viewer?raceSearch=LUC&race=724');
  const inspector = page.getByTestId('race-inspector');
  await expect(inspector.getByRole('heading', { name: 'Luclin', exact: true })).toBeVisible();
  await expect(inspector).toContainText('luc.eqg');
  await expect(inspector).toContainText('Plane of Shadow');
  await expect(page.getByText('Image previews are unavailable.', { exact: false })).toBeVisible();
  await page.getByLabel('Search models', { exact: true }).fill('PG3');
  await expect(inspector).toContainText('pg3.eqg');
  await expect(inspector).toContainText('No zones listed for this archive.');
  await expect(inspector.locator('.availability-tag').filter({ hasText: 'Global' })).toHaveCount(0);
  await page.getByLabel('Search models', { exact: true }).fill('2253');
  await expect(inspector).toContainText('PPOINT');
  await expect(inspector).toContainText('No source archives or zone locations are listed');
});

test('failed inventory requests expose retry while existing previews remain usable', async ({ page }) => {
  await mockViewer(page, { inventoryFails: true });
  await page.goto('/race-viewer');
  await expect(page.getByRole('alert')).toContainText('Model locations could not be loaded');
  await expect(page.getByRole('button', { name: 'Human, race 1. View model locations', exact: true })).toBeVisible();
  await page.route('**/api/v1/static-map/race-inventory-map.json', route => route.fulfill({ json: inventory }));
  await page.getByRole('button', { name: 'Retry model inventory' }).click();
  await expect(page.getByRole('alert')).toHaveCount(0);
  await expect(page.getByTestId('race-inspector')).toContainText('globalhum_chr.s3d');
});

test('mobile layout fits without horizontal overflow and preserves model locations', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await mockViewer(page);
  await page.goto('/race-viewer?raceSearch=AVI&race=13');
  await expect(page.getByTestId('race-inspector').getByRole('heading', { name: 'Aviak', exact: true })).toBeVisible();
  await expect(page.getByTestId('race-inspector')).toContainText('southkarana_chr.s3d');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBeTruthy();
  await page.screenshot({ path: test.info().outputPath('race-viewer-mobile.png'), fullPage: true });
  await page.getByRole('button', { name: 'Aviak, race 558. View model locations', exact: true }).click();
  await expect(page.getByTestId('race-inspector')).toContainText('RACE 558');
  await expect(page.getByTestId('race-inspector').getByRole('heading', { name: 'Aviak', exact: true })).toBeInViewport();
});

test('location copy includes the gender model and archive name without its extension', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await mockViewer(page);
  await page.goto('/race-viewer?raceSearch=clock&race=88');
  const male = page.getByRole('button', { name: 'Copy male model location CLM,steamfont_chr', exact: true });
  await male.click();
  await expect(male).toContainText('Copied');
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe('CLM,steamfont_chr');
  const female = page.getByRole('button', { name: 'Copy female model location CLF,steamfont_chr', exact: true });
  await female.click();
  await expect(female).toContainText('Copied');
  await expect(male).not.toContainText('Copied');
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe('CLF,steamfont_chr');
  await page.getByLabel('Search models', { exact: true }).fill('LUC');
  await page.getByRole('button', { name: 'Luclin, race 724. View model locations', exact: true }).click();
  await page.getByRole('button', { name: 'Copy female model location LUC,luc', exact: true }).click();
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe('LUC,luc');
});

test('clipboard failure provides the exact value at the location', async ({ page }) => {
  await mockViewer(page);
  await page.addInitScript(() => {
    Object.defineProperty(navigator, 'clipboard', { value: {
      writeText: async () => { throw new Error('Clipboard unavailable'); },
    } });
  });
  await page.goto('/race-viewer?raceSearch=clock&race=88');
  const source = page.locator('.model-source').filter({ hasText: 'steamfont_chr.s3d' });
  await source.getByRole('button', { name: 'Copy male model location CLM,steamfont_chr', exact: true }).click();
  await expect(source.getByRole('alert')).toHaveText('Could not copy. Select: CLM,steamfont_chr');
  await expect(source).not.toContainText('Copied');
});

test('locations is the default and classic gallery preserves filters and selection', async ({ page }) => {
  await mockViewer(page);
  await page.goto('/race-viewer?raceSearch=AVI&zoneSearch=14&race=13');
  await expect(page.getByRole('button', { name: 'Model locations', exact: true })).toHaveAttribute('aria-pressed', 'true');
  await expect(page.getByTestId('race-inspector').getByRole('img')).toHaveCount(1);
  await page.getByRole('button', { name: 'Classic gallery', exact: true }).click();
  await expect(page).toHaveURL(/view=classic/);
  await expect(page.getByTestId('race-inspector')).toHaveCount(0);
  await expect(page.getByTestId('classic-gallery').getByRole('img')).toHaveCount(2);
  await page.reload();
  await expect(page.getByRole('button', { name: 'Classic gallery', exact: true })).toHaveAttribute('aria-pressed', 'true');
  await expect(page.getByLabel('Search models', { exact: true })).toHaveValue('AVI');
  await expect(page.getByLabel('Available in zone')).toHaveValue('14');
  await page.getByRole('button', { name: 'Model locations', exact: true }).click();
  await expect(page.getByTestId('race-inspector')).toContainText('RACE 13');
  await expect(page).not.toHaveURL(/view=classic/);
  await page.goBack();
  await expect(page.getByTestId('classic-gallery')).toBeVisible();
  await page.setViewportSize({ width: 390, height: 844 });
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBeTruthy();
});

test('desktop panes fill the remaining viewport and scroll independently', async ({ page }) => {
  await mockViewer(page);
  await page.setViewportSize({ width: 1347, height: 940 });
  await page.goto('/race-viewer?race=88');
  await expect(page.getByTestId('race-inspector')).toContainText('Clockwork Gnome');
  for (const viewport of [{ width: 1347, height: 940 }, { width: 1024, height: 768 }, { width: 1347, height: 700 }]) {
    await page.setViewportSize(viewport);
    await expect.poll(() => page.evaluate(() => {
      const gallery = document.querySelector('.race-gallery-window')!.getBoundingClientRect();
      const inspector = document.querySelector('.race-inspector')!.getBoundingClientRect();
      return Math.abs(gallery.bottom - inspector.bottom) < 1 &&
        Math.abs(gallery.height - inspector.height) < 1 &&
        Math.abs(gallery.bottom - (window.innerHeight - 24)) < 2 &&
        document.documentElement.scrollHeight <= window.innerHeight;
    })).toBeTruthy();
  }
  for (const selector of ['.race-gallery', '.inspector-scroll']) {
    await page.locator(selector).hover();
    await page.mouse.wheel(0, 600);
    await expect.poll(() => page.locator(selector).evaluate(element => element.scrollTop)).toBeGreaterThan(0);
    expect(await page.evaluate(() => window.scrollY)).toBe(0);
  }
});
