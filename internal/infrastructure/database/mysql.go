package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GiaBao0510/CleanArchitecture-templates/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDB(
	ctx context.Context,
	cfg config.DatabaseConfig,
) (*sql.DB, error) {
	
}