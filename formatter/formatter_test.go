package formatter_test

import (
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/mgechev/revive/formatter"
	"github.com/mgechev/revive/lint"
)

func TestFormatter(t *testing.T) {
	lintFailure := lint.Failure{
		Failure:  "test failure",
		RuleName: "rule",
		Category: "cat",
		Position: lint.FailurePosition{
			Start: token.Position{
				Filename: "test.go",
				Line:     2,
				Column:   5,
			},
			End: token.Position{
				Filename: "test.go",
				Line:     2,
				Column:   10,
			},
		},
	}
	for _, td := range []struct {
		formatter lint.Formatter
		want      string
	}{
		{
			formatter: &formatter.Checkstyle{},
			want: `
<?xml version='1.0' encoding='UTF-8'?>
<checkstyle version="5.0">
    <file name="test.go">
      <error line="2" column="5" message="test failure (confidence 0)" severity="warning" source="revive/rule"/>
    </file>
</checkstyle>
`,
		},
		{
			formatter: &formatter.Default{},
			want:      `test.go:2:5: test failure`,
		},
		{
			formatter: &formatter.Friendly{},
			want: `
⚠  https://revive.run/r#rule  test failure  
  test.go:2:5

⚠ 1 problem (0 errors, 1 warning)

Warnings:
  1  rule
`,
		},
		{
			formatter: &formatter.JSON{},
			want:      `[{"Severity":"warning","Failure":"test failure","RuleName":"rule","Category":"cat","Position":{"Start":{"Filename":"test.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"test.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}]`,
		},
		{
			formatter: &formatter.NDJSON{},
			want:      `{"Severity":"warning","Failure":"test failure","RuleName":"rule","Category":"cat","Position":{"Start":{"Filename":"test.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"test.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}`,
		},
		{
			formatter: &formatter.Plain{},
			want:      `test.go:2:5: test failure https://revive.run/r#rule`,
		},
		{
			formatter: &formatter.Sarif{},
			want: `
{
  "runs": [
    {
      "results": [
        {
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "test.go"
                },
                "region": {
                  "startColumn": 5,
                  "startLine": 2
                }
              }
            }
          ],
          "message": {
            "text": "test failure"
          },
          "ruleId": "rule"
        }
      ],
      "tool": {
        "driver": {
          "informationUri": "https://revive.run",
          "name": "revive"
        }
      }
    }
  ],
  "version": "2.1.0"
}
`,
		},
		{
			formatter: &formatter.Stylish{},
			want: `
test.go
  (2, 5)  https://revive.run/r#rule  test failure  


 ✖ 1 problem (0 errors) (1 warnings)
`,
		},
		{
			formatter: &formatter.Unix{},
			want:      `test.go:2:5: [rule] test failure`,
		},
	} {
		t.Run(td.formatter.Name(), func(t *testing.T) {
			dir := t.TempDir()
			realStdout := os.Stdout
			fakeStdout, err := os.Create(dir + "/fakeStdout")
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = fakeStdout
			defer func() {
				os.Stdout = realStdout
			}()
			failures := make(chan lint.Failure, 10)
			failures <- lintFailure
			close(failures)
			output, err := td.formatter.Format(failures, lint.Config{})
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = realStdout
			err = fakeStdout.Close()
			if err != nil {
				t.Fatal(err)
			}
			stdout, err := os.ReadFile(fakeStdout.Name())
			if err != nil {
				t.Fatal(err)
			}
			if len(stdout) > 0 {
				t.Errorf("formatter wrote to stdout: %q", stdout)
			}
			got := strings.TrimSpace(output)
			want := strings.TrimSpace(td.want)
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
