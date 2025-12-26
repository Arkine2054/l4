package workers

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/arkine/l4/3/internal/calendar"
)

func Cleaner(ctx context.Context, svc *calendar.Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cleaner stopped")
			return
		case <-ticker.C:
			border := time.Now().AddDate(0, 0, -30)
			svc.CleanupOldEvents(border)
		}
	}
}
