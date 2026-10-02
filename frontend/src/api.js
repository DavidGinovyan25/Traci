const sessionKey = 'traci.session';

const messages = {
  FORBIDDEN: 'Недостаточно прав для этого действия.',
  USER_BLOCKED: 'Аккаунт заблокирован администратором.',
  NOT_FOUND: 'Запись не найдена. Возможно, её уже удалили.',
  ALREADY_EXISTS: 'Такая запись уже существует.',
  SELF_MODIFICATION: 'Нельзя заблокировать или удалить собственный аккаунт.',
  BAD_REQUEST: 'Проверь заполненные поля.',
  INTERNAL_ERROR: 'Ошибка сервера. Попробуй ещё раз позже.',
};

export class ApiError extends Error {
  constructor(message, status = 0, code = '') {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

export function tokenSubject(token) {
  try {
    const part = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    const claims = JSON.parse(atob(part));
    if (!claims.sub || !Number.isFinite(claims.exp) || claims.exp * 1000 <= Date.now()) return null;
    return claims.sub;
  } catch {
    return null;
  }
}

export function createApi({ storage = sessionStorage, fetcher = fetch, onSessionEnd = () => {} } = {}) {
  let session = null;
  try {
    const saved = JSON.parse(storage.getItem(sessionKey));
    if (saved?.token && tokenSubject(saved.token)) session = saved;
    else storage.removeItem(sessionKey);
  } catch {
    storage.removeItem(sessionKey);
  }

  function setSession(value) {
    session = value;
    if (value) storage.setItem(sessionKey, JSON.stringify(value));
    else storage.removeItem(sessionKey);
  }

  async function request(path, { method = 'GET', body, anonymous = false } = {}) {
    const token = anonymous ? null : session?.token;
    const headers = { Accept: 'application/json' };
    if (token) headers.Authorization = `Bearer ${token}`;
    if (body !== undefined) headers['Content-Type'] = 'application/json';
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 15000);
    try {
      const response = await fetcher(`/api${path}`, {
        method, headers, signal: controller.signal,
        ...(body === undefined ? {} : { body: JSON.stringify(body) }),
      });
      const data = response.status === 204 ? null : await response.json();
      if (!response.ok) {
        const ended = !anonymous && (response.status === 401 || data?.code === 'USER_BLOCKED');
        if (ended && session?.token === token) {
          setSession(null);
          onSessionEnd();
        }
        const message = response.status === 401
          ? (anonymous ? 'Неверная почта или пароль.' : 'Сессия истекла. Войди ещё раз.')
          : messages[data?.code] || 'Не удалось выполнить запрос.';
        throw new ApiError(message, response.status, data?.code);
      }
      return data;
    } catch (error) {
      if (error instanceof ApiError) throw error;
      throw new ApiError(error.name === 'AbortError'
        ? 'Сервер долго не отвечает. Попробуй ещё раз.'
        : 'Не удалось связаться с сервером. Проверь соединение.');
    } finally {
      clearTimeout(timeout);
    }
  }

  return { request, setSession, get session() { return session; } };
}
