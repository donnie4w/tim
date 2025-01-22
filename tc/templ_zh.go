// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

const (
	loginText = `
    <html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
</head>

<body class="container">
    <div class="container-fluid text-right">
        <span>
            <h4 style="display: inline;">Tim 管理后台</h4>
        </span>
        <span style="text-align:right">
            <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[EN]</a></h6>
        </span>
        <hr>
        <div id="login">
            <h5>登录</h5>
            <form class="form-control" id="loginform" action="/login" method="post">
                <input name="type" value="1" hidden />
                <input name="name" placeholder="用户名" />
                <input name="pwd" placeholder="密码" type="password" />
                <input type="submit" class="btn btn-primary" value="登录" />
            </form>
        </div>
        <hr>
    </div>
</html>
    `
	initText = `
    <html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-align-top" viewBox="0 0 16 16">
                    <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)" />
                    <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <div class="navbar-nav">
                    <a class="nav-link active" href='/init'>账号管理</a>
                    <a class="nav-link" href='/sysvar'>集群环境</a>
                    <a class="nav-link" href='/data'>数据监控</a>
                    <a class="nav-link" href='/monitor'>性能监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: small;">
        {{if .ShowCreate }}
        <div class="container-fluid card mt-1 p-1">
            </h6>
            <form class="form-control" id="createAdminform" action="/init?type=1" method="post">
                <h6>新建管理员 <h6 class="important">{{ .Show }}</h6>
                    <input name="adminName" placeholder="用户名" />
                    <input name="adminPwd" placeholder="密码" type="password" />
                    管理员<input name="adminType" type="radio" value="1" checked />
                    {{if not .Init}}
                    观察员<input name="adminType" type="radio" value="2" />
                    {{end}}
                    <input type="submit" class="btn btn-primary" value="新建管理员" />
            </form>
        </div>
        {{end}}
        {{if not .Init}}
        <div class="container-fluid card mt-1 p-1">
            <div class="m-2">
                <h6>后台管理员</h6>
                {{range $k,$v := .AdminUser}}
                <form class="form-control" id="adminform" action="/init?type=2" method="post">
                    <input name="adminName" value='{{ $k }}' readonly style="border:none;" /> 权限:{{ $v }}
                    <input type="button" value="删除用户" class="btn btn-danger"
                        onclick="javascipt:if (confirm('确定删除?')){this.parentNode.submit();};" />
                </form>
                {{end}}
            </div>
        </div>
        <hr>
        {{end}}
    </div>
</html>
    `

	sysvarText = `
    <html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
    <meta http-equiv="refresh" content="30">
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-align-top" viewBox="0 0 16 16">
                    <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)" />
                    <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <div class="navbar-nav">
                    <a class="nav-link" href='/init'>账号管理</a>
                    <a class="nav-link active" href='/sysvar'>集群环境</a>
                    <a class="nav-link" href='/data'>数据监控</a>
                    <a class="nav-link" href='/monitor'>性能监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: xx-small;">
        <table class="table table-bordered table-Info">
            <tr>
                <td style="width: 150px;">节点启动时间</td>
                <td colspan="2">{{ .SYS.StartTime }}</td>
            </tr>
            <tr>
                <td>当前时间</td>
                <td colspan="2">{{ .SYS.Time }}</td> 
            </tr>
            <tr>
                <td>节点UUID</td>
                <td class="text-danger" colspan="2">{{ .SYS.UUID }} [{{ .SYS.CSNUM }}]</td>
            </tr>
            <tr>
                <td>集群节点</td>
                <td class="text-danger" colspan="2">{{ .SYS.ALLUUIDS }}</td>
            </tr>

            <tr>
                <td>集群监听地址</td>
                <td colspan="2">{{ .SYS.ADDR }}</td>
            </tr>
        </table>
    </div>
</body>
</html>
    `

	dataText = `
    <html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-align-top" viewBox="0 0 16 16">
                    <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)" />
                    <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <div class="navbar-nav">
                    <a class="nav-link" href='/init'>账号管理</a>
                    <a class="nav-link" href='/sysvar'>集群环境</a>
                    <a class="nav-link active" href='/data'>数据监控</a>
                    <a class="nav-link" href='/monitor'>性能监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>用户数据监控</h3>
            <div class="input-group">
                <span class="input-group-text">监控时间间隔(单位:秒)</span>
                <input id="stime" placeholder="输入时间间隔" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">开始</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">停止</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">清除数据</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>当前节点在线用户数</th>
                <th>节点输入数据(B)</th>
                <th>节点输出数据(B)</th>
                <th>未连接节点</th>
            </tr>
            <tbody id="monitorBody">
            </tbody>
        </table>
    </div>
</body>
<script type="text/javascript">
    var pro = window.location.protocol;
    var wspro = "ws:";
    if (pro === "https:") {
        wspro = "wss:";
    }
    var wsmnt = null;
    var id = 1;
    function WS() {
        this.ws = null;
    }

    WS.prototype.monitor = function () {
        let obj = this;
        this.ws = new WebSocket(wspro + "//" + window.location.host + "/ddmonitorData");
        this.ws.onopen = function (evt) {
            obj.ws.send(document.getElementById("stime").value);
        }
        this.ws.onmessage = function (evt) {
            if (evt.data != "") {
                var json = JSON.parse(evt.data);
                var tr = document.createElement('tr');
                var d = '<td style="font-weight: bold;">' + id++ + '</td>'
                    + '<td>' + json.Online + '</td>'
                    + '<td>' + json.Input  + '</td>'
                    + '<td>' + json.Output  + '</td>'
                    + '<td>' + json.Unaccess + '</td>'
                tr.innerHTML = d;
                document.getElementById("monitorBody").appendChild(tr);
            }
        }
    }

    WS.prototype.close = function () {
        this.ws.close();
    }

    function monitorLoad() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
        wsmnt = new (WS);
        wsmnt.monitor();
    }

    function stop() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
    }

    function clearData() {
        document.getElementById("monitorBody").innerHTML = "";
    }

</script>
</html>
    `

	monitorText = `
    <html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
    <script src="/bootstrap.min.js" type="text/javascript"></script>
</head>

<body class="container">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid">
            <a class="navbar-brand" href="">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-align-top" viewBox="0 0 16 16">
                    <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)" />
                    <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z" />
                </svg>
            </a>
            <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">

                <div class="navbar-nav">
                    <a class="nav-link" href='/init'>账号管理</a>
                    <a class="nav-link" href='/sysvar'>集群环境</a>
                    <a class="nav-link" href='/data'>数据监控</a>
                    <a class="nav-link active" href='/monitor'>性能监控</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>登录</a>
                    <a class="nav-link" href="/lang?lang=en">[EN]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>性能数据监控</h3>
            <div class="input-group">
                <span class="input-group-text">监控时间间隔(单位:秒)</span>
                <input id="stime" placeholder="输入时间间隔" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">开始</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">停止</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">清除数据</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>内存使用(MB)</th>
                <th>内存已分配(MB)</th>
                <th>内存释放次数</th>
                <th>业务任务数</th>
                <th>集群任务数</th>
                <th>协程数</th>
                <th>CPU核数</th>
                <th>磁盘剩余(GB)</th>
                <th>内存使用率</th>
                <th>CPU使用率</th>
            </tr>
            <tbody id="monitorBody">
            </tbody>
        </table>
    </div>
</body>
<script type="text/javascript">
    var pro = window.location.protocol;
    var wspro = "ws:";
    if (pro === "https:") {
        wspro = "wss:";
    }
    var wsmnt = null;
    var id = 1;
    function WS() {
        this.ws = null;
    }

    WS.prototype.monitor = function () {
        let obj = this;
        this.ws = new WebSocket(wspro + "//" + window.location.host + "/monitorData");
        this.ws.onopen = function (evt) {
            obj.ws.send(document.getElementById("stime").value);
        }
        this.ws.onmessage = function (evt) {
            if (evt.data != "") {
                var json = JSON.parse(evt.data);
                var tr = document.createElement('tr');
                var d = '<td style="font-weight: bold;">' + id++ + '</td>'
                    + '<td>' + Math.round(json.Alloc / (1 << 20)) + '</td>'
                    + '<td>' + Math.round(json.TotalAlloc / (1 << 20)) + '</td>'
                    + '<td>' + json.NumGC + '</td>'
                    + '<td>' + json.NumTx + '</td>'
                    + '<td>' + json.CluserLoad + '</td>'
                    + '<td>' + json.NumGoroutine + '</td>'
                    + '<td>' + json.NumCPU + '</td>'
                    + '<td>' + json.DiskFree + '</td>'
                    + '<td>' + Math.round(json.RamUsage * 10000) / 100 + '%</td>'
                    + '<td>' + Math.round(json.CpuUsage * 100) / 100 + '%</td>';
                tr.innerHTML = d;
                document.getElementById("monitorBody").appendChild(tr);
            }
        }
    }

    WS.prototype.close = function () {
        this.ws.close();
    }

    function monitorLoad() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
        wsmnt = new (WS);
        wsmnt.monitor();
    }

    function stop() {
        if (typeof wsmnt != "undefined" && wsmnt != null && wsmnt != "") {
            wsmnt.close();
        }
    }

    function clearData() {
        document.getElementById("monitorBody").innerHTML = "";
    }

</script>
</html>
    `
)
