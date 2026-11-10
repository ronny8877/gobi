import { test, expect } from '@playwright/test';
import { mkdirSync } from 'node:fs';
const capture = process.env.GOBI_CAPTURE_DIR;
if (capture) mkdirSync(capture, { recursive: true });
test('guides render, navigate, and filter without broken links', async ({
  page,
}) => {
  await page.goto('/docs/quick-start/');
  await expect(
    page.getByRole('heading', { name: 'Quick start', exact: true }),
  ).toBeVisible();
  await page.getByLabel('Find a guide').fill('array');
  await expect(page.locator('[data-guide-link]:visible')).toHaveCount(1);
  await page.getByRole('link', { name: 'Arrays', exact: true }).first().click();
  await expect(
    page.getByRole('heading', { name: 'Arrays', exact: true }),
  ).toBeVisible();
  for (const slug of [
    'configuration',
    'routes',
    'generators',
    'references',
    'validation-auth',
    'pocketbase',
    'deployment',
  ]) {
    await page.goto(`/docs/${slug}/`);
    await expect(page.locator('article h1')).toBeVisible();
    await expect(page.locator('article pre').first()).toBeVisible();
  }
});
test('playground switches schemas and reports live failures', async ({
  page,
}) => {
  await page.goto('/playground/');
  await page.getByRole('button', { name: 'References', exact: true }).click();
  await expect(page.locator('pre').first()).toContainText('Ref(ref.profile)');
  await page.getByRole('button', { name: 'HTTP methods', exact: true }).click();
  await expect(page.getByLabel('Request body (JSON)')).toBeVisible();
  await page.getByLabel('Response source').selectOption('live');
  await page.getByLabel('Gobi server / proxy URL').fill('not-a-url');
  await page.getByRole('button', { name: 'Send request', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText(
    'Check that the Gobi demo server',
  );
});
test('builder exports and validates a runnable schema', async ({ page }) => {
  await page.goto('/builder/');
  await page.getByLabel('Route path').fill('/team');
  await page.getByRole('button', { name: 'Add field', exact: true }).click();
  await page.getByLabel('Field name', { exact: true }).last().fill('owner');
  await page
    .getByLabel('Value type', { exact: true })
    .last()
    .selectOption('reference');
  await expect(page.locator('pre')).toContainText('Ref(ref.profile)');
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download api.gobi.json' }).click();
  const download = await downloadPromise;
  expect(download.suggestedFilename()).toBe('api.gobi.json');
  const stream = await download.createReadStream();
  const chunks: Buffer[] = [];
  for await (const chunk of stream!) chunks.push(chunk);
  const config = JSON.parse(Buffer.concat(chunks).toString());
  expect(config.api[0].path).toBe('/team');
  expect(config.ref.profile.email).toBe('User(email)');
  await page.getByLabel('Route path').fill('invalid');
  await expect(page.getByRole('alert')).toContainText('Route path must start');
  await expect(
    page.getByRole('button', { name: 'Download api.gobi.json' }),
  ).toBeDisabled();
  await page.getByRole('button', { name: 'Reset starter' }).click();
  await expect(page.getByLabel('Route path')).toHaveValue('/users');
});
test('dashboard search and schema reveal work', async ({ page }) => {
  await page.goto('/dashboard/');
  await expect(page.getByRole('status')).toHaveText('Sample fixture');
  await page.getByLabel('Search customers').fill('no-such-person');
  await expect(page.getByText('No customers match')).toBeVisible();
  await page.getByLabel('Search customers').fill('');
  await page.getByRole('button', { name: 'View the schema' }).click();
  await expect(page.locator('pre')).toContainText('Ref(ref.customer)');
});
test('live API drives arrays, refs, methods and dashboard', async ({
  page,
}) => {
  test.skip(!process.env.GOBI_E2E_LIVE, 'Requires bundled Gobi API');
  await page.goto('/playground/');
  await page.getByLabel('Response source').selectOption('live');
  for (const label of ['Arrays', 'References', 'Generators', 'HTTP methods']) {
    await page.getByRole('button', { name: label, exact: true }).click();
    await page
      .getByRole('button', { name: 'Send request', exact: true })
      .click();
    await expect(page.getByRole('status')).toHaveText('HTTP 200');
  }
  await page.goto('/dashboard/');
  await page.getByLabel('Data source').selectOption('live');
  await page.getByRole('button', { name: 'Refresh data' }).click();
  await expect(page.getByRole('status')).toHaveText('Live Gobi response');
  await expect(page.locator('.customer-name')).toHaveCount(5);
});
test('desktop and mobile layouts remain usable', async ({ page }) => {
  test.setTimeout(120000);
  const errors: string[] = [];
  page.on('pageerror', (error) => errors.push(error.message));
  for (const viewport of [
    { width: 1440, height: 1050 },
    { width: 390, height: 844 },
  ]) {
    await page.setViewportSize(viewport);
    for (const [route, name] of [
      ['/', 'home'],
      ['/docs/quick-start/', 'guide'],
      ['/playground/', 'playground'],
      ['/dashboard/', 'dashboard'],
      ['/builder/', 'builder'],
    ]) {
      await page.goto(route);
      await page.evaluate(() => document.fonts.ready);
      await expect(page.locator('main h1')).toBeVisible();
      const overflow = await page.evaluate(
        () => document.documentElement.scrollWidth > innerWidth + 1,
      );
      expect(overflow, `${name} overflows ${viewport.width}`).toBe(false);
      if (capture)
        await page.screenshot({
          path: `${capture}/${name}-${viewport.width}.png`,
          fullPage: true,
        });
    }
    if (viewport.width === 390) {
      await page.getByText('Menu', { exact: true }).click();
      await expect(
        page.getByRole('navigation', { name: 'Mobile navigation' }),
      ).toBeVisible();
    }
  }
  expect(errors).toEqual([]);
});
