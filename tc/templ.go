// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

import (
	"fmt"
)

/***********************************************************************************/
func resultHtml(s ...any) (_r string) {
	_r = `<html>
		<body style="text-align:center;">
			<h2>` + fmt.Sprint(s...) + `</h2>
			<h3><span id='s'></span></h3>
			<h4><a href="javascript:window.history.go(-1)"> click here </a> go back。<h4>
		</body>
		<script type="text/javascript">
			var t = 5; 
			function trans() {
				if (t == 0) {window.history.go(-1);}
				document.getElementById('s').innerHTML = "go back after " + t + " seconds"; 
				t--;
			}
			setInterval("trans()", 1000); 
		</script>
	</html>`
	return
}

/***********************************************************************************/
func resultHtmlAndClose(s ...any) (_r string) {
	_r = `<html>
		<body style="text-align:center;">
			<h2>` + fmt.Sprint(s...) + `</h2>
			<h3><span id='s'></span></h3>
			<h4><a href="javascript:window.close();"> click here </a> close page。<h4>
		</body>
		<script type="text/javascript">
			var t = 5; 
			function trans() {
				if (t == 0) {window.close();}
				document.getElementById('s').innerHTML = "close page after " + t + " seconds"; 
				t--;
			}
			setInterval("trans()", 1000); 
		</script>
	</html>`
	return
}
