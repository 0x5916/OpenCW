import { randomChar, randomInt } from '$lib/random';

export const LESSONS: string[] = [
  'KM',
  'R',
  'S',
  'U',
  'A',
  'P',
  'T',
  'L',
  'O',
  'W',
  'I',
  '.',
  'N',
  'J',
  'E',
  'F',
  '0',
  'Y',
  ',',
  'V',
  'G',
  '5',
  '/',
  'Q',
  '9',
  'Z',
  'H',
  '3',
  '8',
  'B',
  '?',
  '4',
  '2',
  '7',
  'C',
  '1',
  'D',
  '6',
  'X'
];
export const MORSE: Record<string, string> = {
  A: '.-',
  B: '-...',
  C: '-.-.',
  D: '-..',
  E: '.',
  F: '..-.',
  G: '--.',
  H: '....',
  I: '..',
  J: '.---',
  K: '-.-',
  L: '.-..',
  M: '--',
  N: '-.',
  O: '---',
  P: '.--.',
  Q: '--.-',
  R: '.-.',
  S: '...',
  T: '-',
  U: '..-',
  V: '...-',
  W: '.--',
  X: '-..-',
  Y: '-.--',
  Z: '--..',
  '0': '-----',
  '1': '.----',
  '2': '..---',
  '3': '...--',
  '4': '....-',
  '5': '.....',
  '6': '-....',
  '7': '--...',
  '8': '---..',
  '9': '----.',
  // Punctuation
  '.': '.-.-.-',
  ',': '--..--',
  '/': '-..-.',
  '?': '..--..'
};

function getGroupSize(groupSize: number | null) {
  return groupSize ?? randomInt(2, 7);
}

function getLessonChars(lesson: number): string {
  return LESSONS.slice(0, lesson).join('');
}

export function getFarnsworthWpmSet(charWpm: number, effWpm: number) {
  const charDot = 1.2 / charWpm;
  const tFarn = (60 / effWpm - charDot * 31) / 19;

  const dash = charDot * 3;
  const symbolSpace = charDot;

  const letterSpace = tFarn * 3;
  const wordSpace = tFarn * 7;

  return { charDot, tFarn, dash, symbolSpace, letterSpace, wordSpace };
}

export function calculateDuration(text: string, charWpm: number, effWpm: number): number {
  const { charDot, dash, symbolSpace, letterSpace, wordSpace } = getFarnsworthWpmSet(
    charWpm,
    effWpm
  );

  let t = 0;
  const upper = text.toUpperCase();

  for (let i = 0; i < upper.length; i++) {
    const ch = upper[i];
    const morse = MORSE[ch];
    if (!morse) {
      t += wordSpace;
      continue;
    }
    for (let j = 0; j < morse.length; j++) {
      const dotOrDash = morse[j];
      const duration = dotOrDash === '.' ? charDot : dash;
      t += duration;
      if (j < morse.length - 1) t += symbolSpace;
    }
    if (i < upper.length - 1) t += letterSpace;
  }

  return t;
}

export function generateTimedLesson(
  lesson: number,
  targetDur: number,
  charWpm: number,
  effWpm: number,
  groupSize: number | null = null
) {
  let result = '';
  let total = 0;

  const lessonChars: string = getLessonChars(lesson);
  const gap = calculateDuration(' ', charWpm, effWpm);

  while (true) {
    const size = getGroupSize(groupSize);
    let group = '';

    for (let i = 0; i < size; i++) {
      group += randomChar(lessonChars);
    }

    const groupDur = calculateDuration(group, charWpm, effWpm);
    const withGapDur = groupDur + gap;
    const overshoot = total + withGapDur - targetDur;

    if (overshoot > 0) {
      const left = targetDur - total;

      if (left < overshoot) break;

      result += group;
      total += groupDur;
      break;
    }
    result += group + ' ';
    total += withGapDur;
  }
  return result;
}
