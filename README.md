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

## Nginx跳转配置

在oss/文件存储/图床的nginx代理做一个配置，实现读取图片文件的url后面加上图片处理参数之后自动跳转到imgv进行处理后返回给客户端。

如： https://mnfvq6jw.oss-hk01.cdncloud.com/02.jpg?x-oss-process=image/resize,m_pad,w_250,h_450,color_FFFFFF/watermark,text_@测试Macro水印,x_20,y_10,t_50,color_FF0000,type_方正楷体,size_22/watermark,image_aHR0cHM6Ly9pbWctaG9tZS5jc2RuaW1nLmNuL2ltYWdlcy8yMDIwMTEyNDAzMjUxMS5wbmc%3D,g_ne,P_80,t_80,x_-30,y_10

实际跳转到 http://oss-hk01.cdncloud.com:8080/image/process?url=https://mnfvq6jw.oss-hk01.cdncloud.com/02.jpg?x-oss-process=image/resize,m_pad,w_250,h_450,color_FFFFFF/watermark,text_@测试Macro水印,x_20,y_10,t_50,color_FF0000,type_方正楷体,size_22/watermark,image_aHR0cHM6Ly9pbWctaG9tZS5jc2RuaW1nLmNuL2ltYWdlcy8yMDIwMTEyNDAzMjUxMS5wbmc%3D,g_ne,P_80,t_80,x_-30,y_10

```nginx
 server {
    listen 80 default_server;
    listen 443 ssl;
    server_name  *.oss.xxx-test.com; #根据实际情况修改泛域名

    ssl_certificate /etc/nginx/cert/oss.xxx-test.com.cer;
    ssl_certificate_key /etc/nginx/cert/oss.xxx-test.com.key;

    location ~ /purge/(.*) {
        allow all;
        proxy_cache_purge cache_zone $1;
   }


    location ~ /image/process {
        allow all;
        proxy_pass http://127.0.0.1:8080;
   }


    location ~* \.(jpg|gif|png|webp)$ {
        if ($request_method = GET) {
            set $test A;
        }
        if ($args ~ "x-oss-process=image") {
            set $test "${test}B";
        }
        if ($test = AB) {
            rewrite ^/ /image/process?url=$scheme://$host$uri last;
        }
        proxy_no_cache 1;
        proxy_pass http://oss-hk01-cdncloud-com;
        include /etc/nginx/default.d/s3rgw.conf;
        add_header Cache-Control no-cache;
   }
}
```



## 图片处理详细说明

### URL添加处理参数格式

```
x-oss-process=image/<操作>,<指标名称>_<值>,...,<单指标值>
```

### 图片缩放

- 操作名： resize

| 指标名 | 是否必选 | 描述                                                    | 取值范围                                                     |
| ------ | -------- | ------------------------------------------------------- | ------------------------------------------------------------ |
| m      | Y        | 指定缩放的模式。                                        | lfit（默认值）：等比缩放，缩放图限制为指定w与h的矩形内的最大图片。<br> mfit：等比缩放，缩放图为延伸出指定w与h的矩形框外的最小图片。<br> fill：将原图等比缩放为延伸出指定w与h的矩形框外的最小图片，然后将超出的部分进行居中裁剪。 <br>pad：将原图缩放为指定w与h的矩形内的最大图片，然后使用指定颜色居中填充空白部分。 fixed：固定宽高，强制缩放。 |
| w      | Y        | 指定目标缩放图的宽度                                    |                                                              |
| h      | Y        | 指定目标缩放图的高度                                    |                                                              |
| color  | N        | 当缩放模式选择为pad（缩放填充）时，可以设置填充的颜色。 |                                                              |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/resize,m_mfit,w_100,h_200
```





### 图片裁剪

- 操作名: corp

| 指标名 | 是否必选 | 描述                                                         | 取值范围                                                     |
| ------ | -------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| x      | Y        | 指定裁剪起点横坐标（默认左上角为原点）                       | [0,图片边界]                                                 |
| y      | Y        | 指定裁剪起点纵坐标（默认左上角为原点）                       | [0,图片边界]                                                 |
| g      | Y        | 设置裁剪的原点位置。原点按照九宫格的形式分布，一共有九个位置可以设置。 | nw：左上(默认)<br> north：中上<br/> ne：右上<br/> west：左中<br/> center：中部<br/> east：右中<br/> sw：左下<br/> south：中下<br/> se：右下 |
| w      | Y        | 指定裁剪宽度                                                 |                                                              |
| h      | Y        | 指定裁剪高度                                                 |                                                              |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/corp,g_nw,x_30,y_50,w_100,h_200
```



