// Package crawlai is the public, importable facade over musu-crawl-ai.
//
// crawl-ai 의 알맹이(오케스트레이터·하베스터·프로세서)는 internal/ 에 봉인돼 외부
// 모듈에서 import 할 수 없다. 이 패키지는 같은 모듈 안에 있어 internal/ 을 호출할 수
// 있고, 외부(예: musu-website-co)가 in-process 로 쓸 수 있는 얇은 공개 API 만 노출한다.
//
// 사용:
//
//	c := crawlai.New(crawlai.Config{WikiDir: "./wiki", Project: "nongjida"})
//	md, _ := c.Fetch(ctx, "web", "https://www.law.go.kr/...")   // 수집(LLM 불필요)
//	hits, _ := c.Search(ctx, "농지법 임대차", false)              // 로컬 wiki 키워드 검색
//	if c.LLMReady(ctx) {                                          // Ollama 살아있을 때만
//	    rpt, _ := c.Research(ctx, "농지 임대 위탁 조건", 2)        // 계획→수집→합성 루프
//	}
//
// 주의(폴백 설계): Fetch 와 키워드 Search 는 LLM(Ollama) 없이 동작한다. Research 와
// semantic Search 는 Ollama 의존이므로 호출 전 LLMReady() 로 게이트하고, 꺼져 있으면
// Fetch+키워드로 폴백하라. 어떤 경우에도 호출측 발행을 막지 않는 게 원칙이다.
package crawlai

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yellowhama/musu-crawl-ai/internal/agent"
)

// Config 는 in-process crawl-ai 의 동작 설정. 빈 필드는 New 가 기본값으로 채운다.
type Config struct {
	WikiDir       string // 수집물 저장 디렉토리 (기본 "./wiki")
	Project       string // 테넌트/프로젝트 스코프 (기본 "default")
	AIBaseURL     string // OpenAI 호환 LLM 엔드포인트 (기본 Ollama "http://localhost:11434/v1")
	AIModel       string // 채팅/계획 모델 (기본 "llama3")
	SearchBaseURL string // 웹 검색 베이스(테스트용 오버라이드, 빈값이면 DuckDuckGo)
}

// Result 는 검색 히트의 공개 표현(내부 processor.IndexEntry 의 외부 노출용 사본).
type Result struct {
	ID      string
	Title   string
	Source  string
	Project string
	Summary string
}

// Client 는 crawl-ai 오케스트레이터를 감싼 in-process 클라이언트.
type Client struct {
	orch *agent.Orchestrator
	cfg  Config
}

// New 는 Config 로 Client 를 만든다. 내부 LoadConfig 가 env(MUSU_*) 로 AI 엔드포인트를
// 읽으므로, 여기서 cfg 값을 해당 env 에 주입한다(프로세스 전역 — 단일 설정 가정).
func New(cfg Config) *Client {
	if cfg.WikiDir == "" {
		cfg.WikiDir = "./wiki"
	}
	if cfg.Project == "" {
		cfg.Project = "default"
	}
	if cfg.AIBaseURL == "" {
		cfg.AIBaseURL = "http://localhost:11434/v1"
	}
	if cfg.AIModel == "" {
		cfg.AIModel = "llama3"
	}
	// internal/utils.LoadConfig 가 viper AutomaticEnv(MUSU_ 접두) 로 읽는 키들.
	os.Setenv("MUSU_AI_URL", cfg.AIBaseURL)
	os.Setenv("MUSU_MODEL", cfg.AIModel)
	os.Setenv("MUSU_OUT", cfg.WikiDir)
	if cfg.SearchBaseURL != "" {
		os.Setenv("MUSU_SEARCH_BASE_URL", cfg.SearchBaseURL)
	}
	return &Client{orch: agent.NewOrchestrator(cfg.WikiDir), cfg: cfg}
}

// Fetch 는 한 소스(web/yt/gh/arxiv/reddit/x/hf)를 마크다운으로 수집해 로컬 wiki 에
// 저장하고 본문을 반환한다. LLM 불필요.
func (c *Client) Fetch(ctx context.Context, source, id string) (string, error) {
	return c.orch.FetchAction(source, id, c.cfg.Project, "ko", c.cfg.AIModel)
}

// Search 는 로컬 wiki 를 검색한다. semantic=false 면 Bleve 키워드(LLM 불필요),
// true 면 벡터 시맨틱(Ollama 임베딩 필요).
func (c *Client) Search(ctx context.Context, query string, semantic bool) ([]Result, error) {
	entries, err := c.orch.SearchAction(query, c.cfg.Project, semantic, c.cfg.AIModel)
	if err != nil {
		return nil, err
	}
	out := make([]Result, 0, len(entries))
	for _, e := range entries {
		out = append(out, Result{ID: e.ID, Title: e.Title, Source: e.Source, Project: e.Project, Summary: e.Summary})
	}
	return out, nil
}

// Research 는 자율 멀티스텝 리서치(계획→발견→수집→합성→재귀)를 돌려 보고서를
// 반환한다. Ollama 의존 — 호출 전 LLMReady 로 게이트 권장. depth<=0 이면 2.
func (c *Client) Research(ctx context.Context, question string, depth int) (string, error) {
	if depth <= 0 {
		depth = 2
	}
	return c.orch.ResearchAction(question, c.cfg.Project, depth, c.cfg.AIModel)
}

// LLMReady 는 AIBaseURL(Ollama) 가 응답하는지 짧게 핑한다. Research/semantic 분기 게이트용.
// 실패(연결 불가/타임아웃/5xx)면 false → 호출측은 Fetch+키워드로 폴백해야 한다.
func (c *Client) LLMReady(ctx context.Context) bool {
	base := strings.TrimSuffix(c.cfg.AIBaseURL, "/v1")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base, nil)
	if err != nil {
		return false
	}
	cl := &http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := cl.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}

// WikiDir 은 수집물 저장 위치를 돌려준다(호출측이 brief 구성 시 경로 참조용).
func (c *Client) WikiDir() string { return c.cfg.WikiDir }
