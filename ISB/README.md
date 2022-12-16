# taskSwitching

## External Library
[qlibc](https://github.com/wolkykim/qlibc.git)
[czmq](http://czmq.zeromq.org/)

## 启动命令示例

启动全局状态存储
```
./ISB_global tcp://0.0.0.0:7053 tcp://0.0.0.0:7054

参数说明:
char* station_state_recieve_endpoint = argv[1];
char* coordination_endpoint = argv[2];

日志地址：
/root/taskswitching/log/global_task_switch.log 代码中写死的
```

启动信息高铁入口
```
./ISB_entry tcp://0.0.0.0:7055 tcp://0.0.0.0:7056 tcp://10.208.104.3:7054 tcp://10.208.104.3:7055
参数说明:
char* entry_task_reciever_endpoint = argv[1];       //任务接收地址
char* entry_worker_reciever_endpoint = argv[2];     //任务执行环境接收地址
char* global_coordination_beijing = argv[3];        //全局状态存储地址
char* entry_name = argv[4];                         //入口名字
日志文件地址：
/root/taskswitching/log/local_task_switch.log 代码中写死的
```

启动站点
```
./ISB_station tcp://0.0.0.0:7057 tcp://0.0.0.0:7058 tcp://0.0.0.0:7059 tcp://159.226.41.229:7057 >& stdout_STS

参数说明:
char* station_task_reciever_endpoint = argv[1];     //任务接收地址
char* station_worker_handler_endpoint = argv[2];    //任务执行环境处理地址
char* station_machines_states_endpoint = argv[3];   //执行机上报状态地址
char* station_name = argv[4];                       //站点名字

日志文件地址
/root/taskswitching/log/station.log 代码中写死的
```

启动执行机
```
./ISB_exec_machine tcp://10.208.104.3:7059 tcp://10.208.104.3:7058 tcp://0.0.0.0:7060 tcp://0.0.0.0:7061 tcp://10.208.104.3:7060 tcp://10.208.104.3:7061 BJ_M1 20 40
参数说明:
char* station_endpoint = argv[1];                   //站点的执行机上报状态地址
char* station_worker_request_endpoint = argv[2];    //向站点请求任务执行环境地址
char* task_reciever_endpoint = argv[3];             //任务接收地址
char* worker_reciever_endpoint = argv[4];           //任务执行环境接收地址
char* machine_name = argv[5];                       //执行机名字
char* worker_endpoint_for_report = argv[6];         //上报给站点的worker接收地址
char* machine_id = argv[7];                         //执行机编号
char* cpu = argv[8];                                //允许使用CPU数（个）
char* mem = argv[9];                                //允许使用内存量（GB）

日志文件地址：
("/root/taskswitching/log/machine_%s.log", machine_id) 代码中写死
```

启动信息高铁的日志服务
```
./logserverd -addr 0.0.0.0:7065 -f local_task_switch.log -f station.log -f machine_BJ_M1.log
源码是logserverd.go
```

## net_config.json 说明
```
config_beijing_ips 北京机器使用的内外网ip

config_nanjing_ips 南京机器使用的内外网ip

beijing_station_endpoints 北京站点内网信息，分别对应name, task_reciever, worker reciever, machine state reciever

beijing_station_endpoints_out 北京站点外网信息，分别对应name, task_reciever, worker reciever, machine state reciever

nanjing_station_endpoints 北京站点内网信息，分别对应name, task_reciever, worker reciever, machine state reciever

nanjing_station_endpoints_out 南京站点外网信息，分别对应name, task_reciever, worker reciever, machine state reciever

global_coordination_beijing 北京全局状态存储地址（获取站点状态的）

global_coordination_nanjing 南京全局状态存储地址（获取站点状态的）

global_endpoints 站点状态上报地址，分别为：北京 name, 公网地址， 内网地址；南京name, 公网， 内网

config_public_intranet_endpoint_map：公网地址到内网地址的映射

```