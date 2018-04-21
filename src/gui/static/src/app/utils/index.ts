export function parseResponseMessage(body: string): string {
  if (/^[0-9]{3}/.test(body)) {
    const parts = body.split(' - ', 2);

    return parts.length === 2
      ? parts[1].charAt(0).toUpperCase() + parts[1].slice(1)
      : body;
  }

  return body;
}
