package config

//
//import (
//	"context"
//	"encoding/json"
//
//	"github.com/it-chep/tutors.git/pkg/postgres"
//
//	"github.com/pkg/errors"
//)
//
//type ConfigKey string
//
//type contextConfigKey string
//
//const contextKey contextConfigKey = "config"
//
//// Config Конфиг реального времени если можно так назвать предназначен, чтобы конфигурировать необходимые для бизнеса параметры
//// например различные фича флаги, лимиты, задержки. Код берет объект конфига из базы и пытается вычитать значение из json поля value
//type Config interface {
//	GetValue(ctx context.Context, key ConfigKey) (Value, error)
//}
//
//// config ...
//type config struct {
//	pool postgres.PoolWrapper
//}
//
//// New ...
//func New(pool postgres.PoolWrapper) Config {
//	return &config{
//		pool: pool,
//	}
//}
//
//// GetValue получает конфиг из таблицы config
//// create table if not exists config (
////
////	key text primary key,
////	value jsonb not null,
////	description text not null
////
//// );
//func (c *config) GetValue(ctx context.Context, key ConfigKey) (Value, error) {
//	sql := `select value from config where key = $1`
//
//	var jsonData []byte
//	err := c.pool.QueryRow(ctx, sql, key).Scan(&jsonData)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to get config value")
//	}
//
//	var val concreteValue
//	if err = json.Unmarshal(jsonData, &val); err != nil {
//		return nil, errors.Wrap(err, "failed to unmarshal config value")
//	}
//
//	return &val, nil
//}
//
//// ContextWithConfig прокинуть конфиг в контекст
//func ContextWithConfig(ctx context.Context, cfg Config) context.Context {
//	return context.WithValue(ctx, contextKey, cfg)
//}
//
//// fromContext получение конфига из контекста
//func fromContext(ctx context.Context) Config {
//	l, ok := ctx.Value(contextKey).(Config)
//	if ok {
//		return l
//	}
//	return nil
//}
//
//// GetValue получает структуру конфига из контекста и возвращает значение из базы
//func GetValue(ctx context.Context, key ConfigKey) (Value, error) {
//	cfg := fromContext(ctx)
//	return cfg.GetValue(ctx, key)
//}
