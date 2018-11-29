package log

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

/// This Formatter allows for custom formatting of
/// log messages based on format patterns.
///
/// The format pattern is used as a template to format the message and
/// is copied character by character except for the following special characters,
/// which are replaced by the corresponding value.
///
///   * %s - message source
///   * %t - message text
///   * %l - message priority level (1 .. 7)
///   * %p - message priority (Fatal, Critical, Error, Warning, Notice, Information, Debug, Trace)
///   * %q - abbreviated message priority (F, C, E, W, N, I, D, T)
///   * %T - message thread name
///   * %I - message thread identifier (numeric)
///   * %O - message thread OS identifier (numeric)
///   * %N - node or host name
///   * %P - message file path
///   * %U - message source file name (empty string if not set)
///   * %u - message source line number (0 if not set)
///   * %w - message date/time abbreviated weekday (Mon, Tue, ...)
///   * %W - message date/time full weekday (Monday, Tuesday, ...)
///   * %b - message date/time abbreviated month (Jan, Feb, ...)
///   * %B - message date/time full month (January, February, ...)
///   * %d - message date/time zero-padded day of month (01 .. 31)
///   * %e - message date/time day of month (1 .. 31)
///   * %f - message date/time space-padded day of month ( 1 .. 31)
///   * %m - message date/time zero-padded month (01 .. 12)
///   * %n - message date/time month (1 .. 12)
///   * %o - message date/time space-padded month ( 1 .. 12)
///   * %y - message date/time year without century (70)
///   * %Y - message date/time year with century (1970)
///   * %H - message date/time hour (00 .. 23)
///   * %h - message date/time hour (00 .. 12)
///   * %a - message date/time am/pm
///   * %A - message date/time AM/PM
///   * %M - message date/time minute (00 .. 59)
///   * %S - message date/time second (00 .. 59)
///   * %i - message date/time millisecond (000 .. 999)
///   * %c - message date/time centisecond (0 .. 9)
///   * %F - message date/time fractional seconds/microseconds (000000 - 999999)
///   * %z - time zone differential in ISO 8601 format (Z or +NN.NN)
///   * %Z - time zone differential in RFC format (GMT or +NNNN)
///   * %L - convert time to local time (must be specified before any date/time specifier; does not itself output anything)
///   * %E - epoch time (UTC, seconds since midnight, January 1, 1970)
///   * %v[width] - the message source (%s) but text length is padded/cropped to 'width'
///   * %[name] - the value of the message parameter with the given name
///   * %% - percent sign

/// json pattern like :{"time"="%y-%m=%d %H:%M:%S", "level"="%q", "file"="%U:%u", "text"="%t"}
/// DefaultPattern set default pattern
const DefaultPattern = "[%q %y-%m-%d %H:%M:%S %U[10]:%u[3]] %t"

func NewPattern() *Pattern {
	p := &Pattern{}
	return p
}

type Action struct {
	Key      rune
	Prepend  string // xxx%
	Property string // %[name]
}

// Pattern 输出格式解析
type Pattern struct {
	actions []*Action
}

func buildWidth(builder *bytes.Buffer, text string, width string) {
	if width != "" {
		format := "%" + width + "s"
		data := fmt.Sprintf(format, text)
		builder.WriteString(data)
	} else {
		builder.WriteString(text)
	}
}

