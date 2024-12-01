package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"

	mylogger "github.com/ritobanbiswas/TDA596_mygologger"
)

type TaskState int

const (
	Idle TaskState = iota
	InProgress
	Completed
)

type Task struct {
	FileName  string
	TaskType  string // "map" or "reduce"
	State     TaskState
	WorkerID  int       // The worker currently assigned (if any)
	StartTime time.Time // Used to track timeouts
}

type Coordinator struct {
	MapTasks       []Task
	ReduceTasks    []Task
	CompletedTasks int
	NReduce        int
	NFiles         int
	Mutex          sync.Mutex // Ensures thread-safe access
	logger         *log.Logger
	AllTasksDone   bool
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

// AssignTask assigns a task to a worker or indicates if all tasks are complete.
func (c *Coordinator) AssignTask(req *TaskRequest, res *TaskResponse) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.logger.Printf("[AssignTask] Worker %d requested a task\n", req.WorkerID)

	if c.AllTasksDone {
		res.AllDone = true
		c.logger.Println("[AssignTask] All tasks are complete. Sending AllDone=true")
		return nil
	}

	// Assign a map task if available
	for i, task := range c.MapTasks {
		if task.State == Idle {
			c.MapTasks[i].State = InProgress
			c.MapTasks[i].WorkerID = req.WorkerID
			c.MapTasks[i].StartTime = time.Now()

			res.TaskType = "map"
			res.FileName = task.FileName
			res.TaskID = i

			res.NReduce = c.NReduce
			res.NFiles = c.NFiles
			c.logger.Printf("[AssignTask] Assigned map task %d to worker %d\n", i, req.WorkerID)
			return nil
		}
	}

	// Assign a reduce task if available
	for i, task := range c.ReduceTasks {
		if task.State == Idle {
			c.ReduceTasks[i].State = InProgress
			c.ReduceTasks[i].WorkerID = req.WorkerID
			c.ReduceTasks[i].StartTime = time.Now()

			res.TaskType = "reduce"
			res.TaskID = i

			res.NReduce = c.NReduce
			res.NFiles = c.NFiles
			c.logger.Printf("[AssignTask] Assigned reduce task %d to worker %d\n", i, req.WorkerID)
			c.logger.Printf("[AssignTask] Assigned, nreduce is %d and nfiles is %d\n", c.NReduce, c.NFiles)
			return nil
		}
	}

	c.logger.Println("[AssignTask] No available tasks at the moment")
	return nil
}

// MarkTaskCompleted is called by workers to indicate task completion.
func (c *Coordinator) MarkTaskCompleted(args *TaskCompleteArgs, reply *struct{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.logger.Printf("[MarkTaskCompleted] Worker %d completed %s task %d\n", args.WorkerID, args.TaskType, args.TaskID)

	if args.TaskType == "map" {
		if c.MapTasks[args.TaskID].WorkerID == args.WorkerID {
			c.MapTasks[args.TaskID].State = Completed
		}
	} else if args.TaskType == "reduce" {
		if c.ReduceTasks[args.TaskID].WorkerID == args.WorkerID {
			c.ReduceTasks[args.TaskID].State = Completed
		}
	}

	c.checkAllTasksDone()
	return nil
}

// checkAllTasksDone verifies if all map and reduce tasks are completed.
func (c *Coordinator) checkAllTasksDone() {
	allMapDone := true
	for _, task := range c.MapTasks {
		if task.State != Completed {
			allMapDone = false
			break
		}
	}

	if allMapDone {
		allReduceDone := true
		for _, task := range c.ReduceTasks {
			if task.State != Completed {
				allReduceDone = false
				break
			}
		}

		if allReduceDone {
			c.AllTasksDone = true
			c.logger.Println("[checkAllTasksDone] All tasks are completed")
		}
	}
}

// monitorTasks reassigns tasks that have been in progress for too long.
func (c *Coordinator) monitorTasks() {
	for {
		time.Sleep(time.Second)

		c.Mutex.Lock()
		for i, task := range c.MapTasks {
			if task.State == InProgress && time.Since(task.StartTime) > 10*time.Second {
				c.logger.Printf("[monitorTasks] Reassigning map task %d due to timeout\n", i)
				c.MapTasks[i].State = Idle
			}
		}
		for i, task := range c.ReduceTasks {
			if task.State == InProgress && time.Since(task.StartTime) > 10*time.Second {
				c.logger.Printf("[monitorTasks] Reassigning reduce task %d due to timeout\n", i)
				c.ReduceTasks[i].State = Idle
			}
		}
		c.Mutex.Unlock()
	}
}

// Done is called by main/mrcoordinator.go to check if all tasks are done.
func (c *Coordinator) Done() bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.logger.Printf("[Done] AllTasksDone=%v\n", c.AllTasksDone)
	return c.AllTasksDone
}

// server starts a thread to listen for RPCs from workers.
func (c *Coordinator) server() {
	c.logger.Println("[server] Starting RPC server")
	rpc.Register(c)
	rpc.HandleHTTP()
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	c.logger.Println("[server] RPC server is running")
}

// MakeCoordinator initializes the coordinator with map tasks and reduce tasks.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	const logdirname = "/var/tmp/"
	logger, flog, _ := mylogger.MyLogger(
		logdirname,
		"mr-coordinatorfile.log",
		os.Args[1],
	)
	defer flog.Close()
	logger.Printf("\n\n\n\n<============================>\n")
	logger.Printf("[MakeCoordinator] Initializing coordinator")
	c := Coordinator{
		MapTasks:    make([]Task, len(files)),
		ReduceTasks: make([]Task, nReduce),
		NReduce:     nReduce,
		NFiles:      len(files),
		logger:      logger,
	}
	c.logger.Printf("[MakeCoordinator] Initializing coordinator with NReduce =%d\n", c.NReduce)
	c.logger.Printf("[MakeCoordinator] Initializing coordinator with NFiles =%d\n", c.NFiles)
	for i, file := range files {
		c.MapTasks[i] = Task{
			FileName: file,
			TaskType: "map",
			State:    Idle,
		}
		c.logger.Printf("[MakeCoordinator] Map task %d created for file %s\n", i, file)
	}

	for i := 0; i < nReduce; i++ {
		c.ReduceTasks[i] = Task{
			TaskType: "reduce",
			State:    Idle,
		}
		c.logger.Printf("[MakeCoordinator] Reduce task %d created\n", i)
	}

	go c.monitorTasks()
	c.server()
	return &c
}
