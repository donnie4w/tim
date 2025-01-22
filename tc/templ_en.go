// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package tc

const (
	loginEnText = `
<html>
<head>
    <title>tim</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/bootstrap.css" rel="stylesheet">
</head>

<body class="container">
    <div class="container-fluid text-right">
        <span>
            <h4 style="display: inline;">Tim Management Platform</h4>
        </span>
        <span style="text-align:right">
            <h6 style="display: inline;">&nbsp;&nbsp;&nbsp;<a href="/lang?lang=en">[中文]</a></h6>
        </span>
        <hr>
        <div id="login">
            <h5>Login</h5>
            <form class="form-control" id="loginform" action="/login" method="post">
                <input name="type" value="1" hidden />
                <input name="name" placeholder="username" />
                <input name="pwd" placeholder="password" type="password" />
                <input type="submit" class="btn btn-primary" value="Login" />
            </form>
        </div>
        <hr>
    </div>
</html>  
    `

	initEnText = `
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
                    <a class="nav-link active" href='/init'>Account</a>
                    <a class="nav-link" href='/sysvar'>Cluster</a>
                    <a class="nav-link" href='/data'>DataMonitor</a>
                    <a class="nav-link" href='/monitor'>Performance</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Login</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: small;">
        {{if .ShowCreate }}
        <div class="container-fluid card mt-1 p-1">
            </h6>
            <form class="form-control" id="createAdminform" action="/init?type=1" method="post">
                <h6>Create Administrator <h6 class="important">{{ .Show }}</h6>
                    <input name="adminName" placeholder="username" />
                    <input name="adminPwd" placeholder="password" type="password" />
                    Administrator<input name="adminType" type="radio" value="1" checked />
                    {{if not .Init}}
                    Observer<input name="adminType" type="radio" value="2" />
                    {{end}}
                    <input type="submit" class="btn btn-primary" value="Create" />
            </form>
        </div>
        {{end}}
        {{if not .Init}}
        <div class="container-fluid card mt-1 p-1">
            <div class="m-2">
                <h6>Manage Accounts</h6>
                {{range $k,$v := .AdminUser}}
                <form class="form-control" id="adminform" action="/init?type=2" method="post">
                    <input name="adminName" value='{{ $k }}' readonly style="border:none;" /> Authority:{{ $v }}
                    <input type="button" value="Delete Account" class="btn btn-danger"
                        onclick="javascipt:if (confirm('confirm delete?')){this.parentNode.submit();};" />
                </form>
                {{end}}
            </div>
        </div>
        <hr>
        {{end}}
    </div>
</html>
    `

	sysvarEnText = `
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
                    <a class="nav-link" href='/init'>Account</a>
                    <a class="nav-link active" href='/sysvar'>Cluster</a>
                    <a class="nav-link" href='/data'>DataMonitor</a>
                    <a class="nav-link" href='/monitor'>Performance</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Login</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="mt-1" style="font-size: xx-small;">
        <table class="table table-bordered table-Info">
            <tr>
                <td style="width: 150px;">System startup time</td>
                <td colspan="2">{{ .SYS.StartTime }}</td>
            </tr>
            <tr>
                <td>current time</td>
                <td colspan="2">{{ .SYS.Time }}</td> 
            </tr>
            <tr>
                <td>Node UUID</td>
                <td class="text-danger" colspan="2">{{ .SYS.UUID }} [{{ .SYS.CSNUM }}]</td>
            </tr>
            <tr>
                <td>Cluster nodes</td>
                <td class="text-danger" colspan="2">{{ .SYS.ALLUUIDS }}</td>
            </tr>

            <tr>
                <td>Cluster listening address</td>
                <td colspan="2">{{ .SYS.ADDR }}</td>
            </tr>
        </table>
    </div>
</body>
</html>
    `
	dataEnText = `
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
                    <a class="nav-link" href='/init'>Account</a>
                    <a class="nav-link" href='/sysvar'>Cluster</a>
                    <a class="nav-link active" href='/data'>DataMonitor</a>
                    <a class="nav-link" href='/monitor'>Performance</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Login</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>data monitor</h3>
            <div class="input-group">
                <span class="input-group-text">time interval(unit:second)</span>
                <input id="stime" placeholder="time interval" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">start</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">stop</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">data clear</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>Online total</th>
                <th>Input data(Byte)</th>
                <th>Output data(Byte)</th>
                <th>Disconnected node</th>
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

	monitorEnText = `
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
                    <a class="nav-link" href='/init'>Account</a>
                    <a class="nav-link" href='/sysvar'>Cluster</a>
                    <a class="nav-link" href='/data'>DataMonitor</a>
                    <a class="nav-link active" href='/monitor'>Performance</a>
                </div>
                <div class="navbar-nav ms-auto">
                    <a class="nav-link" href='/login'>Login</a>
                    <a class="nav-link" href="/lang?lang=zh">[中文]</a>
                </div>
            </div>
        </div>
    </nav>
    <div class="container mt-1 card" style="font-size: small;">
        <div class="container mt-1" style="font-size: small;">
            <h3>performance data monitoring</h3>
            <div class="input-group">
                <span class="input-group-text">time interval(unit:second)</span>
                <input id="stime" placeholder="time interval" value="3" />
                <button class="btn btn-primary" onclick="monitorLoad();">start</button>&nbsp;
                <button class="btn btn-primary" onclick="stop();">stop</button>&nbsp;
                <button class="btn btn-primary" onclick="clearData();">data clear</button>
            </div>
        </div>

        <table class="table table-striped " style="font-size: smaller;">
            <tr>
                <th></th>
                <th>memory usage(MB)</th>
                <th>Memory allocated(MB)</th>
                <th>Memory release times</th>
                <th>Number of business tasks</th>
                <th>Number of cluster tasks</th>
                <th>Coroutine number</th>
                <th>CPU number</th>
                <th>Disk Free Space(GB)</th>
                <th>Memory usage</th>
                <th>CPU usage</th>
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
