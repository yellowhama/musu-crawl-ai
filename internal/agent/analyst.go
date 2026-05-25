package agent

import (
	"fmt"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

type Analyst struct {
	Client        *OllamaClient
	CustomPersona string
}

func (a *Analyst) Synthesize(query string, contents []string) (string, error) {
	fullContext := ""
	for i, c := range contents {
		content := c
		// If single source is huge, summarize it first locally
		if len(content) > 4000 {
			content = utils.Summarize(content, 10) // Get top 10 sentences
		}

		fullContext += fmt.Sprintf("--- Source %d ---\n%s\n\n", i+1, content)
		if len(fullContext) > 12000 { // Increased limit but still guarded
			break
		}
	}

	prompt := fmt.Sprintf(`You are a Research Analyst for "musu-crawl-ai". 

%s

Below are several sources collected for the query: "%s"

Sources:
%s

Your task:
1. Synthesize a comprehensive answer based ONLY on the provided sources.
2. Identify any specific information that is still missing to fully answer the query.

Format your output as:
ANSWER: [Your synthesized response]
MISSING: [List of specific information gaps, or "None"]`, a.CustomPersona, query, fullContext)

	return a.Client.Ask(prompt, false)
}
