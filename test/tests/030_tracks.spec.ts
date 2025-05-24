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

test('should clone a track, verify it exists, delete it, and verify it was deleted', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click();
  await expect(page.getByTestId(globalSelectors.tracksTab)).toContainText("Strecken");

  await page.getByTestId(trackSelectors.toggler).nth(0).click();
  const trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(5);
  await trackNodes.nth(1).click(); // Select "Farmrunde"

  await expect(page.getByLabel('Streckenname')).toHaveValue('Farmrunde');
  await expect(page.getByLabel('Kommentar')).toHaveValue('bei Regen sehr matschig');

  await page.getByLabel('Klonen').click();

  await expect(page.url()).toContain('/tracks/new');

  await expect(page.getByLabel('Streckenname')).toHaveValue('');
  await expect(page.getByLabel('Kommentar')).toHaveValue('bei Regen sehr matschig');

  await page.getByLabel('Streckenname').fill('Cloned Farmrunde');
  await page.getByLabel('Speichern').click();
  await page.getByTestId(trackSelectors.createTrackButton).click();

  await page.getByTestId(globalSelectors.tracksTab).click();
  await expect(trackNodes).toHaveCount(3); //because it was closed after saving the clone

  const clonedTrack = page.getByTestId(trackSelectors.trackNode).filter({ hasText: 'Cloned Farmrunde' });

  await expect(clonedTrack).toBeVisible();

  await clonedTrack.click();

  await expect(page.getByLabel('Streckenname')).toHaveValue('Cloned Farmrunde');
  await expect(page.getByLabel('Kommentar')).toHaveValue('bei Regen sehr matschig');

  await page.getByLabel('Löschen').click();
  await expect(page.getByTestId(trackSelectors.deleteTrackConfirmation)).toBeVisible();
  await page.getByLabel('Löschen').nth(1).click();

  const deletedTrack = page.getByTestId(trackSelectors.trackNode).filter({ hasText: 'Cloned Farmrunde' });

  await expect(deletedTrack).not.toBeVisible();
});

test('should filter tracks based on search input', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click();
  await expect(page.getByTestId(globalSelectors.tracksTab)).toContainText("Strecken");

  let trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2); // Root nodes: Höchberg and Dummy

  await page.getByTestId(trackSelectors.toggler).nth(0).click();
  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(5); // After expansion: Höchberg, Farmrunde, Stadtrunde, Waldweg, Dummy

  await page.getByTestId('track-filter-input').fill('Farm');

  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2); // Höchberg and Farmrunde
  await expect(trackNodes.nth(0)).toContainText('Höchberg');
  await expect(trackNodes.nth(1)).toContainText('Farmrunde');

  await page.getByTestId('track-filter-input').fill('farm');

  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2); // Höchberg and Farmrunde
  await expect(trackNodes.nth(0)).toContainText('Höchberg');
  await expect(trackNodes.nth(1)).toContainText('Farmrunde');

  await page.getByTestId('track-filter-input').fill('Stadt');

  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(2);
  await expect(trackNodes.nth(0)).toContainText('Höchberg');
  await expect(trackNodes.nth(1)).toContainText('Stadtrunde');

  await page.getByTestId('track-filter-input').fill('');

  trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(5);
});
