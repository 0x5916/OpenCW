<script lang="ts">
  import { ArrowLeft, Send } from '@lucide/svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { user } from '$lib/auth';
  import { createForumThread, getForumCategories, type ForumCategory } from '$lib/api';
  import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { localizeApiError } from '$lib/errorLocalization';
  import { localizedHref as href } from '$lib/i18n.svelte';
  import * as m from '$lib/paraglide/messages';

  let categories = $state<ForumCategory[]>([]);
  let categoryId = $state('');
  let title = $state('');
  let body = $state('');
  let loading = $state(true);
  let submitting = $state(false);
  let error = $state('');

  $effect(() => {
    void loadCategories();
  });

  async function loadCategories() {
    if (!$user) {
      loading = false;
      return;
    }
    try {
      categories = await getForumCategories();
      categoryId = page.url.searchParams.get('category') ?? categories[0]?.id ?? '';
    } catch (cause) {
      error = localizeApiError(cause, () => m.forum_load_error());
    } finally {
      loading = false;
    }
  }

  async function submitThread(event: SubmitEvent) {
    event.preventDefault();
    if (!categoryId || !title.trim() || !body.trim()) {
      error = m.forum_required_fields();
      return;
    }
    submitting = true;
    error = '';
    try {
      const result = await createForumThread(categoryId, title.trim(), body.trim());
      await goto(href(`/forum/thread/${result.thread.id}`));
    } catch (cause) {
      error = localizeApiError(cause, () => m.api_error_forum_thread_create_failed());
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head><title>{m.forum_new_thread()} | OpenCW</title></svelte:head>

<main class="composer-page page-narrow">
  <a class="back-link" href={href('/forum')}><ArrowLeft size={16} /> {m.forum_title()}</a>
  <h1 class="page-title">{m.forum_new_thread()}</h1>
  {#if loading}<div class="state"><LoadingSpinner /></div>
  {:else if !$user}<div class="card state">
      <p>{m.forum_login_to_post()}</p>
      <a class="btn-primary" href={href('/login')}>{m.nav_login()}</a>
    </div>
  {:else if error && categories.length === 0}<ErrorAlert message={error} />
  {:else}
    <form class="composer card" onsubmit={submitThread}>
      {#if error}<ErrorAlert message={error} />{/if}
      <label for="category">{m.forum_category_label()}</label>
      <select id="category" class="select" bind:value={categoryId} required>
        <option value="" disabled>{m.forum_category_placeholder()}</option>
        {#each categories as category (category.id)}<option value={category.id}
            >{category.name}</option
          >{/each}
      </select>
      <label for="title">{m.forum_thread_title_label()}</label>
      <input
        id="title"
        bind:value={title}
        maxlength="200"
        required
        placeholder={m.forum_thread_title_placeholder()}
      />
      <label for="body">{m.forum_post_body_label()}</label>
      <textarea
        id="body"
        bind:value={body}
        rows="9"
        maxlength="10000"
        required
        placeholder={m.forum_post_body_placeholder()}></textarea>
      <button class="btn-primary" type="submit" disabled={submitting}
        ><Send size={16} />{submitting ? m.forum_publishing() : m.forum_publish()}</button
      >
    </form>
  {/if}
</main>

<style>
  .composer-page {
    padding-top: 2.25rem;
    padding-bottom: 3rem;
  }
  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.85rem;
  }
  .back-link:hover {
    color: var(--accent);
  }
  .page-title {
    margin: 1.5rem 0;
  }
  .composer {
    display: grid;
    gap: 0.75rem;
  }
  label {
    color: var(--text-label);
    font-weight: 650;
  }
  input,
  textarea {
    width: 100%;
    padding: 0.75rem;
    color: var(--text-primary);
    background: var(--bg-inset);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    font: inherit;
  }
  textarea {
    resize: vertical;
  }
  input:focus,
  select:focus,
  textarea:focus {
    outline: none;
    box-shadow: var(--focus-ring);
    border-color: var(--accent);
  }
  .state {
    display: grid;
    place-items: center;
    gap: 1rem;
    min-height: 12rem;
    text-align: center;
  }
  .state p {
    margin: 0;
  }
</style>
