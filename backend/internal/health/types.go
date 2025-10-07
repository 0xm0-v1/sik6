package health

import "context"

// Checker is defined once for the package.
type Checker func(ctx context.Context) error
