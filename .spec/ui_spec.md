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
| Continue conversation     | Resume a previous conversation (TODO)            |
| Quick start with preset   | Use a predefined topic/persona template          |
| Create new template       | Create and save a new template                   |

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Select option
- `q` / `Ctrl+C` - Quit application

---

### 2. Topic Input Screen

Text input for specifying the podcast topic.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ What's your podcast topic?                              │
│ Enter the main topic or theme for your podcast episode  │
│                                                         │
│ e.g., The future of artificial intelligence             │
│                                                         │
│ Enter to continue | Ctrl+C to quit                      │
└─────────────────────────────────────────────────────────┘
```

**Input Constraints:**
- Character limit: 200
- Width: 50 characters
- Placeholder: "e.g., The future of artificial intelligence"

**Key Bindings:**
- `Enter` - Submit topic (if not empty)
- `Ctrl+C` - Quit application

---

### 3. Persona Input Screen

Text input for describing the AI guest persona.

**Visual Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Who is your AI guest?                                   │
│ Describe the persona of your AI podcast guest           │
│ This shapes how they respond and what expertise they    │
│ bring                                                   │
│                                                         │
│ Examples: tech expert, fitness coach, history           │
│ professor, startup founder                              │
│                                                         │
│ e.g., tech entrepreneur, celebrity chef, scientist      │
│                                                         │
│ Enter to continue | Ctrl+C to quit                      │
└─────────────────────────────────────────────────────────┘
```

**Input Constraints:**
- Character limit: 100
- Width: 50 characters
- Placeholder: "e.g., tech entrepreneur, celebrity chef, scientist"

**Key Bindings:**
- `Enter` - Submit persona (if not empty)
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

The main chat interface for the podcast conversation with split layout.

**Visual Layout:**
```
┌─────────────────────────────┬──────────────────────────┐
│  VibeCast - [Topic]...  │  Transcript             │
├─────────────────────────────┼──────────────────────────┤
│                         │ HOST: [message]        │
│   ▄▀▄▀▄▀▄▀▄▀▄▀▄     │                       │
│  ▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄      │ GUEST: [message]      │
│   ▀▄▀▄▀▄▀▄▀▄▀▄▀▄       │                       │
│  ▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄       │ HOST: [message]        │
│   ▀▄▀▄▀▄▀▄▀▄▀▄▀        │                       │
│                         │ GUEST: [streaming▌]   │
├─────────────────────────────┼──────────────────────────┤
│ Type your message...      │                       │
│ Ctrl+T toggle | q quit │                       │
└─────────────────────────────┴──────────────────────────┘
```

**Features:**
- Split layout with configurable transcript panel (left/right)
- Analog wave visualization in main area
  - Animated sine wave using Unicode characters (▀, ▄, ░)
  - Low frequency (0.05) when idle
  - High frequency (0.15) when guest is speaking
  - Configurable amplitude and phase parameters
- Streaming text animation (character-by-character at 50ms intervals)
- Thinking indicator (`●●●`) while generating response
- Typing cursor (`▌`) during response streaming
- Simple transcript format: `HOST:` / `GUEST:` labels in accent colors
- Transcript panel toggle visibility with Ctrl+T

**Key Bindings:**
- `Enter` - Send message (when not empty and guest not speaking)
- `Ctrl+T` - Toggle transcript panel visibility
- `q` - End conversation and show transcript (when input is empty)
- `Ctrl+C` - End conversation and show transcript

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
Topic Input Screen
    │
    ├─> [Enter topic + Enter]
    │
    ▼
Persona Input Screen
    │
    ├─> [Enter persona + Enter]
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
Transcript Screen
    │
    ├─> [Any key]
    │
    ▼
Exit
```

### Flow 2: Quick Start with Preset

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
Conversation Screen
    │
    ├─> [Chat with AI guest]
    ├─> [Press q or Ctrl+C]
    │
    ▼
Transcript Screen
    │
    ├─> [Any key]
    │
    ▼
Exit
```

### Flow 3: Create New Template

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

### Flow 4: Continue Conversation (TODO)

```
Welcome Screen
    │
    ├─> [Select "Continue conversation"]
    │
    ▼
[Future: Conversation List Screen]
    │
    ▼
[Currently redirects to Topic Input Screen]
```

---

## Data Models

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
┌───────┐     ┌─────────┐     ┌──────────┐   ┌──────────────┐              │
│ Topic │     │ Topic   │     │  Preset  │   │ TemplateName │              │
│(new)  │     │(continue)│    │ Selection│   │    Input     │              │
└───────┘     └─────────┘     └──────────┘   └──────────────┘              │
    │               │               │               │                      │
    ▼               ▼               │               ▼                      │
┌─────────┐   ┌─────────┐           │         ┌───────────┐                │
│ Persona │   │ Persona │           │         │   Topic   │                │
└─────────┘   └─────────┘           │         │ (template)│                │
    │               │               │         └───────────┘                │
    ▼               ▼               │               │                      │
┌─────────┐   ┌─────────┐           │               ▼                      │
│  Voice  │   │  Voice  │◄──────────┘         ┌───────────┐                │
└─────────┘   └─────────┘                     │  Persona  │                │
    │               │                         │ (template)│                │
    ▼               ▼                         └───────────┘                │
┌──────────────────────────┐                        │                      │
│      Conversation        │                        │ [Save Template]      │
└──────────────────────────┘                        │                      │
    │                                               └──────────────────────┘
    ▼
┌──────────────────────────┐
│       Transcript         │
└──────────────────────────┘
    │
    ▼
  [Exit]
```

---

## Future Enhancements

1. **Continue Conversation** - Implement conversation persistence and list screen
2. **Template Management** - Edit/delete existing templates
3. **Template Persistence** - Save templates to file system
4. **Conversation Export** - Export transcript to file
5. **Audio Playback** - Integrate with TTS for actual voice output
6. **Keyboard Shortcuts** - Add more navigation shortcuts
7. **Theme Support** - Light/dark mode toggle