func (pattern *Pattern) Format(msg *Message) string {
	// fmt.Printf("%+v, %+v\n", msg.Level, LEVEL_ALL)
	builder := bytes.Buffer{}

	timestamp := msg.Timetamp
	sec := timestamp / 1000
	nsec := (timestamp - sec*1000) * 1000
	ts := time.Unix(sec, nsec)

	for _, action := range pattern.actions {
		builder.WriteString(action.Prepend)
		switch action.Key {
		case 's': // TODO:source
		case 't':
			builder.WriteString(msg.Text)
		case 'l':
			builder.WriteString(strconv.Itoa(msg.Level))
		case 'p':
			builder.WriteString(gLevelName[msg.Level])
		case 'q':
			builder.WriteByte(gLevelName[msg.Level][0])
		case 'T': // TODO:thread
		case 'I': // TODO:tid
		case 'O': // ostid
		case 'N': // node name
		case 'P': // path,in poco is pid
			buildWidth(&builder, msg.File, action.Property)
		case 'U':
			buildWidth(&builder, msg.FileName, action.Property)
		case 'u':
			buildWidth(&builder, strconv.Itoa(msg.Line), action.Property)
			// builder.WriteString(strconv.Itoa(msg.Line))
		case 'w':
			builder.WriteString(ts.Weekday().String()[0:3])
		case 'W':
			builder.WriteString(ts.Weekday().String())
		case 'b':
			builder.WriteString(ts.Month().String()[0:3])
		case 'B':
			builder.WriteString(ts.Month().String())
		case 'd':
			builder.WriteString(fmt.Sprintf("%02d", ts.Day()))
		case 'e':
			builder.WriteString(fmt.Sprintf("%d", ts.Day()))
		case 'f':
			builder.WriteString(fmt.Sprintf("%2d", ts.Day()))
		case 'm':
			builder.WriteString(fmt.Sprintf("%02d", ts.Month()))
		case 'n':
			builder.WriteString(fmt.Sprintf("%d", ts.Month()))
		case 'o':
			builder.WriteString(fmt.Sprintf("%2d", ts.Month()))
		case 'y':
			builder.WriteString(fmt.Sprintf("%02d", ts.Year()%100))
		case 'Y':
			builder.WriteString(fmt.Sprintf("%04d", ts.Year()))
		case 'H':
			builder.WriteString(fmt.Sprintf("%02d", ts.Hour()))
		case 'h':
			hour := ts.Hour()
			if hour < 1 {
				hour = 12
			} else if hour > 12 {
				hour -= 12
			}
			builder.WriteString(fmt.Sprintf("%02d", hour))
		case 'a':
			if ts.Hour() < 12 {
				builder.WriteString("am")
			} else {
				builder.WriteString("pm")
			}
		case 'A':
			if ts.Hour() < 12 {
				builder.WriteString("AM")
			} else {
				builder.WriteString("PM")
			}
		case 'M':
			builder.WriteString(fmt.Sprintf("%02d", ts.Minute()))
		case 'S':
			builder.WriteString(fmt.Sprintf("%02d", ts.Second()))
		case 'i':
			builder.WriteString(fmt.Sprintf("%03d", ts.Nanosecond()/int(time.Millisecond)))
		case 'c':
		case 'F':
		case 'z': // tzdISO
		case 'Z': // tzdRFC
		case 'E':
			builder.WriteString(fmt.Sprintf("%d", ts.Unix()))
		case 'v': // source width
		case 'x':
			// property
			builder.WriteString(msg.Data[action.Property])
		case 'L':
		case '%':
			builder.WriteByte('%')
		case '[':
			builder.WriteByte('[')
		}
	}
	builder.WriteByte('\n')
	return builder.String()
}

// Parse 通用解析规则%xx[prop]
func (pattern *Pattern) Parse(format string) {
	actions := make([]*Action, 0)

	end := len(format)
	cur := 0
	for cur < end {
		act := &Action{}
		// parse prepend
		for beg := cur; ; cur++ {
			if cur >= end || format[cur] == '%' {
				if beg < cur {
					act.Prepend = format[beg:cur]
				}
				break
			}
		}

		// check end
		if cur == end {
			actions = append(actions, act)
			break
		}
		cur++ // ignore %

		// parse key
		if format[cur] == '[' {
			act.Key = 'x'
		} else {
			act.Key = rune(format[cur])
			cur++
		}
		// parse property
		if cur < end && format[cur] == '[' {
			cur++
			for beg := cur; cur < end; cur++ {
				if format[cur] == ']' {
					act.Property = format[beg:cur]
					cur++
					break
				}
			}
		}
		actions = append(actions, act)
	}

	pattern.actions = actions
}
