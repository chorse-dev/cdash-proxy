// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package algorithm

func FindIf[T any](ts []T, pred func(T) bool) *T {
	for _, t := range ts {
		if pred(t) {
			return &t
		}
	}
	return nil
}
