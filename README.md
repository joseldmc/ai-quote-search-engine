# ai-quote-search-engine

# Approach Documentation

## Ranking Strategy

The Movie Quote Search Engine uses a **dynamic semantic matching approach** based on cosine similarity between feature vectors extracted from both user queries and movie quotes. Rather than relying on hardcoded quote profiles, the system employs a universal emotional lexicon that analyzes text on-the-fly to detect 19 emotion categories, 20+ thematic elements, sentiment polarity, and tonal characteristics. The ranking algorithm applies weighted feature matching (emotions: 3.0x, themes: 2.5x) combined with aggressive tone-compatibility filtering to ensure quotes resonate emotionally with user situations, preventing inappropriate matches like showing threatening quotes to worried users or philosophical quotes to joyful celebrations.

## Implementation Approach
 
1. In the challenge, I already have the base requirements, so I summarized and provide the same information to Claude.ai
2. I start testing the application, and as soon as I noticed that someting was not getting very "helpful" quotes I provided this input to Claude.ai.
3. The ranking approach was pretty much adjusted based on the prompts and the feedback I provided to the AI.

## Key Prompts Used

Below are the most relevant prompts that shaped the development of this application:

### 1. Initial Requirements
```
I want to build a Movie Quote Search Engine for a mental wellness app. 
The way it should works is: The user describe a situation or feeling, and 
the tool should find the most relevant movie quotes when they're going 
through difficult moments. The tool needs to understand the emotional 
intent behind messy, real-world queries and surface quotes that genuinely 
resonate, not just keyword matches.
```

**Impact**: Established the core requirement for semantic understanding over simple keyword matching, leading to the feature vector approach.

---

### 2. Hardcoded Profile Concern
```
I think the current logic is very tied to the query.json I provided, 
but what if the json changes?
```

**Impact**: This critical feedback triggered a complete redesign from hardcoded quote profiles to a dynamic, universal lexicon-based system that works with ANY quotes JSON file. This led to implementing feature extraction, cosine similarity, and the `EmotionalLexicon` structure.

---

### 3. Crisis Safety Feature
```
I tested with "I don't want to live" and I feel concerned that people 
that is dealing with that kind of thoughts about not want to live, 
maybe we can recommend to go with a professional?
```

**Impact**: Led to implementing the crisis detection system with pattern matching for suicidal ideation and a dedicated crisis response that displays mental health resources (988 Suicide & Crisis Lifeline, Crisis Text Line, etc.) instead of movie quotes.

---

### 4. Emotional Mismatch Issue
```
I tried with "I'm very happy because I will meet with my family tonight"
✨ Here are some quotes that might resonate with you:
1. [0.56] "It's not who I am underneath, but what I do that defines me."
   — Batman (Batman Begins)
2. [0.49] "I'm going to make him an offer he can't refuse."
   — Don Vito Corleone (The Godfather)

I don't think that quote number 1 or 2 match the emotion
```

**Impact**: This feedback exposed that feature matching alone wasn't enough. It led to implementing the `areTonesCompatible` pre-filtering function that blocks quotes with incompatible emotional contexts (serious/philosophical quotes for joyful queries, threatening quotes for worried queries, etc.).

---

### 5. Worry Context Mismatches
```
For "My dog is sick, and I'm worried" I got:
1. [0.66] "I'm going to make him an offer he can't refuse."
   — Don Vito Corleone (The Godfather)
2. [0.32] "After all, tomorrow is another day!"
   — Scarlett O'Hara (Gone with the Wind)

The feeling is about a person that is concerned about the dog's health, 
I don't think that any of the quotes listed here are worthy for someone 
that is worried
```

**Impact**: Revealed that worried/health contexts needed special handling. This led to expanding the tone filtering to specifically block threatening, aggressive, and dismissive quotes for worried queries, and to enhance the lexicon with "support", "hope", and "difficulty" themes.

---

## Reflections

### What worked well?
- Generating documentation was super easy.
- Providing the base requirements was easy to get the baseline project.
- There is a "conversation" going on, and the agent is getting the feedback and adjust the implementation in seconds.


### What was challenging or required multiple attempts?
- When testing, I was trying to get more "accurate" or "helpful" quotes, but at the end we were limited by the data in the .json, so even though I was asking to make adjustments in the implementation, never gets where I was trying to go.

## Instructions
Go to the [Project README](./ai-quote-engine/README.md)

## Note
- I didn't work in the user engagement feature, I totally forgot, but I decided to add the Crisis Safety Feature