# cdash-proxy

[![PkgGoDev](https://pkg.go.dev/badge/github.com/chorse-dev/cdash-proxy)](https://pkg.go.dev/github.com/chorse-dev/cdash-proxy)
[![Coverage Status](https://coveralls.io/repos/github/chorse-dev/cdash-proxy/badge.svg?branch=master)](https://coveralls.io/github/chorse-dev/cdash-proxy?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/chorse-dev/cdash-proxy)](https://goreportcard.com/report/github.com/chorse-dev/cdash-proxy)

This code is used on https://chorse.dev/ to perform additional processing of the
data that is sent from CTest. The longterm goal of cdash-proxy is to become
obsolete, because ideally, all processing should be performed by CTest itself.

## Difference to CDash

While CDash has separate tables for `configure`, `build`, and `test`, we prefer
to insert most data into tables `commands` and `diagnostics`. **cdash-proxy**
converts the XML files from CTest into JSON of a simpler schema.

The schema for CTest's XML files can be found [here](ctestxml/model.rnc)
(It still lacks the recent extensions for instrumentation, though).

The structure for **cdash-proxy**'s JSON can be seen [here](model/model.go).

### Configure

Even though CDash has a column for configure warnings, neither CTest nor CDash
actually parse warnings from the configure log. **cdash-proxy** parses
`Site>Configure>Log` from `Configure.xml` using regular expressions and stores
errors and warnings as diagnostics.

Ideally, CTest should invoke `cmake` with the `--sarif-output` option and then
parse the given file and store diagnostics in `Configure.xml`.

### Build

CTest parses the build output for errors and warnings using expressions that are
customizable (`CTEST_CUSTOM_ERROR_MATCH`, etc) and stores the results in
`Site>Build>Error` and/or `Site>Build>Warning` in `Build.xml`.

However, if `CTEST_USE_LAUNCHERS` is set, it uses different matching rules (that
are not extensible) and it stores errors and warnings under
`Site>Build>Failure` **instead**. It does not split the output into separate
diagnostics, but inserts markers like `[CTest: warning matched]`.

If instrumentation is enabled, individual build commands are written to
`Site>Build>Commands` and `Site>Build>Targets>Target>Commands`.

**cdash-proxy** combines the information from lanchers and instrumentation if
available. It removes the `[CTest: warning matched]` markers and extracts
diagnostics.

Ideally, CTest should store `stdout` and `stderr` when instrumentation is
enabled (this would make `CTEST_USE_LAUNCHERS` obsolete). When wrapping the
compiler (using launchers or instrumentaton), CTest should instruct the compiler
to output diagnostics in SARIF or JSON, parse that, and then store them in
`Build.xml`.

CTest also produces some bogus like [this](https://github.com/Kitware/CMake/blob/3d3d3f94703e23d3d2cbff67537057474e3e0ff1/Source/CTest/cmCTestBuildHandler.cxx#L636) or [that](https://github.com/Kitware/CMake/blob/3d3d3f94703e23d3d2cbff67537057474e3e0ff1/Source/CTest/cmCTestBuildHandler.cxx#L645-L648) which probably is not useful for anyone, but no one dares to remove in fear of breaking CDash.

### Test

CTest reports all test names in `Site>Testing>TestList>Test`, followed by all
tests in `Site>Testing>Test`, where each test is named in `FullName`. Also, the
test's command line is stored both in `FullCommandLine` and in a measurement
named "Command Line".

Basically, the information stored about a test is not much different from a
configure or build command, only the XML structure is completely different.

**cdash-proxy** converts each test into a command.

### DynamicAnalysis

The information in `DynamicAnalysis.xml` is largely redundant with
`DynamicAnalysis-Test.xml`. **cdash-proxy** merges the information from both
files (actually the merge happens during insertion into the database, as CTest
uploads the two files separately).

### Coverage / CoverageLog

CTest stores the line number, line coverage, and the line content under
`Site>CoverageLog>File>Report>Line` in `CoverageLog.xml`. Only the line coverage
is interesting, because the file content can be retrieved from the repository
and the line number is just the index itself. **cdash-proxy** stores line
coverage in a simple array of ints.

### Notes / Upload

**cdash-proxy** treats both the same, as file attachments to the job.

### GCovTar

As an alternative to `ctest_coverage()`, this is used as:

```cmake
include(CTestCoverageCollectGCOV)
ctest_coverage_collect_gcov(TARBALL "${CTEST_BINARY_DIRECTORY}/gcov.tbz2"
  GCOV_COMMAND "${CTEST_COVERAGE_COMMAND}"
  )
ctest_submit(
  CDASH_UPLOAD "${CTEST_BINARY_DIRECTORY}/gcov.tbz2"
  CDASH_UPLOAD_TYPE GcovTar
  )
```

This derives from CTest's approach of doing processing on the client.

**cdash-proxy** opens the uploaded tar ball and extracts the coverage
information as if it was submitted as `Coverage.xml` and `CoverageLog.xml`.
Uncovered branches are recorded as diagnostics.

Ideally, CTest should not derive from the approach of doing processing on the
client. The server should be kept dumb.
