# Movie Quote Search Engine

A CLI-based mental wellness tool that finds relevant movie quotes when you're going through difficult moments. Simply describe your situation or feelings, and the engine will surface quotes that genuinely resonate with your emotional state.

## Features

- **Semantic Understanding**: Goes beyond keyword matching to understand emotional intent and context
- **Emotional Profiling**: Dynamically analyzes any quote for emotions, themes, and tones
- **Crisis Detection**: Identifies suicidal ideation and provides immediate mental health resources
- **Clean Architecture**: Separation of concerns with repository pattern and service layer
- **Confidence Scores**: See how well each quote matches your situation (0.0-1.0 scale)
- **Two Modes**: Interactive conversation or single-query mode for quick searches
- **Smart Matching**: Understands related emotions (happy ↔ joyful, worried ↔ anxious) and themes (family ↔ home)
- **Sentiment Filtering**: Strong penalties prevent tone-mismatched quotes (no threatening quotes for happy moments)
- **Universal Lexicon**: Works with ANY quotes JSON file - no hardcoded quote profiles needed
- **Error Handling**: Graceful handling of empty queries, missing files, and no matches
- **Interactive CLI**: User-friendly command-line interface with backspace support
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
> I'm very happy because I will meet with my family tonight

✨ Here are some quotes that might resonate with you:

1. [0.78] "There's no place like home."
   — Dorothy (The Wizard of Oz)

2. [0.64] "You had me at hello."
   — Dorothy Boyd (Jerry Maguire)

3. [0.58] "To infinity and beyond!"
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
- **Detects emotions**: 19 emotion categories including happy, excited, worried, sad, grateful, loved, nostalgic, and more
- **Identifies themes**: 20+ universal themes like family, connection, celebration, home, persistence, challenge, health, etc.
- **Understands context**: Differentiates between "excited about moving" vs "worried about health" vs "happy to meet family"
- **Analyzes sentiment**: Positive, negative, or neutral tone detection

### 2. Dynamic Quote Analysis
The engine analyzes ANY quote in your JSON file without hardcoded profiles:
- **Feature extraction**: Automatically detects emotions, themes, sentiment, and tone from quote text
- **Universal lexicon**: Uses a comprehensive emotional vocabulary that works with any movie quote
- **No manual setup**: Just add quotes to your JSON - the system handles the rest

Example analysis:
```
Query: "I'm very happy because I will meet with my family tonight"
→ Extracts: emotion:happy, emotion:excited, theme:family, theme:connection, 
            theme:time, sentiment:positive

Quote: "There's no place like home."
→ Extracts: theme:home, theme:family, theme:belonging, sentiment:positive

Result: Strong match! ✓
```

### 3. Intelligent Scoring System
The engine uses **cosine similarity** (industry-standard ML technique) with weighted features:

**Feature Weights:**
- **Emotions**: 3.0x weight (highest priority)
- **Themes**: 2.5x weight
- **Sentiment/Tone**: 1.0x weight

**Sentiment Filtering:**
- **Opposite sentiments**: 60-70% penalty
    - Happy query + threatening quote = heavily penalized
- **Mismatched emotional context**: 70% penalty
    - Joyful query + conflict/struggle quote = blocked
- **Neutral mismatches**: 20% penalty

**Related Emotion Bonuses:**
- "happy" relates to "excited", "grateful", "joyful", "content"
- "worried" relates to "anxious", "uncertain", "stressed"
- System automatically boosts scores for emotionally related matches

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
- **Crisis detection**: Special handling for suicidal ideation (see Safety Features below)

## Safety Features

The application includes **crisis detection** to protect user wellbeing:

### Crisis Indicators
Automatically detects phrases indicating suicidal ideation or self-harm:
- "don't want to live", "want to die", "kill myself"
- "end my life", "suicide", "suicidal"
- "hurt myself", "harm myself"
- "better off dead", "end it all"
- "no reason to live", "can't go on"

### Crisis Response
When crisis language is detected, the app:
1. **Does not show movie quotes** (inappropriate for crisis situations)
2. **Displays compassionate message** acknowledging their difficult time
3. **Provides immediate resources**:
    - **988 Suicide & Crisis Lifeline** (US - call or text, 24/7)
    - **Crisis Text Line** (text HOME to 741741)
    - **International Association for Suicide Prevention** (global resources)
    - **Emergency Services** (911 or local emergency number)

Example output:
```
⚠️  It sounds like you might be going through a really difficult time.

While movie quotes can be inspiring, what you're experiencing
may need professional support. Please consider reaching out:

🆘 CRISIS RESOURCES:
   • National Suicide Prevention Lifeline (US)
     Call or Text: 988
   ...
```

This ensures the application acts **responsibly** when users are in mental health crisis.

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

The lexicon is comprehensive but you can extend it by modifying `NewEmotionalLexicon()`:

**Current Emotion Categories (19 total):**
- Negative: overwhelmed, worried, sad, tired, stuck, uncertain, struggling, lonely, rejected, angry
- Positive: motivated, excited, happy, hopeful, proud, grateful, peaceful, loved
- Neutral: nostalgic

**Current Theme Categories (20+ total):**
- persistence, change, future, challenge, opportunity, home, family, journey, truth
- action, choice, life, time, support, beginning, loss, health, moving, preparation
- connection, celebration

Add new emotions or themes:

```go
EmotionKeywords: map[string][]string{
    "overwhelmed": {"overwhelm", "too much", "swamp", ...},
    "confident": {"confident", "sure", "certain", "assured"}, // New emotion
    // ... more emotions
}

ThemeKeywords: map[string][]string{
    "home": {"home", "belong", "place", ...},
    "friendship": {"friend", "buddy", "companion", "pal"}, // New theme
    // ... more themes
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
- **Dynamic Analysis**: No hardcoded quote profiles - works with any quotes JSON
- **Cosine Similarity**: Industry-standard ML technique for semantic matching

### How the Matching Algorithm Works

1. **Tokenization**: Removes stop words and punctuation from text
2. **Feature Extraction**: Analyzes both query and quotes for emotions, themes, sentiment
3. **Feature Vectors**: Creates weighted vectors (emotions: 3.0x, themes: 2.5x)
4. **Cosine Similarity**: Calculates angle between query and quote vectors
5. **Sentiment Filtering**: Applies penalties for mismatched emotional contexts
6. **Normalization**: Scores normalized to 0.0-1.0 range

## Future Enhancements

- User feedback loop to improve matching accuracy
- Custom emotional lexicon configuration files
- Multi-language support for international quotes
- Quote categories and advanced filtering options
- User favorites and search history
- Expanded crisis resources for different countries
- API endpoint for programmatic access
- Web-based interface option
- Machine learning model training on user preferences

## License

This is a demonstration project for educational purposes.

## Contributing

Feel free to submit issues and enhancement requests!