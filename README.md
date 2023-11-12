# `glean`: Lean 4 镜像适配工具

使用上海交通大学的 https://mirror.sjtu.edu.cn 镜像服务，软件源镜像托管在 `elan`
和 `git/lean4-packages` 下。

请自行修改命令中的版本号，可用版本参见：
http://mirror.sjtu.edu.cn/elan/?mirror_intel_list

也可以通过这个链接下载 glean。

## 安装 Elan

Elan 是 Lean 的版本管理工具，在 Lake 调用时根据项目 `lean-toolchain` 文件下载安装 Lean 并切换到对应的版本。

```
glean -install elan -version 3.0.0
```

## 安装 Lean

以下操作会安装 Lean 与 Lean 工具链，包含语言服务器、构建工具等。

```
glean -install lean --version 4.1.0
```

## 在构建项目前下载依赖

每当下载完一个 Lean 项目后，在启动 VSCode 或命令行运行 `lake build` 前，可以提前通过镜像下载依赖。

```
glean -lake-manifest-path ~/EG/lake-manifest.json
```
