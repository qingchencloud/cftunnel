package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/quickshare"
	"github.com/spf13/cobra"
)

// shareFlags are deliberately available on both quick and preset so the
// terminal workflow matches the desktop client's one-click sharing actions.
type shareFlags struct {
	qr       bool
	telegram bool
	copy     bool
	open     bool
	share    bool
}

func (f *shareFlags) bind(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&f.qr, "qr", false, "在终端显示公网地址二维码")
	cmd.Flags().BoolVar(&f.telegram, "telegram", false, "打开 Telegram 分享链接")
	cmd.Flags().BoolVar(&f.copy, "copy", false, "复制公网地址到剪贴板")
	cmd.Flags().BoolVar(&f.open, "open", false, "用系统浏览器打开公网地址")
	cmd.Flags().BoolVar(&f.share, "share", false, "一键显示二维码、复制地址并输出分享链接")
}

func (f shareFlags) enabled() bool {
	return f.qr || f.telegram || f.copy || f.open || f.share
}

func (f shareFlags) normalized() shareFlags {
	if f.share {
		f.qr = true
		f.copy = true
	}
	return f
}

// presentShare prints a stable, copy/paste-friendly share block. Action
// failures are warnings: they must not stop a running tunnel.
func presentShare(target string, flags shareFlags) {
	flags = flags.normalized()
	text := "cftunnel 分享地址: " + target
	telegramURL := quickshare.TelegramShareURL(target, text)

	fmt.Printf("分享地址: %s\n", target)
	fmt.Printf("Telegram 分享: %s\n", telegramURL)
	fmt.Printf("Telegram 群: %s\n", quickshare.TelegramGroupURL)
	fmt.Printf("推荐链接: %s\n", quickshare.InviteURL)

	if flags.qr {
		qr, err := quickshare.QRText(target)
		if err != nil {
			fmt.Printf("二维码生成失败: %v\n", err)
		} else {
			fmt.Println("二维码（手机扫码即可访问）:")
			fmt.Print(qr)
		}
	}
	if flags.copy {
		if err := quickshare.Copy(target); err != nil {
			fmt.Printf("剪贴板复制失败（可手动复制上面的地址）: %v\n", err)
		} else {
			fmt.Println("已复制公网地址到剪贴板")
		}
	}
	if flags.telegram {
		if err := quickshare.OpenURL(telegramURL); err != nil {
			fmt.Printf("Telegram 自动打开失败，请手动打开上面的分享链接: %v\n", err)
		} else {
			fmt.Println("已打开 Telegram 分享页面")
		}
	}
	if flags.open {
		if err := quickshare.OpenURL(target); err != nil {
			fmt.Printf("浏览器自动打开失败，请手动打开上面的地址: %v\n", err)
		} else {
			fmt.Println("已打开公网地址")
		}
	}
}

var shareCmdFlags shareFlags

func init() {
	shareFlagsCmd := &cobra.Command{
		Use:   "share <公网地址>",
		Short: "为公网地址生成二维码并分享到 Telegram",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := strings.TrimSpace(args[0])
			parsed, err := url.Parse(target)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return fmt.Errorf("公网地址格式错误，应为完整的 http(s) URL")
			}
			presentShare(target, shareCmdFlags)
			return nil
		},
	}
	shareCmdFlags.bind(shareFlagsCmd)
	rootCmd.AddCommand(shareFlagsCmd)
}
