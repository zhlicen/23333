<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="static/css/bootstrap.min.css" rel="stylesheet">
    <link href="static/css/style_1.css" rel="stylesheet">
    <title>Customer_login</title>
</head>

<body class="login-body">
<div class="container">
    <form class="form-signin" id="login-form" method="post" action="/login">
        <div class="form-signin-heading text-center">
            <h1 class="sign-title">用户登录</h1>
            <img src="static/img/logo.png" style="width:120px;"/></div>
        <div class="login-wrap">
            <input type="text" class="form-control" name="username" placeholder="用户名" autofocus>
            <input type="password" class="form-control" name="password" placeholder="密码">
            <button class="btn btn-lg btn-login btn-block" type="submit"> 登录</button>
        </div>
    </form>
</div>
</body>
</html>
