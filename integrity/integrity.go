package integrity

import (
	"context"
	"encoding/json"
)

type QualityCheck interface {
	Fetch(ctx context.Context, db *bun.DB) error
	Check()
}

type QualityResult interface {
	Report(encoder *json.Encoder) error
	Fix() iter.Seq[any]
}
