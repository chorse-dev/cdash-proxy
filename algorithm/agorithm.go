// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package algorithm

func Map[T, R any](ts []T, fn func(T) R) []R {
	r := make([]R, len(ts))
	for i, t := range ts {
		r[i] = fn(t)
	}
	return r
}
