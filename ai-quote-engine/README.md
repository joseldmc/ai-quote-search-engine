# Movie Quote Search Engine

A CLI-based mental wellness tool that finds relevant movie quotes when you're going through difficult moments. Simply describe your situation or feelings, and the engine will surface quotes that genuinely resonate with your emotional state.

## Features

- **Semantic Understanding**: Goes beyond keyword matching to understand emotional intent and context
- **Emotional Profiling**: Each quote is mapped to specific emotions, themes, and tones
- **Clean Architecture**: Separation of concerns with repository pattern and service layer
- **Confidence Scores**: See how well each quote matches your situation (0.0-1.0 scale)
- **Two Modes**: Interactive conversation or single-query mode for quick searches
- **Smart Matching**: Understands related emotions (worried ↔ anxious) and themes (new beginning ↔ change)
- **Negative Filtering**: Prevents tone-mismatched quotes (no playful quotes during a crisis)
- **Error Handling**: Graceful handling of empty queries, missing files, and no matches
- **Interactive CLI**: User-friendly command-line interface with visual feedback
- **Flexible Input**: Support for custom quotes files and command-line queries

## Architecture

The application follows clean architecture principles:

```
├── Domain Models (Quote, QuoteData, SearchResult)
├── Repository Layer (QuoteRepository interface, FileQuoteRepository)
├── Service Layer (QuoteService interface, SemanticQuoteService)
└── Presentation Layer (CLI)
```

### Components

- **Domain Models**: Core business entities without dependencies
- **Repository**: Handles data access and file operations
- **Service**: Contains business logic for semantic search and scoring
- **CLI**: User interface for interaction

## Prerequisites

- Go 1.16 or higher
- `quotes.json` file in the same directory as the executable

## Installation

1. Clone or download the repository
2. Ensure you have Go installed:
   ```bash
   go version
   ```

3. Build the application:
   ```bash
   go build -o quote-search main.go
   ```

## Usage

The application supports two modes: **Interactive Mode** and **Single Query Mode**.

### Interactive Mode (Default)

1. Make sure `quotes.json` is in the same directory as the executable

2. Run the application:
   ```bash
   ./quote-search
   ```
   Or with Go:
   ```bash
   go run main.go
   ```

3. Describe your situation or feelings when prompted:
   ```
   How are you feeling? Describe your situation:
   > I'm feeling overwhelmed and need motivation to keep going
   ```

4. The engine will return the top 3 most relevant quotes with confidence scores

5. Type `exit` or `quit` to close the application

### Single Query Mode

For quick searches without entering interactive mode:

```bash
# Using default quotes.json
go run main.go --query "I just got rejected and feel like giving up"

# Short flag version
go run main.go -q "My dog is sick, I'm very worried"

# With custom quotes file
go run main.go my_quotes.json --query "I need motivation"

# After building
./quote-search --query "I'm moving to a new city"
```

### Command Line Options

```bash
go run main.go [quotes_file] [options]

Arguments:
  quotes_file    Path to quotes JSON file (default: quotes.json)

Options:
  --query, -q    Custom query to search (skips interactive mode)
  --help, -h     Show help message

Examples:
  go run main.go                                    # Interactive mode
  go run main.go my_quotes.json                     # Custom quotes file
  go run main.go --query "feeling overwhelmed"      # Single query
  go run main.go quotes.json -q "need motivation"   # Combined
```

## Example Interactions

### Interactive Mode

```
How are you feeling? Describe your situation:
> I need motivation to keep going when things are tough

✨ Here are some quotes that might resonate with you:

1. [0.92] "Just keep swimming."
   — Dory (Finding Nemo)

2. [0.87] "Get busy living, or get busy dying."
   — Andy Dufresne (The Shawshank Redemption)

3. [0.81] "The only way out is through."
   — John Ottway (The Grey)
```

```
How are you feeling? Describe your situation:
> Life feels unpredictable and I don't know what's next

✨ Here are some quotes that might resonate with you:

1. [0.85] "Life is like a box of chocolates. You never know what you're gonna get."
   — Forrest Gump (Forrest Gump)

2. [0.76] "Our lives are defined by opportunities, even the ones we miss."
   — Benjamin Button (The Curious Case of Benjamin Button)

3. [0.68] "To infinity and beyond!"
   — Buzz Lightyear (Toy Story)
```

