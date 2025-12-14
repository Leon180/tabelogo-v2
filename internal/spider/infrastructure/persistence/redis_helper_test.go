package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func setupTestRedis(t *testing.T) (*RedisHelper, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)

	client := redisclient.NewClient(&redisclient.Options{
		Addr: mr.Addr(),
	})

	logger := zap.NewNop()
	helper := NewRedisHelper(client, logger)

	return helper, mr
}

func TestRedisHelper_SetJSON_GetJSON(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:key"
	original := &testStruct{Name: "test", Value: 42}

	// Act - Set
	err := helper.SetJSON(ctx, key, original, 1*time.Hour)

	// Assert - Set
	require.NoError(t, err)

	// Act - Get
	var retrieved testStruct
	err = helper.GetJSON(ctx, key, &retrieved)

	// Assert - Get
	require.NoError(t, err)
	assert.Equal(t, original.Name, retrieved.Name)
	assert.Equal(t, original.Value, retrieved.Value)
}

func TestRedisHelper_GetJSON_NotFound(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()

	// Act
	var result testStruct
	err := helper.GetJSON(ctx, "nonexistent:key", &result)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}

func TestRedisHelper_Delete(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:key"

	// Set a value first
	err := helper.SetJSON(ctx, key, &testStruct{Name: "test"}, 1*time.Hour)
	require.NoError(t, err)

	// Act
	err = helper.Delete(ctx, key)

	// Assert
	require.NoError(t, err)

	// Verify deletion
	var result testStruct
	err = helper.GetJSON(ctx, key, &result)
	assert.Error(t, err)
}

func TestRedisHelper_Exists(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:key"

	// Act - Before setting
	exists, err := helper.Exists(ctx, key)
	require.NoError(t, err)
	assert.False(t, exists)

	// Set value
	err = helper.SetJSON(ctx, key, &testStruct{Name: "test"}, 1*time.Hour)
	require.NoError(t, err)

	// Act - After setting
	exists, err = helper.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisHelper_SetOperations(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:set"

	// Act - Add members
	err := helper.SetAdd(ctx, key, "member1", "member2", "member3")
	require.NoError(t, err)

	// Act - Get members
	members, err := helper.SetMembers(ctx, key)
	require.NoError(t, err)
	assert.Len(t, members, 3)
	assert.Contains(t, members, "member1")
	assert.Contains(t, members, "member2")
	assert.Contains(t, members, "member3")

	// Act - Remove member
	err = helper.SetRemove(ctx, key, "member2")
	require.NoError(t, err)

	// Verify removal
	members, err = helper.SetMembers(ctx, key)
	require.NoError(t, err)
	assert.Len(t, members, 2)
	assert.NotContains(t, members, "member2")
}

func TestRedisHelper_Expire(t *testing.T) {
	// Arrange
	helper, mr := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:key"

	// Set value
	err := helper.SetJSON(ctx, key, &testStruct{Name: "test"}, 0) // No TTL initially
	require.NoError(t, err)

	// Act
	err = helper.Expire(ctx, key, 1*time.Hour)

	// Assert
	require.NoError(t, err)
}
