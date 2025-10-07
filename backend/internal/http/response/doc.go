// Package response provides utilities for building consistent HTTP responses,
// primarily in JSON.
//
// It includes:
//   - writing helpers (`WriteJSON`, `WriteNoBody`) to ensure proper headers
//     and uniform JSON encoding,
//   - a small envelope type (`Envelope`) to structure responses with a status,
//     data, or error message,
//   - middleware such as `MethodGuard` to restrict allowed HTTP methods,
//   - and helpers like `HeadAware` to adapt handlers to HTTP HEAD semantics.
//
// These utilities aim to simplify API development by reducing boilerplate
// and encouraging HTTP-compliant, idiomatic response handling.
package response
