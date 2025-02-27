// Package cpstory cp短打
package cpstory

import (
	"strings"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/math"

	"github.com/FloatTech/zbputils/control/order"
)

func init() {
	engine := control.Register("cpstory", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help:             "cp短打\n- 组cp[@xxx][@xxx]\n- 磕cp大老师 雪乃",
		PublicDataFolder: "CpStory",
	})

	go func() {
		dbpath := engine.DataFolder()
		db.DBPath = dbpath + "cp.db"
		// os.RemoveAll(dbpath)
		_, _ = file.GetLazyData(db.DBPath, false, true)
		err := db.Create("cp_story", &cpstory{})
		if err != nil {
			panic(err)
		}
		n, err := db.Count("cp_story")
		if err != nil {
			panic(err)
		}
		logrus.Printf("[cpstory]读取%d条故事", n)
	}()

	engine.OnRegex("^组cp.*?(\\d+).*?(\\d+)", zero.OnlyGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cs := getRandomCpStory()
		gong := ctxext.CardOrNickName(ctx, math.Str2Int64(ctx.State["regex_matched"].([]string)[1]))
		shou := ctxext.CardOrNickName(ctx, math.Str2Int64(ctx.State["regex_matched"].([]string)[2]))
		text := strings.ReplaceAll(cs.Story, "<攻>", gong)
		text = strings.ReplaceAll(text, "<受>", shou)
		text = strings.ReplaceAll(text, cs.Gong, gong)
		text = strings.ReplaceAll(text, cs.Shou, gong)
		ctx.SendChain(message.Text(text))
	})
	engine.OnPrefix("磕cp").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		cs := getRandomCpStory()
		params := strings.Split(ctx.State["args"].(string), " ")
		if len(params) < 2 {
			ctx.SendChain(message.Text(ctx.Event.MessageID), message.Text("请用空格分开两个人名"))
		} else {
			gong := params[0]
			shou := params[1]
			text := strings.ReplaceAll(cs.Story, "<攻>", gong)
			text = strings.ReplaceAll(text, "<受>", shou)
			text = strings.ReplaceAll(text, cs.Gong, gong)
			text = strings.ReplaceAll(text, cs.Shou, gong)
			ctx.SendChain(message.Text(text))
		}
	})
}
