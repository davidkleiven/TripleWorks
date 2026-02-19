package integrity

import (
	"encoding/json"
)

type QualityCheck interface {
	Check() QualityResult
}

type QualityResult interface {
	Report(encoder *json.Encoder) error
}
