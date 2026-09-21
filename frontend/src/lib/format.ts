import { LESSONS } from '$lib/morse';
import { scoreGrade } from '$lib/score';
import * as m from '$lib/paraglide/messages';

/** Format an ISO timestamp as a short local date, e.g. "Sep 17, 2026". */
export function formatDate(value?: string): string {
  if (!value) return '';
  const date = new Date(value);
  return Number.isNaN(date.getTime())
    ? ''
    : date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

/** Display name for a forum author, falling back to a localized community label. */
export function authorLabel(author: { username?: string } | null | undefined): string {
  return author?.username ?? m.forum_community_member();
}

/** Format a lesson number or cumulative character string as "3 - K, M". */
export function formatLesson(lessonStr: string): string {
  const numericLesson = Number.parseInt(lessonStr, 10);
  if (Number.isInteger(numericLesson) && numericLesson >= 1 && numericLesson <= LESSONS.length) {
    return `${numericLesson} - ${LESSONS[numericLesson - 1].split('').join(', ')}`;
  }

  let cumulative = '';
  for (let i = 0; i < LESSONS.length; i++) {
    cumulative += LESSONS[i];
    if (cumulative === lessonStr.toUpperCase()) {
      return `${i + 1} - ${LESSONS[i].split('').join(', ')}`;
    }
  }

  // Fallback: truncate raw string
  return lessonStr.length > 12 ? lessonStr.slice(0, 12) + '…' : lessonStr;
}

/** Format a 0..1 ratio as a rounded percentage. */
export function percentage(value: number): string {
  return Math.round(value * 100) + '%';
}

/** CSS modifier for an accuracy ratio (0..1). */
export function accuracyClass(value: number): string {
  return `acc-${scoreGrade(value)}`;
}
