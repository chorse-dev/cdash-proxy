// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package algorithm

func NewPointer[T any](value T) *T {
	return &value
}
