// collage_go.go — создатель коллажей с автоматической компоновкой на Go (image/png, image/jpeg)

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type CollageMaker struct {
	images   []image.Image
	names    []string
	cols     int
	padding  int
	border   int
	shadow   bool
	labels   bool
	bgColor  color.Color
	output   string
	autoCols bool
}

func NewCollageMaker(paths []string, cols, padding, border int, shadow, labels, autoCols bool, bgColor color.Color, output string) (*CollageMaker, error) {
	m := &CollageMaker{
		cols:     cols,
		padding:  padding,
		border:   border,
		shadow:   shadow,
		labels:   labels,
		autoCols: autoCols,
		bgColor:  bgColor,
		output:   output,
	}
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Printf("Ошибка загрузки %s: %v\n", p, err)
			continue
		}
		m.images = append(m.images, img)
		m.names = append(m.names, filepath.Base(p))
	}
	return m, nil
}

func (m *CollageMaker) Make() error {
	if len(m.images) == 0 {
		return fmt.Errorf("нет изображений")
	}
	count := len(m.images)
	colsUsed := m.cols
	if m.autoCols {
		colsUsed = int(math.Ceil(math.Sqrt(float64(count))))
		if colsUsed == 0 {
			colsUsed = 1
		}
	}
	rows := int(math.Ceil(float64(count) / float64(colsUsed)))

	// Определяем размер ячейки
	var cellW, cellH int
	for _, img := range m.images {
		b := img.Bounds()
		if b.Dx() > cellW {
			cellW = b.Dx()
		}
		if b.Dy() > cellH {
			cellH = b.Dy()
		}
	}

	labelHeight := 0
	if m.labels {
		labelHeight = 30
	}
	totalW := colsUsed*(cellW+m.padding) + m.padding + 2*m.border
	totalH := rows*(cellH+m.padding+labelHeight) + m.padding + 2*m.border

	collage := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	draw.Draw(collage, collage.Bounds(), &image.Uniform{m.bgColor}, image.Point{}, draw.Src)

	idx := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < colsUsed; col++ {
			if idx >= count {
				break
			}
			img := m.images[idx]
			b := img.Bounds()
			// Масштабирование
			scale := math.Min(float64(cellW)/float64(b.Dx()), float64(cellH)/float64(b.Dy()))
			newW := int(float64(b.Dx()) * scale)
			newH := int(float64(b.Dy()) * scale)
			// Центрирование
			xOffset := (cellW - newW) / 2
			yOffset := (cellH - newH) / 2
			x := m.padding + col*(cellW+m.padding) + m.border
			y := m.padding + row*(cellH+m.padding+labelHeight) + m.border

			// Рамка
			if m.border > 0 {
				for i := 0; i < m.border; i++ {
					for j := 0; j < cellW; j++ {
						collage.Set(x+j, y+i, color.Black)
						collage.Set(x+j, y+cellH-1-i, color.Black)
					}
					for j := 0; j < cellH; j++ {
						collage.Set(x+i, y+j, color.Black)
						collage.Set(x+cellW-1-i, y+j, color.Black)
					}
				}
			}
			// Тень (простая: рисуем серый прямоугольник)
			if m.shadow {
				for i := 0; i < cellH; i++ {
					for j := 0; j < cellW; j++ {
						collage.Set(x+5+j, y+5+i, color.RGBA{180, 180, 180, 255})
					}
				}
			}
			// Вставка с масштабированием
			scaled := image.NewRGBA(image.Rect(0, 0, newW, newH))
			draw.Draw(scaled, scaled.Bounds(), img, b.Min, draw.Src)
			draw.Draw(collage, image.Rect(x+xOffset, y+yOffset, x+xOffset+newW, y+yOffset+newH), scaled, image.Point{}, draw.Src)
			// Подпись
			if m.labels {
				// Простой текст, в Go сложно без внешних библиотек, поэтому просто выведем в консоль
				// Для демонстрации пропускаем
				// В реальном проекте можно использовать библиотеку для текста
			}
			idx++
		}
	}
	// Сохранение
	out, err := os.Create(m.output)
	if err != nil {
		return err
	}
	defer out.Close()
	if filepath.Ext(m.output) == ".png" {
		err = png.Encode(out, collage)
	} else {
		err = jpeg.Encode(out, collage, nil)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Коллаж сохранён в %s\n", m.output)
	return nil
}

func main() {
	var (
		cols     = flag.Int("cols", 3, "Количество колонок")
		padding  = flag.Int("padding", 10, "Отступ")
		border   = flag.Int("border", 0, "Рамка")
		shadow   = flag.Bool("shadow", false, "Тени")
		labels   = flag.Bool("labels", false, "Подписи")
		autoCols = flag.Bool("auto-cols", false, "Авто-колонки")
		output   = flag.String("output", "collage.png", "Выходной файл")
	)
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		fmt.Println("Использование: go run collage_go.go image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]")
		return
	}
	bgColor := color.White
	maker, err := NewCollageMaker(paths, *cols, *padding, *border, *shadow, *labels, *autoCols, bgColor, *output)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if err := maker.Make(); err != nil {
		fmt.Println("Ошибка создания коллажа:", err)
	}
}
