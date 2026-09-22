import {expect, Page, test} from '@playwright/test';

async function openKinbound(page: Page, overrides: Record<string, any> = {}) {
  let state = {supported: true, id: 100, kinbound: 0, nodrop: 0, notransfer: 0, epic: false, can_always: true, ...overrides};
  const patches: any[] = [];
  let conflict = false;
  await page.route('**/api/v1/**', route => route.request().isNavigationRequest() ? route.continue() : route.fulfill({json: []}));
  await page.route('https://api.github.com/**', route => route.abort());
  await page.route('**/api/v1/app/env**', route => route.fulfill({json: {data: {
    is_spire_initialized: true, env: 'local', version: '1.0.0', features: {}, settings: [], os: 'windows'
  }}}));
  await page.route('**/api/v1/item/100**', route => route.fulfill({json: {
    id: 100, name: 'Kinbound Blade', lore: '', color: 0, itemtype: 0, itemclass: 0,
    nodrop: 0, notransfer: 0, races: 65535, classes: 65535, slots: 0,
    clickeffect: -1, proceffect: -1, worneffect: -1, focuseffect: -1, scrolleffect: -1, bardeffect: -1,
    evoitem: 0, evoid: 0, evolvinglevel: 0, evomax: 0
  }}));
  await page.route('**/api/v1/item-kinbound/100', route => {
    if (route.request().method() === 'PATCH') {
      const patch = route.request().postDataJSON();
      patches.push(patch);
      if (conflict) return route.fulfill({status: 409, json: {error: 'The Kinbound setting changed since it was loaded. Refresh before saving.'}});
      state = {...state, kinbound: patch.kinbound};
    }
    return route.fulfill({json: state});
  });
  await page.goto('/');
  await page.waitForSelector('#sidebar');
  await page.evaluate(() => {
    window.history.pushState({}, '', '/item/100');
    window.dispatchEvent(new PopStateEvent('popstate'));
  });
  await page.getByRole('listitem').filter({hasText: 'Kinbound'}).click();
  return {patches, setConflict: () => { conflict = true; }, setState: (next: Record<string, any>) => {state = {...state, ...next};}};
}

test('Kinbound saves explicit selections separately and preserves item dirty fields', async ({page}) => {
  const state = await openKinbound(page);
  await expect(page.locator('#kinbound')).toHaveValue('0');
  await page.locator('#kinbound').selectOption('1');
  await expect(page.locator('#kinbound')).toHaveClass(/pulsate-highlight-modified/);
  await expect(page.getByText('Unsaved Kinbound selection')).toBeVisible();
  await page.evaluate(() => document.querySelector('#id')!.classList.add('pulsate-highlight-modified'));
  await page.locator('#save-item-kinbound').click();
  await expect.poll(() => state.patches).toEqual([{kinbound: 1, expected_kinbound: 0}]);
  await expect(page.locator('#save-item-kinbound')).toBeDisabled();
  await expect(page.locator('#kinbound')).not.toHaveClass(/pulsate-highlight-modified/);
  await expect(page.locator('#id')).toHaveClass(/pulsate-highlight-modified/);
  await page.locator('#kinbound').selectOption('2');
  await page.locator('#save-item-kinbound').click();
  await expect.poll(() => state.patches[1]).toEqual({kinbound: 2, expected_kinbound: 1});
  await page.locator('#kinbound').selectOption('0');
  await page.locator('#save-item-kinbound').click();
  await expect.poll(() => state.patches[2]).toEqual({kinbound: 0, expected_kinbound: 2});
});

test('Kinbound reports absent extension schema without blocking the item editor', async ({page}) => {
  const state = await openKinbound(page, {supported: false, message: 'Kinbound is unavailable on this database. Install migrations 9397 and 9398.'});
  await expect(page.getByText('Kinbound is unavailable on this database.', {exact: false})).toBeVisible();
  await expect(page.locator('#kinbound')).toHaveCount(0);
  await expect(page.getByRole('button', {name: /Save Item/})).toBeVisible();
  expect(state.patches).toHaveLength(0);
});

test('Kinbound blocks Always for excluded items while allowing Never', async ({page}) => {
  const state = await openKinbound(page, {can_always: false, epic: true, blocked_reason: 'This item has a known epic exclusion in the Kinbound catalog.'});
  await expect(page.locator('#kinbound option[value="1"]')).toBeDisabled();
  await expect(page.getByText('This item has a known epic exclusion', {exact: false})).toBeVisible();
  await page.locator('#kinbound').selectOption('2');
  await page.locator('#save-item-kinbound').click();
  await expect.poll(() => state.patches).toEqual([{kinbound: 2, expected_kinbound: 0}]);
});

test('Kinbound stale save keeps the selection dirty until refresh', async ({page}) => {
  const state = await openKinbound(page);
  state.setConflict();
  await page.locator('#kinbound').selectOption('1');
  await page.locator('#save-item-kinbound').click();
  await expect(page.getByText('The Kinbound setting changed since it was loaded.', {exact: false})).toBeVisible();
  await expect(page.getByText('Unsaved Kinbound selection')).toBeVisible();
  state.setState({kinbound: 2});
  await page.getByRole('button', {name: 'Refresh Kinbound'}).click();
  await expect(page.locator('#kinbound')).toHaveValue('2');
  await expect(page.locator('#save-item-kinbound')).toBeDisabled();
});
