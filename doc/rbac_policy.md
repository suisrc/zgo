# 策略

## 说明
一开始， 我想使用casbin + （method, domain, path, netmask）方式做权限判断， 但是整体来说， 想法太简单， 虽然说简单就是好的， 是稳定的， 如果是做一个应用权限管理是好的， 如果做的是一个平台的，无法提供想阿里云一样复杂的授权体系。所以套用“拿来主义”精神， 让我们借助casbin + (阿里云授权规则)构建一个我们自己的授权系统

授权策略
https://help.aliyun.com/document_detail/93739.html?spm=a2c4g.11186623.2.5.4b8610d45RVjOg
```json
policy  = {
     <version_block>,
     <statement_block>
}
<version_block> = "Version" : ("1")
<statement_block> = "Statement" : [ <statement>, <statement>, ... ]
<statement> = { 
    <effect_block>,
    <action_block>,
    <resource_block>,
    <condition_block?>
}
<effect_block> = "Effect" : ("Allow" | "Deny")  
<action_block> = "Action" : 
    ("*" | [<action_string>, <action_string>, ...])
<resource_block> = "Resource" : 
    ("*" | [<resource_string>, <resource_string>, ...])
<condition_block> = "Condition" : <condition_map>
<condition_map> = {
  <condition_type_string> : { 
      <condition_key_string> : <condition_value_list>,
      <condition_key_string> : <condition_value_list>,
      ...
  },
  <condition_type_string> : {
      <condition_key_string> : <condition_value_list>,
      <condition_key_string> : <condition_value_list>,
      ...
  }, ...
}  
<condition_value_list> = [<condition_value>, <condition_value>, ...]
<condition_value> = ("String" | "Number" | "Boolean" | "Date and time" | "IP address")
```

权限策略语法说明：

元素名称                描述
效力（Effect）          授权效力包括两种：允许（Allow）和拒绝（Deny）。
操作（Action）          操作是指对具体资源的操作。
资源（Resource）        资源是指被授权的具体对象。
限制条件（Condition）   限制条件是指授权生效的限制条件。

版本：当前支持的权限策略版本，固定为1，不允许修改。
授权语句：一个权限策略可以有多条授权语句。
每条授权语句的效力为：Allow或Deny 。
每条授权语句中，操作（Action）和资源（Resource）都支持多值。
每条授权语句都支持独立的限制条件（Condition）。

一个条件块支持多个条件的组合，每个条件的操作类型可以不同。
Deny优先原则： 一个用户可以被授予多个权限策略。当这些权限策略同时包含Allow和 Deny时，遵循Deny优先原则。

元素取值：
当元素取值为字符串类型（String）、数字类型（Number）、日期类型（Date and time）、布尔类型（Boolean）和IP地址类型（IP address）时，需要使用双引号。
当元素取值为字符串值（String）时，支持使用*和?进行模糊匹配。
*代表0个或多个任意的英文字母。 例如：ecs:Describe* 表示ECS的所有以Describe开头的操作。
?代表1个任意的英文字母。

操作（Action）：云服务所定义的API操作名称
格式 <service-name>:<action-name>
service-name：产品名称。
action-name： 相关的API操作接口名称
样例："Action": ["lhdg2:ListBuckets", "lhdg2:Describe*", "lhdg2:Describe*"]

资源（Resource）：资源是指被授权的具体对象。
格式：遵循（Resource Name）的统一规范，fmes:<service>:<role>:<account>:<relative>
样例："Resource": ["fmes:lhdg2:*:*:instance/inst-001"]
