package cmd

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	blackjack "github.com/onemorebsmith/blackjack-solver/src"
	"github.com/onemorebsmith/blackjack-solver/src/blackjack/strategies"
)

type BJConfig struct {
	Decks         int                            `json:"decks"`
	Penetration   float32                        `json:"pen"`
	MaxSplits     int                            `json:"maxSplits"`
	IsH17         bool                           `json:"h17"`
	IsRSA         bool                           `json:"rsa"`
	IsDAS         bool                           `json:"das"`
	ShoesToSim    int                            `json:"shoesToSim"`
	RoundsPerHour float32                        `json:"rph"`
	Bidspread     map[int]strategies.BidStrategy `json:"bidspread"`
	Strategy      string                         `json:"strategy"`
}

func (cfg BJConfig) BuildGameDescription() string {
	game := fmt.Sprintf("%d deck ", cfg.Decks)
	if cfg.IsH17 {
		game += "H17 "
	} else {
		game += "S17 "
	}
	if cfg.IsDAS {
		game += "DAS "
	}
	if cfg.IsRSA {
		game += "RSA "
	}
	return strings.TrimRight(game, " ")
}

var threads = runtime.NumCPU()

func Run(cfg BJConfig) {
	start := time.Now()
	bjRules := blackjack.NewBlackjackGameRules(blackjack.InitGame(blackjack.H17Rules, blackjack.H17Splits))
	bjRules.SetDealerHitsSoft17(cfg.IsH17)
	bjRules.SetDoubleAfterSplit(cfg.IsDAS)
	bjRules.SetResplitAces(cfg.IsRSA)
	bjRules.SetMaxPlayerSplits(cfg.MaxSplits)
	bjRules.SetUseSimpleDeviations(false)
	bjRules.SetPenetration(cfg.Penetration)

	log.Printf("simming %d shoes of %s w/ %f pen", cfg.ShoesToSim, cfg.BuildGameDescription(), cfg.Penetration)
	switch cfg.Strategy {
	case "hilo":
		log.Println("using HiLo strategy")
		bjRules.TrackingStrategy = strategies.InitHighLow(cfg.Bidspread)
	case "flatbet":
		log.Println("Using flatbet strategy")
		bjRules.TrackingStrategy = strategies.InitFlatbetStrategy()
	default:
		log.Println("Using flatbet strategy")
		bjRules.TrackingStrategy = strategies.InitFlatbetStrategy()
	}

	overallResults := make([]blackjack.GameResults, threads)
	wg := sync.WaitGroup{}

	shoesPerThread := cfg.ShoesToSim / threads

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			overallResults[idx] = blackjack.PlayGame(*bjRules, cfg.Decks, shoesPerThread, 10000, cfg.RoundsPerHour)
		}(i)
	}
	wg.Wait()

	variance := float32(0)
	hourlyVariance := float32(0)
	// extract variance before we re-aggregate
	for _, r := range overallResults {
		variance += r.EVVariance
		hourlyVariance += r.HourlyEVVariance
	}
	variance /= float32(len(overallResults))
	hourlyVariance /= float32(len(overallResults))

	game := fmt.Sprintf("%d deck ", cfg.Decks)
	if bjRules.DealerHitsSoft17 {
		game += "H17 "
	} else {
		game += "S17 "
	}
	if bjRules.DoubleAfterSplit {
		game += "DAS "
	}
	if bjRules.ReSplitAces {
		game += "RSA "
	}
	game = strings.TrimRight(game, " ")

	aggregatedResults := blackjack.AggregateResults(overallResults...)
	winPct := float32(aggregatedResults.Wins) / float32(aggregatedResults.Hands)
	losePct := float32(aggregatedResults.Losses) / float32(aggregatedResults.Hands)
	pushPct := float32(aggregatedResults.Pushes) / float32(aggregatedResults.Hands)
	bjPct := float32(aggregatedResults.Blackjacks) / float32(aggregatedResults.Hands)

	log.Println("====================================")
	log.Printf("   Threads %d, elapsed: %s", threads, time.Since(start).Truncate(time.Millisecond))
	log.Println("====================================")
	log.Printf("%s, %f pen, %d hands, %f rph", game, bjRules.Penetration, aggregatedResults.Hands, cfg.RoundsPerHour)
	log.Printf("   EV (units):         %f units", aggregatedResults.EV)
	log.Printf("   EV (hand):          %f units", aggregatedResults.EV/float32(aggregatedResults.Hands))
	log.Printf("   EV (hourly):        %f units", aggregatedResults.EV/float32(aggregatedResults.Hands)*cfg.RoundsPerHour)
	log.Printf("   W/L/P:              %f/%f/%f", winPct, losePct, pushPct)
	log.Printf("   Blackjacks:         %d, %f%%", aggregatedResults.Blackjacks, bjPct)
	log.Printf("   1 STD (hand):     +-%f units", variance)
	log.Printf("   1 STD (hourly):   +-%f units", hourlyVariance)
	log.Printf("TC Stats --- ")
	log.Printf("   HighTC (avg)        %f ", aggregatedResults.HighTC/float32(aggregatedResults.Hands))
	log.Printf("   LowTC  (avg)        %f ", aggregatedResults.LowTC/float32(aggregatedResults.Hands))
	log.Printf("   AvgTC  (avg)        %f ", aggregatedResults.AvgTC/float32(aggregatedResults.Hands))

}
