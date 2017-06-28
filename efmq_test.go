package efmq_test

import (
	"testing"

	"github.com/olliephillips/efmq"
)

func TestNewEMFQ(t *testing.T) {
	const badInterface = "bad1"
	if _, err := efmq.NewEFMQ(badInterface); err == nil {
		t.Error(err)
	}
}
