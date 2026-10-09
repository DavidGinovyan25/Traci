import { ApiError, createApi, tokenSubject } from './api.js';
import { genres, statuses, escapeHTML as e, options, initials, duration, query, collectionPatch } from './ui.js';

const main = document.querySelector('#main');
const navigation = document.querySelector('#navigation');
const account = document.querySelector('#account');
const dialog = document.querySelector('#dialog');
const notice = document.querySelector('#notice');
const limit = 12;
const filters = {
  catalog: { search: '', genre: '', offset: 0 },
  collection: { status: '', offset: 0 },
  users: { offset: 0 },
};
let view = 'catalog';
let authMode = 'login';
let loadID = 0;
let noticeTimer;
let currentPage = { items: [], total: 0 };
let modalID = 0;
const api = createApi({ onSessionEnd: () => {
  dialog.close();
  renderAuth();
} });

function notify(message, error = false) {
  clearTimeout(noticeTimer);
  notice.textContent = message;
  notice.classList.toggle('error', error);
  notice.hidden = false;
  noticeTimer = setTimeout(() => { notice.hidden = true; }, 5000);
}

function renderHeader() {
  const session = api.session;
  navigation.innerHTML = session ? `
    <a href="#catalog" ${view === 'catalog' ? 'aria-current="page"' : ''}>Каталог</a>
    <a href="#collection" ${view === 'collection' ? 'aria-current="page"' : ''}>Моя коллекция</a>
    ${session.role === 'admin' ? `<a href="#users" ${view === 'users' ? 'aria-current="page"' : ''}>Пользователи</a>` : ''}` : '';
  account.innerHTML = session ? `<span class="account-email" title="${e(session.email)}">${e(session.name || session.email)}</span>
    ${session.role === 'admin' ? '<span class="badge">Админ</span>' : ''}<button class="text-button" data-action="logout">Выйти ↗</button>`
    : '<span class="header-note">Для тех, кто любит кино</span>';
}

function field(label, name, { value = '', type = 'text', required = false, extra = '' } = {}) {
  return `<label class="field">${label}<input name="${name}" type="${type}" value="${e(value)}" ${required ? 'required' : ''} ${extra}></label>`;
}

function renderAuth(email = '') {
  loadID++;
  renderHeader();
  const register = authMode === 'register';
  main.innerHTML = `<section class="auth-layout">
    <div class="auth-intro"><div class="eyebrow">ТВОЯ ИСТОРИЯ КИНО</div>
      <h1>Хорошее кино.<br>В твоей<br><span>коллекции.</span></h1>
      <p>Сохраняй то, что хочется посмотреть.<br>Отмечай то, что уже стало любимым.</p>
      <div class="cinema-art" aria-hidden="true"><div class="art-orbit"></div><div class="art-frame"><span>TRACI ORIGINAL</span><strong>Ещё один<br>хороший<br>вечер.</strong><small>▶ &nbsp; ТВОЙ СЛЕДУЮЩИЙ ФИЛЬМ</small></div><span class="art-caption">COLLECT MOMENTS / NOT JUST MOVIES</span></div>
    </div>
    <div class="auth-panel"><div class="eyebrow">НАЧНИ СО СВОЕГО АККАУНТА</div>
      <h2>${register ? 'Знакомство с Traci' : 'С возвращением'}</h2>
      <p class="muted">${register ? 'Твои фильмы, оценки и впечатления — в одном месте.' : 'Твоя коллекция ждёт продолжения.'}</p>
      <div class="auth-tabs"><button data-action="auth-login" class="${!register ? 'active' : ''}">Войти</button><button data-action="auth-register" class="${register ? 'active' : ''}">Зарегистрироваться</button></div>
      <form id="auth-form" data-mode="${authMode}">
        ${register ? field('Никнейм', 'username', { required: true, extra: 'maxlength="50" autocomplete="username"' }) + `<div class="form-row">${field('Имя', 'first_name', { required: true, extra: 'autocomplete="given-name"' })}${field('Фамилия', 'second_name', { required: true, extra: 'autocomplete="family-name"' })}</div>` : ''}
        ${field('Почта', 'email', { type: 'email', value: email, required: true, extra: 'maxlength="255" autocomplete="email" placeholder="you@example.com"' })}
        ${field('Пароль', 'password', { type: 'password', required: true, extra: `minlength="8" maxlength="128" autocomplete="${register ? 'new-password' : 'current-password'}" placeholder="Не менее 8 символов"` })}
        <p class="form-error" role="alert" hidden></p>
        <button type="submit" class="button full">${register ? 'Создать аккаунт' : 'Войти в Traci'} <span aria-hidden="true">↗</span></button>
      </form>
      <p class="auth-footnote">${register ? 'После регистрации ты сможешь войти и собрать свою коллекцию.' : 'Маленькое место для большой любви к кино.'}</p>
    </div>
  </section>`;
}

