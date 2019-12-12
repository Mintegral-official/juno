package builder

import "github.com/Mintegral-official/juno/index"

type Builder interface {
	Build() *index.IndexImpl
}
