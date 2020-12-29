package main

import (
	"flag"
	"fmt"
	"github.com/penguin-statistics/db-benchmark/psql"
	"github.com/penguin-statistics/db-benchmark/redis"
	"github.com/penguin-statistics/db-benchmark/utils"
	"time"
)

func main() {
	rowCount := flag.Int("datacount", 5000, "data size to init")
	readCount := flag.Int("readcount", 5000, "times to test read")
	psqldir := flag.String("psqldir", "", "directory to the db file of your testing db")
	redispid := flag.Int("redispid", 0, "pid of the redis instance")
	flag.Parse()

	fmt.Println("testing DBs with", *rowCount, "records.")

	fmt.Println("initializing DBs")

	fmt.Println("\n==== psql ====")

	// https://www.postgresql.org/docs/9.0/storage-file-layout.html
	p := psql.New(*psqldir)
	p.PrintStats()
	t := time.Now()
	errs := p.Init(*rowCount)
	if len(errs) != 0 {
		fmt.Println("error when #Init psql. amount:", len(errs))
	}
	fmt.Println("#Init-ed in", utils.MillisecondsFrom(t), "ms")
	p.PrintStats()

	fmt.Println("\n==== redis ====")

	r := redis.New(*redispid)
	r.PrintStats()
	t = time.Now()
	errs = r.Init(*rowCount)
	if len(errs) != 0 {
		fmt.Println("error when #Init redis. amount:", len(errs))
	}
	fmt.Println("#Init-ed in", utils.MillisecondsFrom(t), "ms")
	r.PrintStats()

	fmt.Println("\n! DBs have now all done initialization")

	fmt.Println("\n! now testing random read speed", *readCount, "times")

	fmt.Println("\n==== psql ====")
	t = time.Now()
	p.Read(*readCount, 20)
	fmt.Println("psql #Read with 20 goroutines done in", utils.MillisecondsFrom(t), "ms")


	fmt.Println("\n==== redis ====")
	t = time.Now()
	r.Read(*readCount)
	fmt.Println("redis #Read done in", utils.MillisecondsFrom(t), "ms")

	fmt.Println("! done test.")
}
