# VibeCast CLI - UI Specification

## Overview

VibeCast is an AI-powered podcast companion CLI application built with Go using the Bubble Tea TUI framework. It enables users to create dynamic conversations with AI personas.

## Technology Stack

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm-architecture TUI)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (text input, key bindings)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) (terminal styling and layout)

## Color Palette

| Color     | Hex Code  | Usage                          |
|-----------|-----------|--------------------------------|
| Primary   | `#7C3AED` | Purple - titles, selections   |
| Secondary | `#10B981` | Green - guest messages        |
| Accent    | `#3B82F6` | Blue - host messages          |
| Muted     | `#6B7280` | Gray - help text, subtitles   |
| Error     | `#EF4444` | Red - error states            |

---

## Screens

### 1. Welcome Screen

The entry point of the application displaying the logo and main menu.

**Visual Layout:**
```
     ╦  ╦╦╔╗ ╔═╗╔═╗╔═╗╔═╗╔╦╗
     ╚╗╔╝║╠╩╗║╣ ║  ╠═╣╚═╗ ║
      ╚╝ ╩╚═╝╚═╝╚═╝╩ ╩╚═╝ ╩

   Your AI-powered podcast companion
   Create dynamic conversations with AI personas

   > Create new conversation
     Continue conversation
     Quick start with preset
     Create new template

   ↑/↓ to navigate • Enter to select • q to quit
```

**Options:**
| Option                    | Description                                      |
|---------------------------|--------------------------------------------------|
| Create new conversation   | Start a new conversation with custom topic/persona |
| Continue conversation     | Resume a previous conversation                   |
| Quick start with preset   | Use a predefined topic/persona template          |
| Create new template       | Create and save a new template                   |

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Select option
- `q` / `Ctrl+C` - Quit application

---

### 2. New Conversation Screen

A single unified screen for creating a new conversation with all required fields.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Create New Conversation                                  │
│ Fill in the details below to start your podcast          │
│ conversation                                             │
│                                                          │
│   Started: Jan 02, 2026 3:04 PM                          │
│                                                          │
│ > Title                                                  │
│   [e.g., AI Future Discussion, Tech Deep Dive...]        │
│                                                          │
│   Topic                                                  │
│   [e.g., The future of artificial intelligence]          │
│                                                          │
│   Persona                                                │
│   [e.g., tech entrepreneur, scientist, chef]             │
│                                                          │
│   Provider                                               │
│   [groq (llama-3.3-70b)] [openai (gpt-4o)]               │
│                                                          │
│ Tab/↓ next • Shift+Tab/↑ prev • ←/→ provider • Enter    │
│ continue • Esc back                                      │
└─────────────────────────────────────────────────────────┘
```

**Fields:**
| Field     | Description                                        |
|-----------|----------------------------------------------------|
| Started   | Auto-detected timestamp (read-only display)        |
| Title     | Name for the conversation (max 100 chars)          |
| Topic     | Main topic or theme for the podcast (max 200 chars)|
| Persona   | AI guest persona description (max 100 chars)       |
| Provider  | AI provider selection from configured providers    |

**Key Bindings:**
- `Tab` / `↓` - Move to next field
- `Shift+Tab` / `↑` - Move to previous field
- `←` / `→` - Select provider (when on provider field)
- `Enter` - Submit form (when on provider field) or move to next field
- `Esc` - Go back to welcome screen
- `Ctrl+C` - Quit application

---

### 3. Conversation List Screen

Displays a list of existing conversations for continuation.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Continue Conversation                                    │
│ Select a conversation to continue                        │
│                                                          │
│ > AI Future Discussion                                   │
│      Jan 02, 2026 3:04 PM                                │
│                                                          │
│   Tech Deep Dive                                         │
│      Jan 01, 2026 10:30 AM                               │
│                                                          │
│   Cooking Show Interview                                 │
│      Dec 31, 2025 2:15 PM                                │
│                                                          │
│ ↑/↓ or j/k navigate | Enter select | Ctrl+I show        │
│ details | Esc back                                       │
└─────────────────────────────────────────────────────────┘
```

