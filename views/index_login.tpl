    <script src="static/js/bootstrap.js"></script>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="static/css/bootstrap.min.css" rel="stylesheet">
    <link href="static/css/font-awesome.min.css" rel="stylesheet">
    <link href="static/css/style.css" rel="stylesheet">
    <script src="static/js/jquery.min.js"></script>
    <script src="static/js/bootstrap.js"></script>
    <title>index_vadication-2</title>
</head>
<body id="myPage" data-spy="scroll" data-target=".navbar" data-offset="60">  
<!--导航-->
<nav class="navbar navbar-default navbar-fixed-top">
    <div class="container-fluid">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#myNavbar">
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="#myPage">
                <img src="static/img/logo.png" style="width:70px;height:35px">
            </a>
        </div>
        <div class="collapse navbar-collapse" id="myNavbar">
            <ul class="nav navbar-nav">
                <li><a href="#About">About</a></li>
                <li><a href="#Service">Service</a></li>
                <li><a href="#Products">Products</a></li>
                <li><a href="#Contact">Contact</a></li>
                <li class="dropdown">
                    <a class="dropdown-toggle" data-toggle="dropdown" href="#">更多
                        <span class="caret"></span></a>
                    <ul class="dropdown-menu">
                        <li><a href="#">Page-1</a></li>
                        <li><a herf="#">Page-2</a></li>
                        <li><a href="#">Page-3</a></li>
                    </ul>
                </li>
            </ul>
            <ul class="nav navbar-nav navbar-right notification-menu">
                <li>
                    <form class="form-inline searchform" action="" method="get">
                        <input type="text" class="col-sm-10 form-control" name="keyword" placeholder="keyword" value="">
                        <button type="submit" class=" col-sm-2 btn btn-default">
                            <i class="glyphicon glyphicon-search" title="Search"></i>
                        </button>
                    </form>
                </li>
                <li><a href="javascript:;" class="btn btn-default dropdown-toggle pull-right" data-toggle="dropdown"><img src="static/img/logo.png" style="width:70px;height:35px" alt="3j"/> 3j <span class="caret"></span></a>
                    <ul class="dropdown-menu pull-right">
                        <li><a href="/custom"><i class="fa fa-cog"></i> 基本资料</a></li>
                        <li><a href="/login"><i class="fa fa-sign-out"></i> 退出</a></li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</nav>
{{template "pro_head_pic.tpl" .}}
{{template "pro_info.tpl" .}}
{{template "pro_foot.tpl" .}}
</body>
</html>

