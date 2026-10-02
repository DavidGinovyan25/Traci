import test from 'node:test';
import assert from 'node:assert/strict';
import { createApi, ApiError, tokenSubject } from '../src/api.js';

function storage() {
  const values = new Map();
  return {
    getItem: key => values.get(key) ?? null,
    setItem: (key, value) => values.set(key, value),
    removeItem: key => values.delete(key),
  };
}

function token(claims = {}) {
  return `header.${Buffer.from(JSON.stringify({ sub: 'user-id', exp: Date.now() / 1000 + 3600, ...claims })).toString('base64url')}.signature`;
}

const session = () => ({ token: token(), id: 'user-id', email: 'test@example.com', role: 'user' });

test('restores a valid session and discards expired or corrupt storage', () => {
  const saved = storage();
  const api = createApi({ storage: saved });
  api.setSession(session());
  assert.equal(createApi({ storage: saved }).session.id, 'user-id');
  api.setSession({ token: token({ exp: 1 }) });
  assert.equal(createApi({ storage: saved }).session, null);
  saved.setItem('traci.session', 'not json');
  assert.equal(createApi({ storage: saved }).session, null);
});

test('rejects malformed tokens and tokens without expiration', () => {
  assert.equal(tokenSubject('broken'), null);
  assert.equal(tokenSubject(token({ exp: null })), null);
  assert.equal(tokenSubject(token({ sub: '' })), null);
});

test('sends authorization and serializes nullable PATCH fields', async () => {
  let request;
  const api = createApi({ storage: storage(), fetcher: async (url, options) => {
    request = { url, ...options };
    return Response.json({ status: 'WATCHED' });
  } });
  const current = session();
  api.setSession(current);
  await api.request('/collection/movie-id', { method: 'PATCH', body: { personal_rating: null } });
  assert.equal(request.url, '/api/collection/movie-id');
  assert.equal(request.headers.Authorization, `Bearer ${current.token}`);
  assert.equal(request.headers['Content-Type'], 'application/json');
  assert.deepEqual(JSON.parse(request.body), { personal_rating: null });
});

test('does not send authorization on public login requests', async () => {
  const api = createApi({ storage: storage(), fetcher: async (url, options) => {
    assert.equal(options.headers.Authorization, undefined);
    return Response.json({ access_token: token() });
  } });
  api.setSession(session());
  await api.request('/auth/login', { method: 'POST', body: {}, anonymous: true });
});

test('DELETE 204 does not attempt to decode an empty body', async () => {
  const api = createApi({ storage: storage(), fetcher: async () => new Response(null, { status: 204 }) });
  assert.equal(await api.request('/catalog/movie-id', { method: 'DELETE' }), null);
});

for (const [status, code] of [[401, 'UNAUTHORIZED'], [403, 'USER_BLOCKED']]) {
  test(`${code} clears a protected session`, async () => {
    let ended = 0;
    const api = createApi({ storage: storage(), onSessionEnd: () => ended++,
      fetcher: async () => Response.json({ code }, { status }) });
    api.setSession(session());
    await assert.rejects(api.request('/catalog'), error => error instanceof ApiError && error.status === status);
    assert.equal(api.session, null);
    assert.equal(ended, 1);
  });
}

test('permission denial keeps the ordinary user signed in', async () => {
  const api = createApi({ storage: storage(), fetcher: async () => Response.json({ code: 'FORBIDDEN' }, { status: 403 }) });
  api.setSession(session());
  await assert.rejects(api.request('/users'), error => error.code === 'FORBIDDEN');
  assert.notEqual(api.session, null);
});

test('a failed old request cannot clear a newer session', async () => {
  let respond;
  const api = createApi({ storage: storage(), fetcher: () => new Promise(resolve => { respond = resolve; }) });
  api.setSession(session());
  const pending = api.request('/catalog');
  const newer = { ...session(), token: token({ sub: 'new-user' }) };
  api.setSession(newer);
  respond(Response.json({ code: 'UNAUTHORIZED' }, { status: 401 }));
  await assert.rejects(pending, ApiError);
  assert.equal(api.session.token, newer.token);
});

test('network failures become readable errors', async () => {
  const api = createApi({ storage: storage(), fetcher: async () => { throw new TypeError('fetch failed'); } });
  await assert.rejects(api.request('/catalog'), error => error instanceof ApiError && error.status === 0);
});
