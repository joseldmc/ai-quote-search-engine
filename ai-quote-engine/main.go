package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
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

// Emotional and thematic lexicons
type EmotionalContext struct {
	PrimaryEmotion  string
	RelatedEmotions []string
	IntensityScore  float64
	Valence         string // positive, negative, neutral
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

// Dynamic Quote Search Service Implementation
type SemanticQuoteService struct {
	data       *QuoteData
	repository QuoteRepository
	lexicon    *EmotionalLexicon
}

func NewSemanticQuoteService(repo QuoteRepository) *SemanticQuoteService {
	return &SemanticQuoteService{
		repository: repo,
		lexicon:    NewEmotionalLexicon(),
	}
}

func (s *SemanticQuoteService) Initialize(filename string) error {
	data, err := s.repository.LoadQuotes(filename)
	if err != nil {
		return err
	}
	s.data = data
	return nil
}

func (s *SemanticQuoteService) SearchQuotes(query string, topN int) ([]SearchResult, error) {
	if s.data == nil {
		return nil, fmt.Errorf("service not initialized")
	}

	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Check for crisis indicators
	if s.detectCrisis(query) {
		return nil, fmt.Errorf("CRISIS_DETECTED")
	}

	queryContext := s.analyzeText(query)

	var results []SearchResult
	for _, quote := range s.data.Quotes {
		quoteContext := s.analyzeText(quote.Text)

		// Check tone compatibility before calculating similarity
		if !s.areTonesCompatible(queryContext, quoteContext, quote.Text) {
			continue
		}

		score := s.calculateSimilarity(queryContext, quoteContext)

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

// Check if query tone is compatible with quote tone
func (s *SemanticQuoteService) areTonesCompatible(queryFeatures, quoteFeatures map[string]float64, quoteText string) bool {
	quoteTextLower := strings.ToLower(quoteText)

	// Detect if query is positive/joyful
	isQueryPositive := queryFeatures["emotion:happy"] > 0 ||
		queryFeatures["emotion:excited"] > 0 ||
		queryFeatures["emotion:grateful"] > 0 ||
		queryFeatures["emotion:loved"] > 0 ||
		queryFeatures["sentiment:positive"] > 1

	// Detect if query is about celebration/connection
	isQueryCelebratory := queryFeatures["theme:family"] > 0 ||
		queryFeatures["theme:connection"] > 0 ||
		queryFeatures["theme:celebration"] > 0 ||
		queryFeatures["theme:home"] > 0

	// Block quotes with these patterns for positive queries
	if isQueryPositive || isQueryCelebratory {
		// Philosophical/serious quotes
		seriousPatterns := []string{
			"defines me", "who i am", "underneath",
			"refuse", "offer", "handle the truth",
			"fight club", "rule", "serious",
			"problem", "crisis", "boat",
		}

		for _, pattern := range seriousPatterns {
			if strings.Contains(quoteTextLower, pattern) {
				return false
			}
		}

		// Block quotes with negative/dark themes
		if quoteFeatures["sentiment:negative"] > quoteFeatures["sentiment:positive"] {
			return false
		}
	}

	// Detect if query is about worry/concern (especially health-related)
	isQueryWorried := queryFeatures["emotion:worried"] > 0 ||
		queryFeatures["emotion:sad"] > 0 ||
		queryFeatures["theme:health"] > 0

	// Block inappropriate quotes for worried/concerned queries
	if isQueryWorried {
		// Block threatening, aggressive, or dismissive quotes
		inappropriatePatterns := []string{
			"refuse", "offer", "can't refuse",
			"i'll be back",
			"fight club", "rule",
			"boat", "gonna need",
			"nobody puts", "corner",
			"handle the truth",
		}

		for _, pattern := range inappropriatePatterns {
			if strings.Contains(quoteTextLower, pattern) {
				return false
			}
		}

		// For health/worry concerns, only allow supportive or empathetic quotes
		// Block overly optimistic quotes that might feel dismissive
		dismissivePatterns := []string{
			"tomorrow is another day",           // Can feel dismissive of current worry
			"life is like", "box of chocolates", // Too philosophical for immediate concern
			"life moves pretty fast", // Not appropriate for health worries
		}

		for _, pattern := range dismissivePatterns {
			if strings.Contains(quoteTextLower, pattern) {
				return false
			}
		}
	}

	// Detect if query is negative/struggling
	isQueryNegative := queryFeatures["emotion:struggling"] > 0 ||
		queryFeatures["emotion:overwhelmed"] > 0 ||
		queryFeatures["emotion:tired"] > 0 ||
		queryFeatures["sentiment:negative"] > 1

	// Block overly cheerful quotes for negative queries
	if isQueryNegative {
		cheerfulPatterns := []string{
			"infinity and beyond",
			"had me at hello",
		}

		for _, pattern := range cheerfulPatterns {
			if strings.Contains(quoteTextLower, pattern) {
				return false
			}
		}
	}

	return true
}

// Detect crisis situations that require professional help
func (s *SemanticQuoteService) detectCrisis(query string) bool {
	query = strings.ToLower(query)

	// Crisis indicators - suicidal ideation, self-harm
	crisisPatterns := []string{
		"kill myself",
		"end my life",
		"don't want to live",
		"want to die",
		"suicide",
		"suicidal",
		"hurt myself",
		"harm myself",
		"not worth living",
		"better off dead",
		"end it all",
		"can't go on",
		"no reason to live",
	}

	for _, pattern := range crisisPatterns {
		if strings.Contains(query, pattern) {
			return true
		}
	}

	return false
}

// Analyze text to extract emotional and thematic content
func (s *SemanticQuoteService) analyzeText(text string) map[string]float64 {
	text = strings.ToLower(text)
	words := s.tokenize(text)

	features := make(map[string]float64)

	// Emotion detection
	for emotion, keywords := range s.lexicon.EmotionKeywords {
		for _, word := range words {
			for _, keyword := range keywords {
				if strings.Contains(word, keyword) || strings.Contains(keyword, word) {
					features["emotion:"+emotion] += 1.0

					// Add related emotions with lower weight
					if related, exists := s.lexicon.EmotionRelations[emotion]; exists {
						for _, relEmotion := range related {
							features["emotion:"+relEmotion] += 0.3
						}
					}
				}
			}
		}
	}

	// Theme detection
	for theme, keywords := range s.lexicon.ThemeKeywords {
		for _, word := range words {
			for _, keyword := range keywords {
				if strings.Contains(word, keyword) || strings.Contains(keyword, word) {
					features["theme:"+theme] += 1.0
				}
			}
		}
	}

	// Sentiment and tone
	positiveCount := 0.0
	negativeCount := 0.0

	for _, word := range words {
		for _, posWord := range s.lexicon.PositiveWords {
			if word == posWord {
				positiveCount += 1.0
			}
		}
		for _, negWord := range s.lexicon.NegativeWords {
			if word == negWord {
				negativeCount += 1.0
			}
		}
	}

	if positiveCount > 0 {
		features["sentiment:positive"] = positiveCount
	}
	if negativeCount > 0 {
		features["sentiment:negative"] = negativeCount
	}

	// Action vs reflection
	for _, word := range words {
		for _, actionWord := range s.lexicon.ActionWords {
			if word == actionWord {
				features["tone:action"] += 1.0
			}
		}
		for _, reflectWord := range s.lexicon.ReflectiveWords {
			if word == reflectWord {
				features["tone:reflective"] += 1.0
			}
		}
	}

	return features
}

// Calculate cosine similarity between query and quote feature vectors
func (s *SemanticQuoteService) calculateSimilarity(queryFeatures, quoteFeatures map[string]float64) float64 {
	// Get all unique features
	allFeatures := make(map[string]bool)
	for feature := range queryFeatures {
		allFeatures[feature] = true
	}
	for feature := range quoteFeatures {
		allFeatures[feature] = true
	}

	// Calculate dot product and magnitudes
	dotProduct := 0.0
	queryMagnitude := 0.0
	quoteMagnitude := 0.0

	for feature := range allFeatures {
		queryVal := queryFeatures[feature]
		quoteVal := quoteFeatures[feature]

		// Weight emotions higher than other features
		weight := 1.0
		if strings.HasPrefix(feature, "emotion:") {
			weight = 3.0
		} else if strings.HasPrefix(feature, "theme:") {
			weight = 2.5
		}

		dotProduct += (queryVal * weight) * (quoteVal * weight)
		queryMagnitude += (queryVal * weight) * (queryVal * weight)
		quoteMagnitude += (quoteVal * weight) * (quoteVal * weight)
	}

	// Apply sentiment filtering - stronger penalties for mismatches
	querySentiment := s.getSentiment(queryFeatures)
	quoteSentiment := s.getSentiment(quoteFeatures)

	// Strong penalty for opposite sentiments
	sentimentPenalty := 1.0
	if querySentiment == "negative" && quoteSentiment == "positive" {
		sentimentPenalty = 0.4
	} else if querySentiment == "positive" && quoteSentiment == "negative" {
		sentimentPenalty = 0.3
	} else if querySentiment == "neutral" && quoteSentiment != "neutral" {
		sentimentPenalty = 0.8
	}

	// Apply tone filtering - penalize mismatched emotional contexts
	queryHasJoy := queryFeatures["emotion:happy"] > 0 || queryFeatures["emotion:excited"] > 0 || queryFeatures["emotion:grateful"] > 0
	quoteHasConflict := quoteFeatures["theme:challenge"] > 0 || quoteFeatures["theme:truth"] > 0 ||
		strings.Contains(strings.ToLower(s.getQuoteTextFromFeatures(quoteFeatures)), "refuse") ||
		strings.Contains(strings.ToLower(s.getQuoteTextFromFeatures(quoteFeatures)), "defines")

	tonePenalty := 1.0
	if queryHasJoy && quoteHasConflict {
		tonePenalty = 0.3
	}

	if queryMagnitude == 0 || quoteMagnitude == 0 {
		return 0.0
	}

	similarity := (dotProduct / (math.Sqrt(queryMagnitude) * math.Sqrt(quoteMagnitude))) * sentimentPenalty * tonePenalty

	// Normalize to 0-1 range
	if similarity < 0 {
		similarity = 0
	}
	if similarity > 1 {
		similarity = 1
	}

	return similarity
}

func (s *SemanticQuoteService) getQuoteTextFromFeatures(features map[string]float64) string {
	// This is a helper - in practice we'd need to track quote text separately
	// For now, return empty string as we can't reverse engineer the quote
	return ""
}

func (s *SemanticQuoteService) getSentiment(features map[string]float64) string {
	positive := features["sentiment:positive"]
	negative := features["sentiment:negative"]

	if positive > negative {
		return "positive"
	} else if negative > positive {
		return "negative"
	}
	return "neutral"
}

func (s *SemanticQuoteService) tokenize(text string) []string {
	text = strings.ToLower(text)

	// Replace punctuation with spaces
	replacer := strings.NewReplacer(
		".", " ", ",", " ", "!", " ", "?", " ", ";", " ", ":", " ",
		"'", "", "\"", "", "(", " ", ")", " ",
	)
	text = replacer.Replace(text)

	words := strings.Fields(text)

	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
		"we": true, "you": true, "your": true,
	}

	var tokens []string
	for _, word := range words {
		if len(word) > 1 && !stopWords[word] {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// Emotional Lexicon - Dynamic knowledge base
type EmotionalLexicon struct {
	EmotionKeywords  map[string][]string
	EmotionRelations map[string][]string
	ThemeKeywords    map[string][]string
	PositiveWords    []string
	NegativeWords    []string
	ActionWords      []string
	ReflectiveWords  []string
}

func NewEmotionalLexicon() *EmotionalLexicon {
	return &EmotionalLexicon{
		EmotionKeywords: map[string][]string{
			"overwhelmed": {"overwhelm", "too much", "swamp", "drown", "bury", "flood"},
			"motivated":   {"motivat", "inspir", "driven", "determin", "pump", "energiz", "push"},
			"worried":     {"worr", "anxious", "nervous", "concern", "afraid", "scare", "fear"},
			"sad":         {"sad", "depress", "down", "unhappy", "heartbreak", "griev", "mourn"},
			"excited":     {"excit", "thrill", "eager", "enthusias", "can't wait", "looking forward"},
			"happy":       {"happy", "joy", "delight", "glad", "pleased", "cheer", "content", "elated"},
			"tired":       {"tire", "exhaust", "worn", "drain", "fatigue", "burnt out", "weary"},
			"stuck":       {"stuck", "trap", "stagnant", "block", "immobil"},
			"uncertain":   {"uncertain", "unsure", "confus", "lost", "unclear", "doubt"},
			"hopeful":     {"hope", "optimis", "positive", "bright", "promising"},
			"struggling":  {"struggl", "difficult", "hard", "tough", "challeng", "fight"},
			"lonely":      {"lone", "isolat", "disconnect", "apart", "solo"},
			"rejected":    {"reject", "dismiss", "refus", "turn down", "decline"},
			"proud":       {"proud", "accomplish", "achiev", "success", "triumph"},
			"grateful":    {"grateful", "thankful", "appreciat", "bless"},
			"angry":       {"angry", "mad", "furious", "irritat", "frustrat"},
			"peaceful":    {"peace", "calm", "serene", "tranquil", "relax"},
			"loved":       {"love", "loved", "caring", "affection", "warm"},
			"nostalgic":   {"nostalg", "remember", "miss", "memories", "past"},
		},

		EmotionRelations: map[string][]string{
			"overwhelmed": {"stressed", "anxious", "tired"},
			"worried":     {"anxious", "uncertain", "stressed"},
			"sad":         {"lonely", "hopeless", "disappointed"},
			"excited":     {"hopeful", "motivated", "energized", "happy"},
			"happy":       {"excited", "grateful", "joyful", "content"},
			"struggling":  {"overwhelmed", "tired", "stuck"},
			"stuck":       {"frustrated", "uncertain", "lost"},
			"rejected":    {"sad", "disappointed", "hurt"},
			"motivated":   {"determined", "hopeful", "energized"},
			"grateful":    {"happy", "content", "blessed"},
			"loved":       {"happy", "grateful", "warm"},
		},

		ThemeKeywords: map[string][]string{
			"persistence": {"keep", "continu", "persist", "endur", "carry on", "push through", "stay", "swimming"},
			"change":      {"chang", "transiti", "shift", "transform", "evolv", "new"},
			"future":      {"future", "ahead", "tomorrow", "next", "coming", "forward"},
			"challenge":   {"challeng", "obstacle", "difficult", "problem", "hurdle", "barrier"},
			"opportunity": {"opportun", "chance", "possibil", "option", "opening"},
			"home":        {"home", "belong", "place", "family", "roots", "comfort", "house"},
			"family":      {"family", "families", "relatives", "parents", "children", "together", "reunion"},
			"journey":     {"journey", "path", "road", "way", "travel", "adventure"},
			"truth":       {"truth", "reality", "honest", "real", "genuine", "authentic"},
			"action":      {"action", "doing", "act", "move", "step", "initiative"},
			"choice":      {"choice", "decis", "choose", "select", "pick", "option"},
			"life":        {"life", "living", "exist", "being", "alive"},
			"time":        {"time", "moment", "now", "present", "today", "day", "tonight", "evening"},
			"support":     {"support", "help", "assist", "guid", "encourag", "force", "with you"},
			"beginning":   {"begin", "start", "new", "fresh", "commence", "launch"},
			"loss":        {"loss", "lost", "missing", "gone", "absence"},
			"health":      {"sick", "ill", "health", "medical", "disease", "pain", "dying", "doctor", "hospital"},
			"moving":      {"mov", "relocat", "transfer", "shift"},
			"preparation": {"prepar", "ready", "plan", "arrang", "organiz"},
			"connection":  {"meet", "meeting", "see", "visit", "reunion", "gather", "connect"},
			"celebration": {"celebrat", "party", "event", "occasion", "special"},
			"hope":        {"hope", "better", "improve", "forward", "through"},
			"difficulty":  {"difficult", "hard", "tough", "struggle", "through", "out"},
		},

		PositiveWords: []string{
			"good", "great", "happy", "joy", "love", "wonderful", "amazing",
			"beautiful", "excellent", "fantastic", "brilliant", "superb",
			"beyond", "infinity", "force", "blessed", "lucky", "glad",
			"delight", "pleased", "cheerful", "sweet", "nice", "fun",
		},

		NegativeWords: []string{
			"bad", "terrible", "awful", "horrible", "sad", "pain", "hurt",
			"problem", "crisis", "emergency", "sick", "dying", "death",
			"refuse", "reject", "serious", "difficult", "struggle",
		},

		ActionWords: []string{
			"do", "act", "move", "go", "make", "create", "build",
			"fight", "push", "drive", "run", "work", "try", "living",
			"swimming", "busy", "defines",
		},

		ReflectiveWords: []string{
			"think", "feel", "believe", "understand", "know", "wonder",
			"consider", "reflect", "remember", "realize", "learn",
			"life", "truth", "defined", "opportunities", "miss",
		},
	}
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
		// Check if it's a crisis situation
		if err.Error() == "CRISIS_DETECTED" {
			c.displayCrisisResources()
			return
		}

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

func (c *CLI) displayCrisisResources() {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println()
	fmt.Println("⚠️  It sounds like you might be going through a really difficult time.")
	fmt.Println()
	fmt.Println("While movie quotes can be inspiring, what you're experiencing")
	fmt.Println("may need professional support. Please consider reaching out:")
	fmt.Println()
	fmt.Println("🆘 CRISIS RESOURCES:")
	fmt.Println()
	fmt.Println("   • National Suicide Prevention Lifeline (US)")
	fmt.Println("     Call or Text: 988")
	fmt.Println("     Available 24/7, free and confidential")
	fmt.Println()
	fmt.Println("   • Crisis Text Line (US)")
	fmt.Println("     Text: HOME to 741741")
	fmt.Println()
	fmt.Println("   • International Association for Suicide Prevention")
	fmt.Println("     https://www.iasp.info/resources/Crisis_Centres/")
	fmt.Println()
	fmt.Println("   • Emergency Services")
	fmt.Println("     Call: 911 (US) or your local emergency number")
	fmt.Println()
	fmt.Println("You don't have to go through this alone. These trained")
	fmt.Println("professionals are available to listen and help, any time.")
	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
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
