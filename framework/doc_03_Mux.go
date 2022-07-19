package framework

/*
每一个请求逻辑,都有一个控制器类型的函数与之对应,那么如何快速的找到指定的控制器呢?
路由简而言之就是让web服务器根据路由规则,理解http请求的信息,匹配对应的控制器,在将请求转发给控制器处理

不同的设计方式,路由器可用性有天壤之别

路由重点关注http的请求头Header部分 由三个部分组成：Method、Request-URI 和 HTTP-Version
如:
	Get /home.html HTTP/1.1
	Host:...
	User-Agent:...
	...

Method是http的方法,如GET POST PUT DELETE等
Request-URI 是请求路径,浏览器请求地址中域名以外的部分
HTTP-Version 是 HTTP 的协议版本，目前常见的有 1.0、1.1、2.0

路由使用的是前两者

*/

/*
路由规则的需求:
简单到复杂排序,主要有4点需求:

需求 1：HTTP 方法匹配
因为RESTful风格流行,为了让URI更可读,框架必须要支持多种http方法输入

需求 2：静态路由匹配
静态路由匹配是一个路由的基本功能，指的是路由规则中没有可变参数，即路由规则地址是固定的，与 Request-URI 完全匹配。
net/http包默认的Mux就是用map实现的静态路由(根据key查找)

需求 3：批量通用前缀
因为业务模块的划分，我们会同时为某个业务模块注册一批路由，所以在路由注册过程中，为了路由的可读性，一般习惯统一定义这批路由的通用前缀。
比如 /user/info、/user/login 都是以 /user 开头，很方便使用者了解页面所属模块。

需求 4：动态路由匹配
根据需求2演进而来,因为uri中的某个字段不一定是固定的,希望路由也能支持这个规则,将动态变化的uri也能匹配出来
/user/:id

*/

/*
需求1,2:
	因为有两个待匹配条件(方法,和路径),自然想到两级hash表,第一级匹配方法,第二级匹配uri
	关键点:
	-匹配处理全部转化为大写,对调用者就是大小写不敏感的路由

需求3:
	一个Group方法,将前缀相同的路径归拢,应该返回一个实例,并且都具有GET POST等方法
	关键点:
	-group封装core的方法,但是 使用接口来解耦,core.Group返回的不是实例,而是一个接口(约定),从而不依赖具体的实现
	如果你觉得这个模块是完整的，而且后续希望有扩展的可能性，那么就应该尽量使用接口来替代实现

需求4:
	如果要支持动态路由,那么之前的哈希规则就会失效,因为有通配符
因为有通配符，在匹配 Request-URI 的时候，请求 URI 的某个字符或者某些字符是动态变化的，无法使用 URI 做为 key 来匹配。
那么，我们就需要其他的算法来支持路由匹配。

	这个问题本质就是一个字符串匹配, 字符串匹配问题比较通用且高效的方法就是字典树,也叫trie树

前缀树:
	-是多叉的树形结构
	-根节点通常为空字符串
	-键值往往在叶子节点或部分内部节点(路径上的前缀组成)

我们可以按照uri的每个段（segment）来切分,每个段在trie树中都能找到对应的节点,每个节点保存一个段
每个叶子节点代表一个URI,有的中间节点也能代表一个uri(是否有对应处理器是另一回事)

因此实现需求4的关键就是实现前缀树
	1.定义树和节点的数据结构
	2.编写函数 增加路由规则
	3.编写函数 查找路由
	4.将以上逻辑添加到框架


*/
