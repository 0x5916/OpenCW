<script lang="ts">
  import { Radio, Smartphone, MessagesSquare } from 'lucide-svelte';
  import { localizeHref } from '$lib/paraglide/runtime';
  import * as m from '$lib/paraglide/messages';

  function href(path: string) {
    return localizeHref(path);
  }
</script>

<!-- Hero -->
<section class="hero">
  <h1 class="hero-title">OpenCW</h1>
  <p class="hero-sub">
    {m.home_hero_subtitle_pre()}<a href={href('/about')} class="hero-link">{m.home_hero_koch()}</a
    >{m.home_hero_subtitle_post()}
  </p>
  <a href={href('/morse/learn')} class="btn-cta">{m.home_cta()}</a>
</section>

<hr class="divider" />

<!-- Feature cards -->
<section class="feature-grid">
  <a href={href('/about')} class="card feature-link-card">
    <div class="feature-icon" aria-hidden="true"><Radio /></div>
    <h2 class="feature-title">{m.home_feature_koch_title()}</h2>
    <p class="feature-body">{m.home_feature_koch_body()}</p>
  </a>
  <div class="card">
    <div class="feature-icon" aria-hidden="true"><Smartphone /></div>
    <h2 class="feature-title">{m.home_feature_mobile_title()}</h2>
    <p class="feature-body">{m.home_feature_mobile_body()}</p>
  </div>
  <div class="card">
    <div class="feature-icon" aria-hidden="true"><MessagesSquare /></div>
    <h2 class="feature-title">{m.home_feature_forum_title()}</h2>
    <p class="feature-body">{m.home_feature_forum_body()}</p>
  </div>
</section>

<hr class="divider" />

<!-- About -->
<section class="about-grid">
  <div>
    <h2 class="section-title">{m.home_about_site_title()}</h2>
    <h3 class="sub-title">{m.home_about_site_q1()}</h3>
    <p class="body-text about-answer">{m.home_about_site_a1()}</p>
    <h3 class="sub-title">{m.home_about_site_q2()}</h3>
    <p class="body-text about-answer">{m.home_about_site_a2()}</p>
  </div>
  <div>
    <h2 class="section-title">{m.home_about_me_title()}</h2>
    <p class="body-text about-bio">{m.home_about_me_p1()}</p>
    <p class="body-text">
      {m.home_about_me_p2_pre()}
      <a href="https://github.com/0x5916" class="link">{m.home_about_me_github()}</a>.
    </p>
  </div>
</section>

<style>
  .hero {
    text-align: center;
    padding-top: 4rem;
    padding-bottom: 4rem;
  }

  .hero-title {
    color: var(--accent);
    font-size: 3rem;
    line-height: 1;
    font-weight: 800;
    letter-spacing: -0.035em;
    margin-bottom: 1rem;
  }

  .hero-sub {
    color: var(--text-secondary);
    font-size: 1.125rem;
    line-height: 1.75rem;
    max-width: 44rem;
    margin-left: auto;
    margin-right: auto;
    margin-bottom: 2rem;
  }

  .hero-link {
    color: var(--accent);
    font-weight: 600;
  }

  .feature-link-card {
    text-decoration: none;
    color: inherit;
  }

  .about-answer {
    margin-bottom: 1rem;
  }

  .about-bio {
    margin-bottom: 0.75rem;
  }

  .about-grid {
    margin-top: 2rem;
    margin-bottom: 2rem;
    display: grid;
    gap: 1.5rem;
  }

  @media (min-width: 640px) {
    .about-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }

  .feature-grid {
    display: grid;
    gap: 1.5rem;
    margin-top: 2rem;
    margin-bottom: 2rem;
  }

  @media (min-width: 640px) {
    .feature-grid {
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
  }

  .feature-icon {
    color: var(--accent);
    font-size: 1.875rem;
    line-height: 2.25rem;
    margin-bottom: 0.75rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .feature-icon :global(svg) {
    width: 1.875rem;
    height: 1.875rem;
  }

  .feature-title {
    color: var(--text-primary);
    font-weight: 700;
    font-size: 1.125rem;
    line-height: 1.75rem;
    margin-bottom: 0.25rem;
  }

  .feature-body {
    color: var(--text-secondary);
    font-size: 0.875rem;
    line-height: 1.25rem;
  }
</style>
