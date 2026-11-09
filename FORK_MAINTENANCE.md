# cddrop 的 Fyne fork

## 基线与仓库

- 独立仓库：`~/dev/go/fyne`；分支 `codex/cddrop-v2.8.2`。
- 官方基线：v2.8.2，`77c1fa8e0a50b6cbd04cd6c4454289e55a307a2d`。
- `upstream`：https://github.com/fyne-io/fyne.git 。目前是浅克隆，升级前按需 fetch 历史。
- Go module 保持 `fyne.io/fyne/v2`，保留全部上游许可证和正常 import 路径。
- `origin`：https://github.com/northsea4/fyne.git ，用户已在 GitHub 创建 fork，本地已关联。
- 首个 fork 版本：`v2.8.2-cddrop.1`，在官方基线之上按问题拆分补丁提交；不改写官方分支或 tag。

## 补丁边界

| 范围 | 主要文件 | 原因 |
| --- | --- | --- |
| 桌面按需事件循环 | `internal/driver/glfw/idle_desktop.go`、`loop.go`、`workqueue.go` | 无工作时 WaitEvents；有动画或可见 dirty 时按截止时间绘制；队列发布完成后唤醒 |
| 动画与 dirty 唤醒 | `internal/animation/runner.go`、`internal/driver/glfw/animation.go`、`internal/driver/common/canvas.go` | 后台任务/动画独立唤醒；最终帧与新动画入队使用同一状态锁 |
| 停机协议 | `workqueue.go`、`loop.go`、`driver.go` | 关闭入队与取消待执行 waiter 同步；停止 native wake 后才 Terminate |
| 设置监听 | `app/settings_desktop.go`、`app/settings.go` | 忽略同目录的其他应用变动；一次 load 只通知一次 |
| 旧字体清理 | `internal/painter/font.go` | ClearFontCache 同时重建 Harfbuzz shaper，释放旧字体 LRU 引用 |
| 可见字形边界 | `internal/painter/cddrop_text.go`、`cddrop/text.go` | 复用 Fyne 实际 shaped runs，应用负责布局对齐 |
| Label 私有对象释放 | `widget/label.go` | renderer 销毁时清理未进入可见树的选择辅助对象与富文本 provider |
| 隐藏控件关闭清理 | `internal/driver/common/canvas.go`、`internal/driver/glfw/window.go` | 遍历已有缓存树，清理隐藏的滚动条、阴影及控件，不创建未渲染对象 |
| 窗口文字纹理释放 | `internal/driver/glfw/window.go` | Close 在所属 GL context 仍存活时删除按 canvas 索引的文字纹理，解除关闭窗口引用 |
| macOS tracking 绘制 | `internal/driver/glfw/cddrop_present_darwin.go` | AppKit 嵌套拖放回调期间及时提交窗口；不重入事件循环 |

不修改 widget 的默认样式/布局。透明窗口、置顶、菜单、拖动与 Finder 区域适配留在 cddrop `internal/windowhost`；上传、任务与设置持久化也留在应用。WASM 保留原 ticker 调度；桌面沿用同步 GLFW callbacks，不移植 #5422 的整套 eventQueue。

`UsesNativeEventWait() bool` 是供平台 host 检测的可选能力；支持时 cddrop 不再安装旧的额外 native runloop timer/observer 与 AX 拦截。`WakeUp()` 仅唤醒，不使任意后台 UI 修改自动线程安全；应用仍需 `fyne.Do` / `fyne.DoAndWait`。

## 版本与依赖

cddrop 的正式依赖固定为：

```go
replace fyne.io/fyne/v2 => github.com/northsea4/fyne/v2 v2.8.2-cddrop.2
```

保留 `require fyne.io/fyne/v2 v2.8.2` 和现有 imports。fork 的 module 路径仍为 `fyne.io/fyne/v2`，通过 Go 的 replace 下载和使用 GitHub fork。正常构建无需本地 `../fyne` 或项目内 `third_party/fyne`。

本地同时开发两个仓库时可以在 cddrop 使用不提交的 Go workspace，包含 `.` 和 `../fyne`；默认 go.mod 始终固定到远程版本。不要把补丁复制回 cddrop 的 `third_party/fyne`，也不要修改 Go module 下载缓存。

后续发布新 tag，更新 cddrop 的 replace 版本和 go.sum，再进行应用构建与回归。不移动已发布 tag，不将 fork 自定义补丁合入官方同步用的分支；通过独立维护分支升级。

## 升级与验证

升级先在独立分支对比官方新 tag；逐项确认补丁是否已有上游等价修复，避免直接覆盖新的窗口生命周期、Wayland frame gate 或线程退出逻辑。上游讨论依据见 cddrop `docs/2026-11-03-fyne-fork-eventloop-review.md`。

必要检查：

1. painter / app race tests、animation race tests。
2. 原生驱动后台唤醒、动画停止、队列并发关闭，以及 `TestWindow_NoDragEndWithoutDraggedEvent`、关闭/resize 回归。
3. cddrop CI 测试、真实驱动构建、`--ui-smoke`、`--hover-smoke`。
4. 无输入 CPU 时间差、AX 查询；开关设置/任务、重复切换主题的正常与 GC 后内存。
5. 用户验证真实鼠标起始拖动与 Finder 多文件投放。自动 smoke 不代替这些体验验收。

当前只在 macOS ARM64 实测。浏览器驱动编译通过；Windows/Linux/Wayland 需要对应环境验证，不能据此保证近 0% CPU。

回退时，cddrop 可以临时把 replace 恢复为 `./third_party/fyne`；其 host 会自动恢复旧事件兼容路径。保留该旧快照直到新版用户验收完成。

## v2.8.2-cddrop.2

修复独立结果窗口反复打开关闭时的文字纹理缓存累积。关闭前在所属窗口的 GL context 中释放普通纹理与文字纹理，再标记关闭，并忽略关闭后的排队鼠标事件，防止重新创建已销毁 renderer。Label renderer 同时释放私有的选择辅助控件和富文本 provider。关闭遍历包括已有 renderer 中的隐藏控件，防止滚动条和阴影脱离可见树后继续留在全局缓存。真实驱动回归覆盖多窗口隔离、重复关闭只释放一次、延迟鼠标事件及标签释放后重建。cddrop 以含多文件与部分失败结果的 60 次展开/收起及 60 次结果窗开关验证 physical footprint，正常生命周期与额外诊断 GC 分开记录。