**With Details (Ctrl+I toggled):**
```
┌─────────────────────────────────────────────────────────┐
│ Continue Conversation                                    │
│ Select a conversation to continue                        │
│                                                          │
│ > AI Future Discussion                                   │
│      Jan 02, 2026 3:04 PM                                │
│      Topic: The future of artificial intelligence        │
│      Persona: Silicon Valley tech founder                │
│                                                          │
│   Tech Deep Dive                                         │
│      Jan 01, 2026 10:30 AM                               │
│      Topic: Building scalable systems                    │
│      Persona: Software architect                         │
│                                                          │
│ ↑/↓ or j/k navigate | Enter select | Ctrl+I hide        │
│ details | Esc back                                       │
└─────────────────────────────────────────────────────────┘
```

**Display:**
- Title displayed in **Primary color** (purple, bold)
- Timestamp displayed in **Muted color** (gray)
- Topic and Persona shown only when Ctrl+I is toggled on

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Select conversation to continue
- `Ctrl+I` - Toggle topic/persona details visibility
- `Esc` - Go back to welcome screen
- `Ctrl+C` - Quit application

---

### 4. Voice Selection Screen

Selection list for choosing the AI guest's voice.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Select a voice for your AI guest                        │
│ Choose the voice that best fits your guest's persona    │
│                                                         │
│ > Alloy (neutral)                                       │
│   Echo (male)                                           │
│   Fable (British)                                       │
│   Onyx (deep male)                                      │
│   Nova (female)                                         │
│   Shimmer (soft female)                                 │
│                                                         │
│ ↑/↓ or j/k to navigate | Enter to select | Ctrl+C quit │
└─────────────────────────────────────────────────────────┘
```

**Available Voices:**
| Voice   | Description  |
|---------|--------------|
| Alloy   | neutral      |
| Echo    | male         |
| Fable   | British      |
| Onyx    | deep male    |
| Nova    | female       |
| Shimmer | soft female  |

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Select voice
- `Ctrl+C` - Quit application

---

### 5. Preset Selection Screen

Selection list for quick-starting with predefined templates.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Quick Start - Select a Template                         │
│ Choose a predefined topic and persona to get started    │
│ quickly                                                 │
│                                                         │
│ > Tech Visionary (The future of artificial intelli...)  │
│   Startup Journey (Building a successful startup fr...)│
│   Wellness Expert (Holistic health and modern welln...)│
│   Creative Mind (The creative process and finding i...)│
│   Culinary Journey (World cuisines and the stories ...) │
│                                                         │
│ ↑/↓ or j/k to navigate • Enter to select • Esc to go   │
│ back • Ctrl+C to quit                                   │
└─────────────────────────────────────────────────────────┘
```

**Default Templates:**
| Name             | Topic                                      | Persona                                                    |
|------------------|--------------------------------------------|------------------------------------------------------------|
| Tech Visionary   | The future of artificial intelligence      | A Silicon Valley tech founder with deep knowledge of AI    |
| Startup Journey  | Building a successful startup from scratch | An experienced entrepreneur who has built and sold companies |
| Wellness Expert  | Holistic health and modern wellness        | A wellness coach with expertise in nutrition and fitness   |
| Creative Mind    | The creative process and finding inspiration | A multi-disciplinary artist across music and visual arts  |
| Culinary Journey | World cuisines and the stories behind food | A chef and food writer who has traveled the world          |

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Select template
- `Esc` - Go back to welcome
- `Ctrl+C` - Quit application

---

### 6. Template Name Input Screen

Text input for naming a new custom template.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Create New Template                                     │
│ Give your template a memorable name                     │
│                                                         │
│ e.g., Tech Interview, Cooking Show...                   │
│                                                         │
│ Enter to continue • Esc to go back • Ctrl+C to quit    │
└─────────────────────────────────────────────────────────┘
```

**Input Constraints:**
- Character limit: 50
- Width: 50 characters
- Placeholder: "e.g., Tech Interview, Cooking Show..."

**Key Bindings:**
- `Enter` - Submit name (if not empty)
- `Esc` - Go back to welcome
- `Ctrl+C` - Quit application

---

### 7. Conversation Screen

The main chat interface for the podcast conversation.

**Visual Layout:**
```
 ╦  ╦╦╔╗ ╔═╗╔═╗╔═╗╔═╗╔╦╗
 ╚╗╔╝║╠╩╗║╣ ║  ╠═╣╚═╗ ║
  ╚╝ ╩╚═╝╚═╝╚═╝╩ ╩╚═╝ ╩

  HOST   Hello, welcome to the show!

  GUEST  Thank you for having me! I'm excited to
         discuss...▌

  ·········█▓░·

  Type your message...
  Enter to send | Ctrl+I show details | q or Ctrl+C to exit
