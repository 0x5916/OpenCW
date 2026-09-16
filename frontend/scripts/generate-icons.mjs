// Generates the full favicon/PWA icon set from the original vector source.
// Source: src/lib/assets/favicon.svg (1024x1024, transparent, green #10b981 dot/dash).
// Outputs into static/:
//   - pwa-{192,512}x{192,512}.png + pwa-maskable-*   (web app manifest icons)
//   - apple-touch-icon.png (180x180)                  (iOS home screen)
//   - favicon.ico (16/32/48 multi-size)               (legacy ICO favicon)
//   - favicon.svg                                     (minified copy of the source)
// NOTE: the default OG image reuses /pwa-512x512.png (see src/lib/seo.ts).
import { mkdir, writeFile } from 'node:fs/promises';
import { readFile } from 'node:fs/promises';
import { Buffer } from 'node:buffer';
import { fileURLToPath } from 'node:url';
import path from 'node:path';
import sharp from 'sharp';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const source = path.join(root, 'src/lib/assets/favicon.svg');
const outDir = path.join(root, 'static');

const MASKABLE_BACKGROUND = '#030712';
const PWA_SIZES = [192, 512];

/** Render the transparent "any" artwork at a given size. */
function renderAny(size) {
  return sharp(source).resize(size, size).png().toBuffer();
}

/** Render the artwork centered on a solid background ("maskable" purpose). */
async function renderMaskable(size) {
  const art = await renderAny(size);
  return sharp({
    create: { width: size, height: size, channels: 4, background: MASKABLE_BACKGROUND }
  })
    .composite([{ input: art, gravity: 'centre' }])
    .png()
    .toBuffer();
}

/** Build a multi-size .ico from PNG buffers (PNG-in-ICO, supported on Vista+). */
function buildIco(entries) {
  const headerSize = 6;
  const entrySize = 16;
  const header = Buffer.alloc(headerSize + entrySize * entries.length);
  header.writeUInt16LE(0, 0); // reserved
  header.writeUInt16LE(1, 2); // type: icon
  header.writeUInt16LE(entries.length, 4); // image count

  let offset = headerSize + entrySize * entries.length;
  entries.forEach(({ size, data }, i) => {
    const base = headerSize + entrySize * i;
    const dim = size >= 256 ? 0 : size; // 0 means 256 in ICO
    header.writeUInt8(dim, base); // width
    header.writeUInt8(dim, base + 1); // height
    header.writeUInt8(0, base + 2); // color count
    header.writeUInt8(0, base + 3); // reserved
    header.writeUInt16LE(1, base + 4); // planes
    header.writeUInt16LE(32, base + 6); // bit count
    header.writeUInt32LE(data.length, base + 8); // bytes in resource
    header.writeUInt32LE(offset, base + 12); // image offset
    offset += data.length;
  });

  return Buffer.concat([header, ...entries.map((entry) => entry.data)]);
}

/** Minify the source SVG into a single line without the XML declaration. */
function minifySvg(svg) {
  return svg
    .replace(/<\?xml[^>]*\?>\s*/i, '')
    .replace(/>\s+</g, '><')
    .replace(/\s{2,}/g, ' ')
    .trim();
}

async function generate() {
  await mkdir(outDir, { recursive: true });

  // PWA manifest icons ("any" + "maskable").
  for (const size of PWA_SIZES) {
    const anyFile = path.join(outDir, `pwa-${size}x${size}.png`);
    await writeFile(anyFile, await renderAny(size));
    console.log(`wrote ${path.relative(root, anyFile)}`);

    const maskableFile = path.join(outDir, `pwa-maskable-${size}x${size}.png`);
    await writeFile(maskableFile, await renderMaskable(size));
    console.log(`wrote ${path.relative(root, maskableFile)}`);
  }

  // iOS home-screen icon (180x180).
  const appleFile = path.join(outDir, 'apple-touch-icon.png');
  await writeFile(appleFile, await renderAny(180));
  console.log(`wrote ${path.relative(root, appleFile)}`);

  // Legacy multi-size ICO favicon (16/32/48).
  const icoEntries = [];
  for (const size of [16, 32, 48]) {
    icoEntries.push({ size, data: await renderAny(size) });
  }
  const icoFile = path.join(outDir, 'favicon.ico');
  await writeFile(icoFile, buildIco(icoEntries));
  console.log(`wrote ${path.relative(root, icoFile)}`);

  // Minified SVG favicon copy, kept in sync with the source.
  const svgFile = path.join(outDir, 'favicon.svg');
  await writeFile(svgFile, minifySvg(await readFile(source, 'utf8')));
  console.log(`wrote ${path.relative(root, svgFile)}`);
}

generate().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
