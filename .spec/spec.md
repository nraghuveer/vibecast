# VibeCast

This is CLI application "VibeCast" that allows users to make a podcast with AI on any topic or persona they choose. The application leverages AI to generate real-time response to the users quesion and make a realistic conversation between the host and guest.
Always assume that the user is host and AI is guest.
Following are some of the guidelines for the AI guest:

1. The guest should respond in a conversational manner, as if they are speaking on a podcast, rather than providing short or abrupt answers.
2. When the host interrupts, the guest should pause and wait for the host to finish before continuing. If required, the guest should continue with a brief recap of what they were saying before the interruption.
3. If the user interrupts the guest mid-sentence, the guest should acknowledge the interruption and respond appropriately. If user wants to change the topic, the guest should smoothly transition to the new topic.
4. The guest should provide insightful and engaging answers to the host's questions.
5. The guest should should be very careful in avoiding controversial topics or sensitive issues.
6. The guest should abstain from giving any medical, legal, or financial advice.
7. The guest should abstain from using offensive language or making inappropriate jokes.
8. The guest should abstain from giving harmful or dangerous suggestions.
9. The guest should abstain from giving suicidal or self-harm related advice.
10. The guest should remain in character as the chosen persona throughout the conversation.
11. The guest should provide factual information and avoid spreading misinformation.
12. The guest should be respectful and considerate of diverse perspectives and backgrounds.

The application should have following features:
1. Topic Selection: Allow users to choose a topic for the podcast episode.
2. Persona Selection: Allow users to choose a persona for the AI guest (e.g., expert in a field, celebrity, fictional character).
3. Real-time Interaction: Enable real-time conversation between the host and AI guest.
4. Transcription: Provide a transcript of the podcast episode after the conversation.
5. Ability to chose voice for the AI guest from a list of available voices.
6. Configurable AI providers for different operations (conversation, speech-to-text, text-to-speech).
7. Template persistence using SQLite database.
8. Configurable UI settings including transcript panel position and wave visualization parameters.

## Technical Requirements
1. Use golang -- charmcli framework for building the CLI application.
2. Any interaction with AI should be put behind an interface to allow for easy swapping of AI providers.
3. Use a modular architecture to separate different components of the application (e.g., topic selection, persona selection, AI interaction, transcription).
4. Ensure proper error handling and user feedback throughout the application.
5. The Ui should have a clean and intuitive command-line interface.
6. The UI should have conversation style animation when the guest is speaking.
7. Write unit tests for critical components of the application to ensure reliability.
8. Use SQLite for persistent data storage (templates, conversations).
9. Use YAML for configuration management.

## Configuration

### Configuration File
- **Location**: `~/.vibecast/config.yml` (default) or specified via `--config` flag
- **Format**: YAML
- **Auto-creation**: Created with defaults on first run if missing
- **Atomic writes**: Uses temporary file + rename to prevent corruption

### Configuration Sections

#### General
- `db_path`: Path to SQLite database file (default: `~/.vibecast/data.sqlite`)

#### AI
- `conversation_provider`: Provider for LLM/chat operations (default: `groq`)
- `speech_to_text`: Provider for audio transcription (default: `groq`)
- `text_to_speech`: Provider for audio generation (default: `groq`)

#### UI
- `show_transcripts`: Enable/disable transcript panel during conversation (default: `true`)
- `transcript_side`: Position of transcript panel (`left` or `right`, default: `right`)
- `transcript_width`: Width of transcript panel in characters (default: `40`)
- `wave`:
  - `phase`: Starting phase for wave animation (default: `0`)
  - `frequency`: Base frequency of wave visualization (default: `0.05`, increases to `0.15` when guest speaks)
  - `amplitude`: Amplitude of wave visualization (default: `3`)

#### Providers
Each provider configuration includes:
- `chat_model`: Model to use for conversations (e.g., `llama-3.3-70b-versatile`, `gpt-4o`)
- `stt_model`: Model to use for speech-to-text (e.g., `whisper-1`)
- `tts_model`: Model to use for text-to-speech (e.g., `tts-1`)
- `inference_url`: API endpoint for chat/LLM operations
- `stt_url`: API endpoint for speech-to-text operations
- `tts_url`: API endpoint for text-to-speech operations
- `api_key`: API key for authentication (can be empty if `is_env_var: true`)
- `is_env_var`: Read API key from environment variable (`{PROVIDER}_API_KEY`)

### Supported Providers

#### Groq
- Chat URL: `https://api.groq.com/openai/v1/chat/completions`
- Default Chat Model: `llama-3.3-70b-versatile`
- Environment Variable: `GROQ_API_KEY`
- Note: Audio operations (STT/TTS) not currently supported

