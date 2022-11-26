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

- | 指标名 | 是否必选 | 描述                                                    | 取值范围                                                     |
  | ------ | -------- | ------------------------------------------------------- | ------------------------------------------------------------ |
  | m      | Y        | 指定缩放的模式。                                        | lfit（默认值）：等比缩放，缩放图限制为指定w与h的矩形内的最大图片。<br> mfit：等比缩放，缩放图为延伸出指定w与h的矩形框外的最小图片。<br> fill：将原图等比缩放为延伸出指定w与h的矩形框外的最小图片，然后将超出的部分进行居中裁剪。 <br>pad：将原图缩放为指定w与h的矩形内的最大图片，然后使用指定颜色居中填充空白部分。 fixed：固定宽高，强制缩放。 |
  | w      | Y        | 指定目标缩放图的宽度                                    |                                                              |
  | h      | Y        | 指定目标缩放图的高度                                    |                                                              |
  | color  | N        | 当缩放模式选择为pad（缩放填充）时，可以设置填充的颜色。 |                                                              |

