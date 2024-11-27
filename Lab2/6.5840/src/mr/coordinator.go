package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type map_task_ds struct {
	map_task_id           int
	map_task_intermediate string
	map_task_status       bool
}

type worker_task_ds struct {
	worker_task_id int
	// worker_task_intermediate []string    // @TODO<ritoban> this might be required let see later
	worker_task_status bool //  this can be 0 -   1-   2-  //@TODO<ritoban> let's see later how it shapes up
}

type Coordinator struct {
	// Your definitions here.
	map_task_list    []map_task_ds
	worker_task_list []worker_task_ds
	//count_worker int
	coordinator_name string // may not be required
	nReduce          int
}

// Your code here -- RPC handlers for the worker to call.
//@TODO<ritoban> Lets see later

// ok := call("Coordinator.Example", &args, &reply)
// if ok {
// 	// reply.Y should be 100.
// 	fmt.Printf("reply.Y %v\n", reply.Y)
// } else {
// 	fmt.Printf("call failed!\n")
// }

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		nReduce:          nReduce,
		map_task_list:    make([]map_task_ds, len(files)),
		worker_task_list: make([]worker_task_ds, nReduce),
	}

	// Your code here.
	// init mapTasks
	for i, file := range files {
		c.map_task_list[i] = map_task_ds{map_task_id: i, map_task_intermediate: file, map_task_status: false}
	}
	// init reduceTasks
	for i := 0; i < nReduce; i++ {
		c.worker_task_list[i] = worker_task_ds{worker_task_id: i, worker_task_status: false}
	}

	c.server()
	return &c
}
