package main

import (
	"fmt"
	"io"
)

func templateIndex(w io.Writer, arg *View) {
	io.WriteString(w, "<!DOCTYPE html>\n<html>\n  <head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <title>Isucon3</title>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap.min.css\">\n    <style>\nbody {\n  padding-top: 60px;\n}\n    </style>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap-responsive.min.css\">\n  </head>\n  <body>\n    <div class=\"navbar navbar-fixed-top\">\n      <div class=\"navbar-inner\">\n        <div class=\"container\">\n          <a class=\"btn btn-navbar\" data-toggle=\"collapse\" data-target=\".nav-collapse\">\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n          </a>\n          <a class=\"brand\" href=\"/\">Isucon3</a>\n          <div class=\"nav-collapse\">\n            <ul class=\"nav\">\n              <li><a href=\"http://localhost/\">Home</a></li>\n              ")
	if arg.User != nil {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/mypage\">MyPage</a></li>\n              <li>\n                <form action=\"/signout\" method=\"post\">\n                  <input type=\"hidden\" name=\"sid\" value=\"")
		io.WriteString(w, string(arg.Token))
		io.WriteString(w, "\">\n                  <input type=\"submit\" value=\"SignOut\">\n                </form>\n              </li>\n              ")
	} else {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/signin\">SignIn</a></li>\n              ")
	}
	io.WriteString(w, "\n            </ul>\n          </div> <!--/.nav-collapse -->\n        </div>\n      </div>\n    </div>\n\n    <div class=\"container\">\n      <h2>Hello ")
	if arg.User != nil {
		io.WriteString(w, string(arg.User.Username))
	}
	io.WriteString(w, "!</h2>\n\n      <h3>public memos</h3>\n      <p id=\"pager\">\n      recent ")
	io.WriteString(w, fmt.Sprintf("%d", arg.PageStart))
	io.WriteString(w, " - ")
	io.WriteString(w, fmt.Sprintf("%d", arg.PageEnd))
	io.WriteString(w, " / total <span id=\"total\">")
	io.WriteString(w, fmt.Sprintf("%d", arg.Total))
	io.WriteString(w, "</span>\n      </p>\n      <ul id=\"memos\">\n        ")
	if len(*arg.Memos) > 0 {
		for _, arg := range *arg.Memos {
			io.WriteString(w, "\n        <li>\n          <a href=\"http://localhost/memo/")
			io.WriteString(w, fmt.Sprintf("%d", arg.Id))
			io.WriteString(w, "\">")
			io.WriteString(w, string(arg.Title))
			io.WriteString(w, "</a> by ")
			io.WriteString(w, string(arg.Username))
			io.WriteString(w, " (")
			io.WriteString(w, string(arg.CreatedAt))
			io.WriteString(w, ")\n        </li>\n        ")
		}
	}
	io.WriteString(w, "\n      </ul>\n\n    </div> <!-- /container -->\n\n    <script type=\"text/javascript\" src=\"http://localhost/js/jquery.min.js\"></script>\n    <script type=\"text/javascript\" src=\"http://localhost/js/bootstrap.min.js\"></script>\n  </body>\n</html>\n")
}

