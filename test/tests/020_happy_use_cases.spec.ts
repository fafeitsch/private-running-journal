import {expect, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {trackSelectors} from '../selectors/track-selectors';
import * as fs from 'fs';
import {journalSelectors} from '../selectors/journal-selectors';

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
  await expect(trackNodes).toHaveCount(4)
  await expect(trackNodes.nth(1)).toContainText('Farmrunde')
  await expect(trackNodes.nth(2)).toContainText('Stadtrunde')

  await trackNodes.nth(1).click()

  await expect(page.getByLabel('Streckenname')).toHaveValue('Farmrunde')
  await expect(page.getByLabel('Verwendungen')).toHaveValue('0')
  await expect(page.getByLabel('Länge')).toHaveValue('10,8')
  await expect(page.getByLabel('Speichern')).toBeDisabled()
});

test('should create new track (in root directory)', async ({page})  => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click()
  await page.getByLabel('Hinzufügen').click()
  await page.getByLabel('Streckenname').fill('Wurzelstrecke')
  await page.getByTestId(trackSelectors.createTrackButton).click()

  let trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(3)

  await expect(page.getByLabel('Streckenname').nth(1)).toHaveValue('Wurzelstrecke')
  await expect(page.getByLabel('Streckenname').nth(1)).toBeVisible()
  await expect(page.getByLabel('Verwendungen')).toHaveValue('0')
  await expect(page.getByLabel('Länge')).toHaveValue('0,0')
  await expect(page.getByLabel('Speichern')).toBeDisabled()

  await page.getByLabel('Vorwärts').click()
  await page.mouse.click(560, 500)
  await page.mouse.click(640, 450)
  await page.mouse.click(720, 390)
  await page.mouse.click(820, 290)
  await page.mouse.click(920, 200)
  await expect(page.getByLabel('Länge')).toHaveValue('0,7')
  await page.mouse.move(910, 210)
  await page.mouse.down()
  await page.mouse.move(560, 500)
  await page.mouse.up()
  await expect(page.getByTestId(trackSelectors.distanceMarker)).not.toBeVisible()
  await page.mouse.click(750, 380)
  await expect(page.getByTestId(trackSelectors.distanceMarker)).toBeVisible()
  await expect(page.getByLabel('Länge')).toHaveValue('1,0')

  await page.getByLabel("Speichern").click()
  await expect(page.getByLabel("Speichern")).toBeDisabled()
  const track  = JSON.parse(fs.readFileSync('testdata/tracks/wurzelstrecke/info.json') as any)
  expect(track.name).toEqual('Wurzelstrecke')
})

test('should use new and old track to create some journal entries', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()

  await expect(page.getByTestId(journalSelectors.emptyState)).toBeVisible()
  await page.getByLabel("Hinzufügen").click()

  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect (page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(3)
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(2).click()
  await page.getByTestId(journalSelectors.todayButton).click()
  await page.getByTestId(journalSelectors.createEntryButton).click()

  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toHaveCount(1)
  const regex = /\d\d.\d\d.\d\d\d\d/
  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toContainText(regex)
  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toContainText('Wurzelstrecke')
  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toContainText('1,0 km')
})
