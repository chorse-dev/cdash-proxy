// SPDX-FileCopyrightText: 2025 Daniel Pfeifer <daniel@pfeifer-mail.de>
// SPDX-License-Identifier: ISC

package ctestxml

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/chorse-dev/cdash-proxy/model"
)

func TestWarnings(t *testing.T) {
	data := `
		<Failure type="Warning">
			<Action>
				<TargetName>mputil</TargetName>
				<Language>C++</Language>
				<SourceFile>src/util.cpp</SourceFile>
				<OutputFile>CMakeFiles/util.dir/src/util.cpp.o</OutputFile>
				<OutputType>object file</OutputType>
			</Action>
			<Command>
				<WorkingDirectory>/home/dp/.cache/build/20ffe0b2269564c31a83df6d79c0bacc</WorkingDirectory>
				<Argument>/usr/bin/ccache</Argument>
				<Argument>/usr/bin/clazy</Argument>
				<Argument>-I/source/include</Argument>
				<Argument>-I/home/dp/.cache/build/20ffe0b2269564c31a83df6d79c0bacc/include</Argument>
				<Argument>-std=gnu++20</Argument>
				<Argument>-MD</Argument>
				<Argument>-MT</Argument>
				<Argument>CMakeFiles/util.dir/src/util.cpp.o</Argument>
				<Argument>-MF</Argument>
				<Argument>CMakeFiles/util.dir/src/util.cpp.o.d</Argument>
				<Argument>-o</Argument>
				<Argument>CMakeFiles/util.dir/src/util.cpp.o</Argument>
				<Argument>-c</Argument>
				<Argument>/source/src/util.cpp</Argument>
			</Command>
			<Result>
				<StdOut/>
				<StdErr>In file included from /source/src/util.cpp:6:
[CTest: warning matched] /source/include/util.h:135:1: warning: mp::UnlockGuard has dtor but not copy-ctor, copy-assignment [-Wclazy-rule-of-three]
  135 | struct UnlockGuard
      | ^
[CTest: warning matched] /source/include/util.h:152:1: warning: mp::DestructorCatcher has dtor but not copy-ctor, copy-assignment [-Wclazy-rule-of-three]
  152 | struct DestructorCatcher
      | ^
[CTest: warning matched] /source/src/util.cpp:82:22: warning: Pass small and trivially-copyable type by value (const kj::ArrayPtr&lt;const char&gt; &amp;) [-Wclazy-function-args-by-value]
   82 |     string.visit([&amp;](const kj::ArrayPtr&lt;const char&gt;&amp; piece) {
      |                      ^~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
      |                      kj::ArrayPtr&lt;const char&gt; piece
[CTest: warning matched] 3 warnings generated.</StdErr>
				<ExitCondition>0</ExitCondition>
			</Result>
		</Failure>
	`
	failure := &Failure{}
	if err := xml.Unmarshal([]byte(data), failure); err != nil {
		t.Errorf("Failed to parse XML: %v\n", err)
		return
	}
	actual := failure.Diagnostics()
	expected := []model.Diagnostic{
		{
			FilePath: "include/util.h",
			Line:     135,
			Column:   1,
			Type:     "Warning",
			Message:  "mp::UnlockGuard has dtor but not copy-ctor, copy-assignment",
			Option:   "-Wclazy-rule-of-three",
		},
		{
			FilePath: "include/util.h",
			Line:     152,
			Column:   1,
			Type:     "Warning",
			Message:  "mp::DestructorCatcher has dtor but not copy-ctor, copy-assignment",
			Option:   "-Wclazy-rule-of-three",
		},
		{
			FilePath: "src/util.cpp",
			Line:     82,
			Column:   22,
			Type:     "Warning",
			Message:  "Pass small and trivially-copyable type by value (const kj::ArrayPtr<const char> &)",
			Option:   "-Wclazy-function-args-by-value",
		},
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Mismatch (-expected +actual):\n%s", diff)
	}
}

