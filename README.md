# `glean`: Lean 4 镜像适配工具

使用上海交通大学的 https://mirror.sjtu.edu.cn 镜像服务，软件源镜像托管在 `elan`
和 `git/lean4-packages` 下。

请自行修改命令中的版本号，可用版本参见：
http://mirror.sjtu.edu.cn/elan/?mirror_intel_list

也可以通过这个链接下载 glean。

## 安装 Elan

Elan 是 Lean 的版本管理工具，在 Lake 调用时根据项目 `lean-toolchain` 文件下载安装 Lean 并切换到对应的版本。

例如要阅读 Mathematics in Lean，可以运行

```sh
git clone --depth 1 https://mirror.sjtu.edu.cn/git/lean4-packages/mathematics_in_lean/
```

然后通过 `cat lean-toolchain` 获取需要安装的版本。

```sh
glean -install elan -version 3.1.1
```

## 安装 Lean

以下操作会安装 Lean 与 Lean 工具链，包含语言服务器、构建工具等。

```sh
glean -install lean --version 4.5.0
```

如需安装 nightly 版本，请以如下例子中的格式编辑命令。

```sh
glean -install lean --version 4.4.0-nightly-2023-11-12
```

nightly 可用版本参见 [lean4_nightly](http://mirror.sjtu.edu.cn/elan/leanprover/lean4_nightly/releases/download?mirror_intel_list)

## 在构建项目前下载依赖

每当下载完一个 Lean 项目后，在启动 VSCode 或命令行运行 `lake build` 前，可以提前通过镜像下载依赖。

```sh
glean -lake-manifest-path ~/EG/lake-manifest.json
```

此处使用选项 `-lake-manifest-path ~/EG/lake-manifest.json` 来手动指定了 `lake-manifest.json` 文件的位置。
也可以在进入一个 Lean4 Lake 项目后，直接运行

```sh
glean
```

命令。这样，glean 会自动找到当前项目中 `lake-manifest.json` 文件的位置。

### 如何判断自己是否在在一个 Lean4 Lake 项目路径工作？

运行 `ls` 命令，如果能找到 `lakefile.lean`，即代表正在 Lean4 Lake 项目路径工作。

使用 `code .` 启动 VSCode 来使用 Lean 时，也需要确保正在 Lean4 Lake 项目路径工作。
