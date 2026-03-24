// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    // interface Locals {}
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

declare module 'svelte/elements' {
  interface HTMLTextareaAttributes {
    autocorrect?: 'on' | 'off';
  }
  // Optionally extend input too, for consistency
  interface HTMLInputAttributes {
    autocorrect?: 'on' | 'off';
  }
}

export {};
