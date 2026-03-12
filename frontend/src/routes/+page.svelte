<script lang="ts">
  import { Radio, Smartphone, MessagesSquare } from 'lucide-svelte';
  import * as m from '$lib/paraglide/messages';
</script>

<!-- Hero -->
<section class="hero">
  <h1 class="hero-title">OpenCW</h1>
  <p class="hero-sub">
    {m.home_hero_subtitle_pre()}<a
      href="/about"
      style="color: var(--accent); font-weight: 600;">{m.home_hero_koch()}</a
    >{m.home_hero_subtitle_post()}
  </p>
  <a href="/morse/learn" class="btn-cta">{m.home_cta()}</a>
</section>

<hr class="divider" />

<!-- Feature cards -->
<section class="feature-grid">
  <a href="/about" class="card" style="text-decoration: none; color: inherit;">
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
    <p class="body-text mb-4">{m.home_about_site_a1()}</p>
    <h3 class="sub-title">{m.home_about_site_q2()}</h3>
    <p class="body-text mb-4">{m.home_about_site_a2()}</p>
  </div>
  <div>
    <h2 class="section-title">{m.home_about_me_title()}</h2>
    <p class="body-text mb-3">{m.home_about_me_p1()}</p>
    <p class="body-text">
      {m.home_about_me_p2_pre()}
      <a href="https://github.com/0x5916" class="link">{m.home_about_me_github()}</a>.
    </p>
  </div>
</section>
