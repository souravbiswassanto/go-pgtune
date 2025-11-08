package pgtune

import (
	"fmt"
	"math"
)

// Constants for OS types
const (
	OSLinux   = "linux"
	OSWindows = "windows"
	OSMac     = "mac"
)

// Constants for DB types
const (
	DBTypeWeb     = "web"
	DBTypeOLTP    = "oltp"
	DBTypeDW      = "dw"
	DBTypeDesktop = "desktop"
	DBTypeMixed   = "mixed"
)

// Constants for storage types
const (
	StorageTypeSSD = "ssd"
	StorageTypeSAN = "san"
	StorageTypeHDD = "hdd"
)

// Size unit constants (in bytes)
const (
	KB = 1024
	MB = 1048576
	GB = 1073741824
	TB = 1099511627776
)

// TuneInput represents the input parameters for PostgreSQL tuning
type TuneInput struct {
	DBVersion      float64 // PostgreSQL version (e.g., 18, 17, 16, etc.)
	OSType         string  // linux, windows, mac
	DBType         string  // web, oltp, dw, desktop, mixed
	TotalMemory    int64   // Total memory in bytes
	CPUNum         int     // Number of CPUs
	MaxConnections int     // Maximum number of connections (0 for default)
	StorageType    string  // ssd, san, hdd
}

// TuneOutput represents the PostgreSQL configuration parameters
type TuneOutput struct {
	MaxConnections               int
	SharedBuffers                string
	EffectiveCacheSize           string
	MaintenanceWorkMem           string
	CheckpointCompletionTarget   float64
	WalBuffers                   string
	DefaultStatisticsTarget      int
	RandomPageCost               float64
	EffectiveIOConcurrency       *int
	WorkMem                      string
	HugePages                    string
	MinWalSize                   string
	MaxWalSize                   string
	MaxWorkerProcesses           *int
	MaxParallelWorkersPerGather  *int
	MaxParallelWorkers           *int
	MaxParallelMaintenanceWorkers *int
	WalLevel                     string
	MaxWalSenders                string
}

// Tune calculates PostgreSQL configuration parameters based on input
func Tune(input TuneInput) (*TuneOutput, error) {
	if input.TotalMemory <= 0 {
		return nil, fmt.Errorf("total memory must be greater than 0")
	}

	totalMemoryKB := input.TotalMemory / KB

	output := &TuneOutput{}

	// Calculate max_connections
	output.MaxConnections = getMaxConnections(input.MaxConnections, input.DBType)

	// Calculate shared_buffers
	sharedBuffersKB := getSharedBuffers(totalMemoryKB, input.DBType, input.OSType, input.DBVersion)
	output.SharedBuffers = formatMemory(sharedBuffersKB)

	// Calculate effective_cache_size
	effectiveCacheSizeKB := getEffectiveCacheSize(totalMemoryKB, input.DBType)
	output.EffectiveCacheSize = formatMemory(effectiveCacheSizeKB)

	// Calculate maintenance_work_mem
	maintenanceWorkMemKB := getMaintenanceWorkMem(totalMemoryKB, input.DBType, input.OSType)
	output.MaintenanceWorkMem = formatMemory(maintenanceWorkMemKB)

	// Calculate checkpoint_completion_target
	output.CheckpointCompletionTarget = 0.9

	// Calculate wal_buffers
	walBuffersKB := getWalBuffers(sharedBuffersKB)
	output.WalBuffers = formatMemory(walBuffersKB)

	// Calculate default_statistics_target
	output.DefaultStatisticsTarget = getDefaultStatisticsTarget(input.DBType)

	// Calculate random_page_cost
	output.RandomPageCost = getRandomPageCost(input.StorageType)

	// Calculate effective_io_concurrency
	if input.OSType == OSLinux {
		ioConcurrency := getEffectiveIOConcurrency(input.StorageType)
		output.EffectiveIOConcurrency = &ioConcurrency
	}

	// Calculate parallel settings
	parallelSettings := getParallelSettings(input.DBVersion, input.DBType, input.CPUNum)
	if parallelSettings != nil {
		output.MaxWorkerProcesses = parallelSettings.MaxWorkerProcesses
		output.MaxParallelWorkersPerGather = parallelSettings.MaxParallelWorkersPerGather
		output.MaxParallelWorkers = parallelSettings.MaxParallelWorkers
		output.MaxParallelMaintenanceWorkers = parallelSettings.MaxParallelMaintenanceWorkers
	}

	// Calculate work_mem
	parallelForWorkMem := getParallelForWorkMem(parallelSettings)
	workMemKB := getWorkMem(totalMemoryKB, sharedBuffersKB, output.MaxConnections, parallelForWorkMem, input.DBType)
	output.WorkMem = formatMemory(workMemKB)

	// Calculate huge_pages
	output.HugePages = getHugePages(totalMemoryKB)

	// Calculate WAL size
	minWalSizeKB, maxWalSizeKB := getCheckpointSegments(input.DBType)
	output.MinWalSize = formatMemory(minWalSizeKB)
	output.MaxWalSize = formatMemory(maxWalSizeKB)

	// Calculate WAL level for desktop
	if input.DBType == DBTypeDesktop {
		output.WalLevel = "minimal"
		output.MaxWalSenders = "0"
	} else {
		output.WalLevel = ""
		output.MaxWalSenders = ""
	}

	return output, nil
}