func templateMemo(w io.Writer, arg *View) {
	io.WriteString(w, "<!DOCTYPE html>\n<html>\n  <head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <title>Isucon3</title>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap.min.css\">\n    <style>\nbody {\n  padding-top: 60px;\n}\n    </style>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap-responsive.min.css\">\n  </head>\n  <body>\n    <div class=\"navbar navbar-fixed-top\">\n      <div class=\"navbar-inner\">\n        <div class=\"container\">\n          <a class=\"btn btn-navbar\" data-toggle=\"collapse\" data-target=\".nav-collapse\">\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n          </a>\n          <a class=\"brand\" href=\"/\">Isucon3</a>\n          <div class=\"nav-collapse\">\n            <ul class=\"nav\">\n              <li><a href=\"http://localhost/\">Home</a></li>\n              ")
	if arg.User != nil {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/mypage\">MyPage</a></li>\n              <li>\n                <form action=\"/signout\" method=\"post\">\n                  <input type=\"hidden\" name=\"sid\" value=\"")
		io.WriteString(w, string(arg.Token))
		io.WriteString(w, "\">\n                  <input type=\"submit\" value=\"SignOut\">\n                </form>\n              </li>\n              ")
	} else {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/signin\">SignIn</a></li>\n              ")
	}
	io.WriteString(w, "\n            </ul>\n          </div> <!--/.nav-collapse -->\n        </div>\n      </div>\n    </div>\n\n    <div class=\"container\">\n      <h2>Hello ")
	if arg.User != nil {
		io.WriteString(w, string(arg.User.Username))
	}
	io.WriteString(w, "!</h2>\n\n      <p id=\"author\">\n      ")
	if arg.Memo.IsPrivate == 1 {
		io.WriteString(w, "\n      Private\n      ")
	} else {
		io.WriteString(w, "\n      Public\n      ")
	}
	io.WriteString(w, "\n      Memo by ")
	io.WriteString(w, string(arg.Memo.Username))
	io.WriteString(w, " (")
	io.WriteString(w, string(arg.Memo.CreatedAt))
	io.WriteString(w, ")\n      </p>\n\n      <hr>\n      ")
	if arg.Older != nil {
		io.WriteString(w, "\n      <a id=\"older\" href=\"http://localhost/memo/")
		io.WriteString(w, fmt.Sprintf("%d", arg.Older.Id))
		io.WriteString(w, "\">&lt; older memo</a>\n      ")
	}
	io.WriteString(w, "\n      |\n      ")
	if arg.Newer != nil {
		io.WriteString(w, "\n      <a id=\"newer\" href=\"http://localhost/memo/")
		io.WriteString(w, fmt.Sprintf("%d", arg.Newer.Id))
		io.WriteString(w, "\">newer memo &gt;</a>\n      ")
	}
	io.WriteString(w, "\n\n      <hr>\n      <div id=\"content_html\">\n        ")
	io.WriteString(w, string(gen_markdown(arg.Memo.Content)))
	io.WriteString(w, "\n      </div>\n\n    </div> <!-- /container -->\n\n    <script type=\"text/javascript\" src=\"http://localhost/js/jquery.min.js\"></script>\n    <script type=\"text/javascript\" src=\"http://localhost/js/bootstrap.min.js\"></script>\n  </body>\n</html>\n")
}

func templateMypage(w io.Writer, arg *View) {
	io.WriteString(w, "<!DOCTYPE html>\n<html>\n  <head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <title>Isucon3</title>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap.min.css\">\n    <style>\nbody {\n  padding-top: 60px;\n}\n    </style>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap-responsive.min.css\">\n  </head>\n  <body>\n    <div class=\"navbar navbar-fixed-top\">\n      <div class=\"navbar-inner\">\n        <div class=\"container\">\n          <a class=\"btn btn-navbar\" data-toggle=\"collapse\" data-target=\".nav-collapse\">\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n          </a>\n          <a class=\"brand\" href=\"/\">Isucon3</a>\n          <div class=\"nav-collapse\">\n            <ul class=\"nav\">\n              <li><a href=\"http://localhost/\">Home</a></li>\n              ")
	if arg.User != nil {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/mypage\">MyPage</a></li>\n              <li>\n                <form action=\"/signout\" method=\"post\">\n                  <input type=\"hidden\" name=\"sid\" value=\"")
		io.WriteString(w, string(arg.Token))
		io.WriteString(w, "\">\n                  <input type=\"submit\" value=\"SignOut\">\n                </form>\n              </li>\n              ")
	} else {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/signin\">SignIn</a></li>\n              ")
	}
	io.WriteString(w, "\n            </ul>\n          </div> <!--/.nav-collapse -->\n        </div>\n      </div>\n    </div>\n\n    <div class=\"container\">\n      <h2>Hello ")
	if arg.User != nil {
		io.WriteString(w, string(arg.User.Username))
	}
	io.WriteString(w, "!</h2>\n\n      <form action=\"http://localhost/memo\" method=\"post\">\n        <input type=\"hidden\" name=\"sid\" value=\"")
	io.WriteString(w, string(arg.Token))
	io.WriteString(w, "\">\n        <textarea name=\"content\"></textarea>\n        <br>\n        <input type=\"checkbox\" name=\"is_private\" value=\"1\"> private\n        <input type=\"submit\" value=\"post\">\n      </form>\n\n      <h3>my memos</h3>\n\n      <ul>\n        ")
	if len(*arg.Memos) > 0 {
		for _, arg := range *arg.Memos {
			io.WriteString(w, "\n        <li>\n          <a href=\"http://localhost/memo/")
			io.WriteString(w, fmt.Sprintf("%d", arg.Id))
			io.WriteString(w, "\">")
			io.WriteString(w, string(arg.Title))
			io.WriteString(w, "</a> by ")
			io.WriteString(w, string(arg.Username))
			io.WriteString(w, " (")
			io.WriteString(w, string(arg.CreatedAt))
			io.WriteString(w, ")\n          ")
			if arg.IsPrivate == 1 {
				io.WriteString(w, "\n          [private]\n          ")
			}
			io.WriteString(w, "\n        </li>\n        ")
		}
	}
	io.WriteString(w, "\n      </ul>\n\n    </div> <!-- /container -->\n\n    <script type=\"text/javascript\" src=\"http://localhost/js/jquery.min.js\"></script>\n    <script type=\"text/javascript\" src=\"http://localhost/js/bootstrap.min.js\"></script>\n  </body>\n</html>\n")
}

