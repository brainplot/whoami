package api

import (
	"context"
)

type Stresser interface {
	Stress(context.Context) error
}