func getMaxConnections(connectionNum int, dbType string) int {
	if connectionNum > 0 {
		return connectionNum
	}

	defaults := map[string]int{
		DBTypeWeb:     200,
		DBTypeOLTP:    300,
		DBTypeDW:      40,
		DBTypeDesktop: 20,
		DBTypeMixed:   100,
	}

	return defaults[dbType]
}

func getSharedBuffers(totalMemoryKB int64, dbType, osType string, dbVersion float64) int64 {
	var sharedBuffers int64

	switch dbType {
	case DBTypeWeb, DBTypeOLTP, DBTypeDW, DBTypeMixed:
		sharedBuffers = totalMemoryKB / 4
	case DBTypeDesktop:
		sharedBuffers = totalMemoryKB / 16
	}

	// Limit shared_buffers to 512MB on Windows for versions < 10
	if dbVersion < 10 && osType == OSWindows {
		winMemoryLimit := int64(512 * MB / KB)
		if sharedBuffers > winMemoryLimit {
			sharedBuffers = winMemoryLimit
		}
	}

	return sharedBuffers
}

func getEffectiveCacheSize(totalMemoryKB int64, dbType string) int64 {
	switch dbType {
	case DBTypeWeb, DBTypeOLTP, DBTypeDW, DBTypeMixed:
		return (totalMemoryKB * 3) / 4
	case DBTypeDesktop:
		return totalMemoryKB / 4
	}
	return (totalMemoryKB * 3) / 4
}

func getMaintenanceWorkMem(totalMemoryKB int64, dbType, osType string) int64 {
	var maintenanceWorkMem int64

	switch dbType {
	case DBTypeWeb, DBTypeOLTP, DBTypeDesktop, DBTypeMixed:
		maintenanceWorkMem = totalMemoryKB / 16
	case DBTypeDW:
		maintenanceWorkMem = totalMemoryKB / 8
	}

	// Cap maintenance RAM at 2GB
	memoryLimit := int64(2 * GB / KB)
	if maintenanceWorkMem >= memoryLimit {
		if osType == OSWindows {
			// 2048MB (2 GB) will raise error on Windows, so we need to remove 1 MB from it
			maintenanceWorkMem = memoryLimit - int64(1*MB/KB)
		} else {
			maintenanceWorkMem = memoryLimit
		}
	}

	return maintenanceWorkMem
}

func getWalBuffers(sharedBuffersKB int64) int64 {
	// Set to 3% of shared_buffers up to a maximum of 16MB
	walBuffers := (3 * sharedBuffersKB) / 100
	maxWalBuffer := int64(16 * MB / KB)

	if walBuffers > maxWalBuffer {
		walBuffers = maxWalBuffer
	}

	// Round upwards to 16MB if near that number
	walBufferNear := int64(14 * MB / KB)
	if walBuffers > walBufferNear && walBuffers < maxWalBuffer {
		walBuffers = maxWalBuffer
	}

	// Set minimum to 32 KB
	if walBuffers < 32 {
		walBuffers = 32
	}

	return walBuffers
}

func getDefaultStatisticsTarget(dbType string) int {
	switch dbType {
	case DBTypeDW:
		return 500
	default:
		return 100
	}
}

