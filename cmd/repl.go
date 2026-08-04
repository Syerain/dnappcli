package cmd

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const prompt = "dnappcli> "

// terminal 保存终端原始配置，用于退出时恢复。
type terminal struct {
	fd       int
	oldState *term.State
}

// newTerminal 将 stdin 设为原始模式（逐字节读取、关闭 echo）。
func newTerminal() (*terminal, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil, fmt.Errorf("stdin 不是交互式终端，无法进入 REPL")
	}
	old, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return &terminal{fd: fd, oldState: old}, nil
}

func (t *terminal) restore() {
	if t.oldState != nil {
		term.Restore(t.fd, t.oldState)
	}
}


// lineEditor 维护当前输入行的文本与光标位置。
type lineEditor struct {
	buf    []rune
	cursor int // [0, len(buf)]
}

func (e *lineEditor) reset() {
	e.buf = e.buf[:0]
	e.cursor = 0
}

func (e *lineEditor) getString() string {
	return string(e.buf)
}

func (e *lineEditor) insert(r rune) {
	e.buf = append(e.buf[:e.cursor], append([]rune{r}, e.buf[e.cursor:]...)...)
	e.cursor++
}

func (e *lineEditor) moveLeft() {
	if e.cursor > 0 {
		e.cursor--
	}
}

func (e *lineEditor) moveRight() {
	if e.cursor < len(e.buf) {
		e.cursor++
	}
}

func (e *lineEditor) moveStart() { e.cursor = 0 }

func (e *lineEditor) moveEnd() { e.cursor = len(e.buf) }

// backspace 删除光标前一字符。
func (e *lineEditor) backspace() bool {
	if e.cursor == 0 {
		return false
	}
	e.buf = append(e.buf[:e.cursor-1], e.buf[e.cursor:]...)
	e.cursor--
	return true
}

// deleteChar 删除光标处字符。
func (e *lineEditor) deleteChar() bool {
	if e.cursor >= len(e.buf) {
		return false
	}
	e.buf = append(e.buf[:e.cursor], e.buf[e.cursor+1:]...)
	return true
}

func (e *lineEditor) clearLine() {
	e.buf = e.buf[:0]
	e.cursor = 0
}

// displayWidth 估算 rune 的终端显示宽度（列）。
// 常见的 CJK 全宽字符计为 2 列，其余计为 1 列。
func runeWidth(r rune) int {
	// 组合零宽字符
	if r == 0x200d || r == 0xfe0f || (r >= 0x200b && r <= 0x200f) ||
		(r >= 0x0300 && r <= 0x036f) { // 组合用音标等
		return 0
	}

	// CJK 统一表意文字、扩展A、兼容表意、假名、谚文音节等
	switch {
	case r >= 0x1100 && r <= 0x11ff: // 谚文 Jamo
		return 2
	case r >= 0x2E80 && r <= 0xA4CF: // CJK 部首、假名、谚文
		return 2
	case r >= 0xAC00 && r <= 0xD7A3: // 谚文音节
		return 2
	case r >= 0xF900 && r <= 0xFAFF: // CJK 兼容表意
		return 2
	case r >= 0xFE10 && r <= 0xFE19: // 竖排形式（全宽）
		return 2
	case r >= 0xFE30 && r <= 0xFE6F: // CJK 兼容形式
		return 2
	case r >= 0xFF00 && r <= 0xFF60: // 全宽形式
		return 2
	case r >= 0xFFE0 && r <= 0xFFE6: // 全宽符号
		return 2
	case r >= 0x1F300 && r <= 0x1FAFF: // emoji 区（大部分宽度需代理对，粗略按 2）
		return 2
	case r >= 0x20000 && r <= 0x3FFFD: // CJK 扩展 B+
		return 2
	}
	return 1
}

// displayRuneWidth 计算 rune 序列的总显示宽度。
func displayRuneWidth(rs []rune) int {
	w := 0
	for _, r := range rs {
		w += runeWidth(r)
	}
	return w
}

