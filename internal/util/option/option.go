package option

import (
	"flag"
	"os"

	"github.com/NevostruevK/GopherMart.git/internal/util/logger"
	"github.com/caarlos0/env/v7"
)

const (
	defRunAddress           = "127.0.0.1:8080"
	defDatabaseUri          = ""
	defAccrualSystemAddress = "127.0.0.1:8081"
)

type Options struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DatabaseUri          string `env:"DATABASE_URI" envDefault:""`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
}

func GetOptions() (*Options, error) {
	runAddressPtr := flag.String("a", defRunAddress, "gophermart server address in format  Host:Port")
	databaseUriPtr := flag.String("d", defDatabaseUri, "data base address")
	accrualSystemAddressPtr := flag.String("r", defAccrualSystemAddress, "accrual system address in format  Host:Port")
	flag.Parse()
	lg := logger.NewLogger("options : ", 0)
	opt := Options{}
	err := env.Parse(&opt)
	if err != nil {
		lg.Printf("ERROR : env.Parse returned error %v", err)
	}

	if _, ok := os.LookupEnv("RUN_ADDRESS"); !ok || err != nil {
		opt.RunAddress = *runAddressPtr
	}
	lg.Printf("RUN_ADDRESS %s", opt.RunAddress)
	if _, ok := os.LookupEnv("DATABASE_URI"); !ok || err != nil {
		opt.DatabaseUri = *databaseUriPtr
	}
	lg.Printf("DATABASE_URI %s", opt.DatabaseUri)
	if _, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); !ok || err != nil {
		opt.AccrualSystemAddress = *accrualSystemAddressPtr
	}
	lg.Printf("ACCRUAL_SYSTEM_ADDRESS %s", opt.AccrualSystemAddress)

	return &opt, err
}
