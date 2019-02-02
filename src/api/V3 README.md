# REST API specification v3

## API Version 2

`/api/v3` endpoints have a standard format.

All `/api/v3` `POST` endpoints accept only `application/json` and return `application/json`.

All `/api/v3` `GET` requires accept data in JSON body.
In the future we may have choose to have `GET` requests also accept `POST` with a JSON body,
to support requests with a large query body, such as when requesting data for a large number
of addresses or transactions.

`/api/v3` responses are always JSON. If there is an error, the JSON object will
look like this:

```json
{
    "error": {
        "code": 400,
        "message": "bad arguments",
    }
}
```

Response data will be included in a `"data"` field, which will always be a JSON object (not an array).

Some endpoints may return both `"error"` and `"data"`. This will be noted in the documentation for that endpoint.

All responses will set an appropriate HTTP status code indicating an error, and it will be equal to the value of `response["error"]["code"]`.

Since `/api/v2` is still under development, there are no guarantees for backwards compatibility.
However, any changes to the API will be recorded in the [changelog](../../CHANGELOG.md).

Under some circumstances an error response body may not be valid JSON.
Any client consuming the API should accomodate this and conditionally parse JSON for non-`200` responses.

## Authentication

Authentication can be enabled with the `-web-interface-username` and `-web-interface-password` options.
The username and password should be provided in an `Authorization: Basic` header.

Authentication can only be enabled when using HTTPS with `-web-interface-https`, unless `-web-interface-plaintext-auth` is enabled.

## CSRF

All `POST`, `PUT` and `DELETE` requests require a CSRF token, obtained with a `GET /api/v1/csrf` call.
The token must be placed in the `X-CSRF-Token` header. A token is only valid
for 30 seconds and it is expected that the client obtains a new CSRF token
for each request. Requesting a CSRF token invalidates any previous CSRF token.

A request rejected for invalid or expired CSRF will respond with `403 Forbidden - invalid CSRF token`
as the response body.

### Get current csrf token

API sets: any

```
URI: /api/v3/csrf
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/csrf
```

Result:

```json
{
    "csrf_token": "klSgXoMOFTvEnt8KptBvHjhlFnW0OIkzyFVn4i8frDvIus9iLsFukqA9sM9Rxf3pLZHRLr82vBQxTq50vbYA8g"
}
```

