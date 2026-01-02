package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Session struct {
	ctx           context.Context
	allocCancel   context.CancelFunc
	browserCancel context.CancelFunc
	timeoutCancel context.CancelFunc
	manager       *Manager
	closed        bool
}

func (s *Session) Navigate(url string) error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}
	return chromedp.Run(s.ctx, chromedp.Navigate(url))
}

func (s *Session) Click(selector string) error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}
	return chromedp.Run(s.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Click(selector, chromedp.ByQuery),
	)
}

func (s *Session) Type(selector, text string) error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}
	return chromedp.Run(s.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Clear(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, text, chromedp.ByQuery),
	)
}

func (s *Session) Wait(duration int) error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}
	time.Sleep(time.Duration(duration) * time.Second)
	return nil
}

func (s *Session) Screenshot() ([]byte, error) {
	if s.closed {
		return nil, fmt.Errorf("session is closed")
	}
	var screenshot []byte
	if err := chromedp.Run(s.ctx, chromedp.CaptureScreenshot(&screenshot)); err != nil {
		return nil, err
	}
	return screenshot, nil
}

func (s *Session) ExtractText(selector string) (string, error) {
	if s.closed {
		return "", fmt.Errorf("session is closed")
	}
	var text string
	if err := chromedp.Run(s.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Text(selector, &text, chromedp.ByQuery),
	); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Session) ExtractAttribute(selector, attribute string) (string, error) {
	if s.closed {
		return "", fmt.Errorf("session is closed")
	}
	var value string
	if err := chromedp.Run(s.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.AttributeValue(selector, attribute, &value, nil, chromedp.ByQuery),
	); err != nil {
		return "", err
	}
	return value, nil
}

func (s *Session) Scroll(x, y int) error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}
	return chromedp.Run(s.ctx, chromedp.Evaluate(fmt.Sprintf("window.scrollBy(%d, %d)", x, y), nil))
}

func (s *Session) Close() error {
	if s.closed {
		return nil
	}
	
	s.closed = true
	s.timeoutCancel()
	s.browserCancel()
	s.allocCancel()
	s.manager.releaseSession()
	
	return nil
}