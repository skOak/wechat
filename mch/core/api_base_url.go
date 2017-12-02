package core

func (clt *Client) APIBaseURL() string {
	// TODO(chanxuehong): 后期做容灾功能
	if clt.mockBaseURL != "" {
		// 优先返回mock server的地址
		return clt.mockBaseURL
	}
	if clt.sandbox {
		return "https://api.mch.weixin.qq.com/sandboxnew"
	}
	return "https://api.mch.weixin.qq.com"
}

func (clt *Client) SetMockBaseURL(url string) {
	clt.mockBaseURL = url
}

func (clt *Client) ClearMockBaseURL() {
	clt.mockBaseURL = ""
}
