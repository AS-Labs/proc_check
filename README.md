# proc_check exporter

A Prometheus exporter written in Go that monitors a process specified by its name via a command-line argument. It exposes metrics such as process existence, CPU usage, memory usage, and command-line arguments, designed for integration with Prometheus and visualization in Grafana, including a table format for arguments.
Features

    Monitors a process by name, provided as a command-line argument.
    Exposes the following metrics:
        process_exists: Indicates if the process exists (1 or 0).
        process_cpu_usage: CPU usage percentage of the process.
        process_memory_usage_bytes: Memory usage in bytes (Resident Set Size, RSS).
        process_arg: Command-line arguments with labels for PID, process name, index, and value, suitable for a Grafana table.
    Runs an HTTP server on port 8081, serving metrics at /metrics.

Prerequisites

    Go: Version 1.16 or later (tested with Go 1.21).
    Dependencies:
        github.com/shirou/gopsutil/process for process monitoring.
        github.com/prometheus/client_golang/prometheus for Prometheus integration.

Setup
1. Clone the Repository
bash
git clone https://github.com/<your-username>/process-exporter.git
cd process-exporter
2. Install Go

    Download and install Go from golang.org/dl/.
    Verify installation:
    ```bash
    go version
    ```

3. Install Dependencies

Run the following commands to fetch required libraries:
```bash
go get github.com/shirou/gopsutil/process
go get github.com/prometheus/client_golang/prometheus
```
4. Build the Exporter

Compile the code into an executable:
```bash
go build -o exporter main.go
```


This creates an executable named exporter (or exporter.exe on Windows).
Usage

Run the exporter with a process name as a command-line argument:
```bash
./proc_check -process <process_name>
```
Examples

```bash
# Monitor Python processes:
./proc_check -process python
# Monitor Nginx processes:
./proc_check -process nginx
```
Note: it's better to add a systemd service file for the exporter.

The exporter will start and serve metrics at http://localhost:8081/metrics. You’ll see output like:
text
Starting exporter on :8081/metrics for process 'python'

To stop the exporter, press Ctrl+C.

## Testing Metrics

Verify the metrics endpoint with curl:
```bash
curl http://localhost:8081/metrics
```

Sample output might look like:
```text
# HELP process_exists Whether the process exists (1 = exists, 0 = does not)
# TYPE process_exists gauge
process_exists{name="python3",pid="1234"} 1
# HELP process_cpu_usage CPU usage percentage of the process
# TYPE process_cpu_usage gauge
process_cpu_usage{name="python3",pid="1234"} 12.5
# HELP process_memory_usage_bytes Memory usage of the process in bytes
# TYPE process_memory_usage_bytes gauge
process_memory_usage_bytes{name="python3",pid="1234"} 52428800
# HELP process_arg Command-line arguments of the process
# TYPE process_arg gauge
process_arg{index="0",name="python3",pid="1234",value="python3"} 1
process_arg{index="1",name="python3",pid="1234",value="script.py"} 1
process_arg{index="2",name="python3",pid="1234",value="--verbose"} 1
```

If no matching process is found:
text
process_exists{name="",pid=""} 0
Integration with Prometheus

Add the exporter to your prometheus.yml:
```yaml
scrape_configs:
  - job_name: 'process_exporter'
    static_configs:
      - targets: ['localhost:8081']
```


## Limitations

Multiple Processes: Reports metrics for all processes matching the name, distinguished by pid.
CPU Usage: Calculated per scrape interval; may require multiple scrapes for accuracy.
Performance: Scans all processes on each scrape, which could be optimized for large systems.

## Contributing

Feel free to open issues or submit pull requests for enhancements or bug fixes.
License

This project is licensed under the MIT License - see the  file for details.