```

**With Details (Ctrl+I toggled):**
```
 ╦  ╦╦╔╗ ╔═╗╔═╗╔═╗╔═╗╔╦╗
 ╚╗╔╝║╠╩╗║╣ ║  ╠═╣╚═╗ ║
  ╚╝ ╩╚═╝╚═╝╚═╝╩ ╩╚═╝ ╩

  Topic: The future of artificial intelligence
  Persona: Silicon Valley tech founder

  HOST   Hello, welcome to the show!

  GUEST  Thank you for having me!▌

  Type your message...
  Enter to send | Ctrl+I hide details | q or Ctrl+C to exit
```

**Features:**
- Streaming text animation (character-by-character at 50ms intervals)
- Flowing dots animation while generating response
- Typing cursor (`▌`) during response streaming
- Simple transcript format: `HOST:` / `GUEST:` labels in accent colors
- Toggle topic/persona details with Ctrl+I

**Key Bindings:**
- `Enter` - Send message (when not empty and guest not speaking)
- `Ctrl+I` - Toggle topic/persona details visibility
- `q` - End conversation and exit (when input is empty)
- `Ctrl+C` - End conversation and exit

---

### 8. Transcript Screen

Displays the full conversation transcript at the end.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Conversation Transcript                                 │
│ Topic: [topic]                                          │
│ Persona: [persona]                                      │
│ Voice: [voice name] ([description])                     │
│                                                         │
│ ──────────────────────────────────────────────────────  │
│                                                         │
│ Guest: [message]                                        │
│                                                         │
│ Host: [message]                                         │
│                                                         │
│ Guest: [message]                                        │
│                                                         │
│ ──────────────────────────────────────────────────────  │
│                                                         │
│ Press any key to exit                                   │
└─────────────────────────────────────────────────────────┘
```

**Key Bindings:**
- Any key - Exit application

---

## User Flows

### Flow 1: Create New Conversation

```
Welcome Screen
    │
    ├─> [Select "Create new conversation"]
    │
    ▼
New Conversation Screen (single screen with all fields)
    │
    ├─> [Enter title, topic, persona, select provider]
    ├─> [Press Enter on provider field]
    │
    ▼
Voice Selection Screen
    │
    ├─> [Select voice + Enter]
    │
    ▼
Conversation Screen
    │
    ├─> [Chat with AI guest]
    ├─> [Press q or Ctrl+C]
    │
    ▼
Exit
```

### Flow 2: Continue Conversation

```
Welcome Screen
    │
    ├─> [Select "Continue conversation"]
    │
    ▼
Conversation List Screen
    │
    ├─> [Browse conversations, Ctrl+I to see details]
    ├─> [Select conversation + Enter]
    │       OR
    ├─> [Esc] ─────────────────────┐
    │                              │
    ▼                              ▼
Conversation Screen          Welcome Screen
    │
    ├─> [Resume chat with AI guest]
    ├─> [Press q or Ctrl+C]
    │
    ▼
Exit
```

### Flow 3: Quick Start with Preset

```
Welcome Screen
    │
    ├─> [Select "Quick start with preset"]
    │
    ▼
Preset Selection Screen
    │
    ├─> [Select template + Enter]
    │       OR
    ├─> [Esc] ─────────────────────┐
    │                              │
    ▼                              ▼
Voice Selection Screen      Welcome Screen
    │
    ├─> [Select voice + Enter]
    │
    ▼
Provider Selection Screen
    │
    ├─> [Select provider + Enter]
    │
    ▼
Conversation Screen
    │
    ├─> [Chat with AI guest]
    ├─> [Press q or Ctrl+C]
    │
    ▼
Exit
```

### Flow 4: Create New Template

