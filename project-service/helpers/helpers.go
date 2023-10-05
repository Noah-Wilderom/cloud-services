package helpers

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"path/filepath"
)

func NavigateAndTakeScreenshot(url string) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func RemoveAllFilesInDirectory(directoryPath string) error {
	dirEntries, err := os.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range dirEntries {
		err := os.RemoveAll(filepath.Join(directoryPath, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to remove file or directory %s: %v", entry.Name(), err)
		}
		fmt.Println("Removed:", entry.Name())
	}

	return nil
}
