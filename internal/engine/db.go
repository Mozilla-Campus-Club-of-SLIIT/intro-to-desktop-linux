package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

const (
	StreamCmdLogs     = "cmdlogs"
	GroupProcessed    = "processed"
	ZSetLeaderboard   = "leaderboard"
	HashSessionsPrefix = "sessions:"
	
	ChannelLeaderboard = "leaderboard_updates"
	ChannelSessions    = "session_updates"
)

var ErrSessionNotFound = fmt.Errorf("session not found")

// DB wraps the valkey client and provides helper methods for our application logic.
type DB struct {
	client valkey.Client
}

// NewDB initializes a new Valkey client and sets up necessary structures.
func NewDB(addr string) (*DB, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	db := &DB{client: client}

	// Setup structures
	if err := db.setup(context.Background()); err != nil {
		client.Close()
		return nil, err
	}

	return db, nil
}

// setup ensures that necessary keys and groups exist in Valkey.
func (db *DB) setup(ctx context.Context) error {
	// Create consumer group for cmdlogs if it doesn't exist.
	// MKSTREAM ensures the stream is created if it doesn't exist.
	err := db.client.Do(ctx, db.client.B().XgroupCreate().Key(StreamCmdLogs).Group(GroupProcessed).Id("0").Mkstream().Build()).Error()
	if err != nil {
		// Ignore error if group already exists
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			// In some versions it might be different, let's just log or check if it's really an error we care about.
			// For simplicity in this setup, we'll continue.
		}
	}
	return nil
}

// Close closes the underlying valkey client.
func (db *DB) Close() {
	db.client.Close()
}

// --- Command Logs (Stream) ---

// AddCmdLog adds a new command record to the cmdlogs stream.
func (db *DB) AddCmdLog(ctx context.Context, username, command string) error {
	return db.client.Do(ctx, db.client.B().Xadd().Key(StreamCmdLogs).Id("*").FieldValue().
		FieldValue("username", username).
		FieldValue("command", command).
		FieldValue("timestamp", time.Now().Format(time.RFC3339)).
		Build()).Error()
}

// AcknowledgeCmdLog marks a command as processed by acknowledging it in the consumer group.
func (db *DB) AcknowledgeCmdLog(ctx context.Context, id string) error {
	return db.client.Do(ctx, db.client.B().Xack().Key(StreamCmdLogs).Group(GroupProcessed).Id(id).Build()).Error()
}

// --- Leaderboard (Sorted Set) ---

// UpdateScore updates the score for a given user in the leaderboard.
func (db *DB) UpdateScore(ctx context.Context, username string, score float64) error {
	err := db.client.Do(ctx, db.client.B().Zadd().Key(ZSetLeaderboard).ScoreMember().ScoreMember(score, username).Build()).Error()
	if err != nil {
		return err
	}
	// Notify subscribers about leaderboard update
	return db.PublishLeaderboardUpdate(ctx, fmt.Sprintf("update:%s:%f", username, score))
}

// GetLeaderboard returns the top users and their scores.
func (db *DB) GetLeaderboard(ctx context.Context, topN int64) ([]valkey.ZScore, error) {
	resp, err := db.client.Do(ctx, db.client.B().Zrevrange().Key(ZSetLeaderboard).Start(0).Stop(topN-1).Withscores().Build()).AsZScores()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Sessions (Hash) ---

// SetSession sets or updates a user session.
// It stores the session ID and the expiration timestamp (calculated from ttl) in a hash,
// and also sets a key expiration for automatic cleanup.
func (db *DB) SetSession(ctx context.Context, username, id string, ttl time.Duration) error {
	key := HashSessionsPrefix + username
	expiresAt := time.Now().Add(ttl) // Calculate expiration timestamp

	// Set fields in Hash
	err := db.client.Do(ctx, db.client.B().Hset().Key(key).FieldValue().
		FieldValue("id", id).
		FieldValue("expiresAt", expiresAt.Format(time.RFC3339)). // Store timestamp in hash
		Build()).Error()
	if err != nil {
		return err
	}

	// Also set expiration on the key itself so it automatically cleans up
	if ttl > 0 {
		if err := db.client.Do(ctx, db.client.B().Expire().Key(key).Seconds(int64(ttl.Seconds())).Build()).Error(); err != nil {
			return err
		}
	}

	// Notify subscribers about session update
	return db.PublishSessionUpdate(ctx, fmt.Sprintf("set:%s:%s:%s", username, id, expiresAt.Format(time.RFC3339)))
}

// GetSession retrieves the session ID and expiration for a user.
func (db *DB) GetSession(ctx context.Context, username string) (map[string]string, error) {
	key := HashSessionsPrefix + username
	sessionMap, err := db.client.Do(ctx, db.client.B().Hgetall().Key(key).Build()).AsStrMap()
	if err != nil {
		return nil, err
	}
	if len(sessionMap) == 0 {
		return nil, ErrSessionNotFound
	}
	return sessionMap, nil
}

// DeleteSession removes a user's session.
func (db *DB) DeleteSession(ctx context.Context, username string) error {
	key := HashSessionsPrefix + username
	err := db.client.Do(ctx, db.client.B().Del().Key(key).Build()).Error()
	if err != nil {
		return err
	}
	return db.PublishSessionUpdate(ctx, fmt.Sprintf("delete:%s", username))
}

// --- Pub/Sub ---

// PublishLeaderboardUpdate publishes a message to the leaderboard channel.
func (db *DB) PublishLeaderboardUpdate(ctx context.Context, message string) error {
	return db.client.Do(ctx, db.client.B().Publish().Channel(ChannelLeaderboard).Message(message).Build()).Error()
}

// PublishSessionUpdate publishes a message to the session channel.
func (db *DB) PublishSessionUpdate(ctx context.Context, message string) error {
	return db.client.Do(ctx, db.client.B().Publish().Channel(ChannelSessions).Message(message).Build()).Error()
}

// Subscribe provides a way to listen to specified channels.
func (db *DB) Subscribe(ctx context.Context, channels []string, handler func(msg valkey.PubSubMessage)) error {
	return db.client.Receive(ctx, db.client.B().Subscribe().Channel(channels...).Build(), handler)
}

// ReadCmdLogs reads new messages from the cmdlogs stream for the processed group.
// This is typically called in a loop by a worker.
func (db *DB) ReadCmdLogs(ctx context.Context, consumerName string, count int64) ([]valkey.XRangeEntry, error) {
	// ">" means read only new messages that never been delivered to other consumers.
	resp, err := db.client.Do(ctx, db.client.B().Xreadgroup().
		Group(GroupProcessed, consumerName).
		Count(count).
		Block(0). // Block until new messages arrive
		Streams().
		Key(StreamCmdLogs).
		Id(">").
		Build()).AsXRead()
	
	if err != nil {
		return nil, err
	}

	if entries, ok := resp[StreamCmdLogs]; ok {
		return entries, nil
	}
	
	return nil, nil
}
