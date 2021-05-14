package judger

type ResultCode = int

const (
	PASS                  ResultCode = iota // 0
	WRONG_ANSWER                            // 1
	RUNTIME_ERROR                           // 2
	COMPILE_ERROR                           // 3
	UNKNOWN_ERROR                           // 4
	TIME_LIMIT_EXCEEDED                     // 5
	MEMORY_LIMIT_EXCEEDED                   // 6
)
