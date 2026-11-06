// Package ut implements a few shared utilities.
package ut

func Is[T any](val any) bool {
	_, res := val.(T)
	return res
}
