// collage_cpp.cpp — создатель коллажей с автоматической компоновкой на C++ (OpenCV)

#include <opencv2/opencv.hpp>
#include <iostream>
#include <vector>
#include <filesystem>
#include <string>
#include <cmath>

using namespace cv;
using namespace std;
namespace fs = std::filesystem;

class CollageMaker {
private:
    vector<Mat> images;
    vector<string> names;
    int cols;
    int padding;
    int border;
    bool shadow;
    bool labels;
    Scalar bg_color;
    string output;
    bool auto_cols;

public:
    CollageMaker(vector<string> paths, int cols = 3, int padding = 10,
                 int border = 0, bool shadow = false, bool labels = false,
                 Scalar bg_color = Scalar(255,255,255), string output = "collage.png",
                 bool auto_cols = false)
        : cols(cols), padding(padding), border(border), shadow(shadow),
          labels(labels), bg_color(bg_color), output(output), auto_cols(auto_cols) {
        for (auto& p : paths) {
            Mat img = imread(p, IMREAD_UNCHANGED);
            if (!img.empty()) {
                images.push_back(img);
                names.push_back(fs::path(p).filename().string());
            } else {
                cerr << "Ошибка загрузки " << p << endl;
            }
        }
    }

    void make() {
        if (images.empty()) {
            cerr << "Нет изображений." << endl;
            return;
        }
        int count = images.size();
        int rows, cols_used;
        if (auto_cols) {
            cols_used = (int)ceil(sqrt(count));
            if (cols_used == 0) cols_used = 1;
        } else {
            cols_used = cols;
        }
        rows = (int)ceil((double)count / cols_used);

        // Определяем размер ячейки (максимальные размеры)
        int cell_w = 0, cell_h = 0;
        for (auto& img : images) {
            if (img.cols > cell_w) cell_w = img.cols;
            if (img.rows > cell_h) cell_h = img.rows;
        }

        int label_height = labels ? 30 : 0;
        int total_w = cols_used * (cell_w + padding) + padding + 2*border;
        int total_h = rows * (cell_h + padding + label_height) + padding + 2*border;

        Mat collage(total_h, total_w, CV_8UC3, bg_color);
        int idx = 0;
        for (int row = 0; row < rows; ++row) {
            for (int col = 0; col < cols_used; ++col) {
                if (idx >= count) break;
                Mat img = images[idx];
                // Масштабирование с сохранением пропорций
                Mat resized;
                double scale = min((double)cell_w / img.cols, (double)cell_h / img.rows);
                resize(img, resized, Size(), scale, scale);
                // Центрирование
                int x_offset = (cell_w - resized.cols) / 2;
                int y_offset = (cell_h - resized.rows) / 2;
                int x = padding + col * (cell_w + padding) + border;
                int y = padding + row * (cell_h + padding + label_height) + border;
                // Рамка
                if (border > 0) {
                    rectangle(collage, Rect(x, y, cell_w, cell_h), Scalar(0,0,0), border);
                }
                // Тень (простая: заливаем серым со смещением)
                if (shadow) {
                    rectangle(collage, Rect(x+5, y+5, cell_w, cell_h), Scalar(180,180,180), FILLED);
                }
                // Вставка
                resized.copyTo(collage(Rect(x + x_offset, y + y_offset, resized.cols, resized.rows)));
                // Подпись
                if (labels) {
                    int label_y = y + cell_h + 5;
                    putText(collage, names[idx], Point(x, label_y+15),
                            FONT_HERSHEY_SIMPLEX, 0.5, Scalar(0,0,0), 1);
                }
                idx++;
            }
        }
        imwrite(output, collage);
        cout << "Коллаж сохранён в " << output << endl;
    }
};

int main(int argc, char** argv) {
    vector<string> paths;
    string output = "collage.png";
    int cols = 3, padding = 10, border = 0;
    bool shadow = false, labels = false, auto_cols = false;
    Scalar bg_color = Scalar(255,255,255);
    // Простой парсинг аргументов (для демонстрации)
    if (argc < 2) {
        cout << "Использование: ./collage image1.jpg image2.jpg ... [--cols N] [--padding N] [--border N] [--shadow] [--labels] [--auto-cols] [--output out.png]" << endl;
        return 0;
    }
    for (int i = 1; i < argc; ++i) {
        string arg = argv[i];
        if (arg == "--cols" && i+1 < argc) cols = stoi(argv[++i]);
        else if (arg == "--padding" && i+1 < argc) padding = stoi(argv[++i]);
        else if (arg == "--border" && i+1 < argc) border = stoi(argv[++i]);
        else if (arg == "--shadow") shadow = true;
        else if (arg == "--labels") labels = true;
        else if (arg == "--auto-cols") auto_cols = true;
        else if (arg == "--output" && i+1 < argc) output = argv[++i];
        else paths.push_back(arg);
    }
    CollageMaker maker(paths, cols, padding, border, shadow, labels, bg_color, output, auto_cols);
    maker.make();
    return 0;
}
