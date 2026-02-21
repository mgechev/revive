package formatter_test

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"

	"github.com/mgechev/revive/formatter"
	"github.com/mgechev/revive/lint"
)

//revive:disable:line-length-limit
func TestFormatter(t *testing.T) {
	for name, td := range map[string]struct {
		formatter lint.Formatter
		want      string
	}{
		"checkstyle": {
			formatter: &formatter.Checkstyle{},
			want: `<?xml version='1.0' encoding='UTF-8'?>
<checkstyle version="5.0">
    <file name="err.go">
      <error line="33" column="4" message="replace fmt.Errorf by errors.New (confidence 0)" severity="error" source="revive/use-errors-new"/>
      <error line="38" column="4" message="replace fmt.Errorf by errors.New (confidence 0)" severity="error" source="revive/use-errors-new"/>
    </file>
    <file name="file.go">
      <error line="2" column="5" message="error var Exp should have name of the form ErrFoo (confidence 0)" severity="warning" source="revive/error-naming"/>
    </file>
</checkstyle>`,
		},
		"default": {
			formatter: &formatter.Default{},
			want: `file.go:2:5: error var Exp should have name of the form ErrFoo
err.go:33:4: replace fmt.Errorf by errors.New
err.go:38:4: replace fmt.Errorf by errors.New`,
		},
		"friendly": {
			formatter: &formatter.Friendly{},
			want: `  ⚠  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo
  file.go:2:5

  ✘  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  err.go:33:4

  ✘  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  err.go:38:4

✘ 3 problems (2 errors, 1 warning)

Errors:
  2  use-errors-new

Warnings:
  1  error-naming

`,
		},
		"json": {
			formatter: &formatter.JSON{},
			want: "[" +
				`{"Severity":"warning","Failure":"error var Exp should have name of the form ErrFoo","RuleName":"error-naming","Category":"naming","Position":{"Start":{"Filename":"file.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"file.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}` +
				"," +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":33,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":33,"Column":8}},"Confidence":0,"ReplacementLine":""}` +
				"," +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":38,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":38,"Column":9}},"Confidence":0,"ReplacementLine":""}` +
				"]",
		},
		"ndjson": {
			formatter: &formatter.NDJSON{},
			want: `{"Severity":"warning","Failure":"error var Exp should have name of the form ErrFoo","RuleName":"error-naming","Category":"naming","Position":{"Start":{"Filename":"file.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"file.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}` +
				"\n" +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":33,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":33,"Column":8}},"Confidence":0,"ReplacementLine":""}` +
				"\n" +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":38,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":38,"Column":9}},"Confidence":0,"ReplacementLine":""}` +
				"\n",
		},
		"plain": {
			formatter: &formatter.Plain{},
			want: `file.go:2:5: error var Exp should have name of the form ErrFoo https://revive.run/r#error-naming` +
				"\n" +
				`err.go:33:4: replace fmt.Errorf by errors.New https://revive.run/r#use-errors-new` +
				"\n" +
				`err.go:38:4: replace fmt.Errorf by errors.New https://revive.run/r#use-errors-new` +
				"\n",
		},
		"sarif": {
			formatter: &formatter.Sarif{},
			want: `{
  "runs": [
    {
      "results": [
        {
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "file.go"
                },
                "region": {
                  "startColumn": 5,
                  "startLine": 2
                }
              }
            }
          ],
          "message": {
            "text": "error var Exp should have name of the form ErrFoo"
          },
          "ruleId": "error-naming"
        },
        {
          "level": "error",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "err.go"
                },
                "region": {
                  "startColumn": 4,
                  "startLine": 33
                }
              }
            }
          ],
          "message": {
            "text": "replace fmt.Errorf by errors.New"
          },
          "ruleId": "use-errors-new"
        },
        {
          "level": "error",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "err.go"
                },
                "region": {
                  "startColumn": 4,
                  "startLine": 38
                }
              }
            }
          ],
          "message": {
            "text": "replace fmt.Errorf by errors.New"
          },
          "ruleId": "use-errors-new"
        }
      ],
      "tool": {
        "driver": {
          "informationUri": "https://revive.run",
          "name": "revive",
          "rules": [
            {
              "helpUri": "https://revive.run/r#use-errors-new",
              "id": "use-errors-new",
              "properties": {
                "severity": "error"
              }
            }
          ]
        }
      }
    }
  ],
  "version": "2.1.0"
}`,
		},
		"stylish": {
			formatter: &formatter.Stylish{},
			want: `err.go
  (33, 4)  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  (38, 4)  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New

file.go
  (2, 5)  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo


 ✖ 3 problems (2 errors) (1 warning)`,
		},
		"unix": {
			formatter: &formatter.Unix{},
			want: "file.go:2:5: [error-naming] error var Exp should have name of the form ErrFoo" +
				"\n" +
				"err.go:33:4: [use-errors-new] replace fmt.Errorf by errors.New" +
				"\n" +
				"err.go:38:4: [use-errors-new] replace fmt.Errorf by errors.New" +
				"\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			previousNoColor := color.NoColor
			color.NoColor = true
			realStdout := os.Stdout
			fakeStdout, err := os.Create(filepath.Join(t.TempDir(), "fakeStdout"))
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = fakeStdout
			t.Cleanup(func() {
				os.Stdout = realStdout
				color.NoColor = previousNoColor
			})
			failures := make(chan lint.Failure, 10)
			failures <- lint.Failure{
				Failure:  "error var Exp should have name of the form ErrFoo",
				RuleName: "error-naming",
				Category: lint.FailureCategoryNaming,
				Position: lint.FailurePosition{
					Start: token.Position{
						Filename: "file.go",
						Line:     2,
						Column:   5,
					},
					End: token.Position{
						Filename: "file.go",
						Line:     2,
						Column:   10,
					},
				},
			}
			failures <- lint.Failure{
				Failure:  "replace fmt.Errorf by errors.New",
				RuleName: "use-errors-new",
				Category: lint.FailureCategoryErrors,
				Position: lint.FailurePosition{
					Start: token.Position{
						Filename: "err.go",
						Line:     33,
						Column:   4,
					},
					End: token.Position{
						Filename: "err.go",
						Line:     33,
						Column:   8,
					},
				},
			}
			failures <- lint.Failure{
				Failure:  "replace fmt.Errorf by errors.New",
				RuleName: "use-errors-new",
				Category: lint.FailureCategoryErrors,
				Position: lint.FailurePosition{
					Start: token.Position{
						Filename: "err.go",
						Line:     38,
						Column:   4,
					},
					End: token.Position{
						Filename: "err.go",
						Line:     38,
						Column:   9,
					},
				},
			}
			close(failures)
			output, err := td.formatter.Format(failures, lint.Config{
				Confidence: 0.8,
				Rules: lint.RulesConfig{
					"use-errors-new": lint.RuleConfig{
						Severity: lint.SeverityError,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = realStdout
			err = fakeStdout.Close()
			if err != nil {
				t.Fatal(err)
			}
			stdout, err := os.ReadFile(fakeStdout.Name()) //nolint:gosec // ignore G703: Path traversal via taint analysis
			if err != nil {
				t.Fatal(err)
			}
			if len(stdout) > 0 {
				t.Errorf("formatter wrote to stdout: %q", stdout)
			}
			if td.want != output {
				t.Errorf("got:\n%s\nwant:\n%s\n", output, td.want)
			}
		})
	}
}