### 图片水印

- 操作名: watermark

| 指标名 | 是否必选 | 描述                                                         | 取值范围                                                     |
| ------ | -------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| t      | N        | 指定图片水印或水印文字的透明度                               | 指定图片水印或水印文字的透明度。                             |
| x      | Y        | 指定水印的水平边距， 即距离图片边缘的水平距离                | 单位：像素（px）                                             |
| y      | Y        | 指定水印的垂直边距，即距离图片边缘的垂直距离                 | 单位：像素（px）                                             |
| g      | Y        | 设置裁剪的原点位置。原点按照九宫格的形式分布，一共有九个位置可以设置。 | nw：左上(左上对齐)<br> north：中上(上对齐，水平居中)<br/> ne：右上(右上对齐)<br/> west：左中(左对齐，垂直居中)<br/> center：中部(全图居中)<br/> east：右中(右对齐，垂直居中)<br/> sw：左下(左下对齐)<br/> south：中下(下对齐，水平居中)<br/> se：右下(右下对齐，默认) |
| image  | N        | 图片水印的水印图片url,绝对或相对url，先Base64编码，后UrlEncode编码 | 水印图片url可带参数，本身也可以使用x-oss-process参数进行处理 |
| P      | N        | 指定图片水印按照原图的比例进行缩放，取值为缩放的百分比。如设置参数值为10，如果原图为100×100， 则图片水印大小为10×10。当原图变成了200×200，则图片水印大小为20×20 | [1,100]                                                      |
| text   | N        | 文本水印内容，用UrlEncode编码                                |                                                              |
| type   | N        | 文本字体名称，用UrlEncode编码                                | 方正仿宋<br>华文宋体<br>方正书宋<br>方正楷体<br>文泉驿正黑<br>文泉驿微米黑 |
| color  | N        | 指定文字水印的文字颜色，参数值为RGB颜色值                    | RGB颜色值，例如：000000表示黑色，FFFFFF表示白色。默认值：000000（黑色） |
| size   | N        | 指定文字水印的文字大小                                       | 默认值：40<br>单位: pt                                       |
| rotate | N        | 指定文字顺时针旋转角度                                       | [-360,360]                                                   |


- 使用范例

```
#文字水印
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/watermark,text_%E6%B5%8B%E8%AF%95%20test%20%E6%B0%B4%E5%8D%B0,type_%E6%96%B9%E6%AD%A3%E4%BB%BF%E5%AE%8B,t_70,x_20,y_30,color_00FFFF

#图片水印
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/watermark,image_aHR0cHM6Ly9pbWctaG9tZS5jc2RuaW1nLmNuL2ltYWdlcy8yMDIyMTAyNjA1MTAwOC5wbmc%3D,t_80,g_center
```



### 图片旋转

- 操作名: rotate

| 指标名 | 是否必选 | 描述           | 取值范围   |
| ------ | -------- | -------------- | ---------- |
| rotate | Y        | 顺时针旋转角度 | [-360,360] |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/rotate,45
```



### 亮度调整

- 操作名: bright

| 指标名  | 是否必选 | 描述           | 取值范围   |
| ------- | -------- | -------------- | ---------- |
| [value] | Y        | 指定图片的亮度 | [-100,100] |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/bright,-30
```



### 对比度调整

- 操作名: contrast

| 指标名  | 是否必选 | 描述             | 取值范围   |
| ------- | -------- | ---------------- | ---------- |
| [value] | Y        | 指定图片的对比度 | [-100,100] |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/contrast,20
```



### 图片锐化

- 操作名: sharpen

| 指标名  | 是否必选 | 描述               | 取值范围 |
| ------- | -------- | ------------------ | -------- |
| [value] | Y        | 设置锐化效果的强度 | [50,399] |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/sharpen,150
```



### 图片模糊

- 操作名: blur

| 指标名 | 是否必选 | 描述                 | 取值范围                        |
| ------ | -------- | -------------------- | ------------------------------- |
| r      | Y        | 设置模糊半径         | [1,50]<br>该值越大，图片越模糊  |
| s      | Y        | 设置正态分布的标准差 | [1,50]<br/>该值越大，图片越模糊 |

- 使用范例

```
https://u3xs13xf.hk03.1112oss.com/image%2FWX20221009-200815@2x.png?x-oss-process=image/blur,r_10,s_20

```

