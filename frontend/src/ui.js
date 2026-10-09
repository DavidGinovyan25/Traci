export const genres = {
  ACTION: 'Боевик', ADVENTURE: 'Приключения', ANIMATION: 'Анимация', COMEDY: 'Комедия',
  CRIME: 'Криминал', DOCUMENTARY: 'Документальный', DRAMA: 'Драма', FAMILY: 'Семейный',
  FANTASY: 'Фэнтези', HISTORY: 'Исторический', HORROR: 'Ужасы', MUSIC: 'Музыкальный',
  MYSTERY: 'Детектив', ROMANCE: 'Мелодрама', SCIENCE_FICTION: 'Фантастика',
  THRILLER: 'Триллер', WAR: 'Военный', WESTERN: 'Вестерн',
};

export const statuses = {
  PLANNED: 'Хочу посмотреть', WATCHING: 'Смотрю', WATCHED: 'Просмотрено', DROPPED: 'Не досмотрел',
};

export function escapeHTML(value) {
  return String(value ?? '').replace(/[&<>"']/g, char => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;',
  })[char]);
}

export function options(values, selected = '', placeholder = null) {
  const empty = placeholder === null ? '' : `<option value="">${escapeHTML(placeholder)}</option>`;
  return empty + Object.entries(values).map(([value, title]) =>
    `<option value="${value}" ${selected === value ? 'selected' : ''}>${escapeHTML(title)}</option>`).join('');
}

export function initials(name) {
  return name.trim().split(/\s+/).slice(0, 2).map(word => Array.from(word)[0] || '').join('').toUpperCase();
}

export function duration(minutes) {
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  return hours ? `${hours} ч${rest ? ` ${rest} мин` : ''}` : `${rest} мин`;
}

export function query(values) {
  return new URLSearchParams(Object.entries(values).filter(([, value]) => value !== '' && value != null)).toString();
}

export function collectionPatch(values) {
  return {
    status: values.status,
    personal_rating: values.personal_rating === '' ? null : Number(values.personal_rating),
    review: values.review.trim() || null,
  };
}
