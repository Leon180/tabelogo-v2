package migrations

import "errors"

var (
	// ErrNilVersion 當沒有 migration 執行過時
	ErrNilVersion = errors.New("no migration has been run")

	// ErrNoChange 當沒有變更需要執行時
	ErrNoChange = errors.New("no change")

	// ErrDirtyState migration 處於 dirty 狀態
	ErrDirtyState = errors.New("migration is in dirty state")

	// ErrInvalidVersion 無效的版本號
	ErrInvalidVersion = errors.New("invalid migration version")

	// ErrDatabaseNotReady 資料庫未就緒
	ErrDatabaseNotReady = errors.New("database is not ready")
)
