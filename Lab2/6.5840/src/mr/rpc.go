package mr

import (
	"fmt"
	"os"
	"strconv"
)

// Example RPC argument and reply structures.
type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	fmt.Printf("[coordinatorSock] Generated socket name: %s\n", s)
	return s
}
