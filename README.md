# DB Benchs

## Preparation
1. Install redis + postgresql
2. Fill in corresponding login details
3. Launch both servers, then fill in corresponding:
   - server data directory (for postgresql)
   - / PID (for redis)

## Run
1. Restart postgresql, as it flushes db and correctly persist them on disk. Make sure you restart it before continuing to get decent results
2. `go run .`
3. Observe results