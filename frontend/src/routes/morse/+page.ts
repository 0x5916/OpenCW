import { redirect } from '@sveltejs/kit';

// The trainer lived at `/morse` before moving to `/learn`. Keeping the route
// means the prerenderer writes a static redirect page (meta refresh + script)
// at `/morse.html`, so old links keep working on the static site.
export const load = () => {
  redirect(308, '/learn');
};
