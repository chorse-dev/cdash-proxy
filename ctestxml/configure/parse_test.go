// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package configure

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/chorse-dev/cdash-proxy/model"
)

func TestParseDeprecationWarning(t *testing.T) {
	log := `CMake Deprecation Warning at CMakeLists.txt:7 (cmake_minimum_required):
  Compatibility with CMake < 3.10 will be removed from a future version of
  CMake.

  Update the VERSION argument <min> value.  Or, use the <min>...<max> syntax
  to tell CMake that the project requires at least <min> but has been updated
  to work with policies introduced by <max> or earlier.


-- The C compiler identification is Clang 18.1.3
-- Detecting C compiler ABI info
-- Detecting C compiler ABI info - done
-- Configuring done (32.7s)
`
	actual := Parse(log, 0)
	expected := []model.Diagnostic{{
		FilePath: "CMakeLists.txt",
		Line:     7,
		Column:   -1,
		Type:     "Warning",
		Option:   "cmake_minimum_required",
		Message: `Compatibility with CMake < 3.10 will be removed from a future version of
CMake.

Update the VERSION argument <min> value.  Or, use the <min>...<max> syntax
to tell CMake that the project requires at least <min> but has been updated
to work with policies introduced by <max> or earlier.`,
	}}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Parse Error:\n%s", diff)
	}
}

func TestParseWarning(t *testing.T) {
	log := `-- The CXX compiler identification is Clang 18.1.3
-- Detecting CXX compiler ABI info
-- Detecting CXX compiler ABI info - done
CMake Warning at examples/CMakeLists.txt:12 (message):
  Missing range support! Skip: identity_as_default_projection


Examples to be built: identity_direct_usage
-- Configuring done (0.9s)
`
	actual := Parse(log, 0)
	expected := []model.Diagnostic{{
		FilePath: "examples/CMakeLists.txt",
		Line:     12,
		Column:   -1,
		Type:     "Warning",
		Option:   "message",
		Message:  `Missing range support! Skip: identity_as_default_projection`,
	}}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Parse Error:\n%s", diff)
	}
}

func TestParseGenerateError(t *testing.T) {
	log :=
		`CMake Error:
  CTEST_USE_LAUNCHERS is enabled, but the RULE_LAUNCH_COMPILE global property
  is not defined.

  Did you forget to include(CTest) in the toplevel CMakeLists.txt ?


`
	actual := Parse(log, 1)
	expected := []model.Diagnostic{{
		FilePath: "CMakeLists.txt",
		Line:     -1,
		Column:   -1,
		Type:     "Error",
		Message: `CTEST_USE_LAUNCHERS is enabled, but the RULE_LAUNCH_COMPILE global property
is not defined.

Did you forget to include(CTest) in the toplevel CMakeLists.txt ?`,
	}}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Parse Error:\n%s", diff)
	}
}
