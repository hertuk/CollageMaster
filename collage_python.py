# collage_python.py — создатель коллажей с автоматической компоновкой на Python (Pillow)

import os
import sys
import argparse
from PIL import Image, ImageDraw, ImageFont
import math
from typing import List, Tuple

class CollageMaker:
    def __init__(self, images: List[str], cols: int = 3, padding: int = 10,
                 border: int = 0, shadow: bool = False, labels: bool = False,
                 bg_color: Tuple[int,int,int] = (255,255,255),
                 output: str = "collage.png", auto_cols: bool = False):
        self.images = images
        self.cols = cols
        self.padding = padding
        self.border = border
        self.shadow = shadow
        self.labels = labels
        self.bg_color = bg_color
        self.output = output
        self.auto_cols = auto_cols
        self.image_objects = []

    def load_images(self):
        for path in self.images:
            try:
                img = Image.open(path).convert("RGBA")
                self.image_objects.append((img, os.path.basename(path)))
            except Exception as e:
                print(f"Ошибка загрузки {path}: {e}")

    def compute_grid(self, count):
        if self.auto_cols:
            # определяем оптимальное количество колонок: квадратная сетка
            self.cols = int(math.ceil(math.sqrt(count)))
            if self.cols == 0: self.cols = 1
        cols = self.cols
        rows = math.ceil(count / cols)
        return rows, cols

    def make(self):
        self.load_images()
        if not self.image_objects:
            print("Нет загруженных изображений.")
            return
        count = len(self.image_objects)
        rows, cols = self.compute_grid(count)

        # Определяем размеры ячейки
        max_width = max(img.size[0] for img, _ in self.image_objects)
        max_height = max(img.size[1] for img, _ in self.image_objects)
        # или можно использовать среднее, но для простоты используем макс
        cell_w = max_width
        cell_h = max_height

        # Итоговый размер коллажа
        total_w = cols * (cell_w + self.padding) + self.padding + 2*self.border
        total_h = rows * (cell_h + self.padding) + self.padding + 2*self.border
        # Дополнительное место для подписей
        label_height = 30 if self.labels else 0
        total_h += rows * label_height

        collage = Image.new('RGB', (total_w, total_h), self.bg_color)
        draw = ImageDraw.Draw(collage)

        # Шрифт для подписей (используем стандартный)
        try:
            font = ImageFont.truetype("arial.ttf", 14)
        except:
            font = ImageFont.load_default()

        idx = 0
        for row in range(rows):
            for col in range(cols):
                if idx >= count:
                    break
                img, name = self.image_objects[idx]
                # Масштабируем изображение под размер ячейки с сохранением пропорций
                img.thumbnail((cell_w, cell_h), Image.Resampling.LANCZOS)
                # Вычисляем позицию для центрирования внутри ячейки
                x_offset = (cell_w - img.width) // 2
                y_offset = (cell_h - img.height) // 2
                # Координаты на холсте
                x = self.padding + col * (cell_w + self.padding) + self.border
                y = self.padding + row * (cell_h + self.padding) + self.border
                # Добавляем смещение для подписей
                if self.labels:
                    y += row * label_height
                # Рамка
                if self.border > 0:
                    draw.rectangle([x, y, x+cell_w, y+cell_h], outline=(0,0,0), width=self.border)
                # Тень (упрощённо: рисуем серый прямоугольник со смещением)
                if self.shadow:
                    shadow_offset = 5
                    shadow_color = (180,180,180)
                    draw.rectangle([x+shadow_offset, y+shadow_offset, x+cell_w+shadow_offset, y+cell_h+shadow_offset],
                                   fill=shadow_color)
                # Вставляем изображение с учётом центрирования
                paste_x = x + x_offset
                paste_y = y + y_offset
                collage.paste(img, (paste_x, paste_y), img)
                # Подпись
                if self.labels:
                    label_y = y + cell_h + 5
                    draw.text((x, label_y), name, fill=(0,0,0), font=font)
                idx += 1
        collage.save(self.output)
        print(f"Коллаж сохранён в {self.output}")

def main():
    parser = argparse.ArgumentParser(description="Создатель коллажей с автоматической компоновкой")
    parser.add_argument("--input", nargs="+", required=True, help="Пути к изображениям")
    parser.add_argument("--output", default="collage.png", help="Выходной файл")
    parser.add_argument("--cols", type=int, default=3, help="Количество колонок")
    parser.add_argument("--padding", type=int, default=10, help="Отступ между изображениями")
    parser.add_argument("--border", type=int, default=0, help="Толщина рамки вокруг каждого изображения")
    parser.add_argument("--shadow", action="store_true", help="Добавить тени")
    parser.add_argument("--labels", action="store_true", help="Добавить подписи с именами файлов")
    parser.add_argument("--bg", default="white", help="Цвет фона (white, black или hex)")
    parser.add_argument("--auto-cols", action="store_true", help="Автоматически определить количество колонок")
    args = parser.parse_args()

    bg_color = (255,255,255)
    if args.bg.lower() == "black":
        bg_color = (0,0,0)
    elif args.bg.startswith("#"):
        hex_color = args.bg.lstrip('#')
        bg_color = tuple(int(hex_color[i:i+2], 16) for i in (0,2,4))

    maker = CollageMaker(args.input, cols=args.cols, padding=args.padding,
                         border=args.border, shadow=args.shadow, labels=args.labels,
                         bg_color=bg_color, output=args.output, auto_cols=args.auto_cols)
    maker.make()

if __name__ == "__main__":
    main()
