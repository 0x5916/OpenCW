export function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function randomChar(arr: string): string {
  return arr[Math.floor(Math.random() * arr.length)];
}
