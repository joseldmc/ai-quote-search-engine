package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Domain Models
type Quote struct {
	Text      string `json:"text"`
	Movie     string `json:"movie"`
	Character string `json:"character"`
}

type QuoteData struct {
	Query  string  `json:"query"`
	Quotes []Quote `json:"quotes"`
}

type SearchResult struct {
	Quote Quote
	Score float64
}

// Emotional Context and Themes
type EmotionalProfile struct {
	Emotions []string
	Themes   []string
	Tone     string
}

// Repository Interface
type QuoteRepository interface {
	LoadQuotes(filename string) (*QuoteData, error)
}

// Service Interface
type QuoteService interface {
	SearchQuotes(query string, topN int) ([]SearchResult, error)
}

// File Repository Implementation
type FileQuoteRepository struct{}

func NewFileQuoteRepository() *FileQuoteRepository {
	return &FileQuoteRepository{}
}

func (r *FileQuoteRepository) LoadQuotes(filename string) (*QuoteData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open quotes file: %w", err)
	}
	defer file.Close()

	var data QuoteData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse quotes file: %w", err)
	}

	if len(data.Quotes) == 0 {
		return nil, fmt.Errorf("no quotes found in file")
	}

	return &data, nil
}

// Quote Search Service Implementation
type SemanticQuoteService struct {
	data          *QuoteData
	repository    QuoteRepository
	quoteProfiles map[string]EmotionalProfile
}

func NewSemanticQuoteService(repo QuoteRepository) *SemanticQuoteService {
	return &SemanticQuoteService{
		repository:    repo,
		quoteProfiles: make(map[string]EmotionalProfile),
	}
}

func (s *SemanticQuoteService) Initialize(filename string) error {
	data, err := s.repository.LoadQuotes(filename)
	if err != nil {
		return err
	}
	s.data = data
	s.buildQuoteProfiles()
	return nil
}