func TestError(t *testing.T) {
	data := `
		<Failure type="Error">
			<Action>
				<TargetName>fpe</TargetName>
				<Language>C</Language>
				<SourceFile>Failures/fpe.c</SourceFile>
				<OutputFile>Failures/CMakeFiles/fpe.dir/fpe.c.o</OutputFile>
				<OutputType>object file</OutputType>
			</Action>
			<Command>
				<WorkingDirectory>/home/dp/.cache/build/83ca88f9fab95a07d56582a37692a8af</WorkingDirectory>
				<Argument>/usr/bin/ccache</Argument>
				<Argument>/usr/bin/gcc</Argument>
				<Argument>--coverage</Argument>
				<Argument>-Wall</Argument>
				<Argument>-Wextra</Argument>
				<Argument>-MD</Argument>
				<Argument>-MT</Argument>
				<Argument>Failures/CMakeFiles/fpe.dir/fpe.c.o</Argument>
				<Argument>-MF</Argument>
				<Argument>Failures/CMakeFiles/fpe.dir/fpe.c.o.d</Argument>
				<Argument>-o</Argument>
				<Argument>Failures/CMakeFiles/fpe.dir/fpe.c.o</Argument>
				<Argument>-c</Argument>
				<Argument>/home/dp/Projects/Example/Failures/fpe.c</Argument>
			</Command>
			<Result>
				<StdOut/>
				<StdErr>/home/dp/Projects/Example/Failures/fpe.c: In function ‘main’:
[CTest: warning matched] /home/dp/Projects/Example/Failures/fpe.c:10:20: warning: format ‘%d’ expects argument of type ‘int’, but argument 2 has type ‘double’ [-Wformat=]
   10 |   printf("Result: %d\n", result);
      |                   ~^     ~~~~~~
      |                    |     |
      |                    int   double
      |                   %f
/home/dp/Projects/Example/Failures/fpe.c:19:10: error: incompatible types when returning type ‘typeof (nullptr)’ but ‘int’ was expected
   19 |   return nullptr;
      |          ^~~~~~~
[CTest: warning matched] /home/dp/Projects/Example/Failures/fpe.c:7:7: warning: unused variable ‘unusedVar’ [-Wunused-variable]
    7 |   int unusedVar = 10;
      |       ^~~~~~~~~</StdErr>
				<ExitCondition>1</ExitCondition>
			</Result>
		</Failure>
	`
	failure := &Failure{}
	if err := xml.Unmarshal([]byte(data), failure); err != nil {
		t.Errorf("Failed to parse XML: %v\n", err)
		return
	}
	actual := failure.Diagnostics()
	expected := []model.Diagnostic{
		{
			FilePath: "Failures/fpe.c",
			Line:     10,
			Column:   20,
			Type:     "Warning",
			Message:  "format ‘%d’ expects argument of type ‘int’, but argument 2 has type ‘double’",
			Option:   "-Wformat=",
		},
		{
			FilePath: "Failures/fpe.c",
			Line:     19,
			Column:   10,
			Type:     "Error",
			Message:  "incompatible types when returning type ‘typeof (nullptr)’ but ‘int’ was expected",
		},
		{
			FilePath: "Failures/fpe.c",
			Line:     7,
			Column:   7,
			Type:     "Warning",
			Message:  "unused variable ‘unusedVar’",
			Option:   "-Wunused-variable",
		},
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Mismatch (-expected +actual):\n%s", diff)
	}
}

func TestMSan(t *testing.T) {
	data := `
		<Failure type="Error">
			<Action>
				<TargetName>TARGET_NAME</TargetName>
				<OutputFile>test/foo.c++</OutputFile>
			</Action>
			<Command>
				<WorkingDirectory>test</WorkingDirectory>
				<Argument>foo</Argument>
				<Argument>test</Argument>
				<Argument>test/foo</Argument>
				<Argument>include</Argument>
			</Command>
			<Result>
				<StdOut/>
				<StdErr>Uninitialized bytes in strlen at offset 4 inside [0x701000000000, 5)
==16852==WARNING: MemorySanitizer: use-of-uninitialized-value
    #0 0x7f959a2e2b46  (/lib64/libfoo.so.1.0.1+0x3bb46) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
    #1 0x7f959a2e31a9  (/lib64/libfoo.so.1.0.1+0x3c1a9) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
    #2 0x7f959a2df6b2  (/lib64/libfoo.so.1.0.1+0x386b2) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
    #3 0x7f959a2fae4f  (/lib64/libfoo.so.1.0.1+0x53e4f) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
    #4 0x7f959a2fb8a2  (/lib64/libfoo.so.1.0.1+0x548a2) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
    #5 0x0000004b1721  (foo+0x4b1721) (BuildId: ddceddf87caa17083f26399219eaeaf1ab41ea87)
    #6 0x7f9599cf75f4  (/lib64/libc.so.6+0x35f4) (BuildId: 2b3c02fe7e4d3811767175b6f323692a10a4e116)
    #7 0x7f9599cf76a7  (/lib64/libc.so.6+0x36a7) (BuildId: 2b3c02fe7e4d3811767175b6f323692a10a4e116)
    #8 0x000000400ab4  (foo+0x400ab4) (BuildId: ddceddf87caa17083f26399219eaeaf1ab41ea87)

SUMMARY: MemorySanitizer: use-of-uninitialized-value (/lib64/libfoo.so.1.0.1+0x3bb46) (BuildId: 0bb29394caeeb1d93fd4116f44090fae725ca4af)
Exiting</StdErr>
				<ExitCondition>1</ExitCondition>
			</Result>
		</Failure>
	`
	failure := &Failure{}
	if err := xml.Unmarshal([]byte(data), failure); err != nil {
		t.Errorf("Failed to parse XML: %v\n", err)
		return
	}
	actual := failure.Diagnostics()
	expected := []model.Diagnostic{
		{
			Line:     -1,
			Column:   -1,
			Type:     "Error",
			Message:  "Command finished with exit code 1",
		},
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Mismatch (-expected +actual):\n%s", diff)
	}
}
