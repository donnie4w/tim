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
                    <a class="nav-link" href='/dashboard'>Dashboard</a>
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
                    <a class="nav-link" href='/dashboard'>Dashboard</a>
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
                    <a class="nav-link" href='/dashboard'>Dashboard</a>
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
                    <a class="nav-link" href='/dashboard'>Dashboard</a>
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

	dashboardEnText = `
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
                <rect width="4" height="12" rx="1" transform="matrix(1 0 0 -1 6 15)"/>
                <path d="M1.5 2a.5.5 0 0 1 0-1v1zm13-1a.5.5 0 0 1 0 1V1zm-13 0h13v1h-13V1z"/>
            </svg>
        </a>
        <button type="button" class="navbar-toggler" data-bs-toggle="collapse" data-bs-target="#navbarCollapse">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarCollapse">

            <div class="navbar-nav">
                <a class="nav-link" href='/init'>Account</a>
                <a class="nav-link" href='/sysvar'>Cluster</a>
                <a class="nav-link" href='/data'>Data Monitor</a>
                <a class="nav-link" href='/monitor'>Performance</a>
                <a class="nav-link active" href='/dashboard'>Dashboard</a>
            </div>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href='/login'>Login</a>
                <a class="nav-link" href="/lang?lang=zh">[中文]</a>
            </div>
        </div>
    </div>
</nav>
<div class="container mt-1 card small">
    <div class="mb-5 p-2 overflow-y-auto">
        <div class="m-1">
            <div class="d-flex justify-content-between align-items-center">
                <h1 class="h4 mb-0">
                    System Resource Monitor
                </h1>
                <div class="d-flex align-items-center">
                    <span id="last-update" class="me-3">Last Update: Loading...</span>
                    <div id="loading-indicator" class="loader d-none"></div>
                </div>
            </div>
        </div>

        <!-- Status Cards -->
        <div class="row row-cols-1 md:row-cols-2 lg:row-cols-4 g-2 mb-1">
            <!-- CPU Usage -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">CPU Usage</h5>
                                <p class="h3 fw-bold mt-1" id="cpu-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="cpu-progress" class="progress-bar bg-primary" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Memory Usage -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">Memory Usage</h5>
                                <p class="h3 fw-bold mt-1" id="mem-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="mem-progress" class="progress-bar bg-success" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Disk Usage -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div>
                                <h5 class="card-title text-secondary mb-0">Disk Usage</h5>
                                <p class="h3 fw-bold mt-1" id="disk-usage">--%</p>
                            </div>
                        </div>
                        <div class="progress" style="height: 6px;">
                            <div id="disk-progress" class="progress-bar bg-warning" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- System Load -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-2 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">System Load</h5>
                                <p class="h3 fw-bold mt-1" id="load-1">--</p>
                            </div>
                        </div>
                        <div class="d-flex gap-2 text-sm">
                            <span class="badge bg-info/20 text-info">1Min: <span id="load-1-min">--</span></span>
                            <span class="badge bg-info/20 text-info">5Min: <span id="load-5-min">--</span></span>
                            <span class="badge bg-info/20 text-info">15Min: <span id="load-15-min">--</span></span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- App Tasks -->
            <div class="col-sm-6 col-md-4 col-lg-3 col-xl-3 col-12 mb-1">
                <div class="bg-white rounded-3 shadow-sm h-70">
                    <div class="card-body p-2">
                        <div class="d-flex justify-content-between align-items-start mb-1">
                            <div>
                                <h5 class="card-title text-secondary mb-0">App Tasks Total</h5>
                                <p class="h3 fw-bold mt-1" id="sdk-task">--</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <hr>
        <!-- Details -->
        <div class="row g-1">
            <!-- System Info Table -->
            <div class="col-12 col-md-6 col-lg-4 col-xl-4">
                <div class="card bg-white rounded-3 shadow-sm">
                    <div class="card-header bg-light">
                        <h5 class="card-title mb-0">System Details</h5>
                    </div>
                    <div class="p-1">
                        <table class="table table-hover">
                            <tbody>
                            <tr>
                                <td class="text-secondary" style="width: 120px;">CPU Model</td>
                                <td id="detail-cpu-model" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">CPU Cores</td>
                                <td id="detail-cpu-cores" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Total Memory</td>
                                <td id="detail-mem-total" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Used Memory</td>
                                <td id="detail-mem-used" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Mount Point</td>
                                <td id="disk-mount" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Total Disk</td>
                                <td id="detail-disk-total" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Free Disk</td>
                                <td id="detail-disk-free" class="fw-bold">Loading...</td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <!-- Go Runtime Info -->
            <div class="col-12 col-md-6 col-lg-4 col-xl-4">
                <div class="card bg-white rounded-3 shadow-sm">
                    <div class="card-header bg-light">
                        <h5 class="card-title mb-0">Go Runtime</h5>
                    </div>
                    <div class="p-1">
                        <table class="table table-hover">
                            <tbody>
                            <tr>
                                <td class="text-secondary" style="width: 150px;">Goroutines</td>
                                <td id="go-goroutines" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Threads</td>
                                <td id="go-threads" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Go Version</td>
                                <td id="go-version" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Heap Alloc</td>
                                <td id="heap_alloc_mb" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Total Alloc</td>
                                <td id="alloc_bytes" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Heap Free</td>
                                <td id="free_bytes" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">GC Count</td>
                                <td id="gc_count_total" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Virtual Memory</td>
                                <td id="heap_sys_mb" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Total Alloc Objects</td>
                                <td id="alloc_objects" class="fw-bold">Loading...</td>
                            </tr>
                            <tr>
                                <td class="text-secondary">Total Free Objects</td>
                                <td id="free_objects" class="fw-bold">Loading...</td>
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
                    document.getElementById('last-update').textContent = ` + "`" + `Last Update: ${data.collect_time}` + "`" + `;
                }
            } catch (error) {
                console.error('Failed to get system data:', error);
                alert('Failed to get system data, please retry later');
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
            document.getElementById('detail-cpu-model').textContent = data.cpu.model_name || 'Unknown';
            document.getElementById('detail-cpu-cores').textContent = data.cpu.logical_cores || 'Unknown';
            document.getElementById('sdk-task').textContent = ` + "`" + `${data.task_number.sdk_task_num}` + "`" + `;
            updateValueWithAnimation('mem-usage', ` + "`" + `${formatNumber(data.memory.used_percent)}%` + "`" + `);
            document.getElementById('mem-progress').style.width = ` + "`" + `${data.memory.used_percent}%` + "`" + `;
            document.getElementById('detail-mem-total').textContent = ` + "`" + `${formatNumber(data.memory.total_gb)} GB` + "`" + `;
            document.getElementById('detail-mem-used').textContent = ` + "`" + `${formatNumber(data.memory.used_gb)} GB` + "`" + `;

            updateValueWithAnimation('disk-usage', ` + "`" + `${formatNumber(data.disk.used_percent)}%` + "`" + `);
            document.getElementById('disk-progress').style.width = ` + "`" + `${data.disk.used_percent}%` + "`" + `;
            document.getElementById('disk-mount').textContent = ` + "`" + `${data.disk.mount_point || 'Unknown'}` + "`" + `;
            document.getElementById('detail-disk-total').textContent = ` + "`" + `${formatNumber(data.disk.total_gb)} GB` + "`" + `;
            document.getElementById('detail-disk-free').textContent = ` + "`" + `${formatNumber(data.disk.free_gb)} GB` + "`" + `;

            updateValueWithAnimation('load-1', formatNumber(data.load_avg.load1));
            updateValueWithAnimation('load-1-min', formatNumber(data.load_avg.load1));
            updateValueWithAnimation('load-5-min', formatNumber(data.load_avg.load5));
            updateValueWithAnimation('load-15-min', formatNumber(data.load_avg.load15));

            updateValueWithAnimation('go-goroutines', data.go_runtime.goroutines || '--');
            updateValueWithAnimation('go-threads', data.go_runtime.threads || '--');
            updateValueWithAnimation('go-version', data.go_runtime.go_version || 'Unknown');

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
</html>
`
)
