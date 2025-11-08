package pgtune

import (
"testing"
)

func TestTuneBasicWeb(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeWeb,
TotalMemory:    2 * GB,
CPUNum:         1,
MaxConnections: 100,
StorageType:    StorageTypeSSD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

// Verify expected values match pgtune website
if output.MaxConnections != 100 {
t.Errorf("MaxConnections = %d; want 100", output.MaxConnections)
}
if output.SharedBuffers != "512MB" {
t.Errorf("SharedBuffers = %s; want 512MB", output.SharedBuffers)
}
if output.EffectiveCacheSize != "1536MB" {
t.Errorf("EffectiveCacheSize = %s; want 1536MB", output.EffectiveCacheSize)
}
if output.MaintenanceWorkMem != "128MB" {
t.Errorf("MaintenanceWorkMem = %s; want 128MB", output.MaintenanceWorkMem)
}
if output.CheckpointCompletionTarget != 0.9 {
t.Errorf("CheckpointCompletionTarget = %f; want 0.9", output.CheckpointCompletionTarget)
}
if output.WalBuffers != "16MB" {
t.Errorf("WalBuffers = %s; want 16MB", output.WalBuffers)
}
if output.DefaultStatisticsTarget != 100 {
t.Errorf("DefaultStatisticsTarget = %d; want 100", output.DefaultStatisticsTarget)
}
if output.RandomPageCost != 1.1 {
t.Errorf("RandomPageCost = %f; want 1.1", output.RandomPageCost)
}
if output.EffectiveIOConcurrency == nil || *output.EffectiveIOConcurrency != 200 {
t.Errorf("EffectiveIOConcurrency = %v; want 200", output.EffectiveIOConcurrency)
}
if output.WorkMem != "4854kB" {
t.Errorf("WorkMem = %s; want 4854kB", output.WorkMem)
}
if output.HugePages != "off" {
t.Errorf("HugePages = %s; want off", output.HugePages)
}
if output.MinWalSize != "1GB" {
t.Errorf("MinWalSize = %s; want 1GB", output.MinWalSize)
}
if output.MaxWalSize != "4GB" {
t.Errorf("MaxWalSize = %s; want 4GB", output.MaxWalSize)
}
}

func TestTuneDW(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeDW,
TotalMemory:    16 * GB,
CPUNum:         8,
MaxConnections: 0, // Use default
StorageType:    StorageTypeSSD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

// DW should have different defaults
if output.MaxConnections != 40 {
t.Errorf("MaxConnections = %d; want 40 for DW", output.MaxConnections)
}
if output.DefaultStatisticsTarget != 500 {
t.Errorf("DefaultStatisticsTarget = %d; want 500 for DW", output.DefaultStatisticsTarget)
}
}

func TestTuneOLTP(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeOLTP,
TotalMemory:    4 * GB,
CPUNum:         2,
MaxConnections: 0,
StorageType:    StorageTypeSSD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

if output.MaxConnections != 300 {
t.Errorf("MaxConnections = %d; want 300 for OLTP", output.MaxConnections)
}
}

func TestTuneHDD(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeWeb,
TotalMemory:    2 * GB,
CPUNum:         1,
MaxConnections: 100,
StorageType:    StorageTypeHDD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

if output.RandomPageCost != 4.0 {
t.Errorf("RandomPageCost = %f; want 4.0 for HDD", output.RandomPageCost)
}
if output.EffectiveIOConcurrency == nil || *output.EffectiveIOConcurrency != 2 {
t.Errorf("EffectiveIOConcurrency = %v; want 2 for HDD", output.EffectiveIOConcurrency)
}
}

func TestTuneWindows(t *testing.T) {
input := TuneInput{
DBVersion:      9.6,
OSType:         OSWindows,
DBType:         DBTypeWeb,
TotalMemory:    4 * GB,
CPUNum:         1,
MaxConnections: 100,
StorageType:    StorageTypeSSD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

// Windows with PG < 10 should limit shared_buffers to 512MB
if output.SharedBuffers != "512MB" {
t.Errorf("SharedBuffers = %s; want 512MB for Windows PG<10", output.SharedBuffers)
}
}

func TestTuneHighCPU(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeWeb,
TotalMemory:    16 * GB,
CPUNum:         8,
MaxConnections: 200,
StorageType:    StorageTypeSSD,
}

output, err := Tune(input)
if err != nil {
t.Fatalf("Tune failed: %v", err)
}

// With 8 CPUs, should set parallel worker settings
if output.MaxWorkerProcesses == nil {
t.Error("MaxWorkerProcesses should be set with 8 CPUs")
}
if output.MaxParallelWorkersPerGather == nil {
t.Error("MaxParallelWorkersPerGather should be set with 8 CPUs")
}
if output.MaxParallelWorkers == nil {
t.Error("MaxParallelWorkers should be set with PG18 and 8 CPUs")
}
if output.MaxParallelMaintenanceWorkers == nil {
t.Error("MaxParallelMaintenanceWorkers should be set with PG18 and 8 CPUs")
}
}

func TestTuneInvalidMemory(t *testing.T) {
input := TuneInput{
DBVersion:      18,
OSType:         OSLinux,
DBType:         DBTypeWeb,
TotalMemory:    0,
CPUNum:         1,
MaxConnections: 100,
StorageType:    StorageTypeSSD,
}

_, err := Tune(input)
if err == nil {
t.Error("Expected error for zero memory, got nil")
}
}

func TestFormatMemory(t *testing.T) {
tests := []struct {
kb   int64
want string
}{
{64, "64kB"},
{1024, "1MB"},
{2048, "2MB"},
{1048576, "1GB"},
{2097152, "2GB"},
{4194304, "4GB"},
}

for _, tt := range tests {
got := formatMemory(tt.kb)
if got != tt.want {
t.Errorf("formatMemory(%d) = %s; want %s", tt.kb, got, tt.want)
}
}
}
