package psql

import (
	"fmt"
	"github.com/penguin-statistics/db-benchmark/models"
	"github.com/penguin-statistics/db-benchmark/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type PSql struct {
	DB             *gorm.DB
	datadir        string
	testCandidates []string
}

func New(datadir string) *PSql {
	db, err := gorm.Open(postgres.Open("host=localhost user=root password=root dbname=penguinnative port=5432 sslmode=disable TimeZone=Asia/Shanghai"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Installation{})
	if err != nil {
		panic(err)
	}
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Installation{})
	return &PSql{
		DB:      db,
		datadir: datadir,
	}
}

func (p *PSql) Init(rowCount int) (err []error) {
	var errors []error
	for i := 0; i < rowCount; i++ {
		inst := utils.RandomInstallation()
		if thisE := p.DB.Create(inst).Error; thisE != nil {
			fmt.Println("error occurred when #Init postgresql:", thisE)
			errors = append(errors, thisE)
			continue
		}

		if utils.SampleHit(0.05) {
			p.testCandidates = append(p.testCandidates, inst.DeviceToken)
		}
	}
	return errors
}

func (p *PSql) Count() int64 {
	count := new(int64)
	err := p.DB.Model(&models.Installation{}).Count(count).Error
	if err != nil {
		return -1
	}
	return *count
}

func (p *PSql) Quota() uint64 {
	i, err := utils.CalculateDirSize(p.datadir)
	if err != nil {
		return 0
	}
	return uint64(i)
}

func (p *PSql) slowRead(times int) {
	l := len(p.testCandidates)
	for i := 0; i < times; i++ {
		candidate := p.testCandidates[utils.RandomIntn(l)]
		i := new(models.Installation)
		err := p.DB.Where("device_token = ?", candidate).First(i).Error
		if err != nil {
			panic(err)
		}
	}
}

func (p *PSql) Read(times int, goroutines int) {
	wg := &sync.WaitGroup{}
	fmt.Println("  + running in goroutines... ")
	for i := 0; i < goroutines; i++ {
		fmt.Print(i, " ")
		if i == goroutines - 1 {
			fmt.Print("\n  + all started. now waiting for goroutines to finish...\n")
		}
		i := i
		go func() {
			wg.Add(1)
			p.slowRead(times / goroutines)
			fmt.Print(i, ":OK ")
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("")
}

func (p *PSql) PrintStats() {
	fmt.Println("===>> PostgreSQL current status" +
		"\n	- count of records:", p.Count(),
		"\n	- db data directory size", p.Quota(), "bytes")
}
