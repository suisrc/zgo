package service

// MobileSender 手机号码发送器
type MobileSender struct {
	GPA // 数据库
}

// EmailSender 邮件发送器
type EmailSender struct {
	GPA // 数据库
}

// ThreeSender 第三方平台消息发送器
type ThreeSender struct {
	GPA // 数据库
}

// SendCaptcha 发送验证码
func (a *MobileSender) SendCaptcha(mobile string) (string, error) {
	return "123456", nil
}

// SendCaptcha 发送验证码
func (a *EmailSender) SendCaptcha(email string) (string, error) {
	return "123456", nil
}

// SendCaptcha 发送验证码
func (a *ThreeSender) SendCaptcha(kid, openid string) (string, error) {
	return "123456", nil
}
