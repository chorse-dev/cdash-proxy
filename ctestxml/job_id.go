// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func GenerateJobID(project, site, stamp, build string) string {
	hasher := md5.New()
	json.NewEncoder(hasher).Encode([]string{project, site, stamp, build})
	return hex.EncodeToString(hasher.Sum(nil))
}
