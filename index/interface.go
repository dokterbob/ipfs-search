package index

import (
	"context"
)

// Index represents an index which stores and retrieves document properties.
type Index interface {
	Index(ctx context.Context, id string, properties map[string]interface{}) error
	Update(ctx context.Context, id string, properties map[string]interface{}) error
	Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error)
}
