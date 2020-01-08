package helpers

type SBuilder interface {
	StringBuilder(cap int, value ...interface{}) string
}
