package main

const indexTemplate = `
<html>

<head>
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta http-equiv="content-type" content="text/html;charset=utf-8" />
<meta http-equiv="refresh" content="60" />
<!--link rel="stylesheet" href="picnic.css" /-->
<style>
/* Picnic CSS v6.4.0 http://picnicss.com/ */
html{font-family:sans-serif;-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%}body{margin:0}article,header,main,nav,section {display:block}a{background:transparent}a:active,a:hover{outline:0}b {font-weight:bold}table{border-collapse:collapse;border-spacing:0}th{padding:0}*{box-sizing:inherit}html,body{font-family:Arial, Helvetica, sans-serif;box-sizing:border-box;height:100%}body{color:#111;font-size:1.1em;line-height:1.5;background:#fff}main{display:block}h3 {margin:0;padding:.6em 0}li{margin:0 0 .3em}a{color:#0074d9;text-decoration:none;box-shadow:none}:checked+.toggle,:checked+.toggle:hover{box-shadow:inset 0 0 0 99em rgba(17,17,17,0.2)}table{text-align:left}th{padding:.3em 2.4em .3em .6em}th{text-align:left;font-weight:900;color:#fff;background-color:#0074d9}tr:nth-child(even){background:rgba(0,0,0,0.05)}.flex{display:-ms-flexbox;display:flex;margin-left:-0.6em;width:calc(100% + .6em);flex-wrap:wrap}.flex>*{box-sizing:border-box;flex:1 1 auto;padding-left:.6em;padding-bottom:.6em}.flex[class*="one"]>*,.flex[class*="two"]>*,.flex[class*="three"]>* {flex-grow:0}.one>*{width:100%}.two>*{width:50%}@media all and (min-width: 800px) {.three-800>*{width:33.33333%}}nav{position:fixed;top:0;left:0;right:0;height:3em;padding:0 .6em;background:#fff;box-shadow:0 0 0.2em rgba(17,17,17,0.2);z-index:10000;transform-style:preserve-3d}nav .brand {float:right;position:relative;top:50%;-webkit-transform:translateY(-50%);transform:translateY(-50%)}nav .brand{font-weight:700;float:left;padding:0 .6em;max-width:50%;white-space:nowrap;color:#111}.card {position:relative;box-shadow:0;border-radius:.2em;border:1px solid #ccc;overflow:hidden;text-align:left;background:#fff;margin-bottom:.6em;padding:0}:checked+.card {font-size:0;padding:0;margin:0;border:0}.card>* {max-width:100%;display:block}.card>*:last-child {margin-bottom:0}.card header,.card section {padding:.6em .8em}.card section {padding:.6em .8em 0}.card header {font-weight:bold;position:relative;border-bottom:1px solid #eee}.card header:last-child {border-bottom:0}
</style>
<title>NetKotH Scorebot</title>
<style>
main {
	margin: 3em auto 0;
	padding: 10px;
}
header.crashed {
	background: #ff4136;
	color: white;
}
header.crashed span {
	float: right;
}
header.offline {
	background: #ff4136;
	color: white;
}
header.offline span {
	float: right;
}
header.online {
	background: #2ecc40;
	color: white;
}
header.online span {
	float: right;
}
</style>
</head>
<body>

<nav><a class="brand" href="/">NetKotH Scoreboard</a></nav>

<main>
	<h3>Scores</h3>
	<div class="flex">
		<div>
			<table style="width: 100%">
				<tr><th>Team</th><th style="text-align: center">Score</th></tr>
{{range .Teams}}
				<tr><td>{{.Name}}</td><td style="text-align: center">{{.Score}}</td></tr>
{{end}}
			</table>
		</div>
	</div>
	<h3>Targets</h3>
	<div class="flex one three-800">
{{range .Victims}}
		<div>
			<article class="card">
				<header class="{{.State}}">{{.IP}}<span>{{.State}}</span></header>
				<section>
					<div class="flex two">
						<div><b>Homepage</b></div>
						<div><a href="http://{{.IP}}" target="_blank">http://{{.IP}}</a></div>
					</div>
					<div class="flex two">
						<div><b>Last Contact</b></div>
						<div>{{.LastSeenRel}} ago</div>
					</div>
					<div class="flex two">
						<div><b>Difficulty</b></div>
						<div>{{if .VM.Meta.difficulty}}{{.VM.Meta.difficulty}}{{else}}unspecified{{end}}</div>
					</div>
					<div class="flex two">
						<div><b>Controlling team</b></div>
						<div>{{.Controller}}</div>
					</div>
				</section>
			</article>
		</div>
{{- end}}
	</div>
	<h3>Rules</h3>
	<div class="flex">
		<div>
			These are the rules:
			<ol>
				<li>Only attack the targets above, not the score board, not the infrastructure, not other players.</li>
				<li>Anything goes on the target systems.</li>
				<li>Plant your flag on the homepage of the target (follow link) to score points.
				<li>Your flag should be your team name inside <tt>&lt;team&gt;$TEAM&lt;/team&gt;</tt> html tags.</li>
				<li>The scoring engine simply retrieves the content of the homepage of the target like a browser. You must prove you can deface this homepage.</li>
			</ol>
		</div>
	</div>
</main>

</body>
</html>
`
