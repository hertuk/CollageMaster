// collage_rs.rs — создатель коллажей с автоматической компоновкой на Rust (image crate)

use std::fs::File;
use std::io::BufReader;
use std::path::Path;
use image::{ImageReader, GenericImageView, ImageBuffer, Rgba, Pixel};
use std::env;
use std::f64;

struct CollageMaker {
    images: Vec<(image::DynamicImage, String)>,
    cols: usize,
    padding: u32,
    border: u32,
    shadow: bool,
    labels: bool,
    bg_color: Rgba<u8>,
    output: String,
    auto_cols: bool,
}

impl CollageMaker {
    fn new(paths: Vec<String>, cols: usize, padding: u32, border: u32,
           shadow: bool, labels: bool, bg_color: Rgba<u8>, output: String, auto_cols: bool) -> Self {
        let mut images = Vec::new();
        for p in paths {
            if let Ok(img) = ImageReader::open(&p).and_then(|r| r.decode()) {
                let name = std::path::Path::new(&p).file_name().unwrap().to_str().unwrap().to_string();
                images.push((img, name));
            } else {
                eprintln!("Ошибка загрузки {}", p);
            }
        }
        CollageMaker { images, cols, padding, border, shadow, labels, bg_color, output, auto_cols }
    }

    fn make(&self) -> Result<(), Box<dyn std::error::Error>> {
        if self.images.is_empty() {
            return Err("Нет изображений".into());
        }
        let count = self.images.len();
        let cols_used = if self.auto_cols {
            let c = (count as f64).sqrt().ceil() as usize;
            if c == 0 { 1 } else { c }
        } else {
            self.cols
        };
        let rows = (count + cols_used - 1) / cols_used;

        // Определяем размер ячейки
        let mut cell_w = 0;
        let mut cell_h = 0;
        for (img, _) in &self.images {
            let (w, h) = img.dimensions();
            if w > cell_w { cell_w = w; }
            if h > cell_h { cell_h = h; }
        }

        let label_height = if self.labels { 30 } else { 0 };
        let total_w = cols_used * (cell_w + self.padding) + self.padding + 2 * self.border;
        let total_h = rows * (cell_h + self.padding + label_height) + self.padding + 2 * self.border;

        let mut collage = ImageBuffer::<Rgba<u8>, Vec<u8>>::new(total_w, total_h);
        for pixel in collage.pixels_mut() {
            *pixel = self.bg_color;
        }

        let mut idx = 0;
        for row in 0..rows {
            for col in 0..cols_used {
                if idx >= count { break; }
                let (img, ref name) = self.images[idx];
                let (img_w, img_h) = img.dimensions();
                let scale = f64::min(cell_w as f64 / img_w as f64, cell_h as f64 / img_h as f64);
                let new_w = (img_w as f64 * scale) as u32;
                let new_h = (img_h as f64 * scale) as u32;
                let x_off = (cell_w - new_w) / 2;
                let y_off = (cell_h - new_h) / 2;
                let x = self.padding + col * (cell_w + self.padding) + self.border;
                let y = self.padding + row * (cell_h + self.padding + label_height) + self.border;

                // Рамка (простая)
                if self.border > 0 {
                    for i in 0..self.border {
                        for j in 0..cell_w {
                            collage.put_pixel(x+j, y+i, Rgba([0,0,0,255]));
                            collage.put_pixel(x+j, y+cell_h-1-i, Rgba([0,0,0,255]));
                        }
                        for j in 0..cell_h {
                            collage.put_pixel(x+i, y+j, Rgba([0,0,0,255]));
                            collage.put_pixel(x+cell_w-1-i, y+j, Rgba([0,0,0,255]));
                        }
                    }
                }
                // Тень
                if self.shadow {
                    for dy in 0..cell_h {
                        for dx in 0..cell_w {
                            collage.put_pixel(x+5+dx, y+5+dy, Rgba([180,180,180,255]));
                        }
                    }
                }
                // Масштабируем изображение
                let scaled = img.resize_exact(new_w, new_h, image::imageops::FilterType::Lanczos3);
                // Вставляем
                for dy in 0..new_h {
                    for dx in 0..new_w {
                        let px = scaled.get_pixel(dx, dy);
                        collage.put_pixel(x + x_off + dx, y + y_off + dy, px);
                    }
                }
                // Подпись (пропускаем, т.к. сложно без шрифтов)
                idx += 1;
            }
        }
        collage.save(&self.output)?;
        println!("Коллаж сохранён в {}", self.output);
        Ok(())
    }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Использование: collage_rs image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]");
        return Ok(());
    }
    let mut paths = Vec::new();
    let mut cols = 3;
    let mut padding = 10;
    let mut border = 0;
    let mut shadow = false;
    let mut labels = false;
    let mut auto_cols = false;
    let mut output = "collage.png".to_string();

    let mut i = 1;
    while i < args.len() {
        match args[i].as_str() {
            "--cols" => { i+=1; cols = args[i].parse()?; }
            "--padding" => { i+=1; padding = args[i].parse()?; }
            "--border" => { i+=1; border = args[i].parse()?; }
            "--shadow" => shadow = true,
            "--labels" => labels = true,
            "--auto-cols" => auto_cols = true,
            "--output" => { i+=1; output = args[i].clone(); }
            _ => paths.push(args[i].clone()),
        }
        i+=1;
    }
    let bg_color = Rgba([255,255,255,255]);
    let maker = CollageMaker::new(paths, cols, padding, border, shadow, labels, bg_color, output, auto_cols);
    maker.make()?;
    Ok(())
}
