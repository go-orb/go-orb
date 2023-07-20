# grpc error codes

Found here: https://harunpeksen.com/advanced-grpc-deadlines-cancellation-error-handling-multiplexing-13ade6ce5c45

| **Code** | **Name**            | **HTTP Mapping**          |
| -------- | ------------------- | ------------------------- |
| 0        | OK                  | 200 OK                    |
| 1        | CANCELLED           | 499 Client Closed Request |
| 2        | UNKNOWN             | 500 Internal Server Error |
| 3        | INVALID_ARGUMENT    | 400 Bad Request           |
| 4        | DEADLINE_EXCEEDED   | 504 Gateway Timeout       |
| 5        | NOT_FOUND           | 404 Not Found             |
| 6        | ALREADY_EXISTS      | 409 Conflict              |
| 7        | PERMISSION_DENIED   | 403 Forbidden             |
| 8        | RESOURCE_EXHAUSTED  | 429 Too Many Requests     |
| 9        | FAILED_PRECONDITION | 400 Bad Request           |
| 10       | ABORTED             | 409 Conflict              |
| 11       | OUT_OF_RANGE        | 400 Bad Request           |
| 12       | UNIMPLEMENTED       | 501 Not Implemented       |
| 13       | INTERNAL            | 500 Internal Server Error |
| 14       | UNAVAILABLE         | 503 Service Unavailable   |
| 15       | DATA_LOSS           | 500 Internal Server Error |
| 16       | UNAUTHENTICATED     | 401 Unauthorized          |
