## How to compile
### With race test
```
go build -race -buildmode=plugin -gcflags="all=-N -l" ../mrapps/wc.go
```

```
go run -race -gcflags="all=-N -l" mrcoordinator.go pg-*.txt
```
> On different session
```
go run -race -gcflags="all=-N -l" mrworker.go wc.so
```

### if no race test

```
go build  -buildmode=plugin -gcflags="all=-N -l" ../mrapps/wc.go

go run -gcflags="all=-N -l" mrcoordinator.go pg-*.txt
```

Again in different terminal.
```
go run -gcflags="all=-N -l" mrworker.go wc.so
```

## Test 

**All test passed**
```
RBISWAS@KW2LV6YD2D main % bash test-mr.sh   
*** Cannot find timeout command; proceeding without timeouts.
*** Starting wc test.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- wc test: PASS
*** Starting indexer test.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- indexer test: PASS
*** Starting map parallelism test.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- reduce parallelism test: PASS
*** Starting job count test.
2024/12/02 00:20:13 [PerformReduceTask] Failed to open input file mr-7-0: open mr-7-0: no such file or directory
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- job count test: PASS
*** Starting early exit test.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- early exit test: PASS
*** Starting crash test.
2024/12/02 00:20:34 [PerformReduceTask] Failed to open input file mr-2-0: open mr-2-0: no such file or directory
2024/12/02 00:20:35 [PerformReduceTask] Failed to open input file mr-2-1: open mr-2-1: no such file or directory
2024/12/02 00:20:36 [PerformReduceTask] Failed to open input file mr-2-2: open mr-2-2: no such file or directory
2024/12/02 00:20:36 [PerformReduceTask] Failed to open input file mr-2-3: open mr-2-3: no such file or directory
2024/12/02 00:20:37 [PerformReduceTask] Failed to open input file mr-2-4: open mr-2-4: no such file or directory
2024/12/02 00:20:37 [PerformReduceTask] Failed to open input file mr-2-5: open mr-2-5: no such file or directory
2024/12/02 00:20:37 [PerformReduceTask] Failed to open input file mr-2-6: open mr-2-6: no such file or directory
2024/12/02 00:20:38 [PerformReduceTask] Failed to open input file mr-2-7: open mr-2-7: no such file or directory
2024/12/02 00:20:38 [PerformReduceTask] Failed to open input file mr-2-8: open mr-2-8: no such file or directory
2024/12/02 00:20:38 [PerformReduceTask] Failed to open input file mr-2-9: open mr-2-9: no such file or directory
2024/12/02 00:20:45 [PerformReduceTask] Failed to open input file mr-2-0: open mr-2-0: no such file or directory
2024/12/02 00:20:46 [PerformReduceTask] Failed to open input file mr-2-1: open mr-2-1: no such file or directory
2024/12/02 00:20:46 [PerformReduceTask] Failed to open input file mr-2-2: open mr-2-2: no such file or directory
2024/12/02 00:20:47 [PerformReduceTask] Failed to open input file mr-2-4: open mr-2-4: no such file or directory
2024/12/02 00:20:47 [PerformReduceTask] Failed to open input file mr-2-3: open mr-2-3: no such file or directory
2024/12/02 00:20:48 [PerformReduceTask] Failed to open input file mr-2-5: open mr-2-5: no such file or directory
2024/12/02 00:20:48 [PerformReduceTask] Failed to open input file mr-2-6: open mr-2-6: no such file or directory
2024/12/02 00:20:49 [PerformReduceTask] Failed to open input file mr-2-7: open mr-2-7: no such file or directory
2024/12/02 00:20:49 [PerformReduceTask] Failed to open input file mr-2-8: open mr-2-8: no such file or directory
2024/12/02 00:20:49 [PerformReduceTask] Failed to open input file mr-2-9: open mr-2-9: no such file or directory
2024/12/02 00:20:55 [PerformReduceTask] Failed to open input file mr-2-0: open mr-2-0: no such file or directory
2024/12/02 00:20:56 [PerformReduceTask] Failed to open input file mr-2-1: open mr-2-1: no such file or directory
2024/12/02 00:20:56 [PerformReduceTask] Failed to open input file mr-2-2: open mr-2-2: no such file or directory
2024/12/02 00:20:57 [PerformReduceTask] Failed to open input file mr-2-3: open mr-2-3: no such file or directory
2024/12/02 00:20:57 [PerformReduceTask] Failed to open input file mr-2-4: open mr-2-4: no such file or directory
2024/12/02 00:20:58 [PerformReduceTask] Failed to open input file mr-2-5: open mr-2-5: no such file or directory
2024/12/02 00:20:59 [PerformReduceTask] Failed to open input file mr-2-6: open mr-2-6: no such file or directory
2024/12/02 00:20:59 [PerformReduceTask] Failed to open input file mr-2-7: open mr-2-7: no such file or directory
2024/12/02 00:21:00 [PerformReduceTask] Failed to open input file mr-2-8: open mr-2-8: no such file or directory
2024/12/02 00:21:00 [PerformReduceTask] Failed to open input file mr-2-9: open mr-2-9: no such file or directory
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
[Worker] All tasks are done. Exiting.
--- crash test: PASS
*** PASSED ALL TESTS
RBISWAS@KW2LV6YD2D main % 

```