function emptyState(title, text, action = '') {
  return `<div class="empty"><div class="empty-icon" aria-hidden="true">◫</div><h2>${title}</h2><p>${text}</p>${action}</div>`;
}

function pager(page) {
  if (!page.total) return '';
  const offset = filters[view].offset;
  return `<div class="pagination"><span>${offset + 1}–${Math.min(offset + limit, page.total)} из ${page.total}</span>
    <div><button class="button secondary" data-action="previous" ${offset === 0 ? 'disabled' : ''}>← Назад</button>
    <button class="button secondary" data-action="next" ${offset + limit >= page.total ? 'disabled' : ''}>Дальше →</button></div></div>`;
}

function movieCard(movie, item = null) {
  const palette = Object.keys(genres).indexOf(movie.genre) % 5;
  return `<article class="movie-card">
    <button class="movie-cover palette-${palette}" data-action="movie" data-id="${e(movie.id)}" aria-label="Открыть фильм ${e(movie.name)}">
      <span class="cover-genre">${e(genres[movie.genre] || movie.genre)}</span>
      <span class="cover-orbit" aria-hidden="true"></span><span class="cover-title" aria-hidden="true">${e(initials(movie.name))}</span>
      <span class="cover-bottom">${movie.release_year}<span aria-hidden="true">↗</span></span>
    </button>
    <div class="movie-info"><div class="movie-meta">${e(genres[movie.genre] || movie.genre)} <span>·</span> ${duration(movie.duration_min)}</div>
      <h2><button class="movie-name" data-action="movie" data-id="${e(movie.id)}">${e(movie.name)}</button></h2>
      ${item ? `<div class="collection-meta"><span class="status status-${item.status.toLowerCase()}">${e(statuses[item.status])}</span>${item.personal_rating != null ? `<span class="rating">★ ${item.personal_rating}/10</span>` : ''}</div>
        <button class="button secondary full" data-action="collection-edit" data-id="${e(movie.id)}">Статус, оценка и отзыв</button>`
      : `<button class="button secondary full" data-action="add" data-id="${e(movie.id)}"><span aria-hidden="true">＋</span> В коллекцию</button>`}
    </div>
  </article>`;
}

function pageFrame() {
  renderHeader();
  const catalog = view === 'catalog';
  const collection = view === 'collection';
  main.innerHTML = `<section class="workspace">
    <div class="page-heading"><div><div class="eyebrow">${catalog ? 'ОТКРОЙ ЧТО-ТО НОВОЕ' : collection ? 'СОХРАНЕНО ДЛЯ ТЕБЯ' : 'УПРАВЛЕНИЕ TRACI'}</div>
    <h1>${catalog ? 'Всё начинается с кино.' : collection ? 'Твоя коллекция.' : 'Пользователи.'}</h1>
    <p class="muted">${catalog ? 'Найди фильм для следующего вечера.' : collection ? 'Планы, впечатления и фильмы, к которым хочется вернуться.' : 'Аккаунты, роли и доступ к приложению.'}</p></div>
    ${catalog && api.session.role === 'admin' ? '<button class="button" data-action="movie-create">＋ Добавить фильм</button>' : ''}</div>
    ${catalog ? `<form id="filters-form" class="filters"><label class="search-field"><span aria-hidden="true">⌕</span><input name="search" aria-label="Поиск по названию" placeholder="Поиск по названию фильма" value="${e(filters.catalog.search)}" maxlength="255"></label>
      <select name="genre" aria-label="Жанр">${options(genres, filters.catalog.genre, 'Все жанры')}</select><button class="button" type="submit">Найти</button><button class="text-button" type="button" data-action="reset-filters">Сбросить</button></form>`
    : collection ? `<form id="filters-form" class="filters"><select name="status" aria-label="Статус просмотра">${options(statuses, filters.collection.status, 'Все статусы')}</select><button class="button" type="submit">Показать</button><button class="text-button" type="button" data-action="reset-filters">Сбросить</button></form>` : ''}
    <div id="results" aria-live="polite" aria-busy="true"><div class="loading"><span class="spinner"></span>Загружаем ${collection ? 'коллекцию' : catalog ? 'каталог' : 'пользователей'}…</div></div>
  </section>`;
}

