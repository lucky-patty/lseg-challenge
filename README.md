# LSEG Challenge
This is the technical task for the position of **Senior Developer Experience (DevEx) Engineer** at LSEG.

## Requirement
1. The program must parse `CSV` log file
2. Identify each job or task and track its start and finish times
3. Calculate the duration of each job from the time it started to the time it finished
4. Produce a report or output that
   - Logs a warning if a job takes longer than 5 minutes
   - Logs an error if a job takes longer than 10 minutes
  
**Note:** I add the summary of all jobs at the end of the program to increase user experience
## Installation Instruction
```
go get
```

## How to run the program
```
go build -o log-analyzer
go run log-analyzer --file path/to/log.log
```
**Notes**: You can change the name `log-analyzer` to anything just make sure to call it the same name when you run the program

## Sample Output
```
==================== Summary =================
 45 Projects in total  
 2 project(s) missing END for job  
 24 passed within 5 minutes  
 9 exceed 5 minutes but not 10  
 10 failed (more than 10 minutes)  
==================== END ====================
```
