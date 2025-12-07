package email

import (
	"fmt"
	"log"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// Config 阿里云邮件服务配置
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	AccountName     string // 发信地址
	FromAlias       string // 发信人昵称
	RegionID        string // 区域，如 cn-hangzhou
}

// AliyunEmailService 阿里云邮件服务
type AliyunEmailService struct {
	config Config
	client *dm.Client
}

// 全局邮件服务实例
var emailService *AliyunEmailService

// InitEmailService 初始化邮件服务
func InitEmailService(cfg Config) {
	service, err := NewAliyunEmailService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize email service: %v", err)
		emailService = &AliyunEmailService{config: cfg}
		return
	}
	emailService = service
	log.Println("Email service initialized")
}

// GetEmailService 获取邮件服务实例
func GetEmailService() *AliyunEmailService {
	return emailService
}

// NewAliyunEmailService 创建阿里云邮件服务实例
func NewAliyunEmailService(cfg Config) (*AliyunEmailService, error) {
	service := &AliyunEmailService{
		config: cfg,
	}

	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return service, nil
	}

	// 创建阿里云客户端配置
	config := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		Endpoint:        tea.String(fmt.Sprintf("dm.%s.aliyuncs.com", cfg.RegionID)),
	}

	client, err := dm.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建邮件客户端失败: %w", err)
	}

	service.client = client
	return service, nil
}

// SendVerificationCode 发送验证码邮件
func (s *AliyunEmailService) SendVerificationCode(to, code, purpose string) error {
	subject := fmt.Sprintf("【VPS Platform】%s验证码", purpose)
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: 'Microsoft YaHei', Arial, sans-serif; background-color: #f5f5f5; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; padding: 40px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; font-size: 24px; margin: 0; }
        .code-box { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 8px; padding: 30px; text-align: center; margin: 30px 0; }
        .code { font-size: 36px; font-weight: bold; color: #ffffff; letter-spacing: 8px; }
        .info { color: #666; font-size: 14px; line-height: 1.8; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #999; font-size: 12px; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>VPS Platform</h1>
        </div>
        <div class="info">
            <p>您好，</p>
            <p>您正在进行<strong>%s</strong>操作，验证码如下：</p>
        </div>
        <div class="code-box">
            <span class="code">%s</span>
        </div>
        <div class="info">
            <p>此验证码 <strong>15 分钟</strong> 内有效，请勿泄露给他人。</p>
            <p>如非本人操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
            <p>© 2025 VPS Platform. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, purpose, code)

	return s.SendEmail(to, subject, htmlBody)
}

// SendEmail 发送邮件
func (s *AliyunEmailService) SendEmail(to, subject, htmlBody string) error {
	if s.config.AccessKeyID == "" || s.config.AccessKeySecret == "" {
		log.Printf("邮件服务未配置，验证码邮件未发送: to=%s, subject=%s", to, subject)
		return nil // 未配置时不报错，只打印日志
	}

	if s.client == nil {
		return fmt.Errorf("邮件客户端未初始化")
	}

	// 构建发送请求
	request := &dm.SingleSendMailRequest{
		AccountName:    tea.String(s.config.AccountName),
		AddressType:    tea.Int32(1),
		ToAddress:      tea.String(to),
		Subject:        tea.String(subject),
		HtmlBody:       tea.String(htmlBody),
		ReplyToAddress: tea.Bool(false),
	}

	if s.config.FromAlias != "" {
		request.FromAlias = tea.String(s.config.FromAlias)
	}

	// 发送邮件
	_, err := s.client.SingleSendMail(request)
	if err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	log.Printf("邮件发送成功: to=%s, subject=%s", to, subject)
	return nil
}

// SendVerificationCode 全局函数，方便调用
func SendVerificationCode(to, code, purpose string) error {
	if emailService == nil {
		// 邮件服务未初始化时仅打印日志，不返回错误（测试环境或未配置时）
		log.Printf("邮件服务未初始化，验证码: %s (发送至 %s，用途: %s)", code, to, purpose)
		return nil
	}
	return emailService.SendVerificationCode(to, code, purpose)
}
