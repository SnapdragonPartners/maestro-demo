# maestro-demo
Testing Maestro development

## Development Environment

This project includes a containerized development environment with the following features:

### Container Specifications
- **Base Image**: Ubuntu 22.04 LTS
- **Languages**: Python 3 with pip
- **Development Tools**: Git, Vim, Nano, curl, wget
- **Build Tools**: GCC, G++, build-essential
- **User**: Non-root developer user with sudo access

### Quick Start

Build and run the development container:

```bash
# Build the container
docker build -t maestro-demo-dev .

# Or use docker-compose
docker-compose up -d dev
docker-compose exec dev bash

# Direct run
docker run -it --rm -v $(pwd):/workspace maestro-demo-dev
```

### Available Ports
- 3000: Web development
- 8000: Python development server  
- 8080: Alternative web server
- 9000: Additional services

### Files Created
- `Dockerfile`: Container definition
- `docker-compose.yml`: Container orchestration
- `.dockerignore`: Build optimization
- `README-container.md`: Detailed container documentation

The development environment is ready for application code deployment and testing.

## Questions File Format

The application requires a `questions.json` file at the repository root containing quiz questions.

### Schema

Each question in the JSON array must follow this structure:

```json
{
  "id": <integer>,           // Unique question identifier
  "question": "<string>",    // The question text
  "answers": [<strings>],    // Array of possible answers
  "correct": <integer>       // Index of the correct answer (0-based)
}
```

### Example

```json
[
  {
    "id": 1,
    "question": "What is the capital of France?",
    "answers": ["London", "Paris", "Berlin", "Madrid"],
    "correct": 1
  },
  {
    "id": 2,
    "question": "Which planet is known as the Red Planet?",
    "answers": ["Venus", "Jupiter", "Mars", "Saturn"],
    "correct": 2
  }
]
```

### Validation

- The `correct` index must be within the bounds of the `answers` array (0 ≤ correct < len(answers))
- The application validates this at startup and when loading questions
- Invalid questions will cause the application to fail with a descriptive error message

### Location

The `questions.json` file must be located at the repository root (same directory as `main.go`).

### Quiz Behavior

- The application loads questions from `questions.json` at startup or when starting a new quiz
- Each quiz session randomly selects exactly 3 questions (defined by `NumQuestions` constant)
- The selection is random for each new quiz session
- If fewer than 3 questions are available, all questions will be used
