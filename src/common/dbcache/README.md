## Doc

### 初始化

```
// 创建一个针对 plants 表的定义
plantTable := NewTable(TableConfig{
    Table:    "plants",                 // 表名
    IdName:   "id",                     // 主键字段名
    Cols:     nil,                      // 留空表示全部字段
    Typ:      reflect.TypeOf(Plant{}),  // 用于 unmarshal 的结构体
    ExpireIn: 3600,                     // 缓存失效时间，默认 60s
    Prefix:   "",                       // 默认空，缓存 key 的前缀
})

// redisPool: redis 链接池对象 *pool.Pool
dcHead := dbcache.NewRedisHeader(redisPool)

// conn: postgresql 数据库连接对象 *sqlx.DB
dc := dbcache.NewCache(conn, dcHead)
```

### 获取数据

```
data := plantTable.Data(id)
err = dc.FetchDatas(data)
if err != nil {
    ...
    return
}

if data.Error != nil {
    ...
    return
}

plant := data.Value.(*proto.PlantModel)
```

### 数据库变动自动更新

因为自动创建会带来一些风险，所以手动创建接收 trigger 和 notify function

```sql
-- 只需要创建一个接收函数
CREATE OR REPLACE FUNCTION trigger_db_update_notify() RETURNS trigger AS $$
  BEGIN
    PERFORM pg_notify('db_update_notify',
      concat(TG_NAME, ',', TG_WHEN, ',', TG_LEVEL, ',', TG_OP, ',', TG_TABLE_SCHEMA, ',', TG_TABLE_NAME, ',', NEW.id));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

-- 每个表要单独创建 trigger
DROP TRIGGER IF EXISTS trigger_db_update_notify ON plants;
CREATE TRIGGER trigger_db_update_notify AFTER INSERT OR UPDATE
ON plants
FOR EACH ROW
EXECUTE PROCEDURE trigger_db_update_notify();
```

确保你的数据表，在读远大于写的情况下才使用 trigger 自动更新

```
trigger := dbcache.NewUpdateTigger(dbcache.TriggerConfig{
    RedisPool:      redisPool,                      // redis 用于分布式锁
    DataSource:     dataSource,                     // postgresql 连接参数
    Tables:         []*dbcache.Table{plantTable},   // dbcache 表定义，用于表名接收哪些通知
    DBCache:        dc,                             // dbcache 实例

    // 以下是参数留空的默认值
    UpdateInterval: 3,                              // 数据批量更新，延迟时间
    BufferSize:     100,                            // 通知缓冲池大小
    LockName:       "db_update_notify",             // 分布式锁，名称
    ChannelName:    "db_update_notify",             // pg_notify channel 名称
})
go trigger.Start() // start

// 日志里看到 lucky, got locked; start receive notify 表示当前进程拿到锁，开始接收通知
```