function usersTable(users) {
  if (!users.length) return emptyState('Пользователей пока нет', 'Здесь появятся зарегистрированные аккаунты.');
  return `<div class="table-wrap"><table><thead><tr><th>Пользователь</th><th>Почта</th><th>Роль</th><th>Статус</th><th>Действия</th></tr></thead><tbody>${users.map(user => `<tr>
    <td><strong>${e(user.username)}</strong><small>${e(user.first_name)} ${e(user.second_name)}</small></td><td>${e(user.email)}</td>
    <td>${user.role === 'admin' ? 'Администратор' : 'Пользователь'}</td><td><span class="status ${user.is_blocked ? 'status-dropped' : 'status-watched'}">${user.is_blocked ? 'Заблокирован' : 'Активен'}</span></td>
    <td>${user.id === api.session.id ? '<span class="muted">Это ты</span>' : `<div class="table-actions"><button class="text-button" data-action="user-block" data-id="${e(user.id)}">${user.is_blocked ? 'Разблокировать' : 'Заблокировать'}</button><button class="text-button danger" data-action="user-delete" data-id="${e(user.id)}">Удалить</button></div>`}</td>
  </tr>`).join('')}</tbody></table></div>`;
}

async function loadPage() {
  if (!api.session) return renderAuth();
  const requested = location.hash.slice(1);
  view = ['catalog', 'collection', 'users'].includes(requested) ? requested : 'catalog';
  if (view === 'users' && api.session.role !== 'admin') view = 'catalog';
  const id = ++loadID;
  const pageView = view;
  pageFrame();
  const results = document.querySelector('#results');
  try {
    let page = await api.request(`/${pageView}?${query({ ...filters[pageView], limit })}`);
    if (id !== loadID || !api.session) return;
    if (!page.items.length && filters[pageView].offset > 0) {
      filters[pageView].offset = Math.max(0, Math.ceil(page.total / limit) - 1) * limit;
      return loadPage();
    }
    currentPage = page;
    results.innerHTML = pageView === 'users' ? usersTable(page.items) : page.items.length
      ? `<div class="section-caption"><span>${pageView === 'catalog' ? 'В КАТАЛОГЕ' : 'В КОЛЛЕКЦИИ'} / ${page.total}</span><span>Твой следующий хороший вечер</span></div><div class="movie-grid">${page.items.map(item => pageView === 'collection' ? movieCard(item.movie, item) : movieCard(item)).join('')}</div>`
      : pageView === 'collection' ? emptyState('Начни свою историю кино', 'Добавь первый фильм из каталога. Здесь будут твои планы, оценки и отзывы.', '<a class="button" href="#catalog">Перейти в каталог ↗</a>')
      : emptyState(filters.catalog.search || filters.catalog.genre ? 'Ничего не нашлось' : 'Каталог ждёт первых фильмов', filters.catalog.search || filters.catalog.genre ? 'Попробуй другое название или сбрось фильтры.' : 'Администратор может добавить фильмы — и здесь появится твой следующий киносеанс.');
    results.innerHTML += pager(page);
    results.setAttribute('aria-busy', 'false');
  } catch (error) {
    if (id !== loadID || !api.session) return;
    results.setAttribute('aria-busy', 'false');
    results.innerHTML = emptyState('Не удалось загрузить страницу', e(error.message), '<button class="button" data-action="retry">Попробовать ещё раз</button>');
  }
}

function openDialog(title, content) {
  modalID++;
  dialog.innerHTML = `<div class="dialog-head"><h2 id="dialog-title">${e(title)}</h2><button class="icon-button" data-action="close-dialog" aria-label="Закрыть">×</button></div>${content}`;
  if (!dialog.open) dialog.showModal();
}

function formFooter(label = 'Сохранить') {
  return `<p class="form-error" role="alert" hidden></p><div class="dialog-actions"><button type="button" class="button secondary" data-action="close-dialog">Отмена</button><button type="submit" class="button">${label}</button></div>`;
}

function movieForm(movie = null) {
  openDialog(movie ? 'Редактировать фильм' : 'Новый фильм', `<form id="movie-form" data-id="${movie?.id || ''}">
    ${field('Название', 'name', { value: movie?.name, required: true, extra: 'maxlength="255"' })}
    <div class="form-row">${field('Длительность, минут', 'duration_min', { value: movie?.duration_min, type: 'number', required: true, extra: 'min="1" max="2147483647" step="1"' })}${field('Год выхода', 'release_year', { value: movie?.release_year, type: 'number', required: true, extra: 'min="1888" max="2147483647" step="1"' })}</div>
    <label class="field">Жанр<select name="genre" aria-label="Жанр" required>${options(genres, movie?.genre || 'DRAMA')}</select></label>
    <label class="field">Описание<textarea name="description" rows="5" maxlength="5000">${e(movie?.description)}</textarea></label>${formFooter(movie ? 'Сохранить изменения' : 'Добавить фильм')}</form>`);
}

