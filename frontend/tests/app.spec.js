import { test, expect } from '@playwright/test';
import { execFileSync } from 'node:child_process';
import { randomUUID } from 'node:crypto';
import { fileURLToPath } from 'node:url';

const project = fileURLToPath(new URL('../../', import.meta.url));
const envFile = process.env.ENV_FILE || `${project}backend/.env`;

function sql(query) {
  execFileSync('docker', ['compose', '--project-directory', project, '--env-file', envFile, 'exec', '-T', 'db',
    'sh', '-c', 'exec psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'],
  { cwd: project, input: query, stdio: ['pipe', 'pipe', 'pipe'], timeout: 15000 });
}

function createAdmin(username, email, password) {
  execFileSync('docker', ['compose', '--project-directory', project, '--env-file', envFile, 'exec', '-T', 'app', '/app/traci', 'create-admin'],
    { cwd: project, input: `${username}\nTest\nAdmin\n${email}\n${password}\n`, stdio: ['pipe', 'pipe', 'pipe'], timeout: 15000 });
}

async function register(page, username, email, password) {
  await page.goto('/');
  await page.getByRole('button', { name: 'Зарегистрироваться', exact: true }).click();
  await page.getByLabel('Никнейм', { exact: true }).fill(username);
  await page.getByLabel('Имя', { exact: true }).fill('Test');
  await page.getByLabel('Фамилия', { exact: true }).fill('User');
  await page.getByLabel('Почта', { exact: true }).fill(email);
  await page.getByLabel('Пароль', { exact: true }).fill(password);
  await page.getByRole('button', { name: 'Создать аккаунт' }).click();
  await expect(page.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
}

async function login(page, email, password) {
  await page.getByLabel('Почта', { exact: true }).fill(email);
  await page.getByLabel('Пароль', { exact: true }).fill(password);
  await page.getByRole('button', { name: 'Войти в Traci' }).click();
  await expect(page.getByRole('heading', { name: 'Всё начинается с кино.' })).toBeVisible();
  await expect(page.locator('#results')).toHaveAttribute('aria-busy', 'false');
}

test('real API: registration, movie CRUD, collection, user administration and mobile layout', async ({ browser }) => {
  const run = randomUUID().replaceAll('-', '');
  const adminEmail = `frontend-admin-${run}@example.com`;
  const userEmail = `frontend-user-${run}@example.com`;
  const movieName = `Ночной сеанс ${run}`;
  const password = `Test-${run}`;
  const adminContext = await browser.newContext({ baseURL: test.info().project.use.baseURL, viewport: { width: 1440, height: 1000 } });
  const userContext = await browser.newContext({ baseURL: test.info().project.use.baseURL, viewport: { width: 1440, height: 1000 } });
  const admin = await adminContext.newPage();
  const user = await userContext.newPage();
  const errors = [];
  admin.on('pageerror', error => errors.push(error.message));
  user.on('pageerror', error => errors.push(error.message));
  try {
    await admin.goto('/');
    await expect(admin.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
    await admin.screenshot({ path: '/tmp/traci-auth.png', fullPage: true });
    createAdmin(`admin_${run}`, adminEmail, password);
    await login(admin, adminEmail, password);
    await expect(admin.getByRole('link', { name: 'Пользователи' })).toBeVisible();
    await admin.getByRole('button', { name: 'Добавить фильм' }).click();
    const adminDialog = admin.getByRole('dialog');
    await adminDialog.getByLabel('Название', { exact: true }).fill(movieName);
    await adminDialog.getByLabel('Длительность, минут').fill('120');
    await adminDialog.getByLabel('Год выхода').fill('2020');
    await adminDialog.getByLabel('Жанр', { exact: true }).selectOption('DRAMA');
    await adminDialog.getByLabel('Описание').fill('<script>window.injected = true</script>');
    await adminDialog.getByRole('button', { name: 'Добавить фильм', exact: true }).click();
    await expect(adminDialog).not.toBeVisible();
    await expect(admin.getByRole('button', { name: movieName, exact: true })).toBeVisible();
    await admin.screenshot({ path: '/tmp/traci-catalog.png', fullPage: true });
    await admin.getByRole('button', { name: movieName, exact: true }).click();
    await expect(adminDialog.locator('.description')).toHaveText('<script>window.injected = true</script>');
    expect(await admin.evaluate(() => window.injected)).toBeUndefined();
    await adminDialog.getByRole('button', { name: 'Редактировать', exact: true }).click();
    await adminDialog.getByLabel('Длительность, минут').fill('130');
    await adminDialog.getByLabel('Описание').fill('');
    await adminDialog.getByRole('button', { name: 'Сохранить изменения' }).click();
    await expect(adminDialog).not.toBeVisible();
    await admin.getByRole('button', { name: movieName, exact: true }).click();
    await expect(adminDialog).toContainText('2 ч 10 мин');
    await expect(adminDialog).toContainText('У этого фильма пока нет описания.');
    await adminDialog.getByRole('button', { name: 'Закрыть', exact: true }).click();

    await register(user, `user_${run}`, userEmail, password);
    await user.getByLabel('Пароль', { exact: true }).fill('incorrect-password');
    await user.getByRole('button', { name: 'Войти в Traci' }).click();
    await expect(user.getByRole('alert')).toHaveText('Неверная почта или пароль.');
    await login(user, userEmail, password);
    await expect(user.getByRole('link', { name: 'Пользователи' })).toHaveCount(0);
    await user.getByRole('link', { name: 'Моя коллекция' }).click();
    await expect(user.getByRole('heading', { name: 'Начни свою историю кино' })).toBeVisible();
    await user.getByRole('link', { name: 'Каталог', exact: true }).click();
    await user.getByLabel('Поиск по названию').fill(movieName);
    await user.getByRole('button', { name: 'Найти', exact: true }).click();
    await expect(user.locator('.movie-card')).toHaveCount(1);
    await user.getByRole('button', { name: 'В коллекцию', exact: true }).click();
    await expect(user.getByRole('status')).toContainText('Фильм теперь в твоей коллекции.');
    await user.getByRole('link', { name: 'Моя коллекция' }).click();
    await user.getByRole('button', { name: 'Статус, оценка и отзыв' }).click();
    const userDialog = user.getByRole('dialog');
    await userDialog.getByLabel('Статус просмотра').selectOption('WATCHED');
    await userDialog.getByLabel('Личная оценка').fill('9');
    await userDialog.getByLabel('Твой отзыв').fill('Отличный фильм');
    await userDialog.getByRole('button', { name: 'Сохранить', exact: true }).click();
    await expect(userDialog).not.toBeVisible();
    await expect(user.locator('.movie-card')).toContainText('★ 9/10');
    await user.getByRole('button', { name: 'Статус, оценка и отзыв' }).click();
    await expect(userDialog.getByLabel('Твой отзыв')).toHaveValue('Отличный фильм');
    await userDialog.getByLabel('Личная оценка').fill('');
    await userDialog.getByLabel('Твой отзыв').fill('');
    await userDialog.getByRole('button', { name: 'Сохранить', exact: true }).click();
    await expect(userDialog).not.toBeVisible();
    await expect(user.locator('.rating')).toHaveCount(0);
    await user.locator('#filters-form').getByLabel('Статус просмотра').selectOption('PLANNED');
    await user.getByRole('button', { name: 'Показать' }).click();
    await expect(user.getByRole('heading', { name: 'Начни свою историю кино' })).toBeVisible();
    await user.getByRole('button', { name: 'Сбросить' }).click();
    await expect(user.locator('.movie-card')).toHaveCount(1);
    await user.reload();
    await expect(user.getByRole('heading', { name: 'Твоя коллекция.' })).toBeVisible();
    await user.setViewportSize({ width: 390, height: 844 });
    await user.screenshot({ path: '/tmp/traci-mobile.png', fullPage: true });
    expect(await user.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBeTruthy();
    await user.getByRole('button', { name: 'Статус, оценка и отзыв' }).click();
    await userDialog.getByRole('button', { name: 'Удалить из коллекции' }).click();
    await userDialog.getByRole('button', { name: 'Удалить', exact: true }).click();
    await expect(userDialog).not.toBeVisible();
    await expect(user.getByRole('heading', { name: 'Начни свою историю кино' })).toBeVisible();

    await admin.getByRole('link', { name: 'Пользователи' }).click();
    const row = admin.getByRole('row').filter({ hasText: userEmail });
    await row.getByRole('button', { name: 'Заблокировать', exact: true }).click();
    await expect(row).toContainText('Заблокирован');
    await user.getByRole('link', { name: 'Каталог', exact: true }).click();
    await expect(user.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
    await row.getByRole('button', { name: 'Разблокировать', exact: true }).click();
    await expect(row).toContainText('Активен');
    await login(user, userEmail, password);
    await row.getByRole('button', { name: 'Удалить', exact: true }).click();
    await adminDialog.getByRole('button', { name: 'Удалить', exact: true }).click();
    await expect(row).toHaveCount(0);
    await user.getByRole('link', { name: 'Моя коллекция' }).click();
    await expect(user.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
    await admin.getByRole('link', { name: 'Каталог', exact: true }).click();
    await admin.getByRole('button', { name: movieName, exact: true }).click();
    await adminDialog.getByRole('button', { name: 'Удалить фильм', exact: true }).click();
    await adminDialog.getByRole('button', { name: 'Удалить', exact: true }).click();
    await expect(adminDialog).not.toBeVisible();
    await expect(admin.getByRole('button', { name: movieName, exact: true })).toHaveCount(0);
    await admin.getByRole('button', { name: 'Выйти' }).click();
    await expect(admin.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
    await admin.reload();
    await expect(admin.getByRole('heading', { name: 'С возвращением' })).toBeVisible();
    expect(errors).toEqual([]);
  } finally {
    try {
      sql(`DELETE FROM movies WHERE name = '${movieName}'; DELETE FROM users WHERE email IN ('${adminEmail}', '${userEmail}');`);
    } finally {
      await Promise.allSettled([adminContext.close(), userContext.close()]);
    }
  }
});
