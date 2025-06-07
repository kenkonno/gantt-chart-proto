/**
 * RGB形式の色文字列をパースして、R, G, B の値を取得する
 * @param rgb RGB形式の色文字列 (例: "rgb(178, 223, 220)")
 * @returns [r, g, b] の配列、パースに失敗した場合は [0, 0, 0]
 */
export const parseRgb = (rgb: string): [number, number, number] => {
    // rgb(r, g, b) 形式の文字列から数値を抽出
    const match = rgb.match(/rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)/);
    if (!match) {
        return [0, 0, 0]; // パースに失敗した場合のデフォルト値
    }

    return [
        parseInt(match[1], 10),
        parseInt(match[2], 10),
        parseInt(match[3], 10)
    ];
};

/**
 * RGB値をRGB形式の文字列に変換する
 * @param r 赤 (0-255)
 * @param g 緑 (0-255)
 * @param b 青 (0-255)
 * @returns RGB形式の色文字列 (例: "rgb(178, 223, 220)")
 */
export const rgbToString = (r: number, g: number, b: number): string => {
    return `rgb(${Math.round(r)}, ${Math.round(g)}, ${Math.round(b)})`;
};

/**
 * RGB値をHSL値に変換する
 * @param r 赤 (0-255)
 * @param g 緑 (0-255)
 * @param b 青 (0-255)
 * @returns [h, s, l] の配列 (h: 0-360, s: 0-100, l: 0-100)
 */
export const rgbToHsl = (r: number, g: number, b: number): [number, number, number] => {
    // RGB値を0-1の範囲に正規化
    r /= 255;
    g /= 255;
    b /= 255;

    const max = Math.max(r, g, b);
    const min = Math.min(r, g, b);
    let h = 0, s = 0;
    const l = (max + min) / 2;

    if (max !== min) {
        const d = max - min;
        s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

        switch (max) {
            case r:
                h = (g - b) / d + (g < b ? 6 : 0);
                break;
            case g:
                h = (b - r) / d + 2;
                break;
            case b:
                h = (r - g) / d + 4;
                break;
        }

        h /= 6;
    }

    // HSL値を一般的な範囲に変換 (h: 0-360, s: 0-100, l: 0-100)
    return [
        Math.round(h * 360),
        Math.round(s * 100),
        Math.round(l * 100)
    ];
};

/**
 * HSL値をRGB値に変換する
 * @param h 色相 (0-360)
 * @param s 彩度 (0-100)
 * @param l 明度 (0-100)
 * @returns [r, g, b] の配列 (r, g, b: 0-255)
 */
export const hslToRgb = (h: number, s: number, l: number): [number, number, number] => {
    // HSL値を0-1の範囲に正規化
    h /= 360;
    s /= 100;
    l /= 100;

    let r, g, b;

    if (s === 0) {
        // 彩度が0の場合はグレースケール
        r = g = b = l;
    } else {
        const hue2rgb = (p: number, q: number, t: number) => {
            if (t < 0) t += 1;
            if (t > 1) t -= 1;
            if (t < 1 / 6) return p + (q - p) * 6 * t;
            if (t < 1 / 2) return q;
            if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
            return p;
        };

        const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
        const p = 2 * l - q;

        r = hue2rgb(p, q, h + 1 / 3);
        g = hue2rgb(p, q, h);
        b = hue2rgb(p, q, h - 1 / 3);
    }

    // RGB値を0-255の範囲に変換
    return [
        Math.round(r * 255),
        Math.round(g * 255),
        Math.round(b * 255)
    ];
};

/**
 * RGB形式の色文字列を強調表示する
 * @param rgb RGB形式の色文字列 (例: "rgb(178, 223, 220)")
 * @returns 強調表示されたRGB形式の色文字列
 */
export const emphasizeColor = (rgb: string): string => {
    // RGB文字列をパース
    const [r, g, b] = parseRgb(rgb);

    // RGBをHSLに変換
    const [h, s, l] = rgbToHsl(r, g, b);

    // 彩度を上げる（最大100%まで）
    const newS = Math.min(s * 1.3, 100);

    // 明度を下げる（最小20%まで）
    const newL = Math.max(l * 0.7, 20);

    // HSLをRGBに戻す
    const [newR, newG, newB] = hslToRgb(h, newS, newL);

    // RGB文字列を返す
    return rgbToString(newR, newG, newB);
};

/**
 * RGB形式の色文字列を明るくする
 * @param rgb RGB形式の色文字列 (例: "rgb(178, 223, 220)")
 * @param factor 明るさの係数 (1より大きい値で明るく、1より小さい値で暗く)
 * @returns 明るさが調整されたRGB形式の色文字列
 */
export const lightenColor = (rgb: string, factor: number): string => {
    // RGB文字列をパース
    const [r, g, b] = parseRgb(rgb);

    // RGBをHSLに変換
    const [h, s, l] = rgbToHsl(r, g, b);

    // 明度を調整（最大100%、最小0%）
    const newL = Math.min(Math.max(l * factor, 0), 100);

    // HSLをRGBに戻す
    const [newR, newG, newB] = hslToRgb(h, s, newL);

    // RGB文字列を返す
    return rgbToString(newR, newG, newB);
};

/**
 * RGB形式の色文字列を暗くする
 * @param rgb RGB形式の色文字列 (例: "rgb(178, 223, 220)")
 * @param factor 暗さの係数 (0-1の範囲、小さいほど暗く)
 * @returns 暗くされたRGB形式の色文字列
 */
export const darkenColor = (rgb: string, factor: number): string => {
    return lightenColor(rgb, factor);
};