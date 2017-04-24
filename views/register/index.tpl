<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="static/css/style_1.css" rel="stylesheet">
    <script src="static/js/bootstrap.js"></script>
    <script src="static/js/jquery.min.js"></script>
    <script src="static/js/validator.js"></script>
    <title>index_regist_register-1</title>
</head>
<body>
<div class="container form-reginster">
<form data-toggle="validator" role="form" id="register-pane" action="/register" method="post">
    <div class="form-signin-heading text-center">
    <h1 class="reginter-title">注册</h1>
    <img src="static/img/logo.png" style="width:120px;"/></div>
        <div class="form-group login-wrap">
            <label for="Username" class="col-sm-2 control-label">用户名</label>
            <div class="col-sm-10">
                <input type="text" name="username" class="form-control" pattern="^[_A-z0-9]{1,}$" id="Username" placeholder="Username" data-error="用户名不能包含特殊字符" required>
				<div class="help-block with-errors"></div>
            </div>
        </div>
        <div class="form-group login-wrap has-feedback">
            <label for="inputEmail" class="col-sm-2 control-label">邮箱</label>
            <div class="col-sm-10">
                <input type="email" class="form-control" id="inputEmail" name="email" placeholder="Email" data-error="输入有效的邮箱'example@test.com'" required>
				<div class="help-block with-errors"></div>
            </div>
        </div>
        <div class="form-group login-wrap">
            <label for="inputTEL" class="col-sm-2 control-label">手机号</label>
            <div class="col-sm-10">
                <input type="tel" name="mobile" pattern="^[0-9]{1,}$" class="form-control" id="inputTEL" placeholder="Phone" required>
            </div>
        </div>
        <div class="form-group login-wrap">
            <label for="inputPassword1" class="col-sm-2 control-label">密码</label>
            <div class="col-sm-10">
                <input type="password" name="password" data-minlength="6" class="form-control" id="inputPassword1" placeholder="Password" data-error="最少6位密码" required>
				<div class="help-block with-errors"></div>
            </div>
        </div>
        <div class="form-group login-wrap">
        <label for="inputPassword2" class="col-sm-2 control-label ">密码确认</label>
        <div class="col-sm-10">
            <input type="password" class="form-control" id="inputPassword2"
                   data-match="#inputPassword1"
                   data-match-error="密码不匹配"
                   placeholder="Comfirm" required>
            <div class="help-block with-errors"></div>
        </div>
        </div>
        <div class="form-group login-wrap">
            <div class="col-sm-offset-2 col-sm-10">
                <button type="submit" class="btn btn-lg btn-login btn-block">Register</button>
            </div>
        </div>
</form>
</div>

</body>
</html>

