package util

func JsRedirectHtmlResp(redirect string) string {
	return `<!DOCTYPE html>
<html>
<head>
  <title></title>
</head>
<body>
<script>
	var redir = "` + redirect + `"; try { window.top.location = redir } catch (k) { window.location = redir };
</script>
</body>
</html>`
}

func GAJsRedirectHtmlResp(ua, redirect string) string {

	return `<!DOCTYPE html>
<html>
<head>
<title></title>
</head>
<body>
<script>
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','//www.google-analytics.com/analytics.js','ga');
ga('create', '` + ua + `', { 'userId': window.user && window.user.id ? window.user.id : 'auto' });
ga('send', 'pageview');
setTimeout(function(){
var redir = "` + redirect + `"; try { window.top.location = redir } catch (k) { window.location = redir };
}, 200);
</script>
</body>`
}
