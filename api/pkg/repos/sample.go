package repos

import (
	"context"
	"time"

	"github.com/fredbi/go-patterns/iterators"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	ItemsIterator = iterators.StructIterator[Item]

	// SampleRepo provides a simple CRUD database access to a table of shippable items.
	SampleRepo interface {
		Get(context.Context, string, ...ItemOption) (Item, error)
		Create(context.Context, Item) (string, error)
		Update(context.Context, Item) error
		Delete(context.Context, string) error
		List(context.Context, ...ItemOption) (ItemsIterator, error)
	}

	// Item represents a shippable item
	Item struct {
		ID                string    `json:"id" db:"id"`                                // NOT NULL
		Name              string    `json:"name" db:"name"`                            // NOT NULL
		WarehouseLocation string    `json:"warehouseLocation" db:"warehouse_location"` // NOT NULL
		LastUpdated       time.Time `json:"lastUpdated" db:"last_updated"`             // NOT NULL
		Weight            float64   `json:"weight" db:"weight"`                        // NOT NULL

		Dimensions   pgtype.FlatArray[float64] `json:"dimensions" db:"dimensions"`                // NULL
		Attributes   map[string]interface{}    `json:"attributes,omitempty" db:"attributes"`      // NULL
		DeliveryTime *time.Duration            `json:"deliveryTime,omitempty" db:"delivery_time"` // NULL
		Description  *string                   `json:"description,omitempty" db:"description"`    // NULL
		Tags         map[string]string         `json:"tags,omitempty" db:"tags"`                  // NULL
	}

	// ItemOption lets the caller specify filters
	ItemOption func(*itemOptions)

	itemOptions struct {
		withTags bool
	}
)

// WithItemTags instructs queries to join items against their tags
func WithItemTags(enabled bool) ItemOption {
	return func(o *itemOptions) {
		o.withTags = enabled
	}
}
