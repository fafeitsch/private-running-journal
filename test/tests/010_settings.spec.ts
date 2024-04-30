import { test, expect } from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {settingsSelectors} from '../selectors/settings-selectors';
import * as fs from 'fs';

test('should show app in correct language', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.settingsTab).click()
  await expect(page.getByTestId(settingsSelectors.saveButton)).toBeDisabled()
  await expect(page.getByTestId(globalSelectors.settingsTab)).toContainText('Einstellungen')

  await expect(page.getByTestId(settingsSelectors.language)).toContainText('Deutsch')
  await page.getByTestId(settingsSelectors.language).click()
  await page.getByLabel('Englisch').click()

  await page.getByTestId(settingsSelectors.saveButton).click()
  await expect(page.getByTestId(globalSelectors.settingsTab)).toContainText('Settings')

  let settings = JSON.parse(fs.readFileSync('testdata/settings.json') as any)
  expect(settings.language).toEqual('en')
  await expect(page.getByTestId(settingsSelectors.language)).toContainText('English')
  await page.getByTestId(settingsSelectors.language).click()
  await page.getByLabel('German').click()
  await page.getByTestId(settingsSelectors.saveButton).click()
  await expect(page.getByTestId(globalSelectors.settingsTab)).toContainText('Einstellungen')
  settings = JSON.parse(fs.readFileSync('testdata/settings.json') as any)
  expect(settings.language).toEqual('de')
});
