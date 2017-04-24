<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="static/css/bootstrap.min.css" rel="stylesheet">
    <link href="static/css/font-awesome.min.css" rel="stylesheet">
    <link href="static/css/style_1.css" rel="stylesheet">
    <script src="static/js/jquery.min.js"></script>
    <script src="static/js/bootstrap.js"></script>
    <title>index_manager_1</title>
</head>
<body>
<section>
{{template "login/ep_nav_left.tpl" .}}
<div class="main-content">
{{template "product/pro_nav_head.tpl" .}}
{{template "product/product_add.tpl" .}}
{{template "login/foot.tpl" .}}
</div>
</section>
</body>
</html>
