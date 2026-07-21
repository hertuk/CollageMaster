// collage_java.java — создатель коллажей с автоматической компоновкой на Java (AWT)

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.nio.file.*;
import java.util.*;
import java.util.List;

public class CollageMaker {
    private List<BufferedImage> images = new ArrayList<>();
    private List<String> names = new ArrayList<>();
    private int cols;
    private int padding;
    private int border;
    private boolean shadow;
    private boolean labels;
    private Color bgColor;
    private String output;
    private boolean autoCols;

    public CollageMaker(String[] paths, int cols, int padding, int border,
                        boolean shadow, boolean labels, Color bgColor,
                        String output, boolean autoCols) throws IOException {
        this.cols = cols;
        this.padding = padding;
        this.border = border;
        this.shadow = shadow;
        this.labels = labels;
        this.bgColor = bgColor;
        this.output = output;
        this.autoCols = autoCols;
        for (String p : paths) {
            try {
                BufferedImage img = ImageIO.read(new File(p));
                images.add(img);
                names.add(Paths.get(p).getFileName().toString());
            } catch (IOException e) {
                System.err.println("Ошибка загрузки " + p + ": " + e.getMessage());
            }
        }
    }

    public void make() throws IOException {
        if (images.isEmpty()) {
            System.out.println("Нет изображений.");
            return;
        }
        int count = images.size();
        int colsUsed = autoCols ? (int)Math.ceil(Math.sqrt(count)) : cols;
        if (colsUsed == 0) colsUsed = 1;
        int rows = (int)Math.ceil((double)count / colsUsed);

        int cellW = 0, cellH = 0;
        for (BufferedImage img : images) {
            cellW = Math.max(cellW, img.getWidth());
            cellH = Math.max(cellH, img.getHeight());
        }

        int labelHeight = labels ? 30 : 0;
        int totalW = colsUsed * (cellW + padding) + padding + 2*border;
        int totalH = rows * (cellH + padding + labelHeight) + padding + 2*border;

        BufferedImage collage = new BufferedImage(totalW, totalH, BufferedImage.TYPE_INT_RGB);
        Graphics2D g = collage.createGraphics();
        g.setColor(bgColor);
        g.fillRect(0, 0, totalW, totalH);

        int idx = 0;
        for (int row = 0; row < rows; ++row) {
            for (int col = 0; col < colsUsed; ++col) {
                if (idx >= count) break;
                BufferedImage img = images.get(idx);
                // Масштабирование
                double scale = Math.min((double)cellW / img.getWidth(), (double)cellH / img.getHeight());
                int newW = (int)(img.getWidth() * scale);
                int newH = (int)(img.getHeight() * scale);
                Image scaled = img.getScaledInstance(newW, newH, Image.SCALE_SMOOTH);
                int xOffset = (cellW - newW) / 2;
                int yOffset = (cellH - newH) / 2;
                int x = padding + col * (cellW + padding) + border;
                int y = padding + row * (cellH + padding + labelHeight) + border;

                // Рамка
                if (border > 0) {
                    g.setColor(Color.BLACK);
                    g.drawRect(x, y, cellW, cellH);
                }
                // Тень
                if (shadow) {
                    g.setColor(new Color(180,180,180));
                    g.fillRect(x+5, y+5, cellW, cellH);
                }
                // Вставка
                g.drawImage(scaled, x + xOffset, y + yOffset, null);
                // Подпись
                if (labels) {
                    g.setColor(Color.BLACK);
                    g.setFont(new Font("Arial", Font.PLAIN, 12));
                    g.drawString(names.get(idx), x, y + cellH + 15);
                }
                idx++;
            }
        }
        g.dispose();
        ImageIO.write(collage, "png", new File(output));
        System.out.println("Коллаж сохранён в " + output);
    }

    public static void main(String[] args) throws IOException {
        // Упрощённый парсинг аргументов
        List<String> paths = new ArrayList<>();
        String output = "collage.png";
        int cols = 3, padding = 10, border = 0;
        boolean shadow = false, labels = false, autoCols = false;
        Color bgColor = Color.WHITE;

        for (int i = 0; i < args.length; ++i) {
            if (args[i].equals("--cols") && i+1 < args.length) cols = Integer.parseInt(args[++i]);
            else if (args[i].equals("--padding") && i+1 < args.length) padding = Integer.parseInt(args[++i]);
            else if (args[i].equals("--border") && i+1 < args.length) border = Integer.parseInt(args[++i]);
            else if (args[i].equals("--shadow")) shadow = true;
            else if (args[i].equals("--labels")) labels = true;
            else if (args[i].equals("--auto-cols")) autoCols = true;
            else if (args[i].equals("--output") && i+1 < args.length) output = args[++i];
            else paths.add(args[i]);
        }

        if (paths.isEmpty()) {
            System.out.println("Использование: java CollageMaker image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]");
            return;
        }
        CollageMaker maker = new CollageMaker(paths.toArray(new String[0]), cols, padding, border, shadow, labels, bgColor, output, autoCols);
        maker.make();
    }
}
