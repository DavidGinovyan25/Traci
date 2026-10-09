import test from 'node:test';
import assert from 'node:assert/strict';
import { escapeHTML, collectionPatch, options, query, duration, initials } from '../src/ui.js';

test('movie titles and reviews cannot inject HTML or attributes', () => {
  assert.equal(escapeHTML('<script>"x" & \'y\'</script>'), '&lt;script&gt;&quot;x&quot; &amp; &#39;y&#39;&lt;/script&gt;');
  assert.equal(escapeHTML(null), '');
  assert.ok(!options({ DRAMA: '<img src=x>' }, 'DRAMA').includes('<img'));
});

test('clearing review and rating sends explicit null values', () => {
  assert.deepEqual(collectionPatch({ status: 'WATCHED', personal_rating: '', review: '   ' }), {
    status: 'WATCHED', personal_rating: null, review: null,
  });
  assert.deepEqual(collectionPatch({ status: 'WATCHED', personal_rating: '9', review: ' Good ' }), {
    status: 'WATCHED', personal_rating: 9, review: 'Good',
  });
});

test('filters omit empty values and URL-encode search terms', () => {
  assert.equal(query({ search: 'a & b', genre: '', offset: 0 }), 'search=a+%26+b&offset=0');
});

test('duration and title initials handle short and non-Latin titles', () => {
  assert.equal(duration(120), '2 ч');
  assert.equal(duration(169), '2 ч 49 мин');
  assert.equal(duration(45), '45 мин');
  assert.equal(initials('  Унесённые   ветром '), 'УВ');
});
