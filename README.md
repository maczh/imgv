# Imgv 仿阿里云图像处理服务器

## 功能说明

针对任意来源的url图片，可以进行缩放、剪裁、旋转、水印、模糊、锐化、圆角、内切圆、亮度调整、对比度调整等功能处理。

- 缩放
- 剪切
- 旋转
- 水印：文字水印、图片水印
- 圆角
- 内切圆
- 模糊
- 锐化
- 亮度
- 对比度

## 安装

### 源码编译安装

```bash
git clone https://github.com/maczh/imgv
go mod tidy -compat=1.17
go build
```

### 下载Release版bin压缩包

## 运行前配置

### 创建图片缓存目录

```bash
#在imgv程序当前目录下
mkdir cache
mkdir fonts
```

### 下载中文字体到fonts/目录

```bash
FZFSK.ttf
FZKTK.ttf
FZSSK.ttf
Songti.ttc
文泉驿微米黑.ttf
文泉驿正黑.ttf
```

## 运行参数

- -p 端口号 ，默认为8080
- -d 缓存图片保存的目录