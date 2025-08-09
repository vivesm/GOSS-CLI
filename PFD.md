⸻

PRD: Gemini CLI with GPT-OSS Backend

1. Overview

This project modifies the existing Gemini CLI to use a locally hosted GPT-OSS model (e.g., gpt-oss-20b in LM Studio) as the inference backend instead of Google’s Gemini API. The goal is to maintain the Gemini CLI’s command-line interface and user workflow while replacing the networked, API-key-based Google calls with local OpenAI-compatible API calls to LM Studio.

⸻

2. Goals
	•	Maintain CLI UX parity with Gemini CLI.
	•	Route all inference requests to LM Studio’s local API.
	•	Remove Google API key dependency.
	•	Support OpenAI-compatible request/response format.
	•	Enable offline usage without internet connectivity.

⸻

3. Non-Goals
	•	Adding new CLI features beyond API swap.
	•	Implementing multimodal support (Gemini’s image/video capabilities) unless GPT-OSS supports them.
	•	Maintaining compatibility with Gemini API after swap (this is a one-way fork).

⸻

4. Users & Use Cases

Users
	•	Developers who like the Gemini CLI interface but want local inference.
	•	Privacy-conscious users avoiding cloud AI services.
	•	AI hobbyists running large models locally.

Use Cases
	•	Local offline AI assistant via CLI.
	•	Prototyping automation scripts with a local model.
	•	Avoiding API costs by leveraging self-hosted models.

⸻

5. Requirements

Functional Requirements
	1.	API Endpoint Change
	•	All requests must be sent to a configurable LM Studio endpoint (default: http://localhost:1234/v1/chat/completions).
	•	Endpoint should be user-configurable via CLI flags or environment variables.
	2.	Request Format Conversion
	•	Convert Gemini request format:

{
  "contents": [
    { "role": "user", "parts": [{ "text": "Hello!" }] }
  ]
}

→ to OpenAI format:

{
  "model": "gpt-oss-20b",
  "messages": [
    { "role": "user", "content": "Hello!" }
  ]
}


	3.	Response Parsing
	•	Gemini: response.candidates[0].content.parts[0].text
	•	OpenAI: response.choices[0].message.content
	4.	Auth Removal
	•	Remove API key checks for Google.
	•	Skip token injection in headers.
	5.	Model Configuration
	•	Allow user to pass LM Studio model name (--model flag, default: gpt-oss-20b).
	6.	Error Handling
	•	Graceful error messages if LM Studio API is unreachable.
	•	Fallback instructions for starting LM Studio’s local server.

⸻

Non-Functional Requirements
	•	Performance: Requests should respond in ~1–3 seconds for local models (dependent on hardware).
	•	Portability: Should run on macOS, Linux, Windows.
	•	Maintainability: Code changes should be isolated so the CLI can be updated with minimal merge conflicts.

⸻

6. Technical Approach
	1.	Fork Gemini CLI repository.
	2.	Identify the API call layer (likely in a client.js or api.ts file).
	3.	Replace fetch() calls to Gemini API with OpenAI-style request to LM Studio.
	4.	Implement a request converter function.
	5.	Implement a response mapper.
	6.	Add CLI flag --api-base for endpoint customization.
	7.	Test against LM Studio running gpt-oss-20b.

⸻

7. Dependencies
	•	LM Studio running with API server enabled.
	•	Node.js (version per Gemini CLI requirements).
	•	node-fetch or equivalent HTTP client.

⸻

8. Risks & Mitigation

Risk	Mitigation
Model mismatch in prompt handling	Test with simple prompts and compare outputs
CLI feature incompatibility	Remove or disable unsupported features
Users running without LM Studio	Add clear error guidance


⸻

9. Deliverables
	•	Modified CLI with LM Studio backend.
	•	Documentation on setup and usage.
	•	Example .env file for local API configuration.

⸻

10. Success Criteria
	•	CLI commands run without requiring a Google API key.
	•	Inference runs entirely offline.
	•	Users can switch models by changing a CLI flag.

⸻
