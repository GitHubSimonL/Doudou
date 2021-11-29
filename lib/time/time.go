package time

import (
	"time"
)

//  时间常量定义
const (
	MsPerSecond        = 1000
	MsPerMinute        = 60000
	MsPerHour          = 3600000
	MsPerDay           = 86400000
	MsPerWeek          = 604800000
	NanoPerMs          = 1000000
	TheThousandthRatio = 10000 //  万分比
	RateBaseNum        = 10000
)

const (
	DateLayout         = "2006-01-02"
	DateTimeLayout     = "2006-01-02 15:04:05"
	DateTimeNanoLayout = "2006-01-02 15:04:05.000000"
)

var (
	timeOffset int64 //  时间偏移，毫秒
)

// 注：
// 以下时间戳都是毫秒级别的

func ResetTimeOffset() {
	timeOffset = 0
}

func AddTimeOffset(offset int64) {
	timeOffset += offset
}

func GetTimeOffset() int64 {
	return timeOffset
}

func MS(t time.Time) int64 {
	return t.UnixNano() / NanoPerMs
}

//  返回unix时间戳。
func CurrentMS() int64 {
	now := time.Now()
	return MS(now) + timeOffset
}

// 返回当前时间
func NowTime() time.Time {
	return Ms2Time(CurrentMS())
}

func CurrentMsTime() string {
	return NowTime().Format(DateTimeNanoLayout)
}

func CurrentTime() string {
	return NowTime().Format(DateTimeLayout)
}

func CurrentDate() string {
	return NowTime().Format(DateLayout)
}

func FormatTime(ms int64) string {
	return Ms2Time(ms).Format(DateTimeLayout)
}

func FormatMsTime(ms int64) string {
	return Ms2Time(ms).Format(DateTimeNanoLayout)
}

func Ms2Time(ms int64) (result time.Time) {
	if ms == 0 {
		return
	}

	sec := ms / 1e3
	nsec := (ms % 1e3) * 1e6
	result = time.Unix(sec, nsec).UTC()
	return
}

//  date format: "2006-01-02 13:04:00"
func S2UnixTime(value string) int64 {
	t, err := time.Parse(DateTimeLayout, value)
	if err != nil {
		return 0
	}
	return t.UnixNano() / 1000000
}

func GetMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//  从一个毫秒时间戳获得当前时区的本日凌晨时间。
func Ms2Midnight(t int64) time.Time {
	midTime := GetMidnight(Ms2Time(t))
	return midTime
}

func CurMidnight() int64 {
	return MS(GetMidnight(NowTime()))
}

func NextMidnight(t int64) int64 {
	midTime := GetMidnight(Ms2Time(t))
	return midTime.UnixNano()/1e6 + MsPerDay
}

//  从一个毫秒时间戳获取下一个准点时间。
func NextHour(t int64) int64 {
	t1 := Ms2Time(t)
	year, month, day := t1.Date()
	hour, _, _ := t1.Clock()
	t2 := time.Date(year, month, day, hour+1, 0, 0, 0, t1.Location())
	return t2.UnixNano() / 1e6
}

//  同一天
func OtherDay(curTime, lstTime int64) bool {
	return curTime > NextMidnight(lstTime)
}

//  是否同一个月
func SameMonth(curTime, lstTime int64) bool {
	t1 := time.Unix(curTime/1000, curTime%1000).UTC()
	y1, m1, _ := t1.Date()

	t2 := time.Unix(lstTime/1000, lstTime%1000).UTC()
	y2, m2, _ := t2.Date()

	return y1 == y2 && m1 == m2
}

func SameWeek(curTime, lstTime int64) bool {
	t1 := time.Unix(curTime/1000, curTime%1000).UTC()
	y1, w1 := t1.ISOWeek()

	t2 := time.Unix(lstTime/1000, lstTime%1000).UTC()
	y2, w2 := t2.ISOWeek()

	return y1 == y2 && w1 == w2
}

func SameDay(curTime, lstTime int64) bool {
	midnight := Ms2Midnight(lstTime)
	begin := MS(midnight)
	end := begin + MsPerDay
	return begin <= curTime && curTime < end
}

func CurDate() (year int, month int, day int) {
	t := NowTime()
	year, m, day := t.Date()
	month = int(m)
	return
}

func CurHour() int {
	t := time.Unix(CurrentMS()/1000, 0).UTC()
	return t.Hour()
}

// 获取当前时间到下次周几凌晨的毫秒数
func GetNextWeekXMs(x time.Weekday) int64 {
	t := NowTime()
	tw := t.Weekday()
	days := 0
	if tw >= x {
		// 计算下周
		days = int(time.Saturday) - int(tw) + int(x) - int(time.Sunday) + 1
	} else {
		// 计算本周
		days = int(x) - int(tw)
	}

	return int64(days*MsPerDay) - int64(t.Sub(GetMidnight(t)).Seconds()*MsPerSecond)
}

// 获取指定时间到下次周X的凌晨毫秒数
func GetCurNextWeekXMs(cur int64, x time.Weekday) int32 {
	t := Ms2Time(cur)
	tw := t.Weekday()
	days := 0
	if tw >= x {
		// 计算下周
		days = int(time.Saturday) - int(tw) + int(x) - int(time.Sunday) + 1
	} else {
		// 计算本周
		days = int(x) - int(tw)
	}

	return int32(days*MsPerDay) - int32(t.Sub(GetMidnight(t)).Seconds()*MsPerSecond)
}

func NextWeekDayMS(now int64, day time.Weekday) int64 {
	tm := Ms2Midnight(now)
	d := tm.Weekday()
	diff := time.Weekday(0)
	if d >= day {
		diff = time.Weekday(6) + day - d
	} else {
		diff = day - d
	}

	return MS(tm.AddDate(0, 0, int(diff)))
}

func GetWeekday(ms int64) int {
	return int(Ms2Time(ms).Weekday())
}

func BeginningOfHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

func BeginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return BeginningOfHour(t).Add(d)
}

func BeginningOfWeek(t time.Time, firstDayMonday bool) time.Time {
	t2 := BeginningOfDay(t)
	weekday := int(t2.Weekday())
	if firstDayMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}
	d := time.Duration(-weekday) * 24 * time.Hour
	return t2.Add(d)
}

func EndOfWeek(t time.Time, firstDayMonday bool) time.Time {
	return BeginningOfWeek(t, firstDayMonday).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func SubTimeDay(t1, t2 int64) int {
	tm1 := time.Unix(t1/1000, 0).UTC()
	tm2 := time.Unix(t2/1000, 0).UTC()
	tm1 = GetMidnight(tm1)
	tm2 = GetMidnight(tm2)
	return int(tm1.Sub(tm2).Hours() / 24)
}

//  从tm到现在是第几天
func WhichDayFrom(tm int64) int32 {
	hours := NowTime().Sub(Ms2Midnight(tm)).Hours()
	return int32(hours/24) + 1
}
