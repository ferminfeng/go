// @Author fermin 2024/6/3 20:42:00
package main

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

func main() {
	fmt.Println("当前时间", GetCurrentTimeFormatted())

	t1 := TimeStringFormatTimeUnix(time.DateTime, "2024-06-03 00:00:00")
	fmt.Println("设置本地之间，将[2024-06-03 00:00:00]转为时间戳", t1)
	fmt.Println("设置本地之间，将[2024-06-03 00:00:00]的时间戳再转为时间", time.Unix(t1, 0).Format(time.DateTime))

	t2, _ := time.Parse(time.DateTime, "2024-06-03 00:00:00")
	t22 := t2.Unix()
	fmt.Println("直接转，将[2024-06-03 00:00:00]转为时间戳", t22)
	fmt.Println("直接转，将[2024-06-03 00:00:00]的时间戳再转为时间", time.Unix(t22, 0).Format(time.DateTime))
}

const (
	TimeOffset = 8 * 3600  // 8 hour offset
	HalfOffset = 12 * 3600 // Half-day hourly offset
)

// Get the current timestamp by Second
func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

// Convert timestamp to time.Time type
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}

// UnixMillSecondToTime convert millSecond to time.Time type
func UnixMillSecondToTime(millSecond int64) time.Time {
	return time.Unix(0, millSecond*1e6)
}

// Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

// Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// Get the timestamp at 0 o'clock of the day
func GetCurDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - TimeOffset
}

// Get the timestamp at 12 o'clock on the day
func GetCurDayHalfTimestamp() int64 {
	return GetCurDayZeroTimestamp() + HalfOffset

}

// Get the formatted time at 0 o'clock of the day, the format is "2006-01-02_00-00-00"
func GetCurDayZeroTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp(), 0).Format("2006-01-02_15-04-05")
}

// Get the formatted time at 12 o'clock of the day, the format is "2006-01-02_12-00-00"
func GetCurDayHalfTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp()+HalfOffset, 0).Format("2006-01-02_15-04-05")
}

// TimeStringFormatTimeUnix convert string to unix timestamp
func TimeStringFormatTimeUnix(timeFormat string, timeSrc string) int64 {
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeFormat, timeSrc, loc)
	return tmp.Unix()
}

// TimeStringFormatTimeUnixMilli convert string to unix timestamp
func TimeStringFormatTimeUnixMilli(timeFormat string, timeSrc string) int64 {
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeFormat, timeSrc, loc)
	return tmp.UnixMilli()
}

// TimeStringToTime convert string to time.Time
func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}

// TimeStringFormatToTime convert string to time.Time
func TimeStringFormatToTime(timeFormat string, timeString string) time.Time {
	t, _ := time.Parse(timeFormat, timeString)
	return t
}

// TimeToString convert time.Time to string
func TimeToString(timeFormat string, t time.Time) string {
	return t.Format(timeFormat)
}

func GetCurrentTimeFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
