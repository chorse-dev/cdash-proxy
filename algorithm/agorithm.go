// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package algorithm

func AnyOf[T any](ts []T, fn func(T) bool) bool {
	for _, t := range ts {
		if fn(t) {
			return true
		}
	}
	return false
}

func Map[T, R any](ts []T, fn func(T) R) []R {
	r := make([]R, len(ts))
	for i, t := range ts {
		r[i] = fn(t)
	}
	return r
}

func MapEx[T, V any](ts []T, fn func(T) (V, error)) ([]V, error) {
	var err error
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i], err = fn(t)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