```
Welcome Screen
    │
    ├─> [Select "Create new template"]
    │
    ▼
Template Name Input Screen
    │
    ├─> [Enter name + Enter]
    │       OR
    ├─> [Esc] ─────────────────────┐
    │                              │
    ▼                              ▼
Topic Input Screen          Welcome Screen
    │
    ├─> [Enter topic + Enter]
    │
    ▼
Persona Input Screen
    │
    ├─> [Enter persona + Enter]
    │
    ▼
[Template Saved]
    │
    ▼
Welcome Screen
```

---

## Data Models

### Conversation
```go
type Conversation struct {
    ID        string     // Unique identifier
    Title     string     // Display name for the conversation
    Topic     string     // Podcast topic
    Persona   string     // AI guest persona
    VoiceID   string     // Selected voice ID
    VoiceName string     // Selected voice name
    Provider  string     // AI provider name
    CreatedAt time.Time  // Started timestamp
    EndedAt   *time.Time // Ended timestamp (nullable)
}
```

### Template
```go
type Template struct {
    ID      string  // Unique identifier
    Name    string  // Display name
    Topic   string  // Podcast topic
    Persona string  // AI guest persona
}
```

### Voice
```go
type Voice struct {
    ID          string  // Voice identifier
    Name        string  // Display name
    Description string  // Voice characteristics
}
```

### Message
```go
type Message struct {
    Content  string  // Message text
    IsHost   bool    // true = host, false = guest
    Complete bool    // true if message is fully rendered
}
```

---

## Screen State Machine

```
                    ┌──────────────────────────────────────────────────────┐
                    │                                                      │
                    ▼                                                      │
              ┌─────────────┐                                              │
              │   Welcome   │◄─────────────────────────────────────────────┤
              └─────────────┘                                              │
                    │                                                      │
    ┌───────────────┼───────────────┬───────────────┐                      │
    │               │               │               │                      │
    ▼               ▼               ▼               ▼                      │
┌───────────┐ ┌─────────────┐ ┌──────────┐   ┌──────────────┐              │
│    New    │ │Conversation │ │  Preset  │   │ TemplateName │              │
│Conversation│ │    List    │ │ Selection│   │    Input     │              │
└───────────┘ └─────────────┘ └──────────┘   └──────────────┘              │
    │               │               │               │                      │
    ▼               │               │               ▼                      │
┌─────────┐         │               │         ┌───────────┐                │
│  Voice  │◄────────┘               │         │   Topic   │                │
└─────────┘                         │         │ (template)│                │
    │                               │         └───────────┘                │
    ▼                               │               │                      │
┌──────────────────────────┐        │               ▼                      │
│      Conversation        │◄───────┤         ┌───────────┐                │
│  (shows title in header) │        │         │  Persona  │                │
└──────────────────────────┘        │         │ (template)│                │
    │                               │         └───────────┘                │
    ▼                               │               │                      │
  [Exit]                            │               │ [Save Template]      │
                                    │               └──────────────────────┘
                                    │
                               ┌─────────┐
                               │  Voice  │
                               └─────────┘
                                    │
                                    ▼
                               ┌─────────┐
                               │Provider │
                               └─────────┘
                                    │
                                    ▼
                            ┌──────────────┐
                            │ Conversation │
                            └──────────────┘
```

---

## Keyboard Shortcuts Summary

| Shortcut    | Context              | Action                           |
|-------------|----------------------|----------------------------------|
| `↑` / `k`   | Lists, Forms         | Navigate up / Previous field     |
| `↓` / `j`   | Lists, Forms         | Navigate down / Next field       |
| `Tab`       | New Conversation     | Next field                       |
| `Shift+Tab` | New Conversation     | Previous field                   |
| `←` / `→`   | New Conversation     | Select provider                  |
| `Enter`     | All screens          | Confirm / Submit                 |
| `Esc`       | Sub-screens          | Go back to previous screen       |
| `Ctrl+I`    | List & Conversation  | Toggle topic/persona details     |
| `q`         | Conversation         | Exit (when input empty)          |
| `Ctrl+C`    | All screens          | Quit application                 |

---

## Future Enhancements

1. **Template Management** - Edit/delete existing templates
2. **Template Persistence** - Save templates to file system
3. **Conversation Export** - Export transcript to file
4. **Audio Playback** - Integrate with TTS for actual voice output
5. **Theme Support** - Light/dark mode toggle
6. **Search** - Search conversations by title, topic, or content
