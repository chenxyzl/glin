### 前置环境
+ 安装protoc  
  下载https://github.com/protocolbuffers/protobuf/releases 放入任意可执行的目录

### dart生成环境
+ 安装proto-gen-dart  
  执行命令 dart pub global activate protoc_plugin
+ 开始生成  
  进入app目录执行make dart

### ts生成环境
+ 安装node.js   
  参考微软官方文档[https://learn.microsoft.com/en-us/windows/dev-environment/javascript/nodejs-on-wsl]
+ 安装ts-protoc-gen插件  
  npm install -g ts-protoc-gen
+ 开始生成  
  进入app目录执行make ts
+ 运行时copy
  运行依赖ts-protoc-gen目录下的@protobuf-ts/runtime/build/es2015和types