// tabComplete 仅当输入是命令前缀时补全为某条命令名。
func (e *lineEditor) tabComplete() bool {
	s := e.getString()
	if strings.Contains(s, " ") {
		return false
	}
	names := []string{"login", "register", "registercode", "help", "exit"}
	for _, n := range names {
		if len(s) > 0 && strings.HasPrefix(n, s) {
			e.buf = []rune(n)
			e.cursor = len(e.buf)
			return true
		}
	}
	return false
}

// replHost 管理 REPL 的编辑与历史。
type replHost struct {
	history   []string
	histIndex int // 上下箭头浏览历史时的索引
}

// StartREPL 启动交互式命令行循环，支持方向键移动与历史命令。
func StartREPL() {
	appInit()
	fmt.Print("dnappcli 交互模式  |  输入 help 查看命令  |  输入 exit 退出\r\n\r\n")

	repl := &replHost{}
	if err := repl.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// run 进入循环。错误返回时退出。
func (r *replHost) run() error {
	term, err := newTerminal()
	if err != nil {
		return err
	}
	defer term.restore()

	ed := &lineEditor{}

	for {
		render(ed)

		key, err := readKey()
		if err != nil {
			return err
		}

		switch {
		case key.isEnter():
			ioLine := ed.getString()
			fmt.Print("\r\n")
			if ioLine == "" {
				continue
			}
			if isExitCommand(ioLine) {
				fmt.Print("bye.\r\n")
				return nil
			}
			if isHelpCommand(ioLine) {
				printREPLHelp()
				ed.reset()
				continue
			}

			r.history = append(r.history, ioLine)
			r.histIndex = len(r.history)

			args := parseArgs(ioLine)
			if len(args) > 0 {
				// 执行命令时临时恢复终端为标准模式，使输出正确
				term.restore()
				ExecuteArgs(args)
				// 重新进入原始模式
				newRaw, err := newTerminal()
				if err != nil {
					return err
				}
				*term = *newRaw
			}
			ed.reset()
			fmt.Print("\r\n")

		case key.isUp():
			r.histPrev(ed)
		case key.isDown():
			r.histNext(ed)
		case key.isLeft():
			ed.moveLeft()
			render(ed)
		case key.isRight():
			ed.moveRight()
			render(ed)
		case key.isHome():
			ed.moveStart()
			render(ed)
		case key.isEnd():
			ed.moveEnd()
			render(ed)
		case key.isDelete():
			if ed.backspace() {
				render(ed)
			}
		case key.isForwardDelete():
			if ed.deleteChar() {
				render(ed)
			}
		case key.isCtrlU():
			ed.clearLine()
			render(ed)
		case key.isTab():
			if ed.tabComplete() {
				render(ed)
			}
		case key.isPrintable():
			ed.insert(key.r)
			render(ed)
		case key.isEOF():
			fmt.Print("\r\n")
			return nil
		default:
		}
	}
}

// render 清行、画提示符与缓冲区、将光标移到正确位置。
// 使用 ANSI 绝对列定位（\x1b[%dG），避免 \x1b[%dC 的参数边界问题。
func render(ed *lineEditor) {
	// 1. \r 回到行首 + \x1b[K 清除到行尾 = 清空整行
	fmt.Print("\r\x1b[K")
	// 2. 打印提示符
	fmt.Print(prompt)
	// 3. 打印输入缓冲区
	fmt.Print(ed.getString())

	// 4. 计算光标应处的绝对列号（1-based）
	//    列号 = 1（第1列） + 提示符显示宽度 + 光标前文本显示宽度
	col := 1 + displayRuneWidth([]rune(prompt)) + displayRuneWidth(ed.buf[:ed.cursor])
	fmt.Printf("\r\x1b[%dG", col)
}

// clearLine 清空当前行（含提示符）。
func clearLine() {
	fmt.Print("\r\x1b[K")
}

// printCursorTo 将光标移到当前输入位置（遗留兼容 - 不再使用）。

// isExitCommand 判断是否为退出命令。
func isExitCommand(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "exit" || s == "quit"
}

// isHelpCommand 判断是否为帮助命令。
func isHelpCommand(s string) bool {
	return strings.TrimSpace(strings.ToLower(s)) == "help"
}

func (r *replHost) histPrev(ed *lineEditor) {
	if len(r.history) == 0 {
		return
	}
	if r.histIndex > 0 {
		r.histIndex--
	}
	if r.histIndex < len(r.history) {
		ed.buf = []rune(r.history[r.histIndex])
		ed.cursor = len(ed.buf)
	}
	render(ed)
}

func (r *replHost) histNext(ed *lineEditor) {
	if r.histIndex >= len(r.history) {
		return
	}
	r.histIndex++
	if r.histIndex >= len(r.history) {
		ed.reset()
	} else {
		ed.buf = []rune(r.history[r.histIndex])
		ed.cursor = len(ed.buf)
	}
	render(ed)
}

// key 表示从终端读到的单个输入。
type key struct {
	r           rune
	escSequence string
}

func (k key) isEnter()        bool { return k.r == '\r' || k.r == '\n' }
func (k key) isPrintable()    bool { return k.r >= 0x20 && k.r != 0x7f }
func (k key) isEOF()          bool { return k.r == 0x04 } // Ctrl-D
func (k key) isTab()          bool { return k.r == 0x09 }
func (k key) isCtrlU()        bool { return k.r == 0x15 }
func (k key) isUp()           bool { return k.escSequence == "[A" }
func (k key) isDown()         bool { return k.escSequence == "[B" }
func (k key) isRight()        bool { return k.escSequence == "[C" }
func (k key) isLeft()         bool { return k.escSequence == "[D" }
func (k key) isHome()         bool { return k.escSequence == "[H" || k.escSequence == "[1~" }
func (k key) isEnd()          bool { return k.escSequence == "[F" || k.escSequence == "[4~" }
func (k key) isDelete()       bool { return k.r == 0x7f }
func (k key) isForwardDelete() bool { return k.escSequence == "[3~" }

// readKey 从 stdin 读取一个键（含 ANSI 转义序列）。
func readKey() (key, error) {
	var one [1]byte
	if _, err := os.Stdin.Read(one[:]); err != nil {
		return key{}, err
	}
	b := one[0]

	// ESC 开头可能是方向键等转义序列
	if b == 0x1b {
		seq := make([]byte, 0, 8)
		for i := 0; i < 8; i++ {
			var c [1]byte
			n, err := os.Stdin.Read(c[:])
			if err != nil || n == 0 {
				break
			}
			seq = append(seq, c[0])
			if (c[0] >= 'A' && c[0] <= 'D') || c[0] == 'H' || c[0] == 'F' || c[0] == '~' {
				break
			}
		}
		return key{escSequence: string(seq)}, nil
	}

	return key{r: rune(b)}, nil
}

// parseArgs 按空格切分，支持引号包裹的含空格的参数。
func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for _, ch := range line {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func printREPLHelp() {
	fmt.Print("\r\n")
	fmt.Print("可用命令：\r\n")
	fmt.Print("  login <username> <password>           用户登录\r\n")
	fmt.Print("  register <username> <password> <nickname> [registercode]  注册新用户（注册码可选，留空自动生成）\r\n")
	fmt.Print("  registercode [--magicword X] [--expire 60m]  生成注册码\r\n")
	fmt.Print("  help                                   显示帮助\r\n")
	fmt.Print("  exit / quit                            退出\r\n")
	fmt.Print("\r\n")
	fmt.Print("方向键：↑↓ 历史命令，←→ 移动光标，Tab 补全命令，Ctrl-U 清空行。\r\n")
	fmt.Print("登录成功后 token 会保存到 data.yaml。\r\n")
	fmt.Print("\r\n")
}

