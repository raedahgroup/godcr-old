package helper 

import (
	"os"
	"image"
	"time"
	"math"
	"strconv"

	"gioui.org/unit"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

const (
	logoPath = "../../gio/assets/decred.png"
	StandaloneScreenPadding = 20
)

var logo material.Image

func InitLogo(theme *Theme) error {
	logoByte, err := os.Open(logoPath)
	if err != nil {
		return err
	}

	src, _, err := image.Decode(logoByte) 
	if err != nil {
		return err
	}

	logo = theme.Image(paint.NewImageOp(src))
	logo.Scale = 1.3

	return nil
}

func DrawLogo(ctx *layout.Context) {
	inset := layout.Inset{
		Left: unit.Dp(StandaloneScreenPadding),
	}
	inset.Layout(ctx, func(){
		logo.Layout(ctx)
	})
}

func TimeAgo(then time.Time) string {
	now := time.Now()

	yearNow, monthNow, dayNow := now.Date()
	hourNow, minuteNow, secondNow := now.Clock()

	yearThen, monthThen, dayThen := then.Date()
	hourThen, minuteThen, secondThen := then.Clock()

	year := math.Abs(float64((int(yearNow - yearThen))))
	month := math.Abs(float64((int(monthNow - monthThen))))
	day := math.Abs(float64((int(dayNow - dayThen))))
	hour := math.Abs(float64((int(hourNow - hourThen))))
	minute := math.Abs(float64((int(minuteNow - minuteThen))))
	second := math.Abs(float64((int(secondNow - secondThen))))

	week := math.Floor(day / 7)

	if year >  0 {
		txt := " years ago"
		if year == 1 {
			txt = " year ago"
		}
		return strconv.Itoa(int(year)) + txt
	} else if month > 0 {
		txt := " months ago"
		if month == 1 {
			txt = " month ago"
		}
		return strconv.Itoa(int(month)) + txt
	} else if week > 0 {
		txt := " weeks ago"
		if week == 1 {
			txt = " week ago"
		}
		return strconv.Itoa(int(week)) + txt
	} else if day > 0 {
		txt := " days ago"
		if day == 1 {
			txt = " day ago"
		}
		return strconv.Itoa(int(day)) + txt
	} else  if hour > 0 {
		txt := " hours ago"
		if hour == 1 {
			txt = " hour ago"
		}
		return strconv.Itoa(int(hour)) + txt
	} else if minute > 0 {
		txt := " minutes ago"
		if minute == 1 {
			txt = " minute ago"
		}
		return strconv.Itoa(int(minute)) + txt
	} 
	
	txt := " seconds ago"
	if second == 1 {
		txt = " second ago"
	}
	return strconv.Itoa(int(second)) + txt
}