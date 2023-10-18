# Lean 4 镜像适配工具

使用上海交通大学的 https://mirror.sjtu.edu.cn 镜像服务，软件源镜像托管在 `elan`
和 `git/lean4-packages` 下。

请自行修改命令中的版本号，可用版本参见：
http://mirror.sjtu.edu.cn/elan/?mirror_intel_list

## 安装 Elan

```
glean -install elan -version 3.0.0
```

## 安装 Lean

```
glean -install lean --version 4.1.0
```

## 在构建项目前下载依赖

```
glean -lake-manifest-path ~/EG/lake-manifest.json
```