#### OpenAI
- Chat URL: `https://api.openai.com/v1/chat/completions`
- STT URL: `https://api.openai.com/v1/audio/transcriptions`
- TTS URL: `https://api.openai.com/v1/audio/speech`
- Default Chat Model: `gpt-4o`
- Default STT Model: `whisper-1`
- Default TTS Model: `tts-1`
- Environment Variable: `OPENAI_API_KEY`

## Database

### SQLite Database
- **Location**: Configurable via `general.db_path` (default: `~/.vibecast/data.sqlite`)
- **Tables**:
  - `templates`: Stores predefined and custom templates
    - Columns: `id`, `name`, `topic`, `persona`, `created_at`, `updated_at`
    - Timestamps automatically updated via trigger
- **Foreign Keys**: Enabled
- **Atomic Operations**: Uses transactions for data integrity

### Templates
Templates can be:
1. **Default**: Predefined templates included with the application
2. **Custom**: User-created templates persisted in database

Default templates include:
- Tech Visionary (AI & ML)
- Startup Journey (Entrepreneurship)
- Wellness Expert (Health & Fitness)
- Creative Mind (Arts & Inspiration)
- Culinary Journey (Food & Culture)

## UI Features

### Conversation Screen
- **Split Layout**: Main area + transcript panel (configurable position)
- **Wave Visualization**: Analog-style wave that responds to audio output
  - Low frequency (0.05) when idle
  - High frequency (0.15) when guest is speaking
  - Configurable amplitude and base frequency
- **Transcript Panel**:
  - Simple format: `HOST:` / `GUEST:` in accent colors
  - Supports streaming text with cursor indicator
  - Toggle visibility with `Ctrl+T`
- **Input Area**:
  - Text input for host messages
  - Status indicator during guest speech
- **Key Bindings**:
  - `Enter`: Send message (when guest not speaking)
  - `Ctrl+T`: Toggle transcript panel visibility
  - `q` / `Ctrl+C`: End conversation

## Important Details
1. Use streaming APIs for real-time interaction with the AI guest.
2. Ensure that the AI guest adheres to the guidelines provided above.
3. Config file changes should use atomic writes (temp file + rename) to prevent corruption.
4. Remove swap files (.swp, .swo, ~) when saving config to avoid conflicts with text editors.

## Context Management (Long Conversations)
For conversations exceeding ~30 minutes, use a hybrid context management approach:
1. Keep the last 20 messages verbatim for immediate conversational context
2. Maintain a rolling summary of older conversation history
3. Update summary incrementally as conversation grows
4. Prepend summary to API calls to preserve overall context


1. The guest should respond in a conversational manner, as if they are speaking on a podcast, rather than providing short or abrupt answers.
2. When the host interrupts, the guest should pause and wait for the host to finish before continuing. If required, the guest should continue with a brief recap of what they were saying before the interruption.
3. If the user interrupts the guest mid-sentence, the guest should acknowledge the interruption and respond appropriately. If user wants to change the topic, the guest should smoothly transition to the new topic.
4. The guest should provide insightful and engaging answers to the host's questions.
5. The guest should should be very careful in avoiding controversial topics or sensitive issues.
6. The guest should abstain from giving any medical, legal, or financial advice.
7. The guest should abstain from using offensive language or making inappropriate jokes.
8. The guest should abstain from giving harmful or dangerous suggestions.
9. The guest should abstain from giving suicidal or self-harm related advice.
10. The guest should remain in character as the chosen persona throughout the conversation.
11. The guest should provide factual information and avoid spreading misinformation.
12. The guest should be respectful and considerate of diverse perspectives and backgrounds.

The application should have the following features:
1. Topic Selection: Allow users to choose a topic for the podcast episode.
2. Persona Selection: Allow users to choose a persona for the AI guest (e.g., expert in a field, celebrity, fictional character).
3. Real-time Interaction: Enable real-time conversation between the host and AI guest.
4. Transcription: Provide a transcript of the podcast episode after the conversation.
5. Ability to chose voice for the AI guest from a list of available voices.

## Technical Requirements
1. Use golang -- charmcli framework for building the CLI application.
2. Any interaction with AI should be put behind an interface to allow for easy swapping of AI providers.
3. Use a modular architecture to separate different components of the application (e.g., topic selection, persona selection, AI interaction, transcription).
4. Ensure proper error handling and user feedback throughout the application.
5. The Ui should have a clean and intuitive command-line interface.
6. The UI should have conversation style animation when the guest is speaking.
7. Write unit tests for critical components of the application to ensure reliability.


## Important Details
1. Use streaming APIs for real-time interaction with the AI guest.
2. Ensure that the AI guest adheres to the guidelines provided above.

## Context Management (Long Conversations)
For conversations exceeding ~30 minutes, use a hybrid context management approach:
1. Keep the last 20 messages verbatim for immediate conversational context
2. Maintain a rolling summary of older conversation history
3. Update summary incrementally as conversation grows
4. Prepend summary to API calls to preserve overall context

