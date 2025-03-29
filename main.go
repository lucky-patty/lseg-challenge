package main 

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "strings"
  "time"
)

type LogEntry struct {
  Timestamp         time.Time `json:"time_stamp"`
  JobDescription    string    `json:"job_description"`
  LogEntry          string    `json:"log_entry"`
  PID               string    `json:"pid"`
}

const (
  ColorReset = "\033[0m"
  ColorRed   = "\033[31m"
  ColorGreen = "\033[32m"
  ColorYellow = "\033[33m"

)

// Summary of the task
var (
  totalJobs       int 
  missingEndCount int 
  passCount       int 
  warnCount       int 
  failCount       int 
)

// Parse the line to JobEntry that we want
func parseLine(line string) (*LogEntry, error) {
  // We split chunks of information with ','
  parts := strings.Split(line, ",")
  if len(parts) != 4 {
    return nil, fmt.Errorf("Invalid log format on line: %s", line)
  }

  // Parse HH:MM:SS time format 
  timestamp, err := time.Parse("15:04:05", strings.TrimSpace(parts[0]))
  if err != nil {
    return nil, fmt.Errorf("Invalid timestamp format: %v", err)
  }

  // Log Entry can either be START or END only 
  entry := strings.TrimSpace(parts[2])
  if entry != "START" && entry != "END" {
    return nil, fmt.Errorf("Invalid log entry type on line: %s  (must be START or END)", entry)
  }

  return &LogEntry{
    Timestamp: timestamp,
    JobDescription: strings.TrimSpace(parts[1]),
    LogEntry: entry,
    PID: strings.TrimSpace(parts[3]),
  }, nil
}


func main() {
  // Put up a flag file
  filePath := flag.String("file", "", "Path to the log file (required)")
  flag.Parse()

  // Make sure that log file must not be EMPTY
  if *filePath == "" {
    fmt.Println("Usage: go run main.go --file path/to/logfile.log")
    os.Exit(1)
  }

  file, err := os.Open(*filePath)
  if err != nil {
    log.Fatalf("Failed to open log file: %v", err)
  }

  defer file.Close()


  startTimes := make(map[string]time.Time)
  endTimes := make(map[string]time.Time)

  scanner := bufio.NewScanner(file)
  var invalidLines []string 
  
  // Scan the file
  for scanner.Scan() {
    line := scanner.Text()
    entry, err := parseLine(line)
    if err != nil {
      fmt.Printf("%s[WARNING]%s Skipping line: %v", ColorYellow, ColorReset , err)
      invalidLines = append(invalidLines, line)
      continue
    }

    switch entry.LogEntry {
    case "START":
          startTimes[entry.PID] = entry.Timestamp 
    case "END":
          endTimes[entry.PID] = entry.Timestamp 
    }
  }

  // Let user know how many line we skip 
  if len(invalidLines) > 0 {
    fmt.Printf("\n %s Skipped %d invalid lines. %s \n", ColorRed, len(invalidLines), ColorReset)
  }

  // Check the time
  for pid, start := range startTimes {
    totalJobs++ 
    end, exists := endTimes[pid]
    if !exists {
      missingEndCount++
      fmt.Printf("%s[WARNING]%s Missing END for job %s \n", ColorYellow, ColorReset, pid)
      continue
    }

    duration := end.Sub(start)
    mins := int(duration.Minutes())
    secs := int(duration.Seconds()) % 60

    fmt.Printf("%s[INFO]%s PID %s took %dm %ds \n", ColorGreen, ColorReset, pid, mins, secs)

    if duration > 10 * time.Minute {
      failCount++
      fmt.Printf("%s[ERROR]%s PID %s exceed 10 minutes \n", ColorRed, ColorReset, pid)
    } else if duration > 5 * time.Minute {
      warnCount++
      fmt.Printf("%s[WARNING]%s Job %s exceeds 5 minutes \n", ColorYellow, ColorReset, pid)
    } else {
      passCount++
    }
  }

  // Summary of the task
  fmt.Println("\n==================== Summary =================")
  fmt.Printf("%s %d Projects in total %s \n", ColorGreen, totalJobs, ColorReset)
  fmt.Printf("%s %d project(s) missing END for job %s \n", ColorYellow, missingEndCount, ColorReset)
  fmt.Printf("%s %d passed within 5 minutes %s \n", ColorGreen, passCount, ColorReset)
  fmt.Printf("%s %d exceed 5 minutes but not 10 %s \n", ColorYellow, warnCount, ColorReset)
  fmt.Printf("%s %d failed (more than 10 minutes) %s \n", ColorRed, failCount, ColorReset)
  fmt.Println("==================== END ====================")


}