func (s *SemanticQuoteService) buildQuoteProfiles() {
	// Define emotional and thematic profiles for each quote
	s.quoteProfiles["Just keep swimming."] = EmotionalProfile{
		Emotions: []string{"overwhelmed", "struggling", "tired", "perseverance", "exhausted"},
		Themes:   []string{"persistence", "continuation", "resilience", "keep going", "don't give up"},
		Tone:     "encouraging",
	}

	s.quoteProfiles["After all, tomorrow is another day!"] = EmotionalProfile{
		Emotions: []string{"disappointed", "setback", "failure", "hopeful", "optimistic"},
		Themes:   []string{"new beginning", "fresh start", "moving forward", "hope", "recovery"},
		Tone:     "hopeful",
	}

	s.quoteProfiles["I'm going to make him an offer he can't refuse."] = EmotionalProfile{
		Emotions: []string{"powerful", "assertive", "confident", "threatening"},
		Themes:   []string{"negotiation", "power", "control", "persuasion", "dominance"},
		Tone:     "assertive",
	}

	s.quoteProfiles["Life is like a box of chocolates. You never know what you're gonna get."] = EmotionalProfile{
		Emotions: []string{"uncertain", "curious", "accepting", "philosophical"},
		Themes:   []string{"uncertainty", "unpredictability", "acceptance", "life lessons", "randomness", "adventure"},
		Tone:     "philosophical",
	}

	s.quoteProfiles["You can't handle the truth!"] = EmotionalProfile{
		Emotions: []string{"angry", "confrontational", "defensive", "intense"},
		Themes:   []string{"reality", "confrontation", "denial", "honesty"},
		Tone:     "confrontational",
	}

	s.quoteProfiles["May the Force be with you."] = EmotionalProfile{
		Emotions: []string{"supportive", "encouraging", "hopeful", "wishing well"},
		Themes:   []string{"good luck", "support", "blessing", "encouragement", "journey", "challenge"},
		Tone:     "supportive",
	}

	s.quoteProfiles["There's no place like home."] = EmotionalProfile{
		Emotions: []string{"homesick", "nostalgic", "longing", "comfort", "belonging"},
		Themes:   []string{"home", "belonging", "comfort", "safety", "family", "roots"},
		Tone:     "nostalgic",
	}

	s.quoteProfiles["I'll be back."] = EmotionalProfile{
		Emotions: []string{"determined", "confident", "persistent", "threatening"},
		Themes:   []string{"return", "persistence", "promise", "determination"},
		Tone:     "determined",
	}

	s.quoteProfiles["Houston, we have a problem."] = EmotionalProfile{
		Emotions: []string{"worried", "anxious", "crisis", "emergency", "stressed"},
		Themes:   []string{"problem", "crisis", "emergency", "trouble", "difficulty", "challenge"},
		Tone:     "urgent",
	}

	s.quoteProfiles["You're gonna need a bigger boat."] = EmotionalProfile{
		Emotions: []string{"overwhelmed", "underprepared", "surprised", "inadequate"},
		Themes:   []string{"underestimation", "surprise", "insufficient", "unprepared", "bigger challenge"},
		Tone:     "humorous",
	}

	s.quoteProfiles["The first rule of Fight Club is: You do not talk about Fight Club."] = EmotionalProfile{
		Emotions: []string{"secretive", "rebellious", "exclusive"},
		Themes:   []string{"secrecy", "rules", "rebellion", "underground"},
		Tone:     "mysterious",
	}

	s.quoteProfiles["Why so serious?"] = EmotionalProfile{
		Emotions: []string{"stressed", "tense", "overthinking", "worried"},
		Themes:   []string{"lighten up", "relax", "perspective", "humor", "don't stress"},
		Tone:     "playful",
	}

	s.quoteProfiles["You had me at hello."] = EmotionalProfile{
		Emotions: []string{"love", "romantic", "convinced", "charmed", "smitten"},
		Themes:   []string{"love", "romance", "connection", "immediate attraction"},
		Tone:     "romantic",
	}

	s.quoteProfiles["To infinity and beyond!"] = EmotionalProfile{
		Emotions: []string{"excited", "adventurous", "optimistic", "ambitious", "enthusiastic"},
		Themes:   []string{"adventure", "limitless", "ambition", "excitement", "new horizons", "exploration"},
		Tone:     "enthusiastic",
	}

	s.quoteProfiles["Life moves pretty fast. If you don't stop and look around once in a while, you could miss it."] = EmotionalProfile{
		Emotions: []string{"rushed", "busy", "overwhelmed", "reflective", "mindful"},
		Themes:   []string{"slow down", "mindfulness", "appreciation", "present moment", "life passing by", "pause"},
		Tone:     "reflective",
	}

	s.quoteProfiles["Nobody puts Baby in a corner."] = EmotionalProfile{
		Emotions: []string{"defensive", "protective", "standing up", "assertive"},
		Themes:   []string{"standing up", "protection", "respect", "worth"},
		Tone:     "protective",
	}

	s.quoteProfiles["It's not who I am underneath, but what I do that defines me."] = EmotionalProfile{
		Emotions: []string{"determined", "purposeful", "identity seeking", "motivated"},
		Themes:   []string{"actions", "identity", "purpose", "character", "doing vs being", "integrity"},
		Tone:     "inspirational",
	}

	s.quoteProfiles["Our lives are defined by opportunities, even the ones we miss."] = EmotionalProfile{
		Emotions: []string{"regretful", "reflective", "philosophical", "contemplative", "accepting"},
		Themes:   []string{"opportunities", "regret", "choices", "life path", "what if", "acceptance", "past decisions"},
		Tone:     "philosophical",
	}

	s.quoteProfiles["Get busy living, or get busy dying."] = EmotionalProfile{
		Emotions: []string{"motivated", "determined", "stuck", "choosing", "decisive"},
		Themes:   []string{"choice", "action", "motivation", "moving forward", "living fully", "purpose"},
		Tone:     "motivational",
	}

	s.quoteProfiles["The only way out is through."] = EmotionalProfile{
		Emotions: []string{"struggling", "difficult", "challenged", "facing hardship", "enduring"},
		Themes:   []string{"perseverance", "endurance", "facing difficulty", "no shortcuts", "courage", "pushing through"},
		Tone:     "resilient",
	}
}