async function movieDetails(id) {
  openDialog('Фильм', '<div class="loading"><span class="spinner"></span>Загружаем…</div>');
  const ticket = modalID;
  const movie = await api.request(`/catalog/${id}`);
  if (!dialog.open || ticket !== modalID) return;
  openDialog(movie.name, `<div class="detail-meta"><span class="badge">${e(genres[movie.genre])}</span><span>${movie.release_year}</span><span>${duration(movie.duration_min)}</span></div>
    <p class="description">${e(movie.description || 'У этого фильма пока нет описания.')}</p>
    <div class="detail-actions"><button class="button" data-action="add" data-id="${e(id)}">＋ В коллекцию</button>
    ${api.session?.role === 'admin' ? `<button class="button secondary" data-action="movie-edit" data-id="${e(id)}">Редактировать</button><button class="text-button danger" data-action="movie-delete" data-id="${e(id)}">Удалить фильм</button>` : ''}</div>`);
}

async function editCollection(id) {
  openDialog('Фильм в коллекции', '<div class="loading"><span class="spinner"></span>Загружаем…</div>');
  const ticket = modalID;
  const item = await api.request(`/collection/${id}`);
  if (!dialog.open || ticket !== modalID) return;
  openDialog(item.movie.name, `<form id="collection-form" data-id="${e(id)}">
    <label class="field">Статус просмотра<select name="status" aria-label="Статус просмотра">${options(statuses, item.status)}</select></label>
    ${field('Личная оценка', 'personal_rating', { value: item.personal_rating, type: 'number', extra: 'min="1" max="10" step="1" placeholder="От 1 до 10, можно оставить пустым"' })}
    <label class="field">Твой отзыв<textarea name="review" rows="5" maxlength="1000" placeholder="Какие впечатления оставил фильм?">${e(item.review)}</textarea></label>${formFooter()}
    <button type="button" class="text-button danger remove-collection" data-action="collection-delete" data-id="${e(id)}">Удалить из коллекции</button></form>`);
}

function confirmAction(title, text, action, id) {
  openDialog(title, `<p class="description">${e(text)}</p><form id="confirm-form" data-operation="${action}" data-id="${e(id)}">${formFooter('Удалить')}</form>`);
}

async function identifyRole() {
  const session = api.session;
  try {
    const user = await api.request(`/users/${session.id}`);
    if (api.session?.token === session.token) api.setSession({ ...session, role: user.role, name: user.username });
  } catch (error) {
    if (error.status === 403 && error.code === 'FORBIDDEN') {
      if (api.session?.token === session.token) api.setSession({ ...session, role: 'user' });
    } else throw error;
  }
}

async function submitAuth(form, values) {
  if (form.dataset.mode === 'register') {
    await api.request('/auth/register', { method: 'POST', body: values, anonymous: true });
    authMode = 'login';
    renderAuth(values.email);
    notify('Аккаунт создан. Теперь войди в Traci.');
    return;
  }
  const login = await api.request('/auth/login', { method: 'POST', body: values, anonymous: true });
  const id = tokenSubject(login.access_token);
  if (!id) throw new ApiError('Сервер вернул некорректную сессию.');
  api.setSession({ token: login.access_token, id, email: values.email, role: 'user' });
  try { await identifyRole(); } catch (error) { api.setSession(null); throw error; }
  location.hash = 'catalog';
  await loadPage();
}

