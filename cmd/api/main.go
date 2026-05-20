package api

import (
	"fmt"
	"log"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-6/pkg/logger"
)

func main() {
	zapLogger, err := logger.NewProductionLogger()
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer zapLogger.Sync()

	err = config.MustLoad("")
	if err != nil {
		zapLogger.Fatal(fmt.Sprintf("warning: loading configuration from .env unreadable: %v\n", err))
	}
	// err = godotenv.Load(".env")
	// if err != nil {
	// 	zapLogger.Fatal(fmt.Sprintf("warning: assuming default configuration: .env unreadable: %v\n", err))
	// 	// logg := logger.NewProductionLogger()
	// 	// fmt.Printf("warning: assuming default configuration: .env unreadable: %v\n", err)
	// }
}
