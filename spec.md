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