### Single Query Mode

```bash
$ go run main.go --query "I just got rejected and feel like giving up"

╔════════════════════════════════════════════════════════════╗
║          Movie Quote Search Engine                         ║
╚════════════════════════════════════════════════════════════╝

Query: I just got rejected and feel like giving up

✨ Here are some quotes that might resonate with you:

1. [0.88] "Get busy living, or get busy dying."
   — Andy Dufresne (The Shawshank Redemption)

2. [0.82] "Just keep swimming."
   — Dory (Finding Nemo)

3. [0.75] "After all, tomorrow is another day!"
   — Scarlett O'Hara (Gone with the Wind)
```

## How It Works

The semantic search engine uses an **Emotional Profiling System** to match quotes with user queries:

### 1. Query Analysis
When you describe your situation, the engine:
- **Detects emotions**: worried, excited, overwhelmed, sad, uncertain, motivated, struggling, etc.
- **Identifies themes**: new beginning, health, challenge, uncertainty, persistence, support, etc.
- **Understands context**: Differentiates between "excited about moving" vs "worried about health"

### 2. Quote Profiling
Each quote has a detailed emotional profile:
- **Emotions**: What feelings does this quote address?
- **Themes**: What situations is it relevant to?
- **Tone**: encouraging, hopeful, motivational, philosophical, etc.

Example:
```go
"Just keep swimming." → {
    Emotions: ["overwhelmed", "struggling", "tired", "perseverance"],
    Themes: ["persistence", "resilience", "keep going"],
    Tone: "encouraging"
}
```

### 3. Relevance Scoring
The engine calculates a normalized score (0.0-1.0) based on:
- **Emotion matches**: 10 points (exact) or 5 points (related)
   - "worried" matches "anxious", "concerned", "crisis"
- **Theme matches**: 8 points (exact) or 4 points (related)
   - "new beginning" relates to "change", "journey", "future"
- **Negative penalties**: -15 points for tone mismatches
   - Won't show playful quotes when you're in crisis
   - Won't show urgent quotes when you're excited

### 4. Results
Returns top 3 quotes with confidence scores showing match quality:
- **[0.90-1.00]**: Excellent match
- **[0.70-0.89]**: Good match
- **[0.50-0.69]**: Moderate match

The displayed score helps you understand how well each quote resonates with your specific situation.

## Error Handling

The application handles several error cases:

- **Missing or invalid quotes file**: Clear error message on startup
- **Empty query**: Prompts user to enter a description
- **No meaningful words**: Asks user to rephrase
- **No matches found**: Suggests trying a different description
- **Malformed JSON**: Reports parsing errors

## Customization

### Adding More Quotes

Edit `quotes.json` and add new quote objects:

```json
{
   "text": "Your quote here",
   "movie": "Movie Name",
   "character": "Character Name"
}
```

### Adjusting Search Results

Change the number of results returned by modifying the `topN` parameter in `CLI.Run()`:

```go
results, err := c.service.SearchQuotes(query, 5) // Returns top 5 instead of 3
```

### Extending Emotional Keywords

Add new emotion-keyword mappings in `calculateRelevanceScore()`:

```go
emotionalKeywords := map[string][]string{
"perseverance": {"swimming", "keep", "going", ...},
"courage": {"brave", "fear", "face", ...}, // New emotion
// ... more mappings
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

- **Interfaces**: Enable dependency injection and testing
- **Error Wrapping**: Using `fmt.Errorf` with `%w` for error chains
- **Separation of Concerns**: Each layer has a single responsibility
- **No External Dependencies**: Uses only Go standard library

## Future Enhancements

- Fuzzy string matching for typo tolerance
- Machine learning-based semantic similarity
- Support for multiple languages
- Quote categories and filtering
- User favorites and history
- API integration for expanded quote database

## License

This is a demonstration project for educational purposes.

## Contributing

Feel free to submit issues and enhancement requests!