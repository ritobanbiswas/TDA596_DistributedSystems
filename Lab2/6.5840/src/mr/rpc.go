package mr

import (
	"os"
	"runtime"
	"strconv"
)

func currentFunction() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// Example RPC argument and reply structures.
type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type TaskRequest struct {
	WorkerID int
}

type TaskResponse struct {
	TaskType string
	FileName string
	TaskID   int
	NReduce  int
	//ReduceID int
	NFiles  int
	AllDone bool
}

type TaskCompleteArgs struct {
	TaskType string
	TaskID   int
	WorkerID int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	// fmt.Printf("[coordinatorSock] Generated socket name: %s\n", s)
	return s
}
