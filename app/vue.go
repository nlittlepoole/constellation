package main

import (
	"fmt"
	"github.com/zserge/webview"
	"html/template"
	"strings"
)

var cssPaths []string = []string{
	"assets/vendors/raleway/raleway.css",
	"assets/vendors/bulma/css/bulma.css",
	"assets/styles/main.css",
	"assets/styles/custom.css",
}

func loadUIFramework(w webview.WebView) {
	cssRaw := make([]string, 0)
	for _, path := range cssPaths {
		cssRaw = append(
			cssRaw,
			string(MustAsset(path)),
		)
	}
	w.Eval(fmt.Sprintf(`(function(css){
            var style = document.createElement('style');
            var head = document.head || document.getElementsByTagName('head')[0];
            style.setAttribute('type', 'text/css');
            if (style.styleSheet) {
                style.styleSheet.cssText = css;
            } else {
                style.appendChild(document.createTextNode(css));
            }
            head.appendChild(style);
        })("%s")`, template.JSEscapeString(strings.Join(cssRaw, "\n"))))

	w.Eval(string(MustAsset("assets/vendors/font-awesome/css/all.js")))
	w.Eval(string(MustAsset("assets/vendors/vue/vue.js")))
	w.Eval(string(MustAsset("assets/vendors/chartjs/Chart.js")))
	w.Eval(string(MustAsset("assets/vendors/chartjs/vue-chartjs.js")))
	w.Eval(string(MustAsset("assets/js/app.js")))

}
