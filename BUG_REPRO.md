# Bug 是什么

缺失通知规则时评论通知被静默放行，随后启用默认规则发生 nil map 写入 panic。

# 如何触发

使用缺省规则构造服务，先检查 comment 通知，再启用 comment 规则。

# 错误信息

`missing rules silently accepted comment notification`；`assignment to entry in nil map`。