func (s *SemanticQuoteService) SearchQuotes(query string, topN int) ([]SearchResult, error) {
	if s.data == nil {
		return nil, fmt.Errorf("service not initialized")
	}

	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	queryProfile := s.analyzeQuery(query)

	var results []SearchResult
	for _, quote := range s.data.Quotes {
		score := s.calculateRelevanceScore(queryProfile, quote)
		if score > 0 {
			results = append(results, SearchResult{
				Quote: quote,
				Score: score,
			})
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no matching quotes found for your situation")
	}

	// Sort by score (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Return top N results
	if topN > len(results) {
		topN = len(results)
	}

	return results[:topN], nil
}

func (s *SemanticQuoteService) analyzeQuery(query string) EmotionalProfile {
	query = strings.ToLower(query)
	profile := EmotionalProfile{
		Emotions: []string{},
		Themes:   []string{},
		Tone:     "neutral",
	}

	// Emotion detection
	emotionKeywords := map[string][]string{
		"excited":     {"excited", "thrilled", "looking forward", "can't wait", "enthusiastic"},
		"overwhelmed": {"overwhelmed", "too much", "stressed", "stressed out", "swamped", "drowning"},
		"worried":     {"worried", "anxious", "nervous", "concerned", "scared", "afraid", "fearful"},
		"sad":         {"sad", "depressed", "down", "unhappy", "heartbroken", "grieving"},
		"uncertain":   {"uncertain", "unsure", "don't know", "confused", "lost", "unclear"},
		"hopeful":     {"hopeful", "optimistic", "positive", "looking up"},
		"tired":       {"tired", "exhausted", "worn out", "drained", "fatigued", "burnt out"},
		"stuck":       {"stuck", "trapped", "can't move", "stagnant", "blocked"},
		"lonely":      {"lonely", "alone", "isolated", "disconnected"},
		"motivated":   {"motivated", "driven", "determined", "ready", "pumped"},
		"struggling":  {"struggling", "difficult", "hard", "tough", "challenging"},
		"nostalgic":   {"miss", "missing", "nostalgia", "remember", "used to"},
		"crisis":      {"emergency", "crisis", "urgent", "sick", "ill", "dying", "serious"},
	}

	for emotion, keywords := range emotionKeywords {
		for _, keyword := range keywords {
			if strings.Contains(query, keyword) {
				profile.Emotions = append(profile.Emotions, emotion)
				break
			}
		}
	}

	// Theme detection
	themeKeywords := map[string][]string{
		"new beginning": {"new city", "moving", "starting", "new job", "new chapter", "fresh start"},
		"preparation":   {"preparing", "preparation", "getting ready", "planning"},
		"health":        {"sick", "ill", "health", "medical", "hospital", "dying", "disease"},
		"pet":           {"dog", "cat", "pet", "animal"},
		"change":        {"change", "transition", "different", "new"},
		"challenge":     {"challenge", "difficult", "hard", "tough", "obstacle"},
		"journey":       {"journey", "path", "road", "way", "going"},
		"home":          {"home", "house", "place", "belong"},
		"future":        {"future", "ahead", "tomorrow", "next", "coming"},
		"uncertainty":   {"don't know", "uncertain", "unsure", "unpredictable"},
		"persistence":   {"keep going", "continue", "don't give up", "push through"},
		"support":       {"need help", "support", "encourage", "guidance"},
		"loss":          {"lost", "losing", "gone", "miss"},
	}

	for theme, keywords := range themeKeywords {
		for _, keyword := range keywords {
			if strings.Contains(query, keyword) {
				profile.Themes = append(profile.Themes, theme)
				break
			}
		}
	}

	return profile
}

func (s *SemanticQuoteService) calculateRelevanceScore(queryProfile EmotionalProfile, quote Quote) float64 {
	quoteProfile, exists := s.quoteProfiles[quote.Text]
	if !exists {
		return 0.0
	}

	score := 0.0
	maxPossibleScore := 100.0

	// Emotion matching (highest weight)
	for _, queryEmotion := range queryProfile.Emotions {
		for _, quoteEmotion := range quoteProfile.Emotions {
			if queryEmotion == quoteEmotion {
				score += 10.0
			}
			// Related emotions
			if s.areEmotionsRelated(queryEmotion, quoteEmotion) {
				score += 5.0
			}
		}
	}

	// Theme matching (high weight)
	for _, queryTheme := range queryProfile.Themes {
		for _, quoteTheme := range quoteProfile.Themes {
			if queryTheme == quoteTheme {
				score += 8.0
			}
			if s.areThemesRelated(queryTheme, quoteTheme) {
				score += 4.0
			}
		}
	}

	// Apply penalties for mismatched emotions
	negativePenaltyPairs := map[string][]string{
		"worried":     {"playful", "humorous", "assertive"},
		"crisis":      {"playful", "romantic", "humorous"},
		"sad":         {"playful", "enthusiastic"},
		"excited":     {"urgent", "crisis"},
		"overwhelmed": {"confrontational", "assertive"},
	}

	for _, queryEmotion := range queryProfile.Emotions {
		if penalties, exists := negativePenaltyPairs[queryEmotion]; exists {
			for _, penaltyEmotion := range penalties {
				for _, quoteEmotion := range quoteProfile.Emotions {
					if quoteEmotion == penaltyEmotion {
						score -= 15.0
					}
				}
			}
		}
	}

	// Normalize to 0.0-1.0 range
	if score < 0 {
		return 0.0
	}
	normalizedScore := score / maxPossibleScore
	if normalizedScore > 1.0 {
		normalizedScore = 1.0
	}

	return normalizedScore
}

func (s *SemanticQuoteService) areEmotionsRelated(emotion1, emotion2 string) bool {
	relatedEmotions := map[string][]string{
		"worried":     {"anxious", "concerned", "crisis", "stressed"},
		"overwhelmed": {"stressed", "tired", "struggling", "exhausted"},
		"excited":     {"enthusiastic", "adventurous", "optimistic", "hopeful"},
		"sad":         {"disappointed", "regretful", "lonely"},
		"uncertain":   {"confused", "lost", "philosophical"},
		"struggling":  {"difficult", "challenged", "tired", "stuck"},
		"motivated":   {"determined", "purposeful", "decisive"},
	}

	if related, exists := relatedEmotions[emotion1]; exists {
		for _, rel := range related {
			if rel == emotion2 {
				return true
			}
		}
	}

	if related, exists := relatedEmotions[emotion2]; exists {
		for _, rel := range related {
			if rel == emotion1 {
				return true
			}
		}
	}

	return false
}

func (s *SemanticQuoteService) areThemesRelated(theme1, theme2 string) bool {
	relatedThemes := map[string][]string{
		"new beginning": {"change", "journey", "future"},
		"health":        {"crisis", "support", "challenge"},
		"challenge":     {"persistence", "journey", "struggle"},
		"uncertainty":   {"future", "change", "journey"},
		"preparation":   {"future", "journey", "challenge"},
	}

	if related, exists := relatedThemes[theme1]; exists {
		for _, rel := range related {
			if rel == theme2 {
				return true
			}
		}
	}

	if related, exists := relatedThemes[theme2]; exists {
		for _, rel := range related {
			if rel == theme1 {
				return true
			}
		}
	}

	return false
}

// CLI Interface
type CLI struct {
	service QuoteService
}

func NewCLI(service QuoteService) *CLI {
	return &CLI{service: service}
}

func (c *CLI) Run() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Movie Quote Search Engine                         ║")
	fmt.Println("║          Finding inspiration in cinema                     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nHow are you feeling? Describe your situation:")
		fmt.Print("> ")

		query, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		query = strings.TrimSpace(query)

		if query == "" {
			continue
		}

		if strings.ToLower(query) == "exit" || strings.ToLower(query) == "quit" {
			fmt.Println("\nTake care! Remember: just keep swimming. 🐠")
			break
		}

		c.displayResults(query)
	}
}

