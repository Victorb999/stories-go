package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"stories-go/internal/models"
	"stories-go/internal/repository"

	"github.com/go-chi/chi/v5"
)

type StoryHandler struct {
	repo *repository.StoryRepository
	log  *slog.Logger
}

func NewStoryHandler(repo *repository.StoryRepository, log *slog.Logger) *StoryHandler {
	return &StoryHandler{repo: repo, log: log}
}

// Routes mounts all story endpoints onto a new chi.Router.
func (h *StoryHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Post("/seed", h.Seed)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// StoryInput is the request body for Create and Update.
type StoryInput struct {
	Title       string      `json:"title"`
	CoverImage  string      `json:"cover_image"`
	Author      string      `json:"author"`
	Content     string      `json:"content"`
	AIGenerated bool        `json:"ai_generated"`
	Size        models.Size `json:"size"`
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func normalizeSize(s models.Size) models.Size {
	if s == models.SizeLarge {
		return models.SizeLarge
	}
	return models.SizeSmall
}

// ── handlers ─────────────────────────────────────────────────────────────────

// List godoc
//
//	GET /api/v1/stories?size=small|large&ai=true|false
func (h *StoryHandler) List(w http.ResponseWriter, r *http.Request) {
	f := repository.ListFilter{}
	if s := r.URL.Query().Get("size"); s != "" {
		f.Size = s
	}
	if ai := r.URL.Query().Get("ai"); ai != "" {
		b := ai == "true"
		f.AIGenerated = &b
	}

	stories, err := h.repo.List(r.Context(), f)
	if err != nil {
		h.log.Error("list stories", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if stories == nil {
		stories = []models.Story{}
	}
	writeJSON(w, http.StatusOK, stories)
}

// Create godoc
//
//	POST /api/v1/stories
func (h *StoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var inp StoryInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if inp.Title == "" || inp.Author == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and author are required")
		return
	}
	inp.Size = normalizeSize(inp.Size)

	story, err := h.repo.Create(r.Context(), &models.Story{
		Title:       inp.Title,
		CoverImage:  inp.CoverImage,
		Author:      inp.Author,
		Content:     inp.Content,
		AIGenerated: inp.AIGenerated,
		Size:        inp.Size,
	})
	if err != nil {
		h.log.Error("create story", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, story)
}

// GetByID godoc
//
//	GET /api/v1/stories/{id}   — also increments the view counter
func (h *StoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	story, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("get story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if story == nil {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	writeJSON(w, http.StatusOK, story)
}

// Update godoc
//
//	PUT /api/v1/stories/{id}
func (h *StoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var inp StoryInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if inp.Title == "" || inp.Author == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and author are required")
		return
	}
	inp.Size = normalizeSize(inp.Size)

	story, err := h.repo.Update(r.Context(), id, &models.Story{
		Title:       inp.Title,
		CoverImage:  inp.CoverImage,
		Author:      inp.Author,
		Content:     inp.Content,
		AIGenerated: inp.AIGenerated,
		Size:        inp.Size,
	})
	if err != nil {
		h.log.Error("update story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if story == nil {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	writeJSON(w, http.StatusOK, story)
}

// Delete godoc
//
//	DELETE /api/v1/stories/{id}
func (h *StoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	deleted, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		h.log.Error("delete story", "id", id, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, "story not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Seed godoc
//
//	POST /api/v1/stories/seed
//
// Inserts sample stories when the database is empty.
func (h *StoryHandler) Seed(w http.ResponseWriter, r *http.Request) {
	existing, err := h.repo.List(r.Context(), repository.ListFilter{})
	if err != nil {
		h.log.Error("seed list", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if len(existing) > 0 {
		writeJSON(w, http.StatusOK, map[string]string{"message": "database already has stories, skipping seed"})
		return
	}

	seeds := []models.Story{
		{
			Title:      "Luffy e o Tesouro dos Sete Mares",
			Author:     "Chat GPT",
			CoverImage: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRl8-FEUSPHKpdfiKIfQwCDDZKXe83QaZvZFw&s",
			Content: `Era uma vez um menino chamado Luffy, que morava em uma vila à beira do oceano. Desde pequeno, ele sonhava em se tornar o Rei dos Piratas e encontrar o lendário Tesouro dos Sete Mares, chamado One Piece.
Um dia, Luffy conheceu um grande pirata chamado Shanks, que sempre contava histórias sobre aventuras incríveis. Mas o que Luffy mais amava nele era seu chapéu de palha, que parecia ter um brilho dourado quando o sol batia. Antes de partir em sua própria jornada, Shanks entregou o chapéu a Luffy e disse: “Seja um grande pirata e um dia nos encontraremos de novo!”
Mas algo mágico aconteceu com Luffy: ele comeu uma Fruta Encantada, chamada Fruta do Borrachudo, e seu corpo ganhou um poder incrível! Agora, ele podia esticar os braços, as pernas e até pular mais alto que os mastros dos navios. Com seu novo poder e sua coragem sem fim, ele decidiu construir sua própria tripulação e zarpar pelos oceanos mágicos.
No caminho, Luffy encontrou grandes amigos: Zoro, o espadachim de três espadas que podia cortar até o vento; Nami, a navegadora que lia os mapas como se falassem com ela; Usopp, o arqueiro que contava histórias tão incríveis que pareciam ganhar vida; Sanji, o cozinheiro que lutava com chutes flamejantes; e muitos outros companheiros leais. Juntos, eles formaram os Piratas do Chapéu de Palha!
Mas os mares eram cheios de desafios! Havia feiticeiros do mar, criaturas gigantescas e até piratas malvados que queriam impedir Luffy de alcançar seu sonho. O mais temido de todos era Barba Negra, um pirata das trevas que queria o One Piece para si.
Luffy e sua tripulação navegaram por ilhas flutuantes, cidades submarinas e até castelos escondidos nas nuvens. Cada lugar guardava segredos sobre o grande tesouro. No caminho, Luffy aprendeu que ser um pirata não significava apenas procurar ouro, mas proteger seus amigos, ajudar quem precisava e nunca desistir dos seus sonhos.
Certa vez, quando tudo parecia perdido e uma tempestade gigantesca ameaçava afundar seu navio, o Chapéu de Palha brilhou como um farol no escuro. Luffy percebeu que aquele não era um chapéu comum – ele carregava a promessa de todos os piratas destemidos antes dele. Com um grande soco de borracha, ele afastou a tempestade e seguiu em frente, pois sabia que sua jornada estava apenas começando.
E assim, Luffy e sua tripulação continuaram navegando, rumo ao horizonte dourado, em busca do maior tesouro de todos: a liberdade de viver suas próprias aventuras.
Fim.
`,
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:      "Ichigo e a Espada dos Espíritos",
			Author:     "Chat GPT",
			CoverImage: "https://static.zerochan.net/Kurosaki.Ichigo.full.2918.jpg",
			Content: `Em uma cidade tranquila, vivia um menino chamado Ichigo. Ele tinha cabelos cor de fogo e um dom especial: conseguia enxergar espíritos! Desde pequeno, ele ajudava fantasmas gentis a encontrarem o caminho para o além, mas nunca imaginou que seu destino seria muito maior.
Certa noite, enquanto as estrelas brilhavam como pequenos faróis no céu, uma guerreira misteriosa apareceu diante dele. Seu nome era Rukia, e ela era uma Guardiã Espiritual, protetora do Reino das Almas. “Nosso mundo está em perigo,” ela disse. “Sombras terríveis chamadas Hollows estão devorando as almas inocentes, e só um verdadeiro guardião pode detê-las.”
Antes que Ichigo pudesse entender tudo, um Hollow gigantesco apareceu, rugindo como uma fera das trevas. Para salvar sua família, Ichigo precisou aceitar um presente mágico de Rukia: uma espada encantada, forjada com a luz das almas puras. No momento em que segurou a lâmina, uma energia dourada brilhou ao seu redor, e ele se tornou um Guardião Espiritual!
Com sua nova força, Ichigo derrotou o monstro e protegeu sua cidade. Mas essa era apenas a primeira de muitas batalhas. Rukia ficou ao seu lado, ensinando-o a usar sua espada e os encantamentos sagrados dos guardiões.
No caminho, ele fez grandes amigos: Orihime, uma menina com poderes de cura brilhantes como o sol; Chad, um guerreiro forte como uma montanha; e Uryuu, um arqueiro que disparava flechas de luz. Juntos, enfrentaram perigos sombrios e desvendaram segredos sobre o Reino das Almas.
Mas um dia, algo terrível aconteceu. Rukia foi capturada por um rei cruel que governava o palácio celestial. Ele queria punir a guerreira por ter dado seus poderes a um humano. Sem hesitar, Ichigo e seus amigos embarcaram em uma jornada mágica para resgatá-la.
Eles cruzaram portais encantados, enfrentaram guardiões poderosos e, com coragem, desafiaram o próprio destino. Ichigo descobriu que sua espada não era apenas uma arma, mas uma chave para liberar um poder ainda maior dentro dele. Quando finalmente chegou ao palácio, sua energia brilhou como um relâmpago dourado, e ele libertou Rukia com um golpe de luz celestial.
Ao voltar para casa, Ichigo entendeu que sua missão nunca terminaria. Sempre haveria almas precisando de ajuda, sombras tentando se espalhar e amigos a proteger. Ele não era apenas um menino com uma espada mágica—era um verdadeiro herói do Reino das Almas.
E assim, com coragem no coração e sua espada brilhando como uma estrela, Ichigo seguiu sua jornada, garantindo que a luz sempre venceria a escuridão.
Fim.
`,
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:      "O Menino e a Raposa de Fogo",
			Author:     "Chat GPT",
			CoverImage: "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/8a0b71a4-d6c1-4653-ad4b-b8e8473bcbf0/db587ir-29f59704-af0c-4177-853e-8590f33168fe.png/v1/fill/w_1024,h_1024/narutinho_sticker_by_hanjorafael_db587ir-fullview.png?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9MTAyNCIsInBhdGgiOiIvZi84YTBiNzFhNC1kNmMxLTQ2NTMtYWQ0Yi1iOGU4NDczYmNiZjAvZGI1ODdpci0yOWY1OTcwNC1hZjBjLTQxNzctODUzZS04NTkwZjMzMTY4ZmUucG5nIiwid2lkdGgiOiI8PTEwMjQifV1dLCJhdWQiOlsidXJuOnNlcnZpY2U6aW1hZ2Uub3BlcmF0aW9ucyJdfQ.aJ5V4NsAUlkfAspfeagcvoFP7OTp1aIv1z9pZNRv_LM",
			Content: `Era uma vez um menino chamado Naruto, que vivia em uma vila encantada chamada Folha Escondida. Mas Naruto não era um menino comum. Dentro dele, adormecia uma criatura mágica: uma raposa de fogo com nove caudas, chamada Kurama. Essa raposa tinha um poder imenso e, no passado, havia causado uma grande tempestade de chamas na vila. Para protegê-los, o sábio líder da aldeia, o Hokage, selou o espírito da raposa dentro do pequeno Naruto quando ele ainda era um bebê.
Mas, por isso, muitas pessoas da vila olhavam para Naruto com medo, sem saber que, dentro dele, também havia um coração cheio de bondade. Ele cresceu sozinho, sem entender por que os outros o evitavam, mas nunca deixou de sorrir. Seu maior desejo era se tornar o Hokage, o guardião da vila, para que todos finalmente o vissem como um herói.
Naruto treinava todos os dias, aprendendo feitiços ninjas e artes secretas. Ele tinha um pergaminho antigo, onde estavam escritas palavras mágicas que o ajudavam a liberar sua energia especial, chamada chakra. No caminho, fez amigos leais: Sakura, uma menina que podia curar feridas com sua magia verde, e Sasuke, um jovem misterioso que controlava chamas negras e relâmpagos. Juntos, eles eram guiados pelo mestre Kakashi, um mago de capa prateada, que sempre os ensinava lições valiosas, como nunca abandonar os companheiros.
Certa noite, quando a lua brilhava azul no céu, Naruto sentiu um chamado dentro de si. Era Kurama, a raposa de fogo, que sussurrava em seus sonhos. "Você tem um grande poder, pequeno ninja, mas precisa aprender a confiar em mim." No começo, Naruto tinha medo, mas aos poucos percebeu que Kurama não era apenas uma fera selvagem – ela era parte dele.
Enquanto treinava, Naruto enfrentou feiticeiros sombrios e criaturas mágicas. Um dia, Sasuke decidiu partir para o reino das sombras, buscando mais poder. Naruto tentou impedi-lo, mas seu amigo estava determinado. Mesmo triste, ele fez uma promessa: “Eu trarei você de volta um dia.”
Com o tempo, Naruto se tornou um grande herói. Ele aprendeu a usar sua magia com sabedoria e, mais importante, conquistou a amizade de Kurama. Quando finalmente enfrentou seus maiores desafios, não estava mais sozinho. Seu coração brilhava como o sol, e sua chama interior se tornava cada vez mais forte.
Após muitas aventuras, Naruto realizou seu sonho e se tornou o grande guardião da vila. Agora, ao invés de medo, todos o olhavam com admiração. E a raposa de fogo? Ela não era mais uma fera selada, mas sua melhor amiga, sempre ao seu lado, protegendo a todos com sua luz dourada.
E assim, Naruto, Kurama e toda a vila viveram em paz, mostrando que a verdadeira magia não vem apenas do poder, mas do amor, da amizade e da coragem de nunca desistir.
Fim.
`,
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:       "A Sereiazinha do Rio",
			Author:      "Mamãe Contadora",
			CoverImage:  "https://www.tvtime.com/_next/image?url=https%3A%2F%2Fartworks.thetvdb.com%2Fbanners%2Fposters%2F265610-1.jpg&w=640&q=75",
			Content:     "Marina era uma sereinha que morava num rio cristalino no meio da floresta. Diferente das outras sereias que cantavam melodias tristes, Marina cantava músicas alegres que faziam os peixes dançarem e as flores das margens balançarem no ritmo. Certo dia, um passarinho perdido seguiu sua música e encontrou o caminho de volta para casa.",
			AIGenerated: true,
			Size:        models.SizeLarge,
		},
		{
			Title:       "Os Ursinhos e o Mel Mágico",
			Author:      "Vovô João",
			CoverImage:  "https://img.freepik.com/vetores-gratis/ilustracao-de-urso-marrom-bonitinho-sentado-com-pote-de-mel_1308-185130.jpg?semt=ais_hybrid&w=740&q=80",
			Content:     "Três ursinhos — Bolinha, Fofa e Pingo — descobriram uma colmeia escondida numa árvore encantada. O mel que saía dela era diferente: cada colherada fazia a pessoa lembrar de uma memória feliz. Juntos, os ursinhos decidiram compartilhar o mel com todos os animais tristes da floresta, espalhando alegria por onde passavam.",
			AIGenerated: false,
			Size:        models.SizeSmall,
		},
	}

	var created []models.Story
	for _, s := range seeds {
		s := s
		story, err := h.repo.Create(r.Context(), &s)
		if err != nil {
			h.log.Error("seed create", "title", s.Title, "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		created = append(created, *story)
	}
	writeJSON(w, http.StatusCreated, map[string]any{"seeded": len(created), "stories": created})
}
