package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// The caller supplies one deadline shared by all attempts. Download without
// console progress so a stalled log pipe cannot block network cancellation.
func downloadUpdate(ctx context.Context, destination, url string) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		err = downloadUpdateAttempt(ctx, destination, url)
		if err == nil {
			return nil
		}
	}
	return err
}

func downloadUpdateAttempt(ctx context.Context, destination, url string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("release download returned HTTP %d", response.StatusCode)
	}
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()
	if _, err := io.Copy(output, response.Body); err != nil {
		return err
	}
	if err := output.Sync(); err != nil {
		return err
	}
	return output.Close()
}
