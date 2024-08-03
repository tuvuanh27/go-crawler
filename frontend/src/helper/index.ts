export function convertToUnixTime(dateString: string): number {
  const [year, month, day] = dateString.split('-').map(Number);

  // Create a Date object (month is 0-indexed, so subtract 1)
  const date = new Date(year, month - 1, day);

  // Get Unix time in seconds
  return Math.floor(date.getTime() / 1000);
}
