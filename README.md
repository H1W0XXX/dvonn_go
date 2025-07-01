# DVONN(火山棋) – Go 实现版

基于 Go 语言和 Ebiten 引擎实现的 DVONN(火山棋) 桌面游戏。

## 介绍

该项目使用 Go 语言重现了经典的 DVONN 游戏，图片素材来自 [Board Game Arena 上的 DVONN](https://www.boardgamearena.com/)。

* 游戏规则与状态管理：`dvonn_go/internal/game`
* GUI 渲染与输入处理：`dvonn_go/internal/ui/ebiten`

核心逻辑参考自 [gautammohan/dvonn](https://github.com/gautammohan/dvonn)。

## 功能特性

* 支持两种游戏模式：人机对战 (PvE) ， 人人对战 (PvP)
* 可选自动填充第一阶段棋子（仅限PvP）
* 30 FPS 限制

## 环境要求

* Go 1.20+
* Ebiten v2

## 构建与运行

1. 克隆仓库并进入项目根目录：

   ```bash
   git clone https://github.com/H1W0XXX/dvonn_go.git
   cd dvonn_go
   ```

2. 构建可执行文件：

   ```bash
   go build -ldflags="-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -o dvonn.exe ./cmd/dvonn-gui/main.go
   ```

3. 运行游戏：

   ```bash
   ./dvonn.exe [flags]
   ```

## 启动参数

| 参数      | 说明                 | 默认值   |
| ------- | ------------------ | ----- |
| `-auto` | 是否自动填充第一阶段棋子       | false |
| `-mode` | 游戏模式：`pvp` 或 `pve` | pve   |

示例：在 PvE 模式下自动放置第一阶段棋子

```bash
./dvonn.exe -mode=pve -auto
```

## 游戏玩法概述

1. **摆放阶段**：棋盘空白，玩家轮流放置自己的棋子，直到所有棋子放置完毕。
2. **移动阶段**：玩家轮流跳跃移动棋子，每次移动距离等于堆叠高度，沿直线方向跳跃。
3. **结束判定**：当无法移动时，按剩余堆叠高度合计统计，最高者获胜。

## 参考资源

* [棋子图片素材来源 Board Game Arena: DVONN](https://www.boardgamearena.com/)
* [游戏规则逻辑 Gautam Mohan 的 DVONN 实现](https://github.com/gautammohan/dvonn)

---

欢迎反馈与贡献！