async function submit(form) {
  if (form.dataset.busy) return;
  const values = Object.fromEntries(new FormData(form));
  const button = form.querySelector('[type="submit"]');
  const errorBox = form.querySelector('.form-error');
  form.dataset.busy = 'true';
  if (button) button.disabled = true;
  if (errorBox) errorBox.hidden = true;
  try {
    if (form.id === 'auth-form') return await submitAuth(form, values);
    if (form.id === 'filters-form') {
      Object.assign(filters[view], values, { offset: 0 });
      return await loadPage();
    }
    if (form.id === 'movie-form') {
      const body = { ...values, duration_min: Number(values.duration_min), release_year: Number(values.release_year) };
      const id = form.dataset.id;
      if (!body.description.trim()) {
        if (id) body.description = null;
        else delete body.description;
      }
      await api.request(id ? `/catalog/${id}` : '/catalog', { method: id ? 'PATCH' : 'POST', body });
      notify(id ? 'Изменения сохранены.' : 'Фильм добавлен в каталог.');
    } else if (form.id === 'collection-form') {
      await api.request(`/collection/${form.dataset.id}`, { method: 'PATCH', body: collectionPatch(values) });
      notify('Коллекция обновлена.');
    } else if (form.id === 'confirm-form') {
      const operation = form.dataset.operation;
      await api.request(`/${operation}/${form.dataset.id}`, { method: 'DELETE' });
      notify(operation === 'users' ? 'Пользователь удалён.' : operation === 'catalog' ? 'Фильм удалён.' : 'Фильм удалён из коллекции.');
    }
    dialog.close();
    await loadPage();
  } catch (error) {
    if (form.isConnected && errorBox) {
      errorBox.textContent = error.message;
      errorBox.hidden = false;
    } else notify(error.message, true);
  } finally {
    delete form.dataset.busy;
    if (button) button.disabled = false;
  }
}

async function action(button) {
  if (button.dataset.busy) return;
  button.dataset.busy = 'true';
  button.disabled = true;
  const id = button.dataset.id;
  try {
    switch (button.dataset.action) {
      case 'logout': api.setSession(null); dialog.close(); renderAuth(); break;
      case 'auth-login': authMode = 'login'; renderAuth(); break;
      case 'auth-register': authMode = 'register'; renderAuth(); break;
      case 'close-dialog': modalID++; dialog.close(); break;
      case 'retry': await loadPage(); break;
      case 'previous': filters[view].offset = Math.max(0, filters[view].offset - limit); await loadPage(); break;
      case 'next': filters[view].offset += limit; await loadPage(); break;
      case 'reset-filters': Object.keys(filters[view]).forEach(key => { filters[view][key] = key === 'offset' ? 0 : ''; }); await loadPage(); break;
      case 'movie-create': movieForm(); break;
      case 'movie': await movieDetails(id); break;
      case 'movie-edit': {
        const movie = await api.request(`/catalog/${id}`);
        movieForm(movie);
        break;
      }
      case 'add':
        await api.request('/collection', { method: 'POST', body: { movie_id: id, status: 'PLANNED' } });
        notify('Фильм теперь в твоей коллекции.');
        button.textContent = '✓ В коллекции';
        return;
      case 'collection-edit': await editCollection(id); break;
      case 'collection-delete': confirmAction('Удалить из коллекции?', 'Твоя оценка и отзыв тоже будут удалены.', 'collection', id); break;
      case 'movie-delete': confirmAction('Удалить фильм?', 'Фильм будет удалён из каталога и всех личных коллекций.', 'catalog', id); break;
      case 'user-delete': confirmAction('Удалить аккаунт?', 'Аккаунт и его личная коллекция будут удалены.', 'users', id); break;
      case 'user-block': {
        const user = currentPage.items.find(item => item.id === id);
        await api.request(`/users/${id}`, { method: 'PATCH', body: { is_blocked: !user.is_blocked } });
        notify(user.is_blocked ? 'Пользователь разблокирован.' : 'Пользователь заблокирован.');
        await loadPage();
        break;
      }
    }
  } catch (error) {
    if (error.status === 409 && button.dataset.action === 'add') notify('Этот фильм уже есть в твоей коллекции.');
    else notify(error.message, true);
  } finally {
    delete button.dataset.busy;
    if (button.textContent !== '✓ В коллекции') button.disabled = false;
  }
}

function showUnexpected(error) {
  notify(error instanceof Error ? error.message : 'Что-то пошло не так. Попробуй ещё раз.', true);
}

document.addEventListener('click', event => {
  const button = event.target.closest('[data-action]');
  if (button) action(button).catch(showUnexpected);
});
document.addEventListener('submit', event => {
  if (event.target.matches('form')) {
    event.preventDefault();
    submit(event.target).catch(showUnexpected);
  }
});
window.addEventListener('hashchange', () => loadPage().catch(showUnexpected));
dialog.addEventListener('close', () => { modalID++; });

async function start() {
  if (!api.session) return renderAuth();
  renderHeader();
  main.innerHTML = '<div class="loading"><span class="spinner"></span>Возвращаемся к твоей коллекции…</div>';
  try {
    await identifyRole();
    await loadPage();
  } catch (error) {
    api.setSession(null);
    renderAuth();
    notify(error.message, true);
  }
}

start().catch(showUnexpected);
