import { MatSnackBar, MatSnackBarConfig } from '@angular/material';

export function parseResponseMessage(body: string): string {
  if (typeof body === 'object') {
    body = body['_body'];
  }

  if (body.startsWith('400') || body.startsWith('403')) {
    const parts = body.split(' - ', 2);

    return parts.length === 2
      ? parts[1].charAt(0).toUpperCase() + parts[1].slice(1)
      : body;
  }

  return body;
}

export function showSnackbarError(snackbar: MatSnackBar, body: string, duration = 300000) {
  const config = new MatSnackBarConfig();
  config.duration = duration;

  snackbar.open(parseResponseMessage(body), null, config);
}