func (c *CLI) RunSingleQuery(query string) {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Movie Quote Search Engine                         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Query: %s\n", query)

	c.displayResults(query)
}

func (c *CLI) displayResults(query string) {
	results, err := c.service.SearchQuotes(query, 3)
	if err != nil {
		fmt.Printf("\n❌ %s\n", err.Error())
		fmt.Println("Try describing your feelings differently.")
		return
	}

	fmt.Println("\n✨ Here are some quotes that might resonate with you:\n")
	for i, result := range results {
		fmt.Printf("%d. [%.2f] \"%s\"\n", i+1, result.Score, result.Quote.Text)
		fmt.Printf("   — %s (%s)\n", result.Quote.Character, result.Quote.Movie)
		if i < len(results)-1 {
			fmt.Println()
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 60))
}

func main() {
	// Parse command line arguments
	args := os.Args[1:]

	var quotesFile string
	var customQuery string

	// Default quotes file
	quotesFile = "quotes.json"

	// Parse arguments
	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "--query" || arg == "-q" {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --query flag requires an argument\n")
				printUsage()
				os.Exit(1)
			}
			customQuery = args[i+1]
			i += 2
		} else if arg == "--help" || arg == "-h" {
			printUsage()
			os.Exit(0)
		} else {
			// Assume it's the quotes file path
			quotesFile = arg
			i++
		}
	}

	// Dependency injection
	repo := NewFileQuoteRepository()
	service := NewSemanticQuoteService(repo)

	// Initialize service with quotes file
	if err := service.Initialize(quotesFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run CLI
	cli := NewCLI(service)

	// If custom query provided, run single query mode
	if customQuery != "" {
		cli.RunSingleQuery(customQuery)
	} else {
		cli.Run()
	}
}

func printUsage() {
	fmt.Println("Movie Quote Search Engine - Find inspiration in cinema")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go [quotes_file] [options]")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  quotes_file    Path to quotes JSON file (default: quotes.json)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --query, -q    Custom query to search (skips interactive mode)")
	fmt.Println("  --help, -h     Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Interactive mode with default file")
	fmt.Println("  go run main.go")
	fmt.Println()
	fmt.Println("  # Interactive mode with custom file")
	fmt.Println("  go run main.go my_quotes.json")
	fmt.Println()
	fmt.Println("  # Single query mode")
	fmt.Println("  go run main.go --query \"I just got rejected and feel like giving up\"")
	fmt.Println()
	fmt.Println("  # Single query with custom file")
	fmt.Println("  go run main.go my_quotes.json --query \"I need motivation\"")
}
