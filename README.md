# PYVS - Python版本切换器

## 什么是PYVS？

PYVS是一个专门为Windows系统设计的Python版本管理工具。简单来说，它可以帮助你：

- 🐍 在一台电脑上安装多个不同版本的Python
- 🔄 轻松切换当前使用的Python版本
- 📦 自动下载和安装Python，无需手动操作
- 🎯 让新手也能轻松管理Python环境
- 🔍 参考项目 [JVMS](https://github.com/ystyle/jvms)
- 🔧 获取👉[Download Now](https://github.com/S1eeeeep/pyvs/releases)!


## 为什么需要PYVS？

想象一下这个场景：
- 你在学习Python，老师要求用Python 3.7
- 但你想尝试最新的Python 3.13功能
- 同时你还要维护一个用Python 3.9写的老项目

如果没有版本管理工具，你需要：
1. 手动下载不同版本的Python安装包
2. 记住每个Python的安装路径
3. 手动修改系统环境变量来切换版本
4. 容易搞混，很麻烦！

有了PYVS，你只需要简单的命令就能搞定一切！

## 系统要求

- ✅ Windows操作系统（Windows 10/11推荐）
- ✅ 管理员权限（首次设置时需要）
- ✅ 网络连接（下载Python时需要）

## 快速开始

### 第一步：获取PYVS

1. 下载`pyvs.exe`文件
2. 将`pyvs.exe`放到一个你喜欢的位置（比如`C:\pyvs`）

### 第二步：初始化（只需要做一次）

打开命令提示符（CMD）或PowerShell，**以管理员身份运行**，然后输入：

```bash
pyvs init
```

这个命令会：
- 自动设置Python的安装目录
- 配置系统环境变量
- 准备好一切后续操作

### 第三步：开始使用！

#### 1. 查看可以安装的Python版本

```bash
pyvs rls
```

你会看到类似这样的输出：
```
    1) 3.13
    2) 3.12
    3) 3.11
    4) 3.10
    5) 3.9
    6) 3.8
    7) 3.7
    8) 3.6
    9) 3.5

use "pyvs rls -a" show all the versions 
```

#### 2. 安装你需要的Python版本

比如你想安装Python 3.9：

```bash
pyvs install 3.9
```

PYVS会自动：
- 从华为云镜像下载Python安装包
- 解压并安装到指定位置
- 显示进度条，让你知道安装进度

#### 3. 查看已安装的Python版本

```bash
pyvs list
```

你会看到类似这样的输出：
```
Installed python (* marks in use):
    1) 3.9
  * 2) 3.7
```

星号（*）表示当前正在使用的Python版本。

#### 4. 切换Python版本

比如你想切换到Python 3.9：

```bash
pyvs switch 3.9
```

或者你也可以使用序号：
```bash
pyvs switch 1
```

切换成功后，你可以在命令行输入：
```bash
python --version
```
来确认当前使用的Python版本。

## 常用命令大全

| 命令 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `pyvs init` | - | 初始化配置（只需一次） | `pyvs init` |
| `pyvs list` | `pyvs ls` | 查看已安装的版本 | `pyvs list` |
| `pyvs rls` | - | 查看可安装的版本 | `pyvs rls` |
| `pyvs install` | `pyvs i` | 安装指定版本 | `pyvs install 3.9` |
| `pyvs switch` | `pyvs s` | 切换到指定版本 | `pyvs switch 3.9` |
| `pyvs remove` | `pyvs rm` | 删除指定版本 | `pyvs remove 3.7` |

## 新手常见问题

### Q: 什么是管理员权限？为什么要用管理员权限？
A: 管理员权限就是电脑的最高权限。因为PYVS需要修改系统环境变量，所以需要管理员权限。就像你要装修房子，需要房主的钥匙一样。

### Q: 如何以管理员身份运行命令提示符？
A: 
1. 点击开始菜单
2. 输入"cmd"或"PowerShell"
3. 右键点击，选择"以管理员身份运行"
4. 如果弹出安全提示，点击"是"

### Q: 安装失败怎么办？
A: 常见原因和解决方法：
1. **网络问题：** 检查网络连接，尝试重新安装
2. **权限问题：** 确保以管理员身份运行
3. **磁盘空间不足：** 检查是否有足够空间

### Q: 切换版本后，原来的Python包还在吗？
A: 每个Python版本都是独立的，切换版本后需要重新安装对应的包。建议使用虚拟环境来管理项目依赖。

### Q: 如何卸载PYVS？
A: 
1. 删除PYVS所在的文件夹
2. 从系统环境变量PATH中移除PYVS路径
3. 删除Python安装目录（默认在`C:\Program Files\python`）

## 高级使用技巧

### 使用序号切换版本
当安装了很多版本时，记住版本号可能很麻烦。你可以使用序号：

```bash
# 先查看列表
pyvs list
# 输出：
# Installed python (* marks in use):
#     1) 3.13
#     2) 3.12
#   * 3) 3.9

# 然后使用序号切换
pyvs switch 1  # 切换到3.13
```

### 查看所有可安装版本
默认`pyvs rls`只显示前9个版本，要查看所有版本：

```bash
pyvs rls -a
```

## 技术支持

如果遇到问题，你可以：
1. 查看本文档的"常见问题"部分
2. 检查网络连接和权限设置
3. 重新运行初始化命令


## 许可证

本项目采用MIT许可证，你可以自由使用、修改和分发。

---

**祝你使用愉快！🎉**