func getRandomPageCost(storageType string) float64 {
	switch storageType {
	case StorageTypeHDD:
		return 4.0
	case StorageTypeSSD, StorageTypeSAN:
		return 1.1
	default:
		return 4.0
	}
}

func getEffectiveIOConcurrency(storageType string) int {
	switch storageType {
	case StorageTypeHDD:
		return 2
	case StorageTypeSSD:
		return 200
	case StorageTypeSAN:
		return 300
	default:
		return 2
	}
}

type ParallelSettings struct {
	MaxWorkerProcesses            *int
	MaxParallelWorkersPerGather   *int
	MaxParallelWorkers            *int
	MaxParallelMaintenanceWorkers *int
}

func getParallelSettings(dbVersion float64, dbType string, cpuNum int) *ParallelSettings {
	if cpuNum < 4 {
		return nil
	}

	workersPerGather := int(math.Ceil(float64(cpuNum) / 2.0))

	if dbType != DBTypeDW && workersPerGather > 4 {
		workersPerGather = 4
	}

	settings := &ParallelSettings{
		MaxWorkerProcesses:          &cpuNum,
		MaxParallelWorkersPerGather: &workersPerGather,
	}

	if dbVersion >= 10 {
		settings.MaxParallelWorkers = &cpuNum
	}

	if dbVersion >= 11 {
		parallelMaintenanceWorkers := int(math.Ceil(float64(cpuNum) / 2.0))
		if parallelMaintenanceWorkers > 4 {
			parallelMaintenanceWorkers = 4
		}
		settings.MaxParallelMaintenanceWorkers = &parallelMaintenanceWorkers
	}

	return settings
}

func getParallelForWorkMem(parallelSettings *ParallelSettings) int {
	if parallelSettings != nil && parallelSettings.MaxWorkerProcesses != nil {
		return *parallelSettings.MaxWorkerProcesses
	}
	return 8 // default value
}

func getWorkMem(totalMemoryKB, sharedBuffersKB int64, maxConnections, parallelForWorkMem int, dbType string) int64 {
	workMemValue := (totalMemoryKB - sharedBuffersKB) / int64((maxConnections+parallelForWorkMem)*3)

	var workMem int64
	switch dbType {
	case DBTypeWeb, DBTypeOLTP:
		workMem = workMemValue
	case DBTypeDW, DBTypeMixed:
		workMem = workMemValue / 2
	case DBTypeDesktop:
		workMem = workMemValue / 6
	default:
		workMem = workMemValue
	}

	// Minimum 64 KB
	if workMem < 64 {
		workMem = 64
	}

	return workMem
}

func getHugePages(totalMemoryKB int64) string {
	// More than 32GB - better to have huge page
	if totalMemoryKB >= 33554432 { // 32GB in KB
		return "try"
	}
	return "off"
}

func getCheckpointSegments(dbType string) (int64, int64) {
	minWalSizeMap := map[string]int64{
		DBTypeWeb:     int64(1024 * MB / KB),
		DBTypeOLTP:    int64(2048 * MB / KB),
		DBTypeDW:      int64(4096 * MB / KB),
		DBTypeDesktop: int64(100 * MB / KB),
		DBTypeMixed:   int64(1024 * MB / KB),
	}

	maxWalSizeMap := map[string]int64{
		DBTypeWeb:     int64(4096 * MB / KB),
		DBTypeOLTP:    int64(8192 * MB / KB),
		DBTypeDW:      int64(16384 * MB / KB),
		DBTypeDesktop: int64(2048 * MB / KB),
		DBTypeMixed:   int64(4096 * MB / KB),
	}

	return minWalSizeMap[dbType], maxWalSizeMap[dbType]
}

func formatMemory(kb int64) string {
	if kb >= GB/KB {
		gb := float64(kb) / float64(GB/KB)
		if gb == float64(int64(gb)) {
			return fmt.Sprintf("%dGB", int64(gb))
		}
		return fmt.Sprintf("%.0fMB", float64(kb)/float64(MB/KB))
	}
	if kb >= MB/KB {
		mb := float64(kb) / float64(MB/KB)
		if mb == float64(int64(mb)) {
			return fmt.Sprintf("%dMB", int64(mb))
		}
		return fmt.Sprintf("%.0fkB", float64(kb))
	}
	return fmt.Sprintf("%dkB", kb)
}
