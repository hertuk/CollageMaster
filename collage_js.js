// collage_js.js — создатель коллажей с автоматической компоновкой на JavaScript (Node.js + canvas)

const fs = require('fs');
const path = require('path');
const { createCanvas, loadImage } = require('canvas');

class CollageMaker {
    constructor(images, cols, padding, border, shadow, labels, bgColor, output, autoCols) {
        this.images = images; // массив путей к изображениям
        this.cols = cols;
        this.padding = padding;
        this.border = border;
        this.shadow = shadow;
        this.labels = labels;
        this.bgColor = bgColor;
        this.output = output;
        this.autoCols = autoCols;
        this.imageData = [];
    }

    async loadImages() {
        for (const p of this.images) {
            try {
                const img = await loadImage(p);
                this.imageData.push({ img, name: path.basename(p) });
            } catch (e) {
                console.error(`Ошибка загрузки ${p}: ${e.message}`);
            }
        }
    }

    async make() {
        await this.loadImages();
        if (this.imageData.length === 0) {
            console.log('Нет изображений.');
            return;
        }
        const count = this.imageData.length;
        const colsUsed = this.autoCols ? Math.ceil(Math.sqrt(count)) : this.cols;
        const rows = Math.ceil(count / colsUsed);

        // Определяем размер ячейки
        let cellW = 0, cellH = 0;
        for (const data of this.imageData) {
            cellW = Math.max(cellW, data.img.width);
            cellH = Math.max(cellH, data.img.height);
        }

        const labelHeight = this.labels ? 30 : 0;
        const totalW = colsUsed * (cellW + this.padding) + this.padding + 2*this.border;
        const totalH = rows * (cellH + this.padding + labelHeight) + this.padding + 2*this.border;

        const canvas = createCanvas(totalW, totalH);
        const ctx = canvas.getContext('2d');
        ctx.fillStyle = this.bgColor;
        ctx.fillRect(0, 0, totalW, totalH);

        let idx = 0;
        for (let row = 0; row < rows; row++) {
            for (let col = 0; col < colsUsed; col++) {
                if (idx >= count) break;
                const { img, name } = this.imageData[idx];
                const scale = Math.min(cellW / img.width, cellH / img.height);
                const newW = img.width * scale;
                const newH = img.height * scale;
                const xOff = (cellW - newW) / 2;
                const yOff = (cellH - newH) / 2;
                const x = this.padding + col * (cellW + this.padding) + this.border;
                const y = this.padding + row * (cellH + this.padding + labelHeight) + this.border;

                // Рамка
                if (this.border > 0) {
                    ctx.strokeStyle = '#000000';
                    ctx.lineWidth = this.border;
                    ctx.strokeRect(x, y, cellW, cellH);
                }
                // Тень
                if (this.shadow) {
                    ctx.fillStyle = '#B4B4B4';
                    ctx.fillRect(x+5, y+5, cellW, cellH);
                }
                // Вставка
                ctx.drawImage(img, x + xOff, y + yOff, newW, newH);
                // Подпись
                if (this.labels) {
                    ctx.fillStyle = '#000000';
                    ctx.font = '12px Arial';
                    ctx.fillText(name, x, y + cellH + 15);
                }
                idx++;
            }
        }
        const buffer = canvas.toBuffer('image/png');
        fs.writeFileSync(this.output, buffer);
        console.log(`Коллаж сохранён в ${this.output}`);
    }
}

// Парсинг аргументов командной строки (упрощённо)
const args = process.argv.slice(2);
if (args.length === 0) {
    console.log('Использование: node collage_js.js image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]');
    process.exit(1);
}
let paths = [];
let cols = 3, padding = 10, border = 0;
let shadow = false, labels = false, autoCols = false;
let output = 'collage.png';
let bgColor = '#ffffff';
for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '--cols') { cols = parseInt(args[++i]); }
    else if (arg === '--padding') { padding = parseInt(args[++i]); }
    else if (arg === '--border') { border = parseInt(args[++i]); }
    else if (arg === '--shadow') { shadow = true; }
    else if (arg === '--labels') { labels = true; }
    else if (arg === '--auto-cols') { autoCols = true; }
    else if (arg === '--output') { output = args[++i]; }
    else if (arg === '--bg') { bgColor = args[++i]; }
    else { paths.push(arg); }
}
if (paths.length === 0) {
    console.log('Не указаны изображения.');
    process.exit(1);
}
const maker = new CollageMaker(paths, cols, padding, border, shadow, labels, bgColor, output, autoCols);
maker.make().catch(console.error);
