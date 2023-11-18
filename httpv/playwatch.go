package httpv

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
)

// 用于监视自动化脚本的变化
// Watcher watches for changes in the automation scripts.

type PlayWatcher interface {
	// 监视自动化脚本的变化
	// Watch watches for changes in the automation scripts.
	Watch() error

	// 停止监视自动化脚本的变化
	// Stop stops watching for changes in the automation scripts.
	Stop()

	// 关闭Page
	Close()
}

//==============================================================================
// 实现Watcher接口的类型

var _ PlayWatcher = (*watcher)(nil)

var (
	SFIC = false     // 截屏图片多份
	SFII = uint32(5) // 截屏默认间隔
)

type watcher struct {
	page  playwright.Page
	fpath string        // 存放监控的位置
	findx int           // 文件索引
	delay time.Duration // 监控的间隔
	done  chan int      // 停止监控
}

func NewWatcher(page playwright.Page, fpath string, delay uint32) PlayWatcher {
	// 创建目录
	if fpath != "" {
		dir := filepath.Dir(fpath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				logrus.Errorf("screenshot mkdir error: %s -> %s", fpath, err.Error())
				return nil
			}
		}
	}
	if page == nil {
		logrus.Errorf("screenshot page is nil: %s", fpath)
		return nil
	}
	if delay == 0 {
		delay = SFII // 默认6秒, 60s, 图片一次循环
	}
	return &watcher{
		page:  page,
		fpath: fpath,
		delay: time.Duration(delay) * time.Second,
		done:  make(chan int, 1),
	}
}

func (aa *watcher) Watch() error {
	if aa.page == nil {
		return fmt.Errorf("screenshot page is nil: %s", aa.fpath)
	}
	if aa.fpath == "" {
		return fmt.Errorf("screenshot path is empty: %s", aa.fpath)
	}
	for {
		// 定期截屏，便于测试和查看
		select {
		case <-aa.done: // 停止监控
			return nil
		case <-time.After(aa.delay):
			// do nothing
		}
		if aa.page == nil || aa.page.IsClosed() {
			return nil // page已经关闭，停止监控
		}
		if SFIC {
			if err := Screenshot(aa.page, fmt.Sprintf("%s_%d.png", aa.fpath, aa.findx)); err != nil {
				logrus.Errorf(err.Error())
			} else {
				aa.findx++
				if aa.findx > 10 {
					aa.findx = 0
				}
			}
		} else {
			if err := Screenshot(aa.page, fmt.Sprintf("%s.png", aa.fpath)); err != nil {
				logrus.Errorf(err.Error())
			}
		}
	}
}

func (aa *watcher) Stop() {
	aa.done <- 1 // 停止监控
}

func (aa *watcher) Close() {
	if aa.page != nil {
		Screenshot(aa.page, aa.fpath+"_e.png")
		aa.page.Close()
	}
	aa.Stop()
}

//==============================================================================

func Screenshot(page playwright.Page, fpath string) error {
	// 截屏
	bts, err := page.Screenshot()
	if err != nil {
		return fmt.Errorf("screenshot error: %s -> %s", fpath, err.Error())
	}
	// 解码截屏
	img, _, err := image.Decode(bytes.NewReader(bts))
	if err != nil {
		return fmt.Errorf("screenshot decode error: %s -> %s", fpath, err.Error())
	}
	// 添加时间戳
	txt := time.Now().Format("2006-01-02 15:04:05") + " <- " + page.URL()
	img, err = AddTextToImage(img, txt)
	if err != nil {
		return fmt.Errorf("screenshot timestamp error: %s -> %s", fpath, err.Error())
	}
	// 编码截屏
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, img)
	if err != nil {
		return fmt.Errorf("screenshot encode error: %s -> %s", fpath, err.Error())
	}
	// 保存截屏, 持久化截屏图片
	err = os.WriteFile(fpath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("screenshot write error: %s -> %s", fpath, err.Error())
	}

	return nil
}

func AddTextToImage(img image.Image, text string) (image.Image, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	x, y := 10, 20
	fcolor := color.RGBA{255, 255, 255, 255} // white
	if IsLightColor(img.At(x, y)) {
		// 如果背景是浅色，就用黑色
		fcolor = color.RGBA{0, 0, 0, 255} // black
	}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(fcolor),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)

	return rgba, nil
}
func IsLightColor(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	luminance := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
	return luminance > 128
}
