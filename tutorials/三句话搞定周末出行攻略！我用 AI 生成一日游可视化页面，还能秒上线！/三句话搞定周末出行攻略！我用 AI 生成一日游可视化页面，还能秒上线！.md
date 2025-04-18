
前言
==

随着 `MCP` 的快速普及和广泛应用，越来越多的工具和服务开始推出自身的 `MCP Server`。借助这些服务，用户在与 `AI` 对话时可以调用相应功能，从而赋予模型更强的执行能力和更丰富的交互体验。

本文将介绍如何结合 `AI` 与高德地图、`EdgeOne Pages Deploy` 提供的 `MCP Server`，通过三句话即可生成一份可视化的一日游行程页面，并部署到 `EdgeOne` 上以便在线访问。

准备好了吗？准备一杯你最喜欢的咖啡或茶，随着本文一探究竟吧。


![image](https://github.com/user-attachments/assets/75166552-b69a-4879-b3aa-4e007e098325)


效果展示
====

在线访问地址：https://mcp.edgeone.site/share/32wIUaOLr2fnUx3vnB2BU。

![image](https://github.com/user-attachments/assets/ba57b52b-1284-452a-bfe2-cf373da2fd66)



AI 生成一日游可视化页面
=============

📐 设计构思
-------

构建一个对话式 `AI` 一日游攻略网页生成的 "助手"，用户通过自然语言与 `AI` 聊天，即可获取所在城市或指定地点的轻量一日游推荐行程。"助手" 由 `AI` 聊天大模型、高德 `MCP` 和 `EdgeOne Pages Deploy MCP` 共同驱动，在无需用户编写代码或操作复杂工具的情况下，完成从旅游信息查询、行程规划、网页构建到公开部署的全流程体验。

通过与 `AI` 对话，用户可以逐步明确需求，例如 "**我现在在深圳，安排一个轻松的一日游路线，推荐几个小吃店和景点---**"。`AI` 根据语义理解调用高德 `MCP` 提供的位置、景点和美食信息，再综合设计成网页展示所需的 `HTML` 页面结构，最终调用 `EdgeOne Pages Deploy MCP` 进行部署，生成一个公开访问的网页链接，用户可直接使用或转发分享。

⚙️ 实现原理
-------

![image](https://github.com/user-attachments/assets/7ed302cd-1b2f-47c2-b553-7d7ba2be37c6)



- **Cursor AI 自然语言解析与上下文管理**

  用户通过 `Cursor AI` 使用自然语言输入需求，`AI`将自动解析用户位置、偏好（如"一日游"、"情侣"、"想吃当地小吃"等），并在对话中逐步明确用户意图。

- **高德地图 MCP 信息聚合**

  - `AI` 调用高德地图 `MCP` 获取用户当前位置附近的：

    - 人气景点（根据距离、热度筛选）

    - 餐饮店铺（如本地特色小吃）

    - 开放时间等信息

  - `AI` 根据时间预算、路线合理性，编排一条 **一日游动线**

- **HTML 页面生成与部署**

  - `AI` 根据行程规划结果，生成标准化 `HTML` 页面，包含：

    - 行程时间表（上午/中午/下午/傍晚）

    - 景点与餐饮介绍（含高德数据字段）

    - 简洁风格的网页结构与样式

  - 最后调用 `EdgeOne Pages Deploy MCP`，将 `HTML` 源码部署至 `EdgeOne Pages`，并返回公开 `URL`，供用户访问、存档或分享。


前置准备
----

### Cursor 下载

`Cursor` 是一种全新的智能 `IDE`，通过与 `AI` 的无缝集成而提供支持。除了 `Cursor`，你也可以选择其他支持 `MCP` 功能的 `IDE`。

![image](https://github.com/user-attachments/assets/6361e10a-b548-415a-ab36-007453c6a7c9)


### 申请高德地图 API Key

本次实践中将使用高德地图的 `MCP Server`，配置包含 `API Key`，因此需要提前完成 `API Key` 的申请。具体步骤可参考：[高德地图 API Key 申请](https://lbs.amap.com/api/mcp-server/create-project-and-key)。

![image](https://github.com/user-attachments/assets/ce0e0889-767c-48bf-861f-6635ae23fca7)


在 Cursor 中配置 MCP Server
-----------------------

接下来，我将介绍如何在 `Cursor` 中配置高德地图`MCP Server` 和 `EdgeOne Pages Deploy MCP Server`。

### 配置高德地图 MCP Server

`Amap Maps` 是一个支持任何 `MCP` 协议客户端的服务器，用户可以轻松地通过 `Amap Maps MCP` 服务器来使用各种基于位置的服务。它是高德地图官方 `MCP Server` 。

- 首先，在 **Cursor Settiongs-MCP** 里，点击 `Add new global MCP server` 按钮跳转到 `MCP Server` 的 `json` 配置文件。

![image](https://github.com/user-attachments/assets/7a91631f-14a2-409b-90f9-c50f20fadbc7)


- 接下来，将下面的配置复制到 `JSON` 文件里。

```
{
  "mcpServers": {
    "amap-maps": {
    "command": "npx",
    "args": [
      "-y",
      "@amap/amap-maps-mcp-server"
    ],
    "env": {
      "AMAP_MAPS_API_KEY": "您在高德官网上申请的key"
    }
    }
  }
}
```

`AMAP_MAPS_API_KEY` 的值替换成我们刚才申请的 `API Key`。

![image](https://github.com/user-attachments/assets/b2daefa5-5b11-4f10-ae77-8574deaf64b4)


通过以上步骤我们就完成了高德地图的 `MCP Server` 配置。

### 配置 EdgeOne Pages Deploy MCP Server

`EdgeOne Pages` 是基于 `Tencent EdgeOne` 基础设施打造的前端开发和部署平台，专为现代 `Web` 开发设计，帮助开发者快速构建、部署静态站点和无服务器应用。通过集成边缘函数能力，实现高效的内容交付和动态功能扩展，支持全球用户的快速访问。

`EdgeOne Pages Deploy MCP` 是一项专用服务，能够将 `HTML` 内容快速部署到 `EdgeOne Pages` 并生成公开访问链接。这使您能够立即预览和分享 `AI` 生成的网页内容。

配置 `EdgeOne Pages Deploy MCP Server` 和刚才的步骤是一样的**，**以下是对应的 `MCP Server` 配置：

```
{
  "mcpServers": {
    "edgeone-pages-mcp-server": {
      "command": "npx",
      "args": ["edgeone-pages-mcp"]
    }
  }
}
```

#### 技术原理

`Pages MCP Server` 利用无服务器边缘计算能力和 `KV` 存储，通过 `API` 接收 `HTML` 内容，即可自动生成即时生效的公共访问链接，实现秒级静态页面部署并内置错误处理机制。

![image](https://github.com/user-attachments/assets/0a2a6b41-49d5-43e8-95ea-0a277c834b2b)


### MCP Server 配置总览

最终的 `mcp.json` 内容为：

![image](https://github.com/user-attachments/assets/267cdc06-e52c-4f71-a602-7175340029cb)


在 **Cursor Settiongs-MCP** 里可以看到我们所配置的 `MCP Server` 列表信息。

![image](https://github.com/user-attachments/assets/377917aa-7a62-427a-a975-e3d540be8d02)


一日游可视化页面生成
----------

首先，在 `Chat` 对话模块，选择 `Agent` 模块，将下面的 `prompt` 发给 `AI`：

> 我现在在深圳，请安排一个轻松的一日游路线，推荐几个小吃店和景点，你可以调用相关的 amap-maps 工具获取相关信息。

![image](https://github.com/user-attachments/assets/fc451fad-59ea-46e1-a8dc-55a3c9ce810b)


从上面的图片信息可以发现，`Cursor AI` 通过 `maps_text_search` 工具查询到相关的景点信息，然后规划路线，在规划的时候还会调用其他工具去计算区域之间的距离以及车程时长。下面是 `AI` 规划的完整路线：

```
深圳一日游完整路线
上午（9:00-12:00）
深圳湾公园（开放时间：06:00-23:00）
地址：滨海大道（近望海路）
活动：沿海边步行道漫步，欣赏海景，拍照留念
中午（12:00-14:00）
前往海上世界（约30分钟车程）
午餐推荐：金稻园砂锅粥、龙记米粉或蔡澜越南粉
地址：望海路1128号
下午（14:00-17:30）
深圳世界之窗（开放时间：10:00-21:30）
地址：深南大道9037号
门票：200元
活动：参观世界各地著名建筑的微缩景观
傍晚至晚上（18:00-21:00）
返回深圳湾公园（约20分钟车程）
活动：欣赏海湾日落美景
晚餐：在红树林附近的ICON天际餐厅或香港茗星坊茶餐厅享用晚餐
这条路线安排紧凑而不匆忙，让您能够体验深圳的自然景观和人文景点，同时品尝到当地特色美食。全程交通便利，景点之间的距离适中，非常适合轻松的一日游。
```

接下来，我们请求 `AI` 根据路线生成可视化的 `HTML` 页面，下面是 `prompt`：

> 请根据完整路线，生成一个好看的的 HTML 页面，然后将源码发给我，如果我确定没有问题，我会请求部署该页面。

![image](https://github.com/user-attachments/assets/38d8bfa0-039e-4487-b299-ac2bc66213fe)


最后，请求 `AI` 将 `html` 页面部署到 `EdgeOne Pages` 上。

![image](https://github.com/user-attachments/assets/e6c81e86-dbc6-4741-b031-dd973db4651a)


部署成功后，`AI` 将会发送访问链接给我们： https://mcp.edgeone.site/share/32wIUaOLr2fnUx3vnB2BU，下面是页面显示效果：

![image](https://github.com/user-attachments/assets/d419f178-246d-46b7-9d60-928de74767da)


HTML 源码
=======

`AI` 根据一日游规划生成的 `HTML` 源码如下：

```
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>深圳一日游 | 轻松游玩攻略</title>
    <link rel="stylesheet" href="https://cdn.bootcdn.net/ajax/libs/bootstrap/5.3.0/css/bootstrap.min.css">
    <link rel="stylesheet" href="https://cdn.bootcdn.net/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <style>
        :root {
            --primary-color: #3498db;
            --secondary-color: #1abc9c;
            --accent-color: #e74c3c;
            --text-color: #333;
            --light-bg: #f9f9f9;
        }

        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background-color: var(--light-bg);
        }

        .navbar-brand {
            font-weight: 700;
            color: var(--primary-color) !important;
        }

        .hero {
            background-image: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)), url('https://images.unsplash.com/photo-1567133804542-21feb3a61979');
            background-size: cover;
            background-position: center;
            color: white;
            padding: 100px 0;
            position: relative;
        }

        .hero-text {
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
        }

        .section-title {
            position: relative;
            margin-bottom: 30px;
            font-weight: 700;
            color: var(--primary-color);
        }

        .section-title:after {
            content: '';
            display: block;
            width: 50px;
            height: 3px;
            background: var(--secondary-color);
            margin-top: 10px;
        }

        .timeline {
            position: relative;
            padding: 0;
            list-style: none;
        }

        .timeline:before {
            content: "";
            position: absolute;
            top: 0;
            bottom: 0;
            left: 50%;
            width: 2px;
            margin-left: -1px;
            background-color: #e9ecef;
        }

        .timeline > li {
            position: relative;
            margin-bottom: 50px;
        }

        .timeline > li:before,
        .timeline > li:after {
            content: " ";
            display: table;
        }

        .timeline > li:after {
            clear: both;
        }

        .timeline > li .timeline-panel {
            float: left;
            position: relative;
            width: 46%;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 1px 6px rgba(0, 0, 0, 0.1);
            background: white;
        }

        .timeline > li .timeline-panel:before {
            right: -15px;
            border-left-width: 15px;
        }

        .timeline > li .timeline-panel:after {
            right: -14px;
            border-left-width: 14px;
        }

        .timeline > li .timeline-badge {
            color: #fff;
            width: 50px;
            height: 50px;
            line-height: 50px;
            font-size: 1.4em;
            text-align: center;
            position: absolute;
            top: 16px;
            left: 50%;
            margin-left: -25px;
            background-color: var(--primary-color);
            border-radius: 50%;
            z-index: 100;
        }

        .timeline > li.timeline-inverted > .timeline-panel {
            float: right;
        }

        .timeline > li.timeline-inverted > .timeline-panel:before {
            left: -15px;
            right: auto;
            border-right-width: 15px;
            border-left-width: 0;
        }

        .timeline > li.timeline-inverted > .timeline-panel:after {
            left: -14px;
            right: auto;
            border-right-width: 14px;
            border-left-width: 0;
        }

        .timeline-title {
            margin-top: 0;
            color: var(--primary-color);
            font-weight: 600;
        }

        .timeline-body > p,
        .timeline-body > ul {
            margin-bottom: 0;
        }

        .card {
            border: none;
            transition: transform 0.3s ease, box-shadow 0.3s ease;
            margin-bottom: 20px;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
        }

        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 12px 20px rgba(0, 0, 0, 0.15);
        }

        .card-img-top {
            height: 200px;
            object-fit: cover;
        }

        .badge-custom {
            background-color: var(--secondary-color);
            color: white;
            font-weight: 500;
            padding: 5px 10px;
            border-radius: 20px;
        }

        .info-box {
            background-color: white;
            border-radius: 8px;
            padding: 25px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            margin-bottom: 30px;
        }

        .info-box i {
            font-size: 24px;
            color: var(--secondary-color);
            margin-right: 10px;
        }

        .footer {
            background-color: #343a40;
            color: white;
            padding: 40px 0;
        }

        .food-item {
            margin-bottom: 20px;
        }

        .food-name {
            font-weight: 600;
            color: var(--accent-color);
        }

        .btn-primary {
            background-color: var(--primary-color);
            border-color: var(--primary-color);
            padding: 10px 20px;
            font-weight: 600;
            border-radius: 30px;
        }

        .btn-primary:hover {
            background-color: #2980b9;
            border-color: #2980b9;
        }

        @media (max-width: 767px) {
            .timeline:before {
                left: 40px;
            }

            .timeline > li .timeline-panel {
                width: calc(100% - 90px);
                float: right;
            }

            .timeline > li .timeline-badge {
                left: 15px;
                margin-left: 0;
                top: 16px;
            }

            .timeline > li.timeline-inverted > .timeline-panel {
                float: right;
            }

            .timeline > li.timeline-inverted > .timeline-panel:before {
                border-right-width: 15px;
                border-left-width: 0;
            }

            .timeline > li.timeline-inverted > .timeline-panel:after {
                border-right-width: 14px;
                border-left-width: 0;
            }
        }
    </style>
</head>
<body>
    <!-- 导航栏 -->
    <nav class="navbar navbar-expand-lg navbar-light bg-white shadow-sm sticky-top">
        <div class="container">
            <a class="navbar-brand" href="#">深圳一日游</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav ms-auto">
                    <li class="nav-item">
                        <a class="nav-link" href="#itinerary">行程安排</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="#attractions">景点介绍</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="#food">美食推荐</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="#tips">旅行贴士</a>
                    </li>
                </ul>
            </div>
        </div>
    </nav>

    <!-- 头部区域 -->
    <section class="hero">
        <div class="container text-center hero-text">
            <h1 class="display-4 fw-bold">深圳轻松一日游</h1>
            <p class="lead">精选景点 - 特色美食 - 完美路线</p>
            <a href="#itinerary" class="btn btn-primary btn-lg mt-3">查看行程</a>
        </div>
    </section>

    <!-- 行程概览 -->
    <section class="py-5" id="overview">
        <div class="container">
            <div class="row justify-content-center">
                <div class="col-lg-8 text-center">
                    <h2 class="section-title text-center">行程概览</h2>
                    <p class="lead">这条精心设计的路线让您在一天内体验深圳的自然景观和人文景点，同时品尝当地特色美食。全程交通便利，景点之间的距离适中，非常适合轻松的一日游。</p>
                </div>
            </div>

            <div class="row mt-5">
                <div class="col-md-3">
                    <div class="info-box text-center">
                        <i class="fas fa-map-marker-alt mb-3"></i>
                        <h5>3个景点</h5>
                        <p class="text-muted">精选深圳代表性景点</p>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="info-box text-center">
                        <i class="fas fa-utensils mb-3"></i>
                        <h5>6家美食店</h5>
                        <p class="text-muted">当地特色小吃与餐厅</p>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="info-box text-center">
                        <i class="fas fa-clock mb-3"></i>
                        <h5>12小时</h5>
                        <p class="text-muted">轻松舒适的游玩时间</p>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="info-box text-center">
                        <i class="fas fa-car mb-3"></i>
                        <h5>40公里</h5>
                        <p class="text-muted">全程车程总计</p>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- 行程时间线 -->
    <section class="py-5 bg-white" id="itinerary">
        <div class="container">
            <h2 class="section-title text-center">精彩行程</h2>
            <div class="row">
                <div class="col-lg-10 mx-auto">
                    <ul class="timeline">
                        <li>
                            <div class="timeline-badge">
                                <i class="fas fa-sun"></i>
                            </div>
                            <div class="timeline-panel">
                                <div class="timeline-heading">
                                    <h4 class="timeline-title">上午: 深圳湾公园</h4>
                                    <p><small class="text-muted"><i class="fas fa-clock"></i> 9:00 - 12:00</small></p>
                                </div>
                                <div class="timeline-body">
                                    <p>在深圳湾公园开始您的一天。这里有美丽的海景和步行道，可以欣赏深圳湾海岸线的美景，呼吸新鲜空气。</p>
                                    <div class="mt-2">
                                        <span class="badge bg-info me-2">开放时间: 06:00-23:00</span>
                                        <span class="badge bg-secondary">免费</span>
                                    </div>
                                </div>
                            </div>
                        </li>
                        <li class="timeline-inverted">
                            <div class="timeline-badge">
                                <i class="fas fa-utensils"></i>
                            </div>
                            <div class="timeline-panel">
                                <div class="timeline-heading">
                                    <h4 class="timeline-title">午餐: 海上世界美食区</h4>
                                    <p><small class="text-muted"><i class="fas fa-clock"></i> 12:00 - 14:00</small></p>
                                </div>
                                <div class="timeline-body">
                                    <p>在海上世界品尝各式特色美食，推荐尝试金稻园砂锅粥、龙记米粉或蔡澜越南粉。</p>
                                    <div class="mt-2">
                                        <span class="badge bg-warning text-dark me-2">车程: 约30分钟</span>
                                        <span class="badge bg-primary">多种价位选择</span>
                                    </div>
                                </div>
                            </div>
                        </li>
                        <li>
                            <div class="timeline-badge">
                                <i class="fas fa-landmark"></i>
                            </div>
                            <div class="timeline-panel">
                                <div class="timeline-heading">
                                    <h4 class="timeline-title">下午: 深圳世界之窗</h4>
                                    <p><small class="text-muted"><i class="fas fa-clock"></i> 14:00 - 17:30</small></p>
                                </div>
                                <div class="timeline-body">
                                    <p>探索深圳世界之窗，这里汇集了世界各地的著名建筑微缩景观，非常适合休闲观光和拍照。</p>
                                    <div class="mt-2">
                                        <span class="badge bg-info me-2">开放时间: 10:00-21:30</span>
                                        <span class="badge bg-danger">门票: 200元</span>
                                    </div>
                                </div>
                            </div>
                        </li>
                        <li class="timeline-inverted">
                            <div class="timeline-badge">
                                <i class="fas fa-moon"></i>
                            </div>
                            <div class="timeline-panel">
                                <div class="timeline-heading">
                                    <h4 class="timeline-title">傍晚及晚餐: 深圳湾日落与晚餐</h4>
                                    <p><small class="text-muted"><i class="fas fa-clock"></i> 18:00 - 21:00</small></p>
                                </div>
                                <div class="timeline-body">
                                    <p>返回深圳湾公园欣赏海湾日落美景，随后在红树林附近的ICON天际餐厅或香港茗星坊茶餐厅享用晚餐。</p>
                                    <div class="mt-2">
                                        <span class="badge bg-warning text-dark me-2">车程: 约20分钟</span>
                                        <span class="badge bg-success">完美结束</span>
                                    </div>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </section>

    <!-- 景点介绍 -->
    <section class="py-5" id="attractions">
        <div class="container">
            <h2 class="section-title text-center">景点介绍</h2>
            <div class="row">
                <div class="col-md-4">
                    <div class="card">
                        <img src="https://images.unsplash.com/photo-1566041510639-8d95a2490bfb" class="card-img-top" alt="深圳湾公园">
                        <div class="card-body">
                            <h5 class="card-title">深圳湾公园</h5>
                            <p class="card-text">深圳湾公园是一个沿海岸线的城市公园，拥有美丽的海景和步行道，是深圳市民休闲放松的好去处。公园沿线可欣赏深圳湾、香港新界的景色，特别是傍晚时分的日落景色格外迷人。</p>
                            <div class="d-flex justify-content-between align-items-center mt-3">
                                <div>
                                    <span class="badge-custom"><i class="fas fa-map-marker-alt"></i> 滨海大道</span>
                                </div>
                                <small class="text-muted">开放时间: 06:00-23:00</small>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="col-md-4">
                    <div class="card">
                        <img src="https://images.unsplash.com/photo-1527224538127-2104bb71c51b" class="card-img-top" alt="海上世界">
                        <div class="card-body">
                            <h5 class="card-title">海上世界</h5>
                            <p class="card-text">海上世界是深圳蛇口的标志性景点，集文化、娱乐、餐饮、购物于一体。这里有多个主题区域，包括游艇中心、特色商店街和美食广场，既可以感受海洋文化，又能品尝各种美食。</p>
                            <div class="d-flex justify-content-between align-items-center mt-3">
                                <div>
                                    <span class="badge-custom"><i class="fas fa-map-marker-alt"></i> 望海路1128号</span>
                                </div>
                                <small class="text-muted">全天开放</small>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="col-md-4">
                    <div class="card">
                        <img src="https://images.unsplash.com/photo-1494059980473-813e73ee784b" class="card-img-top" alt="深圳世界之窗">
                        <div class="card-body">
                            <h5 class="card-title">深圳世界之窗</h5>
                            <p class="card-text">深圳世界之窗是一个汇集了世界著名建筑和景观微缩复制品的主题公园。公园内按不同国家和地区分区，让您足不出户就能领略世界各地的文化特色和建筑风情。</p>
                            <div class="d-flex justify-content-between align-items-center mt-3">
                                <div>
                                    <span class="badge-custom"><i class="fas fa-map-marker-alt"></i> 深南大道9037号</span>
                                </div>
                                <small class="text-muted">门票: 200元</small>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- 美食推荐 -->
    <section class="py-5 bg-white" id="food">
        <div class="container">
            <h2 class="section-title text-center">美食推荐</h2>
            <div class="row">
                <div class="col-lg-6">
                    <div class="info-box">
                        <h4>午餐推荐 <i class="fas fa-utensils float-end text-primary"></i></h4>
                        <hr>
                        <div class="food-item">
                            <p class="food-name">金稻园砂锅粥</p>
                            <p>位于海上世界C区101，这里的砂锅粥香浓可口，食材新鲜，是当地人气美食。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 望海路海上世界C区101</p>
                        </div>
                        <div class="food-item">
                            <p class="food-name">龙记米粉</p>
                            <p>位于海上世界地铁站A口附近，是深圳很有名的特色米粉，口感劲道，汤底鲜美。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 海上世界地铁站A口步行90米</p>
                        </div>
                        <div class="food-item">
                            <p class="food-name">蔡澜越南粉</p>
                            <p>位于海上世界船前广场，由著名美食家蔡澜推荐，提供正宗的越南风味米粉。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 招商街道水湾社区海上世界船前广场2栋103</p>
                        </div>
                    </div>
                </div>

                <div class="col-lg-6">
                    <div class="info-box">
                        <h4>晚餐推荐 <i class="fas fa-moon float-end text-primary"></i></h4>
                        <hr>
                        <div class="food-item">
                            <p class="food-name">ICON天际餐厅</p>
                            <p>位于尚美-红树湾1号，可以欣赏到海湾美景，提供精致的中西融合料理，环境优雅。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 尚美-红树湾1号</p>
                        </div>
                        <div class="food-item">
                            <p class="food-name">香港茗星坊茶餐厅</p>
                            <p>位于沙嘴中心街，提供正宗的港式美食，如菠萝包、奶茶、鸳鸯饭等经典港式料理。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 沙嘴中心街7号</p>
                        </div>
                        <div class="food-item">
                            <p class="food-name">潮味海鲜餐厅</p>
                            <p>位于福荣路，专注于提供新鲜的海鲜料理，以粤菜和潮汕风味为主，价格适中。</p>
                            <p><i class="fas fa-map-marker-alt text-danger"></i> 地址: 沙头街道福荣路沙嘴村三坊32-33号</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- 旅行贴士 -->
    <section class="py-5" id="tips">
        <div class="container">
            <h2 class="section-title text-center">旅行贴士</h2>
            <div class="row">
                <div class="col-md-6">
                    <div class="info-box">
                        <h5><i class="fas fa-car"></i> 交通出行</h5>
                        <ul class="mt-3">
                            <li>深圳公共交通发达，可以使用深圳通卡或手机支付搭乘公交、地铁</li>
                            <li>也可以使用滴滴出行、高德打车等APP叫车前往各景点</li>
                            <li>世界之窗附近地铁站为"世界之窗站"，十分便利</li>
                            <li>深圳湾公园附近地铁站为"红树湾站"，可步行前往</li>
                            <li>海上世界有专门的"海上世界站"，交通便利</li>
                        </ul>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="info-box">
                        <h5><i class="fas fa-sun"></i> 游玩建议</h5>
                        <ul class="mt-3">
                            <li>深圳天气温暖，建议携带防晒用品和遮阳伞</li>
                            <li>世界之窗面积较大，建议穿着舒适的鞋子</li>
                            <li>深圳湾公园傍晚人较多，建议提前到达找到好位置欣赏日落</li>
                            <li>世界之窗门票可以提前在官网或旅游平台购买，避免排队</li>
                            <li>建议携带充足的饮用水，特别是在夏季游玩时</li>
                            <li>海上世界美食种类丰富，可以少量多尝试不同美食</li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </section>

    <!-- 页脚 -->
    <footer class="footer">
        <div class="container">

            <div class="text-center">
                <p class="mb-0">© AI & 陈明勇 | 版权所有</p>
            </div>
        </div>
    </footer>

    <script src="https://cdn.bootcdn.net/ajax/libs/bootstrap/5.3.0/js/bootstrap.bundle.min.js"></script>
</body>
</html>
```

小结
==

在本次实践中，我们不仅体验了 `AI` 通过 `MCP` 协议连接外部数据服务的能力，还深入探索了 EdgeOne Pages Deploy 在部署环节中的关键作用。作为 **最后一公里** 的部署引擎，`EdgeOne Pages Deploy MCP` 让 `AI` 生成的 `HTML` 页面能够以秒级速度部署至线上，并立即生成公开链接，极大提升了内容的可访问性与传播效率。

相比传统的部署流程，`EdgeOne Pages Deploy` 摒弃了复杂的构建配置与发布流程，借助边缘计算和无服务器架构，实现了真正的"所见即所得，生成即上线"。这不仅提升了 `AI` 应用的落地效率，也为开发者、内容创作者、产品设计师等非工程背景用户打开了快速上线创意的全新通道。
