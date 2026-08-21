# Bug 是什么

重试成功后，首轮延迟回调将详情状态回退为处理中，列表与详情不一致。

# 如何触发

阻塞首轮回调，完成第二次转码并写成功状态，再放行首轮回调。

# 错误信息

detail status=processing version=1；list status=succeeded version=2。
