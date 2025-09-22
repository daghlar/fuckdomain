package screenshot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type ScreenshotConfig struct {
	Width     int
	Height    int
	Quality   int
	FullPage  bool
	Timeout   time.Duration
	UserAgent string
}

type ScreenshotResult struct {
	URL       string
	FilePath  string
	Width     int
	Height    int
	Size      int64
	Timestamp time.Time
	Success   bool
	Error     string
}

type ScreenshotCapture struct {
	config ScreenshotConfig
}

func NewScreenshotCapture(config ScreenshotConfig) *ScreenshotCapture {
	return &ScreenshotCapture{
		config: config,
	}
}

func (sc *ScreenshotCapture) Capture(url string) (*ScreenshotResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sc.config.Timeout)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.UserAgent(sc.config.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	var width, height int

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			width, height, err = sc.getViewportSize(ctx)
			return err
		}),
		chromedp.FullScreenshot(&buf, sc.config.Quality),
	)

	if err != nil {
		return &ScreenshotResult{
			URL:       url,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}, err
	}

	filePath := sc.saveScreenshot(url, buf)

	return &ScreenshotResult{
		URL:       url,
		FilePath:  filePath,
		Width:     width,
		Height:    height,
		Size:      int64(len(buf)),
		Timestamp: time.Now(),
		Success:   true,
	}, nil
}

func (sc *ScreenshotCapture) getViewportSize(ctx context.Context) (int, int, error) {
	var width, height int

	err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			Math.max(
				document.body.scrollWidth,
				document.body.offsetWidth,
				document.documentElement.clientWidth,
				document.documentElement.scrollWidth,
				document.documentElement.offsetWidth
			)
		`, &width),
		chromedp.Evaluate(`
			Math.max(
				document.body.scrollHeight,
				document.body.offsetHeight,
				document.documentElement.clientHeight,
				document.documentElement.scrollHeight,
				document.documentElement.offsetHeight
			)
		`, &height),
	)

	if err != nil {
		return sc.config.Width, sc.config.Height, err
	}

	if width == 0 {
		width = sc.config.Width
	}
	if height == 0 {
		height = sc.config.Height
	}

	return width, height, nil
}

func (sc *ScreenshotCapture) saveScreenshot(url string, data []byte) string {
	dir := "screenshots"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}

	filename := fmt.Sprintf("%s_%d.png",
		filepath.Base(url),
		time.Now().Unix())

	filePath := filepath.Join(dir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return ""
	}

	return filePath
}

func (sc *ScreenshotCapture) CaptureMultiple(urls []string) map[string]*ScreenshotResult {
	results := make(map[string]*ScreenshotResult)

	for _, url := range urls {
		result, err := sc.Capture(url)
		if err != nil {
			result = &ScreenshotResult{
				URL:       url,
				Success:   false,
				Error:     err.Error(),
				Timestamp: time.Now(),
			}
		}
		results[url] = result
	}

	return results
}

func (sc *ScreenshotCapture) CaptureWithCustomSize(url string, width, height int) (*ScreenshotResult, error) {
	originalWidth := sc.config.Width
	originalHeight := sc.config.Height

	sc.config.Width = width
	sc.config.Height = height

	result, err := sc.Capture(url)

	sc.config.Width = originalWidth
	sc.config.Height = originalHeight

	return result, err
}

func (sc *ScreenshotCapture) CaptureElement(url, selector string) (*ScreenshotResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sc.config.Timeout)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.UserAgent(sc.config.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	var width, height int

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			width, height, err = sc.getElementSize(ctx, selector)
			return err
		}),
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible),
	)

	if err != nil {
		return &ScreenshotResult{
			URL:       url,
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}, err
	}

	filePath := sc.saveScreenshot(url, buf)

	return &ScreenshotResult{
		URL:       url,
		FilePath:  filePath,
		Width:     width,
		Height:    height,
		Size:      int64(len(buf)),
		Timestamp: time.Now(),
		Success:   true,
	}, nil
}

func (sc *ScreenshotCapture) getElementSize(ctx context.Context, selector string) (int, int, error) {
	var width, height int

	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(`
			document.querySelector('%s').offsetWidth
		`, selector), &width),
		chromedp.Evaluate(fmt.Sprintf(`
			document.querySelector('%s').offsetHeight
		`, selector), &height),
	)

	return width, height, err
}
