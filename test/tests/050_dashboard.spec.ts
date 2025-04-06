import {expect, Page, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {journalSelectors} from '../selectors/journal-selectors';
import * as fs from 'fs';
import {dashboardSelectors} from '../selectors/dashboard-selectors';

test('should display dashboard data', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.monthInput.input).nth(0).fill("March 2025")
  await page.keyboard.press("Enter")
  await page.getByTestId(globalSelectors.monthInput.input).nth(1).fill("April 2025")
  await page.keyboard.press("Enter")
  const generalStatistics = page.getByTestId(dashboardSelectors.generalStatistics);
  await expect(generalStatistics).toContainText("Gesamtlänge: 30 km")
  await expect(generalStatistics).toContainText("Läufe: 2")
  await expect(generalStatistics).toContainText("Durchschnitt: 15 km")
  await expect(generalStatistics).toContainText("Median: 21 km")
  const topTracks = page.getByTestId(dashboardSelectors.topTrack);
  await expect(topTracks).toHaveCount(3)
  await expect(topTracks.nth(0)).toContainText("Farmrunde")
  await expect(topTracks.nth(0)).toContainText("1⨰")
  await expect(topTracks.nth(0)).toContainText("Höchberg")
  await expect(topTracks.nth(1)).toContainText("Stadtrunde")
  await expect(topTracks.nth(2)).toContainText("Waldrunde (kurz)")

  const monthlyPanel = page.getByTestId(dashboardSelectors.monthlyPanel);
  await expect(monthlyPanel).toHaveCount(2)
  await expect(monthlyPanel.nth(0)).toContainText("03/2025")
  await expect(monthlyPanel.nth(0)).toContainText("Gesamtlänge: 21 km")
  await expect(monthlyPanel.nth(0)).toContainText("Läufe: 1")
  await expect(monthlyPanel.nth(0)).toContainText("Durchschnitt: 21 km")
  await expect(monthlyPanel.nth(0)).toContainText("Median: 21 km")

  await page.getByTestId(globalSelectors.monthInput.input).nth(1).fill("March 2025")
  await page.keyboard.press("Enter")
  await expect(generalStatistics).toContainText("Gesamtlänge: 21 km")
  await expect(generalStatistics).toContainText("Läufe: 1")
  await expect(generalStatistics).toContainText("Durchschnitt: 21 km")
  await expect(generalStatistics).toContainText("Median: 21 km")
  await expect(topTracks).toHaveCount(2)
  await expect(topTracks.nth(0)).toContainText("Farmrunde")
  await expect(topTracks.nth(0)).toContainText("1⨰")
  await expect(topTracks.nth(0)).toContainText("Höchberg")
  await expect(topTracks.nth(1)).toContainText("Stadtrunde")
});
