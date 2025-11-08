package main

import (
	"encoding/json"
	"fmt"
	"log"

	pgtune "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"
)

func main() {
	// Sample input matching the example in README
	input := pgtune.TuneInput{
		DBVersion:      18,
		OSType:         pgtune.OSLinux,
		DBType:         pgtune.DBTypeWeb,
		TotalMemory:    2 * pgtune.GB,
		CPUNum:         1,
		MaxConnections: 100,
		StorageType:    pgtune.StorageTypeSSD,
	}

	output, err := pgtune.Tune(input)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("# DB Version:", input.DBVersion)
	fmt.Println("# OS Type:", input.OSType)
	fmt.Println("# DB Type:", input.DBType)
	fmt.Printf("# Total Memory (RAM): %d GB\n", input.TotalMemory/pgtune.GB)
	fmt.Println("# CPUs num:", input.CPUNum)
	fmt.Println("# Connections num:", output.MaxConnections)
	fmt.Println("# Data Storage:", input.StorageType)
	fmt.Println()

	fmt.Printf("max_connections = %d\n", output.MaxConnections)
	fmt.Printf("shared_buffers = %s\n", output.SharedBuffers)
	fmt.Printf("effective_cache_size = %s\n", output.EffectiveCacheSize)
	fmt.Printf("maintenance_work_mem = %s\n", output.MaintenanceWorkMem)
	fmt.Printf("checkpoint_completion_target = %.1f\n", output.CheckpointCompletionTarget)
	fmt.Printf("wal_buffers = %s\n", output.WalBuffers)
	fmt.Printf("default_statistics_target = %d\n", output.DefaultStatisticsTarget)
	fmt.Printf("random_page_cost = %.1f\n", output.RandomPageCost)

	if output.EffectiveIOConcurrency != nil {
		fmt.Printf("effective_io_concurrency = %d\n", *output.EffectiveIOConcurrency)
	}

	fmt.Printf("work_mem = %s\n", output.WorkMem)
	fmt.Printf("huge_pages = %s\n", output.HugePages)
	fmt.Printf("min_wal_size = %s\n", output.MinWalSize)
	fmt.Printf("max_wal_size = %s\n", output.MaxWalSize)

	fmt.Println("\n--- JSON Output ---")
	jsonOutput, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonOutput))
}
