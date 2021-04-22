# ptm
## 介绍

`ptm` 即`(Packaging tools mirrors)`即全自动配置包管理工具的镜像网站，灵感来源于`pacman-mirrors` ptm提供镜像站延迟测试，自动挑选并设置最优镜像站 目前支持: `yum` `apt` </br>
当前版本号：`0.1.2`
## 使用教程
### 下载
可从右侧发行版下载最新版本<br>
```shell
sudo wget https://gitee.com/Stitchtor/ptm/attach_files/677499/download/ptm_v0.1.0-next_Linux_x86_64.tar.gz
```
```
sudo tar xvf ptm_v0.1.0-next_Linux_x86_64.tar.gz
```
### 运行
`sudo ./ptm`
注：`ptm`需要root权限才能修改本地镜像文件，请以root用户运行
## 使用帮助
配置项|含义|默认参数
--|:--|--
`--api`|镜像数据文件地址|仓库内[/raw/mirrors.json](https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json)
`--auto` or `--at`|全自动配置|false
`--interactive` or `--it`|启用交互式配置|true
`--mirrorCount` or `--mc `|写入的镜像站数目|3
`--mirrorSites` or `--ms`|自定义镜像地址，若启用则只会写入此一个镜像地址| 无
`--onlyShowMirror` or `--osm`|只显示镜像站点信的信息,不进行操作| false

**默认镜像数据文件地址** `https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json`
### 目前提供的开源镜像站点
站点名|地址|高校|企业|
--|:--:|--:|--:
清华|https://mirrors.tuna.tsinghua.edu.cn| ✅|
中国科学技术大学|https://mirrors.ustc.edu.cn| ✅ |
大连东软信息学院|http://mirrors.neusoft.edu.cn| ✅ |
东北大学|http://mirror.neu.edu.cn| ✅ |
浙江大学|http://mirrors.zju.edu.cn| ✅ |
华中科技大学|http://mirrors.hust.edu.cn| ✅ |
哈尔滨工业大学|http://mirrors.hust.edu.cn| ✅ |
重庆大学|http://mirrors.cqu.edu.cn| ✅ |
南京大学|https://mirrors.nju.edu.cn| ✅ |
兰州大学|http://mirror.lzu.edu.cn| ✅ |
东莞理工学院|https://mirrors.dgut.edu.cn| ✅ |
阿里云|https://mirrors.aliyun.com| |✅
网易|http://mirrors.163.com| |✅
华为|https://mirrors.huaweicloud.com| |✅
腾讯|https://mirrors.cloud.tencent.com| |✅
搜狐|http://mirrors.sohu.com| |✅

### 其他

代码水平有限，如有Bug或疑难，欢迎提交PR或issue🎉🎉🎉
