package params

import "time"

type BuildParams struct {
	SourcePath string // path of source code
	TimeLimit  time.Duration
	SizeLimit  int
}

var SourceSuffix = map[string]string{
	"go":  ".go",
	"c":   ".c",
	"cpp": ".cpp",
	"c++": ".cpp",
}

func NewBuild() *BuildParams {
	return &BuildParams{
		SourcePath: "",
		TimeLimit:  5000 * time.Microsecond,
		SizeLimit:  10,
	}
}
