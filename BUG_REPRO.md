# Bug 是什么

明确拒绝的视频被重复发送并留下记录；临时失败后发送成功的视频仍返回失败，同时记录重复。

# 如何触发

先提交 `video-rejected`，再提交首次返回临时错误、第二次成功的 `video-temporary`。

# 错误信息

拒绝项 calls=2、records=2；临时项 status=failed、records=2。
