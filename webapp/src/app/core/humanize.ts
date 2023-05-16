const prefixes = ['', 'ki', 'Mi', 'Gi', 'Ti', 'Pi'];

export function humanize_bytes(value: number, units: string = 'B') {
  let prefix: string = '';
  for (prefix of prefixes) {
    if (Math.abs(value) < 1024) {
      break;
    }
    value /= 1024;
  }
  return value.toFixed(1) + ' ' + prefix + units;
}
