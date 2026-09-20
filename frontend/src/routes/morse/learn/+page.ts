import { redirect } from '@sveltejs/kit';

// The trainer moved from `/morse/learn` to `/learn`; see src/routes/morse/+page.ts.
export const load = () => {
  redirect(308, '/learn');
};
