package miscellaneous

import (
	"context"
	"fmt"
	"time"
)

type StoreInterface interface {
	GetOrphanedURLs(ctx context.Context) ([]string, error)
	DeleteURLByID(ctx context.Context, id string) error
}

type CleanupStatus struct {
	LastRunTime        time.Time `json:"last_run_time"`
	TotalURLsRemoved   int       `json:"total_urls_removed"`
	IsRunning          bool      `json:"is_running"`
	RunIntervalMinutes int       `json:"run_interval_minutes"`
}

var cleanupStatus = CleanupStatus{
	IsRunning:          false,
	RunIntervalMinutes: 0,
	TotalURLsRemoved:   0,
}

func CleanupOrphanedURLs(ctx context.Context, store StoreInterface) (int, error) {
	orphanedURLs, err := store.GetOrphanedURLs(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get orphaned URLs: %w", err)
	}

	deletedCount := 0
	for _, urlID := range orphanedURLs {
		err = store.DeleteURLByID(ctx, urlID)
		if err != nil {
			fmt.Printf("Failed to delete orphaned URL %s: %v\n", urlID, err)
			continue
		}
		deletedCount++
	}

	return deletedCount, nil
}

func StartCleanupJob(store StoreInterface, interval time.Duration) chan bool {
	stopChan := make(chan bool)
	cleanupStatus.IsRunning = true
	cleanupStatus.RunIntervalMinutes = int(interval.Minutes())

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				fmt.Printf("[Cleanup] Running cleanup job at %s\n", time.Now().Format(time.RFC3339))
				count, err := CleanupOrphanedURLs(ctx, store)
				cleanupStatus.LastRunTime = time.Now()

				if err != nil {
					fmt.Printf("[Cleanup ERROR] Job failed: %v\n", err)
				} else {
					cleanupStatus.TotalURLsRemoved += count
					fmt.Printf("[Cleanup] Job completed - Removed %d orphaned URLs (total: %d)\n",
						count, cleanupStatus.TotalURLsRemoved)
				}
			case <-stopChan:
				fmt.Println("[Cleanup] Job stopped")
				cleanupStatus.IsRunning = false
				return
			}
		}
	}()

	return stopChan
}
