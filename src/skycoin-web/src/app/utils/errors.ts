export function parseResponseMessage(body: string): string {
  if (body.startsWith('400') || body.startsWith('403')) {
    const parts = body.split(' - ', 2);

    return parts.length === 2
      ? parts[1].charAt(0).toUpperCase() + parts[1].slice(1)
      : body;
  }

  return body;
}
