package rocket

import (
	"fmt"
	"github.com/lestrrat/go-xslate"
)


type XslateView struct {
	View
	xt *xslate.Xslate // TODO ViewごとにXslateオブジェクト作っていいのかな
}

func NewXslateView() *XslateView {
	view := new(XslateView)

	var err error
	view.xt, err = xslate.New(xslate.Args{
		"Parser": xslate.Args{"Syntax": "TTerse"},
	})
	if err != nil {
		panic(fmt.Sprintf("Xslate initiate error: err=%v", err))
	}
	return view
}

// override
func (v *XslateView) Render(tmplFile string, bind RenderVars) string {
	result, err := v.xt.Render(tmplFile, xslate.Vars(bind))
	if err != nil {
		panic(fmt.Sprintf("Xslate render error: err=%v", err))
	}

	return result
}



