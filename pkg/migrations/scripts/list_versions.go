// 列出所有可用的 migration 版本
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type MigrationFile struct {
	Version     uint
	Description string
	HasUp       bool
	HasDown     bool
}

func main() {
	var migrationsPath string
	flag.StringVar(&migrationsPath, "path", "", "Path to migrations directory")
	flag.Parse()

	if migrationsPath == "" {
		fmt.Println("Usage: go run list_versions.go -path migrations/auth")
		os.Exit(1)
	}

	migrations, err := listMigrations(migrationsPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return
	}

	// 顯示標題
	fmt.Printf("\n%-10s %-50s %-6s %-6s\n", "Version", "Description", "Up", "Down")
	fmt.Println("---------------------------------------------------------------------------------------------")

	// 顯示每個 migration
	for _, m := range migrations {
		upMark := "❌"
		if m.HasUp {
			upMark = "✅"
		}
		downMark := "❌"
		if m.HasDown {
			downMark = "✅"
		}

		fmt.Printf("%-10d %-50s %-6s %-6s\n",
			m.Version,
			truncate(m.Description, 50),
			upMark,
			downMark,
		)
	}

	// 統計資訊
	fmt.Println("\nSummary:")
	fmt.Printf("Total migrations: %d\n", len(migrations))

	complete := 0
	for _, m := range migrations {
		if m.HasUp && m.HasDown {
			complete++
		}
	}
	fmt.Printf("Complete (up+down): %d\n", complete)
	fmt.Printf("Incomplete: %d\n", len(migrations)-complete)

	// 檢查版本連續性
	fmt.Println("\nVersion sequence check:")
	if isSequential(migrations) {
		fmt.Println("✅ Versions are sequential")
	} else {
		fmt.Println("⚠️  Warning: Versions have gaps or are not sequential")
		printGaps(migrations)
	}
}

func listMigrations(dir string) ([]MigrationFile, error) {
	pattern := regexp.MustCompile(`^(\d+)_([a-z0-9_]+)\.(up|down)\.sql$`)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Map to store migrations by version
	migrationMap := make(map[uint]*MigrationFile)

	for _, file := range files {
		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) == 0 {
			continue
		}

		version, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			continue
		}

		description := matches[2]
		direction := matches[3]

		v := uint(version)
		if _, exists := migrationMap[v]; !exists {
			migrationMap[v] = &MigrationFile{
				Version:     v,
				Description: description,
			}
		}

		if direction == "up" {
			migrationMap[v].HasUp = true
		} else {
			migrationMap[v].HasDown = true
		}
	}

	// Convert map to sorted slice
	versions := make([]uint, 0, len(migrationMap))
	for v := range migrationMap {
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	result := make([]MigrationFile, 0, len(versions))
	for _, v := range versions {
		result = append(result, *migrationMap[v])
	}

	return result, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func isSequential(migrations []MigrationFile) bool {
	if len(migrations) == 0 {
		return true
	}

	// 檢查是否使用序列號 (000001, 000002, ...)
	if migrations[0].Version < 100 {
		for i := 0; i < len(migrations); i++ {
			expected := uint(i + 1)
			if migrations[i].Version != expected {
				return false
			}
		}
		return true
	}

	// 對於時間戳格式，只要遞增即可
	return true
}

func printGaps(migrations []MigrationFile) {
	if len(migrations) < 2 {
		return
	}

	// 只檢查序列號格式
	if migrations[0].Version >= 100 {
		return
	}

	for i := 0; i < len(migrations)-1; i++ {
		current := migrations[i].Version
		next := migrations[i+1].Version
		if next-current > 1 {
			fmt.Printf("  Gap: %d -> %d (missing %d)\n", current, next, next-current-1)
		}
	}
}

func init() {
	// 確保 migrations path 存在
	wd, _ := os.Getwd()
	fmt.Printf("Current directory: %s\n", wd)
}
