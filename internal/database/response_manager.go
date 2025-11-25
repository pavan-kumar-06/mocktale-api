package database

import (
	"database/sql"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type MovieResponseCounts struct {
	Option0 int
	Option1 int  
	Option2 int
	Option3 int
	Total   int
}

type VoteDelta struct {
	UserToken    string
	MovieSlug    string
	OptionChosen int
}

type ResponseManager struct {
	db         *sql.DB
	newVotes   chan VoteDelta
	voteCount  int64
	active     int32
	mu         sync.Mutex
}

var ResponseManagerInstance *ResponseManager

func InitResponseManager(db *sql.DB) {
	ResponseManagerInstance = &ResponseManager{
		db:       db,
		newVotes: make(chan VoteDelta, 5000),
		active:   1,
	}

	// Start multiple flushers
	for i := 0; i < 3; i++ {
		go ResponseManagerInstance.periodicFlusher()
	}
	log.Println("✅ Response manager initialized - No cache mode")
}

// AddResponse - Just push to channel
func (rm *ResponseManager) AddResponse(userToken, movieSlug string, optionChosen int) bool {
	if atomic.LoadInt32(&rm.active) == 0 {
		return false
	}

	select {
	case rm.newVotes <- VoteDelta{
		UserToken:    userToken,
		MovieSlug:    movieSlug,
		OptionChosen: optionChosen,
	}:
		atomic.AddInt64(&rm.voteCount, 1)
		return true
	default:
		// Channel full - apply backpressure
		return false
	}
}

// GetMovieCounts - Direct DB query (reads are fast anyway)
func (rm *ResponseManager) GetMovieCounts(movieSlug string) *MovieResponseCounts {
	var counts MovieResponseCounts
	err := rm.db.QueryRow(`
		SELECT option_0, option_1, option_2, option_3, total_votes 
		FROM movie_responses WHERE movie_slug = ?
	`, movieSlug).Scan(&counts.Option0, &counts.Option1, &counts.Option2, &counts.Option3, &counts.Total)
	
	if err != nil {
		return &MovieResponseCounts{}
	}
	return &counts
}

// HasUserVoted - Direct DB query (no cache)
func (rm *ResponseManager) HasUserVoted(userToken, movieSlug string) (bool, int) {
	var option int
	err := rm.db.QueryRow(`
		SELECT option_chosen FROM user_responses 
		WHERE user_token = ? AND movie_slug = ?
	`, userToken, movieSlug).Scan(&option)
	
	if err == nil {
		return true, option
	}
	
	return false, 0
}

// periodicFlusher - Just flush votes to DB
func (rm *ResponseManager) periodicFlusher() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for atomic.LoadInt32(&rm.active) == 1 {
		select {
		case <-ticker.C:
			rm.flushAvailableVotes()
		}
	}
}

// flushAvailableVotes - Process available votes
func (rm *ResponseManager) flushAvailableVotes() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Collect up to 1000 votes
	votesBatch := make([]VoteDelta, 0, 1000)
	
	for i := 0; i < 1000; i++ {
		select {
		case vote := <-rm.newVotes:
			votesBatch = append(votesBatch, vote)
		default:
			break
		}
	}

	if len(votesBatch) > 0 {
		rm.flushBatchToDB(votesBatch)
	}
}

// flushBatchToDB - Batch insert without cache
func (rm *ResponseManager) flushBatchToDB(votes []VoteDelta) {
	if len(votes) == 0 {
		return
	}

	tx, err := rm.db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userResponseStmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO user_responses (user_token, movie_slug, option_chosen) 
		VALUES (?, ?, ?)
	`)
	if err != nil {
		log.Printf("❌ Failed to prepare user response statement: %v", err)
		return
	}
	defer userResponseStmt.Close()

	movieResponseStmt, err := tx.Prepare(`
		INSERT INTO movie_responses (movie_slug, option_0, option_1, option_2, option_3, total_votes)
		VALUES (?, 0, 0, 0, 0, 0)
		ON CONFLICT(movie_slug) DO UPDATE SET
			option_0 = option_0 + ?,
			option_1 = option_1 + ?,
			option_2 = option_2 + ?,
			option_3 = option_3 + ?,
			total_votes = total_votes + ?
	`)
	if err != nil {
		log.Printf("❌ Failed to prepare movie response statement: %v", err)
		return
	}
	defer movieResponseStmt.Close()

	moviesToUpdate := make(map[string][5]int)
	successfulVotes := 0

	for _, vote := range votes {
		_, err := userResponseStmt.Exec(vote.UserToken, vote.MovieSlug, vote.OptionChosen)
		if err != nil {
			log.Printf("❌ Failed to update user response for %s: %v", vote.MovieSlug, err)
			continue
		}

		// No cache update - just update aggregation
		delta := moviesToUpdate[vote.MovieSlug]
		switch vote.OptionChosen {
		case 0: delta[0]++
		case 1: delta[1]++
		case 2: delta[2]++
		case 3: delta[3]++
		}
		delta[4]++
		moviesToUpdate[vote.MovieSlug] = delta
		successfulVotes++
	}

	for movieSlug, delta := range moviesToUpdate {
		_, err := movieResponseStmt.Exec(movieSlug, delta[0], delta[1], delta[2], delta[3], delta[4])
		if err != nil {
			log.Printf("❌ Failed to update movie response for %s: %v", movieSlug, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("❌ Failed to commit: %v", err)
		return
	}

	log.Printf("📤 Flushed %d/%d votes (%d movies updated)", successfulVotes, len(votes), len(moviesToUpdate))
}

// Shutdown - Clean shutdown
func (rm *ResponseManager) Shutdown() {
	atomic.StoreInt32(&rm.active, 0)
	time.Sleep(2 * time.Second)
	rm.flushAvailableVotes()
}