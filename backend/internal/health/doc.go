// Package health exposes HTTP handlers for process liveness and service readiness.
//
// /livez   : liveness - process is up and running
// /readyz  : readiness - dependencies are ready to serve traffic
//
// Note: Kubernetes docs recommend /livez and /readyz; /healthz is deprecated.
// See: https://kubernetes.io/docs/reference/using-api/health-checks/
package health
