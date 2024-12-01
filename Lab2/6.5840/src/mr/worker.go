package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sort"

	mylogger "github.com/ritobanbiswas/TDA596_mygologger"
)

type WorkerPrivate struct {
	logger     *log.Logger
	logdirname string
}
type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	const logdirname = "/var/tmp/"
	logger, flog, _ := mylogger.MyLogger(
		logdirname,
		"mr-workerfile.log",
		os.Args[1],
	)
	defer flog.Close()
	w := WorkerPrivate{
		logger:     logger,
		logdirname: logdirname,
	}
	w.logger.Printf("\n\n\n\n<============================>\n")
	w.logger.Printf("[Worker] Starting worker...")

	for {
		task := w.RequestTask()
		w.logger.Printf("[Worker] Received task: %+v\n", task)

		if task.AllDone {
			fmt.Println("[Worker] All tasks are done. Exiting.")
			break
		}

		if task.TaskType == "map" {
			w.logger.Printf("[Worker] Starting map task: ID=%d, File=%s\n", task.TaskID, task.FileName)
			w.PerformMapTask(task, mapf)
			w.logger.Printf("[Worker] Completed map task: ID=%d\n", task.TaskID)
		} else if task.TaskType == "reduce" {
			w.logger.Printf("[Worker] Starting reduce task: ID=%d\n", task.TaskID)
			w.PerformReduceTask(task, reducef)
			w.logger.Printf("[Worker] Completed reduce task: ID=%d\n", task.TaskID)
		}
		w.ReportTaskCompletion(task.TaskType, task.TaskID)
		w.logger.Printf("[Worker] Reported completion of task: Type=%s, ID=%d\n", task.TaskType, task.TaskID)
	}
}

// RequestTask asks the coordinator for a task.
func (w *WorkerPrivate) RequestTask() TaskResponse {
	args := TaskRequest{}
	reply := TaskResponse{}

	w.logger.Printf("[RequestTask] Requesting a task from the coordinator...")
	ok := call("Coordinator.AssignTask", &args, &reply)
	if !ok {
		log.Fatalf("[RequestTask] Worker failed to contact coordinator")
	}
	w.logger.Printf("[RequestTask] Received response: %+v\n", reply)
	return reply
}

// PerformMapTask executes a map task assigned by the coordinator.
func (w *WorkerPrivate) PerformMapTask(task TaskResponse, mapf func(string, string) []KeyValue) {
	w.logger.Printf("[PerformMapTask] Reading input file: %s\n", task.FileName)
	content, err := ioutil.ReadFile(task.FileName)
	if err != nil {
		log.Fatalf("[PerformMapTask] Failed to read input file %v: %v", task.FileName, err)
	}

	intermediate := mapf(task.FileName, string(content))
	w.logger.Printf("[PerformMapTask] Map function produced %d key-value pairs\n", len(intermediate))

	buckets := make([][]KeyValue, task.NReduce)
	for _, kv := range intermediate {
		bucket := ihash(kv.Key) % task.NReduce
		buckets[bucket] = append(buckets[bucket], kv)
	}

	for i := 0; i < task.NReduce; i++ {
		outputFile := fmt.Sprintf("mr-%d-%d", task.TaskID, i)
		outputFile = filepath.Join(w.logdirname, outputFile)
		w.logger.Printf("[PerformMapTask] Writing to file: %s\n", outputFile)
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("[PerformMapTask] Failed to create output file %v: %v", outputFile, err)
		}
		enc := json.NewEncoder(file)
		for _, kv := range buckets[i] {
			err := enc.Encode(&kv)
			if err != nil {
				log.Fatalf("[PerformMapTask] Failed to write key-value pair to file %v: %v", outputFile, err)
			}
		}
		file.Close()
	}
}

// PerformReduceTask executes a reduce task assigned by the coordinator.
func (w *WorkerPrivate) PerformReduceTask(task TaskResponse, reducef func(string, []string) string) {
	w.logger.Printf("[PerformReduceTask] Reading intermediate files for ReduceID=%d\n", task.TaskID)
	//w.logger.Printf("NREDUCE IS %d\n", task.NReduce)
	var intermediate []KeyValue
	w.logger.Printf("task.NFiles is = %d\n", task.NFiles)
	for i := 0; i < task.NFiles; /* len(files) /*task.NReduce*/ i++ {
		inputFile := fmt.Sprintf("mr-%d-%d", i, task.TaskID)
		inputFile = filepath.Join(w.logdirname, inputFile)
		w.logger.Printf("[PerformReduceTask] Opening file: %s\n", inputFile)
		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatalf("[PerformReduceTask] Failed to open input file %v: %v", inputFile, err)
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
		//w.logger.Printf("CLosing file!")
		file.Close()
		//w.logger.Printf("Intermediates are", intermediate)
	}

	sort.Sort(ByKey(intermediate))
	w.logger.Printf("[PerformReduceTask] Sorted %d key-value pairs\n", len(intermediate))

	outputFile := fmt.Sprintf("mr-out-%d", task.TaskID)
	outputFile = filepath.Join(w.logdirname, outputFile)
	w.logger.Printf("[PerformReduceTask] Writing reduce output to file: %s\n", outputFile)
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("[PerformReduceTask] Failed to create output file %v: %v", outputFile, err)
	}

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)
		fmt.Fprintf(file, "%v %v\n", intermediate[i].Key, output)
		i = j
	}
	file.Close()
}

// ReportTaskCompletion informs the coordinator that a task is complete.
func (w *WorkerPrivate) ReportTaskCompletion(taskType string, taskID int) {
	args := TaskCompleteArgs{TaskType: taskType, TaskID: taskID}
	reply := struct{}{}

	w.logger.Printf("[ReportTaskCompletion] Reporting task completion: Type=%s, ID=%d\n", taskType, taskID)
	ok := call("Coordinator.MarkTaskCompleted", &args, &reply)
	if !ok {
		log.Fatalf("[ReportTaskCompletion] Failed to report task completion to coordinator")
	}
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("[call] Dialing failed:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println("[call] RPC call failed:", err)
	return false
}
