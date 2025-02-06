// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

func getSubproject(subprojects []Subproject, labels []string) string {
	for _, label := range labels {
		for _, sub := range subprojects {
			if sub.Label == label {
				return sub.Name
			}
		}
	}
	return ""
}
