1. ContainerList接口从docker daemon中可以获取所有 容器信息(包括容器ID)
2. ContainerTop接口根据容器ID获取容器内运行的所有进程ID(系统进程ID,非容器内虚拟PID)
3. 根据PID获取进程CPU/内存等信息
    若容器内只有一个运行进程,则可认为容器网络流量即改进程流量
4. ContainerInspect 接口获取容器信息,包括k8s中设置的label(podname containerName)
5. Events接口可以向docker daemon注册事件监听，当有事件发生时通知(容器创建/启动/退出)
    容器退出时有退出码,可以知道进程执行结果

6. NetworkList接口可以列出所有的网络设备ID NetworkID,每个网络设备ID下可以关联多个容器
7. NetworkInspect接口可以查询该设备下关联了哪些容器及容器的EndpointID



启动docker事件触发一下
msg:{"status":"create","id":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","from":"af5b7f010a65","Type":"container","Action":"create","Actor":{"ID":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","Attributes":{"describe":"centos with common command","image":"af5b7f010a65","name":"goofy_keller","org.label-schema.schema-version":"= 1.0     org.label-schema.name=CentOS Base Image     org.label-schema.vendor=CentOS     org.label-schema.license=GPLv2     org.label-schema.build-date=20180531","version":"7.5.1804"}},"scope":"local","time":1531214697,"timeNano":1531214697636579099}
msg:{"status":"attach","id":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","from":"af5b7f010a65","Type":"container","Action":"attach","Actor":{"ID":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","Attributes":{"describe":"centos with common command","image":"af5b7f010a65","name":"goofy_keller","org.label-schema.schema-version":"= 1.0     org.label-schema.name=CentOS Base Image     org.label-schema.vendor=CentOS     org.label-schema.license=GPLv2     org.label-schema.build-date=20180531","version":"7.5.1804"}},"scope":"local","time":1531214697,"timeNano":1531214697638127047}
msg:{"Type":"network","Action":"connect","Actor":{"ID":"e5f8ab4a46c653cdcdd2327bbefb1eadd179b6bba7cac6d98d73058aaee7d373","Attributes":{"container":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","name":"bridge","type":"bridge"}},"scope":"local","time":1531214697,"timeNano":1531214697708874745}
msg:{"status":"start","id":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","from":"af5b7f010a65","Type":"container","Action":"start","Actor":{"ID":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","Attributes":{"describe":"centos with common command","image":"af5b7f010a65","name":"goofy_keller","org.label-schema.schema-version":"= 1.0     org.label-schema.name=CentOS Base Image   org.label-schema.vendor=CentOS     org.label-schema.license=GPLv2     org.label-schema.build-date=20180531","version":"7.5.1804"}},"scope":"local","time":1531214698,"timeNano":1531214698301679012}
msg:{"status":"resize","id":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","from":"af5b7f010a65","Type":"container","Action":"resize","Actor":{"ID":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","Attributes":{"describe":"centos with common command","height":"39","image":"af5b7f010a65","name":"goofy_keller","org.label-schema.schema-version":"= 1.0     org.label-schema.name=CentOS Base Image     org.label-schema.vendor=CentOS     org.label-schema.license=GPLv2     org.label-schema.build-date=20180531","version":"7.5.1804","width":"209"}},"scope":"local","time":1531214698,"timeNano":1531214698325853274}
退出docker事件
msg:{"status":"die","id":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","from":"af5b7f010a65","Type":"container","Action":"die","Actor":{"ID":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","Attributes":{"describe":"centos with common command","exitCode":"0","image":"af5b7f010a65","name":"goofy_keller","org.label-schema.schema-version":"= 1.0     org.label-schema.name=CentOS Base Image     org.label-schema.vendor=CentOS     org.label-schema.license=GPLv2     org.label-schema.build-date=20180531","version":"7.5.1804"}},"scope":"local","time":1531214901,"timeNano":1531214901165904922}
msg:{"Type":"network","Action":"disconnect","Actor":{"ID":"e5f8ab4a46c653cdcdd2327bbefb1eadd179b6bba7cac6d98d73058aaee7d373","Attributes":{"container":"f40288ebba3b4a64eb850c1880605e7307158a9018c539716248e68d264b9556","name":"bridge","type":"bridge"}},"scope":"local","time":1531214901,"timeNano":1531214901258431022}

监控项
节点信息包括：hostname/cpu个数/memory总量
容器信息：一系列label

1. 节点整体CPU使用率 label: 节点信息  值:使用率
2. 节点内存整体使用情况 label: free/avaliable/节点信息 值: total-avaliable
3. PID CPU使用率 label: pid/cmdline/节点信息/容器信息  值: 使用率(多少个CPU)
4. PID 内存使用情况 label：pid/cmdline/节点信息/容器信息 值:内存使用量
5. 容器内PID数量 label：节点信息/容器信息 值：进程数量

6. 容器状态  label: 节点信息/容器信息/启动时间/状态码 值：1/0

测试：
每秒钟请求一次，运行15个小时，内存信息如下
VmSize:	  240104 kB
VmRSS:	   12504 kB