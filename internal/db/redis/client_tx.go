package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Tx is the interface that provides methods to add to a transaction.
type Tx interface {
	Get(ctx context.Context, key string)
	HGet(ctx context.Context, key string)
	Set(ctx context.Context, key string, value []byte, exp time.Duration)
	HSet(ctx context.Context, key string, value map[string]any)
	Delete(ctx context.Context, key string)
	Expire(ctx context.Context, key string, exp time.Duration)
	LastCommand() command
}

// TxFunc is a function that contains all steps in a transaction.
type TxFunc func(tx Tx)

// tx contains the transaction pipe and the results.
type tx struct {
	pipe redis.Pipeliner
	cmds []command
}

// Get returns the value for the key.
func (t *tx) Get(ctx context.Context, key string) {
	cmd := t.pipe.Get(ctx, key)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// HGet returns the structured data for the key as a map of strings.
func (t *tx) HGet(ctx context.Context, key string) {
	cmd := t.pipe.HGetAll(ctx, key)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// Set the value for the key with an expiration time. If the expiration
// time is zero, the key will not expire.
func (t *tx) Set(ctx context.Context, key string, value []byte, exp time.Duration) {
	cmd := t.pipe.Set(ctx, key, value, exp)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// HSet sets the value for the key with an expiration time. If the expiration
// time is zero, the key will not expire.
func (t *tx) HSet(ctx context.Context, key string, value map[string]any) {
	cmd := t.pipe.HSet(ctx, key, value)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// Delete the key.
func (t tx) Delete(ctx context.Context, key string) {
	cmd := t.pipe.Del(ctx, key)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// Expire sets expire time for the key.
func (t tx) Expire(ctx context.Context, key string, exp time.Duration) {
	cmd := t.pipe.Expire(ctx, key, exp)
	t.cmds = append(t.cmds, command{key: key, cmd: cmd})
}

// LastCommand returns the last command in the transaction pipe.
func (t tx) LastCommand() command {
	if len(t.cmds) == 0 {
		return command{}
	}
	return t.cmds[len(t.cmds)-1]
}

// TxResults contains the results from a transaction.
type TxResult struct {
	b [][]byte
	m []map[string]string
}

// AllBytes returns all bytes retreived with successful Get operations
// in the order they where set to the transaction.
func (r TxResult) AllBytes() [][]byte {
	return r.b
}

// FirstBytes returns the first bytes retreived with the transaction.
func (r TxResult) FirstBytes() []byte {
	if len(r.b) == 0 {
		return nil
	}
	return r.b[0]
}

// LastBytes returns the bytes retreived with the transaction.
func (r TxResult) LastBytes() []byte {
	if len(r.b) == 0 {
		return nil
	}
	return r.b[len(r.b)-1]
}

// IndexBytes returns the data at the provided index.
func (r TxResult) IndexBytes(i int) []byte {
	if len(r.b) == 0 || i >= len(r.b) {
		return nil
	}
	return r.b[i]
}

// AllMaps returns all maps retreived with successful HGet operations
func (r TxResult) AllMaps() []map[string]string {
	return r.m
}

// command contains the key that the command target and the redis.Cmder
// containing the issued command.
type command struct {
	key string
	cmd redis.Cmder
}

// FirstMap returns the first map retreived with the transaction.
func (r TxResult) FirstMap() map[string]string {
	if len(r.m) == 0 {
		return nil
	}
	return r.m[0]
}

// LastMap returns the last map retreived with the transaction.
func (r TxResult) LastMap() map[string]string {
	if len(r.m) == 0 {
		return nil
	}
	return r.m[len(r.m)-1]
}

// IndexMap returns the data at the provided index.
func (r TxResult) IndexMap(i int) map[string]string {
	if len(r.m) == 0 || i >= len(r.m) {
		return nil
	}
	return r.m[i]
}

// execCommands runs the commands set to the pipeline.
func execCommands(ctx context.Context, tx *tx) (TxResult, error) {
	var result TxResult
	if _, err := tx.pipe.Exec(ctx); err != nil {
		if errors.Is(err, redis.Nil) {
			return result, ErrKeyNotFound
		}
		return result, err
	}
	for _, cmd := range tx.cmds {
		switch c := cmd.cmd.(type) {
		case *redis.StringCmd:
			b, err := c.Bytes()
			if err != nil {
				return result, err
			}
			result.b = append(result.b, b)
		case *redis.MapStringStringCmd:
			m, err := c.Result()
			if err != nil {
				return result, err
			}
			result.m = append(result.m, m)
		case *redis.StatusCmd, *redis.InfoCmd:
			if err := c.Err(); err != nil {
				return result, err
			}
		}
	}
	return result, nil
}
