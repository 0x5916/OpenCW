// Validates the message catalogue: translation coverage and dead keys.
//
// English (the inlang base locale) is the source of truth. Other locales may
// lag, but the drift has to be visible rather than discovered in production, so
// the checks are:
//   1. every base key exists in every locale, with the same placeholders;
//   2. no locale carries keys the base locale does not know about;
//   3. no base key is unused by `m.<key>()` anywhere under `src/`;
//   4. no referenced key is missing from the base catalogue.
// Run with: npm run messages:validate

import { readFileSync, readdirSync, statSync } from 'node:fs';
import { dirname, extname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

type Catalogue = Record<string, string>;

const scriptsDir = dirname(fileURLToPath(import.meta.url));
const projectDir = join(scriptsDir, '..');

const { baseLocale, locales } = JSON.parse(
  readFileSync(join(projectDir, 'project.inlang/settings.json'), 'utf8')
) as { baseLocale: string; locales: string[] };

function readCatalogue(locale: string): Catalogue {
  return JSON.parse(readFileSync(join(projectDir, 'messages', `${locale}.json`), 'utf8'));
}

/** The input names a message interpolates, order-independent. */
function placeholders(value: string): string {
  return [...value.matchAll(/\{([A-Za-z0-9_]+)\}/g)]
    .map((match) => match[1])
    .sort()
    .join(',');
}

/** Every source file that can reference a message, excluding generated output. */
function sourceFiles(dir: string): string[] {
  const files: string[] = [];

  for (const entry of readdirSync(dir)) {
    const path = join(dir, entry);
    if (statSync(path).isDirectory()) {
      // The generated paraglide runtime declares every key, so it can never
      // tell us whether a key is actually used.
      if (entry !== 'paraglide') files.push(...sourceFiles(path));
    } else if (['.svelte', '.ts', '.js'].includes(extname(entry))) {
      files.push(path);
    }
  }

  return files;
}

/** Keys referenced by `m.<key>` in app code. */
function referencedKeys(): Set<string> {
  const used = new Set<string>();

  for (const file of sourceFiles(join(projectDir, 'src'))) {
    const source = readFileSync(file, 'utf8');
    for (const match of source.matchAll(/\bm\.([a-z][a-z0-9_]*)/g)) used.add(match[1]);
  }

  return used;
}

const base = readCatalogue(baseLocale);
const baseKeys = Object.keys(base);
const used = referencedKeys();
const problems: string[] = [];

for (const locale of locales) {
  if (locale === baseLocale) continue;

  const catalogue = readCatalogue(locale);
  const missing: string[] = [];
  const empty: string[] = [];
  const placeholderDrift: string[] = [];
  const orphans = Object.keys(catalogue).filter((key) => !(key in base));

  for (const key of baseKeys) {
    if (!(key in catalogue)) {
      missing.push(key);
      continue;
    }
    if (!catalogue[key].trim()) empty.push(key);
    else if (placeholders(catalogue[key]) !== placeholders(base[key])) {
      placeholderDrift.push(key);
    }
  }

  // Group missing keys into runs, so the report names the neighbour a
  // translation belongs next to instead of listing 70 keys flat.
  const runs: string[] = [];
  let run: string[] = [];
  let anchor: string | null = null;
  for (const key of baseKeys) {
    if (missing.includes(key)) {
      run.push(key);
    } else if (run.length > 0) {
      runs.push(`after "${anchor}": ${run.join(', ')}`);
      run = [];
    }
    if (!missing.includes(key)) anchor = key;
  }
  if (run.length > 0) runs.push(`after "${anchor}": ${run.join(', ')}`);

  for (const entry of runs) problems.push(`[${locale}] untranslated — ${entry}`);
  for (const key of orphans) {
    problems.push(`[${locale}] unknown key (not in ${baseLocale}): ${key}`);
  }
  for (const key of empty) problems.push(`[${locale}] empty value: ${key}`);
  for (const key of placeholderDrift)
    problems.push(`[${locale}] placeholder mismatch vs ${baseLocale}: ${key}`);
}

for (const key of baseKeys) {
  if (!used.has(key)) problems.push(`dead key (referenced by no source file): ${key}`);
}

for (const key of used) {
  if (!(key in base)) {
    problems.push(`missing from ${baseLocale}.json but referenced in src/: ${key}`);
  }
}

if (problems.length > 0) {
  for (const problem of problems) console.error(problem);
  console.error(
    `\nmessages:validate failed — ${problems.length} problem(s) across ${locales.length} locales and ${baseKeys.length} keys`
  );
  process.exit(1);
}

console.log(
  `messages:validate passed — ${baseKeys.length} keys, ${locales.length} locales, every key referenced`
);
