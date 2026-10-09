import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  testMatch: '*.spec.js',
  timeout: 90000,
  expect: { timeout: 10000 },
  workers: 1,
  use: {
    baseURL: process.env.FRONTEND_URL || 'http://localhost:3000',
    browserName: 'chromium',
    actionTimeout: 10000,
    viewport: { width: 1440, height: 1000 },
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
});