func templateSignin(w io.Writer, arg *View) {
	io.WriteString(w, "<!DOCTYPE html>\n<html>\n  <head>\n    <meta http-equiv=\"Content-Type\" content=\"text/html\" charset=\"utf-8\">\n    <title>Isucon3</title>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap.min.css\">\n    <style>\nbody {\n  padding-top: 60px;\n}\n    </style>\n    <link rel=\"stylesheet\" href=\"http://localhost/css/bootstrap-responsive.min.css\">\n  </head>\n  <body>\n    <div class=\"navbar navbar-fixed-top\">\n      <div class=\"navbar-inner\">\n        <div class=\"container\">\n          <a class=\"btn btn-navbar\" data-toggle=\"collapse\" data-target=\".nav-collapse\">\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n            <span class=\"icon-bar\"></span>\n          </a>\n          <a class=\"brand\" href=\"/\">Isucon3</a>\n          <div class=\"nav-collapse\">\n            <ul class=\"nav\">\n              <li><a href=\"http://localhost/\">Home</a></li>\n              ")
	if arg.User != nil {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/mypage\">MyPage</a></li>\n              <li>\n                <form action=\"/signout\" method=\"post\">\n                  <input type=\"hidden\" name=\"sid\" value=\"")
		io.WriteString(w, string(arg.Token))
		io.WriteString(w, "\">\n                  <input type=\"submit\" value=\"SignOut\">\n                </form>\n              </li>\n              ")
	} else {
		io.WriteString(w, "\n              <li><a href=\"http://localhost/signin\">SignIn</a></li>\n              ")
	}
	io.WriteString(w, "\n            </ul>\n          </div> <!--/.nav-collapse -->\n        </div>\n      </div>\n    </div>\n\n    <div class=\"container\">\n      <h2>Hello ")
	if arg.User != nil {
		io.WriteString(w, string(arg.User.Username))
	}
	io.WriteString(w, "!</h2>\n\n      <form action=\"http://localhost/signin\" method=\"post\">\n        username <input type=\"text\" name=\"username\" size=\"20\">\n        <br>\n        password <input type=\"password\" name=\"password\" size=\"20\">\n        <br>\n        <input type=\"submit\" value=\"signin\">\n      </form>\n\n    </div> <!-- /container -->\n\n    <script type=\"text/javascript\" src=\"http://localhost/js/jquery.min.js\"></script>\n    <script type=\"text/javascript\" src=\"http://localhost/js/bootstrap.min.js\"></script>\n  </body>\n</html>\n")
}
