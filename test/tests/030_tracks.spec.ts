import {expect, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {trackSelectors} from '../selectors/track-selectors';

test('should load and display all available tracks and their information', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click()
  await expect(page.getByTestId(globalSelectors.tracksTab)).toContainText("Strecken")

  let trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2)
  await expect(trackNodes.nth(0)).toContainText('Höchberg')
  await expect(trackNodes.nth(1)).toContainText('Dummy')

  await page.getByTestId(trackSelectors.toggler).nth(0).click()
  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(5)
  await expect(trackNodes.nth(1)).toContainText('Farmrunde')
  await expect(trackNodes.nth(2)).toContainText('Stadtrunde')

  await trackNodes.nth(1).click()

  await expect(page.getByLabel('Streckenname')).toHaveValue('Farmrunde')
  await expect(page.getByLabel('Verwendungen')).toHaveValue('2')
  await expect(page.getByLabel('Streckenlänge')).toHaveValue('10,8')
  await expect(page.getByLabel('Kommentar')).toHaveValue('bei Regen sehr matschig')
  await expect(page.getByLabel('Speichern')).toBeDisabled()
});

test('should show usages correctly', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click()
  await expect(page.getByTestId(globalSelectors.tracksTab)).toContainText("Strecken")

  let trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2)
  await expect(trackNodes.nth(0)).toContainText('Höchberg')
  await expect(trackNodes.nth(1)).toContainText('Dummy')

  await page.getByTestId(trackSelectors.toggler).nth(0).click()
  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(5)
  await expect(trackNodes.nth(1)).toContainText('Farmrunde')
  await expect(trackNodes.nth(2)).toContainText('Stadtrunde')

  await trackNodes.nth(2).click()

  await expect(page.getByLabel('Streckenname')).toHaveValue('Stadtrunde')
  await expect(page.getByLabel('Verwendungen')).toHaveValue('2')
  await expect(page.getByLabel('Streckenlänge')).toHaveValue('10,3')
  await expect(page.getByLabel('Speichern')).toBeDisabled()
})
