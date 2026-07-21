// collage_cs.cs — создатель коллажей с автоматической компоновкой на C# (System.Drawing)

using System;
using System.Collections.Generic;
using System.Drawing;
using System.Drawing.Imaging;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;

class CollageMaker
{
    private List<Image> images = new List<Image>();
    private List<string> names = new List<string>();
    private int cols;
    private int padding;
    private int border;
    private bool shadow;
    private bool labels;
    private Color bgColor;
    private string output;
    private bool autoCols;

    public CollageMaker(string[] paths, int cols, int padding, int border,
                        bool shadow, bool labels, Color bgColor,
                        string output, bool autoCols)
    {
        this.cols = cols;
        this.padding = padding;
        this.border = border;
        this.shadow = shadow;
        this.labels = labels;
        this.bgColor = bgColor;
        this.output = output;
        this.autoCols = autoCols;
        foreach (var p in paths)
        {
            try
            {
                Image img = Image.FromFile(p);
                images.Add(img);
                names.Add(Path.GetFileName(p));
            }
            catch (Exception e)
            {
                Console.WriteLine($"Ошибка загрузки {p}: {e.Message}");
            }
        }
    }

    public void Make()
    {
        if (images.Count == 0)
        {
            Console.WriteLine("Нет изображений.");
            return;
        }
        int count = images.Count;
        int colsUsed = autoCols ? (int)Math.Ceiling(Math.Sqrt(count)) : cols;
        if (colsUsed == 0) colsUsed = 1;
        int rows = (int)Math.Ceiling((double)count / colsUsed);

        int cellW = images.Max(img => img.Width);
        int cellH = images.Max(img => img.Height);

        int labelHeight = labels ? 30 : 0;
        int totalW = colsUsed * (cellW + padding) + padding + 2 * border;
        int totalH = rows * (cellH + padding + labelHeight) + padding + 2 * border;

        using (Bitmap collage = new Bitmap(totalW, totalH))
        using (Graphics g = Graphics.FromImage(collage))
        {
            g.Clear(bgColor);
            int idx = 0;
            for (int row = 0; row < rows; ++row)
            {
                for (int col = 0; col < colsUsed; ++col)
                {
                    if (idx >= count) break;
                    Image img = images[idx];
                    // Масштабирование
                    double scale = Math.Min((double)cellW / img.Width, (double)cellH / img.Height);
                    int newW = (int)(img.Width * scale);
                    int newH = (int)(img.Height * scale);
                    int xOffset = (cellW - newW) / 2;
                    int yOffset = (cellH - newH) / 2;
                    int x = padding + col * (cellW + padding) + border;
                    int y = padding + row * (cellH + padding + labelHeight) + border;

                    // Рамка
                    if (border > 0)
                    {
                        g.DrawRectangle(Pens.Black, x, y, cellW, cellH);
                    }
                    // Тень
                    if (shadow)
                    {
                        g.FillRectangle(new SolidBrush(Color.FromArgb(180, 180, 180)), x + 5, y + 5, cellW, cellH);
                    }
                    // Вставка
                    g.DrawImage(img, x + xOffset, y + yOffset, newW, newH);
                    // Подпись
                    if (labels)
                    {
                        g.DrawString(names[idx], new Font("Arial", 10), Brushes.Black, x, y + cellH + 5);
                    }
                    idx++;
                }
            }
            collage.Save(output, ImageFormat.Png);
            Console.WriteLine($"Коллаж сохранён в {output}");
        }
    }

    public static void Main(string[] args)
    {
        var paths = new List<string>();
        string output = "collage.png";
        int cols = 3, padding = 10, border = 0;
        bool shadow = false, labels = false, autoCols = false;
        Color bgColor = Color.White;

        for (int i = 0; i < args.Length; ++i)
        {
            if (args[i] == "--cols" && i+1 < args.Length) cols = int.Parse(args[++i]);
            else if (args[i] == "--padding" && i+1 < args.Length) padding = int.Parse(args[++i]);
            else if (args[i] == "--border" && i+1 < args.Length) border = int.Parse(args[++i]);
            else if (args[i] == "--shadow") shadow = true;
            else if (args[i] == "--labels") labels = true;
            else if (args[i] == "--auto-cols") autoCols = true;
            else if (args[i] == "--output" && i+1 < args.Length) output = args[++i];
            else paths.Add(args[i]);
        }

        if (paths.Count == 0)
        {
            Console.WriteLine("Использование: collage_cs.exe image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]");
            return;
        }
        var maker = new CollageMaker(paths.ToArray(), cols, padding, border, shadow, labels, bgColor, output, autoCols);
        maker.Make();
    }
}
