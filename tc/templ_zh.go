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
 					<a class="nav-link" href='/dashboard'>性能看板</a>
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
 					<a class="nav-link" href='/dashboard'>性能看板</a>
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
 					<a class="nav-link" href='/dashboard'>性能看板</a>
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
 					<a class="nav-link" href='/dashboard'>性能看板</a>
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

	dashboardText = `<html>
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
                <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)"/>
                <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z"/>
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
                <a class="nav-link" href='/monitor'>性能监控</a>
                <a class="nav-link active" href='/dashboard'>性能看板</a>
            </div>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href='/login'>登录</a>
                <a class="nav-link" href="/lang?lang=en">[EN]</a>
            </div>
        </div>
    </div>
</nav>
<div class="container mt-1 card small">
    <div class="mb-5 p-2 overflow-y-auto">
        <div class="m-1">
            <div class="d-flex justify-content-between align-items-center">
                <h1 class="h4 mb-0">
                    系统资源监控面板
                </h1>
                <div class="d-flex align-items-center">
                    <span id="last-update" class="me-3">最后更新：加载中...</span>
                    <div id="loading-indicator" class="loader d-none"></div>
                </div>
            </div>
        </div>

        <!-- 状态卡片行 -->
        <div class="row row-cols-1 md:row-cols-2 lg:row-cols-4 g-2 mb-1">
            <!-- CPU使用率卡片 -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">CPU使用率</h5>
                                <p class="h3 fw-bold mt-1" id="cpu-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="cpu-progress" class="progress-bar bg-primary" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 内存使用率卡片 -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">内存使用率</h5>
                                <p class="h3 fw-bold mt-1" id="mem-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="mem-progress" class="progress-bar bg-success" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 磁盘使用率卡片 -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div>
                                <h5 class="card-title text-secondary mb-0">磁盘使用率</h5>
                                <p class="h3 fw-bold mt-1" id="disk-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="disk-progress" class="progress-bar bg-warning" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 系统负载卡片 -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">系统负载</h5>
                                <p class="h3 fw-bold mt-1" id="load-1">--</p>
                            </div>
                        </div>
                        <div class="d-flex gap-2 text-sm">
                            <span class="badge bg-info/20 text-info">1分钟: <span id="load-1-min">--</span></span>
                            <span class="badge bg-info/20 text-info">5分钟: <span id="load-5-min">--</span></span>
                            <span class="badge bg-info/20 text-info">15分钟: <span id="load-15-min">--</span></span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- APP 任务卡片 -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-3 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">APP当前任务总数</h5>
                                <p class="h3 fw-bold mt-1" id="sdk-task">--</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <hr>
        <!-- 详细信息行 -->
        <div class="row g-1">
            <!-- 系统详情表格 -->
            <div class="col-12 col-md-6 col-lg-4 col-xl-4">
                <div class="card bg-white rounded-3 shadow-sm">
                    <div class="card-header bg-light">
                        <h5 class="card-title mb-0">系统资源详情</h5>
                    </div>
                    <div class="p-1">
                        <table class="table table-hover">
                            <tbody>
                            <tr>
                                <td class="text-secondary" style="width: 120px;">CPU型号</td>
                                <td id="detail-cpu-model" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">CPU核心数</td>
                                <td id="detail-cpu-cores" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">总内存</td>
                                <td id="detail-mem-total" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">已用内存</td>
                                <td id="detail-mem-used" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">挂载点</td>
                                <td id="disk-mount" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">磁盘总容量</td>
                                <td id="detail-disk-total" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">磁盘可用空间</td>
                                <td id="detail-disk-free" class="fw-bold">加载中...</td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <!-- Go运行时信息 -->
            <div class="col-12 col-md-6 col-lg-4 col-xl-4">
                <div class="card bg-white rounded-3 shadow-sm">
                    <div class="card-header bg-light">
                        <h5 class="card-title mb-0">Go运行时信息</h5>
                    </div>
                    <div class="p-1">
                        <table class="table table-hover">
                            <tbody>
                            <tr>
                                <td class="text-secondary" style="width: 150px;">协程数</td>
                                <td id="go-goroutines" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">线程数</td>
                                <td id="go-threads" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Go版本</td>
                                <td id="go-version" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">当前堆内存</td>
                                <td id="heap_alloc_mb" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">累计内存占用</td>
                                <td id="alloc_bytes" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">堆内存释放</td>
                                <td id="free_bytes" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">GC次数</td>
                                <td id="gc_count_total" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">虚拟内存总量</td>
                                <td id="heap_sys_mb" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">累计分配的对象数</td>
                                <td id="alloc_objects" class="fw-bold">加载中...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">累计释放的对象数</td>
                                <td id="free_objects" class="fw-bold">加载中...</td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

    </div>
</div>
</body>
<script type="text/javascript">

    document.addEventListener('DOMContentLoaded', async function () {
        const formatNumber = (num) => {
            try {
                return num.toFixed(2);
            } catch (e) {
                return 0
            }
        };

        const updateValueWithAnimation = (elementId, value) => {
            const element = document.getElementById(elementId);
            if (element) {
                element.textContent = value;
                element.classList.add('value-pulse');
                setTimeout(() => {
                    element.classList.remove('value-pulse');
                }, 500);
            }
        };

        const fetchSystemData = async () => {
            const loadingIndicator = document.getElementById('loading-indicator');
            loadingIndicator.classList.remove('d-none');
            try {
                const data = await getDashboardData()
                if (data) {
                    updateUI(data);
                    document.getElementById('last-update').textContent = ` + "`" + `最后更新：${data.collect_time}` + "`" + `;
                }
            } catch (error) {
                console.error('获取系统数据失败:', error);
                alert('获取系统数据失败，请稍后重试');
            } finally {
                loadingIndicator.classList.add('d-none');
            }
        };


        async function getDashboardData() {
            try {
                const response = await fetch('/dashboardData');
                if (response.ok) {
                    return await response.json()
                }
            } catch (e) {
                console.error(e)
            }
        }

        const updateUI = (data) => {
            updateValueWithAnimation('cpu-usage', ` + "`" + `${formatNumber(data.cpu.usage_percent)}%` + "`" + `);
            document.getElementById('cpu-progress').style.width = ` + "`" + `${data.cpu.usage_percent}%` + "`" + `;
            document.getElementById('detail-cpu-model').textContent = data.cpu.model_name || '未知';
            document.getElementById('detail-cpu-cores').textContent = data.cpu.logical_cores || '未知';
            document.getElementById('sdk-task').textContent = ` + "`" + `${data.task_number.sdk_task_num}` + "`" + `;
            updateValueWithAnimation('mem-usage', ` + "`" + `${formatNumber(data.memory.used_percent)}%` + "`" + `);
            document.getElementById('mem-progress').style.width = ` + "`" + `${data.memory.used_percent}%` + "`" + `;
            document.getElementById('detail-mem-total').textContent = ` + "`" + `${formatNumber(data.memory.total_gb)} GB` + "`" + `;
            document.getElementById('detail-mem-used').textContent = ` + "`" + `${formatNumber(data.memory.used_gb)} GB` + "`" + `;

            updateValueWithAnimation('disk-usage', ` + "`" + `${formatNumber(data.disk.used_percent)}%` + "`" + `);
            document.getElementById('disk-progress').style.width = ` + "`" + `${data.disk.used_percent}%` + "`" + `;
            document.getElementById('disk-mount').textContent = ` + "`" + `${data.disk.mount_point || '未知'}` + "`" + `;
            document.getElementById('detail-disk-total').textContent = ` + "`" + `${formatNumber(data.disk.total_gb)} GB` + "`" + `;
            document.getElementById('detail-disk-free').textContent = ` + "`" + `${formatNumber(data.disk.free_gb)} GB` + "`" + `;

            updateValueWithAnimation('load-1', formatNumber(data.load_avg.load1));
            updateValueWithAnimation('load-1-min', formatNumber(data.load_avg.load1));
            updateValueWithAnimation('load-5-min', formatNumber(data.load_avg.load5));
            updateValueWithAnimation('load-15-min', formatNumber(data.load_avg.load15));

            updateValueWithAnimation('go-goroutines', data.go_runtime.goroutines || '--');
            updateValueWithAnimation('go-threads', data.go_runtime.threads || '--');
            updateValueWithAnimation('go-version', data.go_runtime.go_version || '未知');

            updateValueWithAnimation('heap_alloc_mb', ` + "`" + `${formatNumber(data.go_runtime.heap_alloc_mb)} MB` + "`" + `);
            updateValueWithAnimation('alloc_bytes', ` + "`" + `${formatNumber(data.go_runtime.alloc_bytes >> 20)} MB` + "`" + `);
            updateValueWithAnimation('free_bytes', ` + "`" + `${formatNumber(data.go_runtime.free_bytes >> 20)} MB` + "`" + `);
            updateValueWithAnimation('gc_count_total', data.go_runtime.gc_count_total || '--');
            updateValueWithAnimation('heap_sys_mb', ` + "`" + `${formatNumber(data.go_runtime.heap_sys_mb)} MB` + "`" + ` || '--');
            updateValueWithAnimation('alloc_objects', ` + "`" + `${data.go_runtime.alloc_objects}` + "`" + ` || '--');
            updateValueWithAnimation('free_objects', ` + "`" + `${data.go_runtime.free_objects}` + "`" + ` || '--')
        };

        fetchSystemData();
        setInterval(fetchSystemData, 10000);
    });

</script>
</html>`
)
