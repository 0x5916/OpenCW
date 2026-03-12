import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
  const stored = cookies.get('learn.lesson');
  const lesson = stored ? Number(stored) || 1 : 1;
  return { lesson };
};
