package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"syscall/js"
)

func addWatermark(img image.Image, watermarkText string) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// 设置水印的字体和颜色
	col := color.RGBA{128, 128, 128, 128} // 灰色半透明水印
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
	}

	// 遍布整个图片的水印
	for y := 0; y < bounds.Dy(); y += 50 {
		for x := 0; x < bounds.Dx(); x += 200 {
			d.Dot = fixed.P(x, y)
			d.DrawString(watermarkText)
		}
	}

	return rgba
}

func processImage(this js.Value, p []js.Value) interface{} {
	if len(p) < 2 {
		js.Global().Get("console").Call("log", "Error: Insufficient arguments")
		return nil
	}

	imgData := p[0]
	watermarkText := p[1].String()

	if imgData.IsUndefined() {
		js.Global().Get("console").Call("log", "Error: Undefined arguments")
		return nil
	}

	imgBytes := make([]byte, imgData.Get("byteLength").Int())
	js.CopyBytesToGo(imgBytes, imgData)

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		js.Global().Get("console").Call("log", "Error decoding image:", err)
		return nil
	}

	resultImg := addWatermark(img, watermarkText)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resultImg, nil)
	if err != nil {
		js.Global().Get("console").Call("log", "Error encoding result image:", err)
		return nil
	}

	result := js.Global().Get("Uint8Array").New(len(buf.Bytes()))
	js.CopyBytesToJS(result, buf.Bytes())

	return result
}

func main() {
	ch := make(chan struct{}, 0)
	// 输出将打印在 console log 中
	fmt.Println("Hello Web Assembly from Go!")
	var a []string
	a = append(a, "1")
	// 将函数导出到全局JavaScript环境，以便在JavaScript中调用
	js.Global().Set("processImage", js.FuncOf(processImage))
	// 从通道 ch 中接收数据。该操作会阻塞，防止主程序退出，直到从通道中接收到数据
	<-ch
}
