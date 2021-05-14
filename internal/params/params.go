package params

type JudgeRequest struct {
	CaseNum   int
	TimeLimit int
	RamLimit  int
	VolumeDir string
}

var SourceSuffix = map[string]string{
	"go":  ".go",
	"c":   ".c",
	"cpp": ".cpp",
	"c++": ".cpp",
}

func NewJudgeRequest() *JudgeRequest {
	return &JudgeRequest{
		TimeLimit: 1000,
		RamLimit:  10,
	}
}
