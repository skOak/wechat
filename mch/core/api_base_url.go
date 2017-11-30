package core

func APIBaseURL(sandbox bool) string {
	// TODO(chanxuehong): 后期做容灾功能
	if sandbox {
		return "https://api.mch.weixin.qq.com/sandboxnew"
	}
	return "https://api.mch.weixin.qq.com"
}
