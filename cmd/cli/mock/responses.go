package mock

import (
	"math/rand"
	"strings"
)

// Response represents a mock AI response
type Response struct {
	Text string
}

// Keywords mapped to specific responses
var keywordResponses = map[string][]string{
	"technology": {
		"That's a fascinating point about technology! I think we're at an inflection point where AI is becoming truly transformative. The key question is how we harness it responsibly.",
		"Technology has always been about augmenting human capability. From the printing press to smartphones, each leap has fundamentally changed how we connect and create.",
		"I've been thinking a lot about how technology shapes our daily routines. It's remarkable how quickly we adapt to new tools and forget life without them.",
	},
	"ai": {
		"AI is genuinely exciting territory. What I find most interesting is not the technology itself, but how it changes human creativity and problem-solving.",
		"The conversation around AI often focuses on replacement, but I see it more as amplification. The best outcomes come from human-AI collaboration.",
		"As someone who works in this space, I think the real magic happens when we use AI to tackle problems that were previously impossible to solve at scale.",
	},
	"food": {
		"Oh, cooking is such a creative outlet! There's something magical about transforming raw ingredients into something that brings people together.",
		"I believe food tells stories. Every dish has a history, a culture, and memories attached to it. That's what makes cooking so meaningful.",
		"The best meals I've had weren't necessarily in fancy restaurants. Sometimes a simple, well-made dish prepared with love is all you need.",
	},
	"music": {
		"Music has this incredible power to transport us through time. A single song can bring back vivid memories from decades ago.",
		"I think what makes music universal is that it speaks to emotions we all share, regardless of language or culture.",
		"The creative process in music fascinates me. How do artists take abstract feelings and translate them into sound?",
	},
	"startup": {
		"Building a startup is like running a marathon at sprint pace. The highs are incredible, but you need resilience for the inevitable lows.",
		"What I've learned from the startup world is that execution matters more than ideas. Everyone has ideas; few can make them reality.",
		"The startup ecosystem has changed dramatically. Today's founders have access to resources and knowledge that would have been unimaginable a decade ago.",
	},
	"health": {
		"I'm a big believer in preventive health. Small daily habits compound over time into significant outcomes.",
		"Mental and physical health are so interconnected. You really can't optimize one without considering the other.",
		"The wellness industry has exploded, but I think the fundamentals remain simple: sleep, movement, nutrition, and connection.",
	},
}

// Generic fallback responses
var genericResponses = []string{
	"That's a really thoughtful question. Let me share my perspective on this...",
	"I love that you brought this up! It's something I've been reflecting on lately.",
	"You know, that reminds me of an interesting experience I had recently.",
	"That's such an important topic. I think there are multiple angles to consider here.",
	"Absolutely! And I'd add that there's often more nuance than people initially realize.",
	"Great point! I think what's often overlooked is the human element in all of this.",
	"I couldn't agree more. And building on that idea, I'd say...",
	"That's exactly the kind of conversation I love having. Here's my take...",
}

// Greeting responses for when conversation starts
var greetingResponses = []string{
	"Thanks so much for having me on the show! I'm really excited to dive into this topic with you.",
	"It's great to be here! I've been looking forward to this conversation.",
	"Hello! I'm thrilled to be on VibeCast. This is going to be fun!",
	"Hey there! Thanks for the invitation. I'm ready to chat about anything you want to explore.",
}

// GetResponse returns a mock response based on the user's input
func GetResponse(input string, isFirstMessage bool) string {
	if isFirstMessage {
		return greetingResponses[rand.Intn(len(greetingResponses))]
	}

	input = strings.ToLower(input)

	// Check for keywords
	for keyword, responses := range keywordResponses {
		if strings.Contains(input, keyword) {
			return responses[rand.Intn(len(responses))]
		}
	}

	// Return generic response
	return genericResponses[rand.Intn(len(genericResponses))]
}

// Voice represents a mock voice option
type Voice struct {
	ID          string
	Name        string
	Description string
}

// GetVoices returns the list of available mock voices
func GetVoices() []Voice {
	return []Voice{
		{ID: "alloy", Name: "Alloy", Description: "neutral"},
		{ID: "echo", Name: "Echo", Description: "male"},
		{ID: "fable", Name: "Fable", Description: "British"},
		{ID: "onyx", Name: "Onyx", Description: "deep male"},
		{ID: "nova", Name: "Nova", Description: "female"},
		{ID: "shimmer", Name: "Shimmer", Description: "soft female"},
	}
}
