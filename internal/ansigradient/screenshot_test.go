package ansigradient

import (
	"fmt"
	"testing"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/gchalk/pkg/ansistyles"
)

func TestScreenshot(t *testing.T) {
	fmt.Println()
	fmt.Println()
	fmt.Println(`  // Background and foreground gradients`)

	g := gchalk.New(gchalk.ForceLevel(gchalk.LevelAnsi16m))

	fmt.Println(`  fg := CSSLinearGradientMust("#ff0, #0ff")`)
	fmt.Println(`  bg := CSSLinearGradientMust("#f00, #00f")`)
	fg := CSSLinearGradientMust("#ff0, #0ff")
	bg := CSSLinearGradientMust("#f00, #00f")

	fmt.Println(`  ApplyGradients("Hello World", fg, nil)`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Hello World", fg, nil, LevelAnsi16m))
	fmt.Println(`  ApplyGradients("Hello World", nil, bg)`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Hello World", nil, bg, LevelAnsi16m))
	fmt.Println(`  ApplyGradients("Hello World", fg, bg)`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Hello World", fg, bg, LevelAnsi16m))
	fmt.Println()

	fmt.Println(`  // CSS Style stops`)

	w := func(str string) string {
		return ansistyles.BrightWhite.Open + str + ansistyles.BrightWhite.Close
	}
	fmt.Println(`  w := func(str string) string {`)
	fmt.Println(`      return ansistyles.WhiteBright.Open + str + ansistyles.WhiteBright.Close`)
	fmt.Println(`  }`)

	fmt.Println(`  CSSLinearGradientMust("#F00, #FF0, #0F0, #0FF, #00F, #F0F").Render(...)`)
	rainbow := CSSLinearGradientMust("#F00, #FF0, #0F0, #0FF, #00F, #F0F")
	fmt.Println(`   => ` + ApplyGradientsRaw("You can use hex colors to create gradients.", rainbow, nil, LevelAnsi16m))
	fmt.Println(`  ApplyGradients("...", CSSLinearGradientMust("#F00, 20%, #0FF"), nil)`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Reach the color midpoint at 20%", CSSLinearGradientMust("#F00, 20%, #0FF"), nil, LevelAnsi16m))
	fmt.Println(`  ApplyGradients(w("..."), nil, CSSLinearGradientMust("#000 20%, #00F 50%, #000 80%")`)
	fmt.Println(`   => ` + ApplyGradientsRaw(w("Stops can be offset from the start and end"), nil, CSSLinearGradientMust("#000 20%, #00F 50%, #000 80%"), LevelAnsi16m))
	fmt.Println(`  ApplyGradients(w("..."), nil, CSSLinearGradientMust("#007, #707 25% 75%, #007")`)
	fmt.Println(`   => ` + ApplyGradientsRaw(w("Multiple stops on the same color makes a big solid area"), nil, CSSLinearGradientMust("#007, #F0F 25% 75%, #007"), LevelAnsi16m))
	fmt.Println(`  ApplyGradients(w("..."), nil, CSSLinearGradientMust("#007 50%, #700 50%")`)
	fmt.Println(`   => ` + ApplyGradientsRaw(w("Stops at same place will make a sharp transition"), nil, CSSLinearGradientMust("#007 50%, #700 50%"), LevelAnsi16m))
	fmt.Println()

	fmt.Println(`  // Nested styles with gchalk`)
	fmt.Println(`  gradient := CSSLinearGradientMust("#F00, #FCC")`)
	gradient := CSSLinearGradientMust("#F00, #FCC")
	fmt.Println(`  ApplyGradients("Colorful "+gchalk.Green("green")+", and colorful again", gradient, nil))`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Colorful "+g.Green("green")+", and colorful again", gradient, nil, LevelAnsi16m))
	fmt.Println(`  gchalk.Green("Green, " + ApplyGradients("colorful", gradient, nil) + ", and green again")`)
	fmt.Println(`   => ` + g.Green("Green, "+ApplyGradientsRaw("colorful", gradient, nil, LevelAnsi16m)+", and green again"))
	fmt.Println()

	fmt.Println(`  // 256 Color Mode`)
	fmt.Println(`  ApplyGradientsRaw("...", CSSLinearGradientMust("#f00, #00f"), nil, LevelAnsi256)`)
	fmt.Println(`   => ` + ApplyGradientsRaw("Support for Ansi 256 color mode", CSSLinearGradientMust("#f00, #00f"), nil, LevelAnsi256))

	fmt.Println()
	fmt.Println()
}
