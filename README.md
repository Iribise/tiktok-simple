# tiktok-simple
实现了所有接口，目前只是能跑，数据库用了Mysql和Redis，不过两者没有同步，除了官方库用到的其他相关框架和库只有gin和gorm

用户注册和登录用bcrypt对password做了哈希，token只是简单的用户名+哈希后的password，鉴权也只是简单地从内存里的map中找

各个相关数据结构定义和数据库表的结构应该可以优化

没有实现根据视频生成封面图

没什么编码规范，部分地方代码重复度还挺高

如果clone了想跑请先根据自身情况修改constant.go中的常量
