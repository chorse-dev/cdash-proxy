// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GenerateJobID(project, site, stamp, build string) string {
	hasher := md5.New()
	fmt.Fprintf(hasher, "%s-%s-%s-%s", project, site, stamp, build)
	return hex.EncodeToString(hasher.Sum(nil))
}
