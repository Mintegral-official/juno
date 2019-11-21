package juno

import "github.com/Mintegral-official/juno/index"

type IndexBuilder interface {
	build() index.Index
}
