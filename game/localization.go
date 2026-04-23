package game

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	LanguageEnglish    = "en"
	LanguagePortuguese = "pt"
	GlobalLanguageKey  = "game_language"
)

func LoadLanguage(store *Store) string {
	if store == nil {
		return LanguageEnglish
	}
	if raw, err := store.GetSetting(GlobalLanguageKey); err == nil {
		lang := NormalizeLanguage(raw)
		if strings.TrimSpace(raw) != "" {
			return lang
		}
	}
	return ""
}

func EnsureLanguageConfigured(reader *bufio.Reader, store *Store) string {
	lang := LoadLanguage(store)
	if strings.TrimSpace(lang) != "" {
		return lang
	}

	fmt.Println("Choose language / Escolha o idioma:")
	fmt.Println("1) English")
	fmt.Println("2) Português")
	fmt.Print("Selection / Selecione: ")
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer == "2" || strings.HasPrefix(answer, "p") {
		lang = LanguagePortuguese
	} else {
		lang = LanguageEnglish
	}
	_ = SaveLanguage(store, lang)
	return lang
}

func SaveLanguage(store *Store, lang string) error {
	if store == nil {
		return nil
	}
	return store.SetSetting(GlobalLanguageKey, NormalizeLanguage(lang))
}

func TranslateText(lang, text string) string {
	if NormalizeLanguage(lang) != LanguagePortuguese {
		return text
	}

	if translated, ok := exactTranslations[text]; ok {
		return translated
	}
	for _, rule := range translationRules {
		if strings.HasPrefix(text, rule.prefix) {
			return rule.apply(text)
		}
	}
	return text
}

func TranslateRoomDescription(lang, roomID, fallback string) string {
	if NormalizeLanguage(lang) != LanguagePortuguese {
		return fallback
	}
	if translated, ok := roomTranslations[roomID]; ok {
		return translated
	}
	return fallback
}

func TranslateActionDescription(lang, desc string) string {
	if NormalizeLanguage(lang) != LanguagePortuguese {
		return desc
	}

	trimmed := strings.TrimSpace(desc)
	for _, rule := range actionRules {
		if strings.HasPrefix(trimmed, rule.prefix) {
			return rule.apply(trimmed)
		}
	}
	if translated, ok := exactTranslations[trimmed]; ok {
		return translated
	}
	return desc
}

func FormatWiseManProvider(lang, provider string) string {
	if NormalizeLanguage(lang) != LanguagePortuguese {
		return provider
	}
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case WiseManProviderGemini:
		return "Gemini"
	case WiseManProviderCloudflare:
		return "Cloudflare"
	case WiseManProviderLegacy:
		return "endereço legado"
	default:
		return provider
	}
}

type translationRule struct {
	prefix string
	apply  func(string) string
}

var exactTranslations = map[string]string{
	"English":                            "Inglês",
	"Português":                          "Português",
	"Settings":                           "Configurações",
	"Selected":                           "Selecionado",
	"-- EMPTY SLOT --":                   "-- SLOT VAZIO --",
	"Lvl":                                "Nv",
	"Inbox":                              "Caixa de entrada",
	"unread":                             "não lido",
	"None":                               "Nenhum",
	"Steel Sword":                        "Espada de Aço",
	"Garden Kit":                         "Kit de Jardim",
	"Meat":                               "Carne",
	"Water":                              "Água",
	"Fur Coat":                           "Casaco de Pele",
	"Sun Amulet":                         "Amuleto Solar",
	"Ice Crystal":                        "Cristal de Gelo",
	"Forsaken Blade":                     "Lâmina Abandonada",
	"Moon Pearl":                         "Pérola da Lua",
	"Sage Blessing":                      "Bênção do Sábio",
	"Ruins Token":                        "Ficha das Ruínas",
	"Moon Charm":                         "Amuleto da Lua",
	"Star Map":                           "Mapa Estelar",
	"Triptych Blessing":                  "Bênção do Tríptico",
	"Monster Hunt":                       "Caçada aos Monstros",
	"Void Expedition":                    "Expedição do Vazio",
	"Frost Pact":                         "Pacto Gélido",
	"Sun Covenant":                       "Aliança Solar",
	"active":                             "ativa",
	"completed":                          "concluída",
	"cooldown":                           "recarga",
	"Movement":                           "Movimento",
	"Room":                               "Sala",
	"Multiplayer":                        "Multijogador",
	"Party":                              "Party",
	"Exploration":                        "Exploração",
	"Guidance":                           "Orientação",
	"Change language":                    "Mudar idioma",
	"Configure wise man AI":              "Configurar IA do sábio",
	"Reconfigure wise man AI":            "Reconfigurar IA do sábio",
	"Ask the wise man for guidance":      "Pedir orientação ao sábio",
	"Send room chat":                     "Enviar chat da sala",
	"Broadcast a room event":             "Anunciar um evento na sala",
	"Whisper a nearby player":            "Sussurrar para um jogador próximo",
	"Create a party":                     "Criar uma party",
	"View party":                         "Ver party",
	"View party quest log":               "Ver registro da quest da party",
	"Follow party leader":                "Seguir líder da party",
	"Enter the party dungeon":            "Entrar na masmorra da party",
	"Share a party heal":                 "Compartilhar cura da party",
	"Rally the party":                    "Reunir a party",
	"Raise party guard":                  "Levantar guarda da party",
	"Open party quest board":             "Abrir quadro de quests da party",
	"Leave party":                        "Sair da party",
	"Invite nearby player to party":      "Convidar jogador próximo para a party",
	"Refresh room chat":                  "Atualizar chat da sala",
	"Check whispers":                     "Ver sussurros",
	"Open notifications inbox":           "Abrir caixa de notificações",
	"Search for a relic":                 "Procurar um relicário",
	"Trade a relic":                      "Trocar um relicário",
	"Ask the wise man":                   "Perguntar ao sábio",
	"Online":                             "Online",
	"connected":                          "conectado",
	"offline":                            "offline",
	"Recent room chat:":                  "Chat recente da sala:",
	"event":                              "evento",
	"Recent whispers:":                   "Sussurros recentes:",
	"unread notification(s)":             "notificação(ões) não lida(s)",
	"led by":                             "liderada por",
	"Members":                            "Membros",
	"Quest":                              "Quest",
	"Phase":                              "Fase",
	"Reward":                             "Recompensa",
	"gold":                               "ouro",
	"Cooldown":                           "Recarga",
	"quest board recovers for":           "o quadro de quests se recupera por",
	"Party support":                      "Apoio da party",
	"damage per hit":                     "dano por golpe",
	"damage for":                         "dano por",
	"turn(s)":                            "turno(s)",
	"HP":                                 "HP",
	"Rally":                              "Reunir",
	"Guard":                              "Guarda",
	"for":                                "por",
	"Selection: ":                        "Seleção: ",
	"Language settings are unavailable.": "As configurações de idioma não estão disponíveis.",
	"Wise man setup failed":              "Falha ao configurar o sábio",
	"The wise man quietly points at":     "O sábio aponta silenciosamente para",
	"The wise man has no guidance right now.":     "O sábio não tem orientação agora.",
	"Language set to %s.":                         "Idioma definido para %s.",
	"Press Enter to continue...":                  "Pressione Enter para continuar...",
	"Enter message: ":                             "Digite a mensagem: ",
	"No message sent.":                            "Nenhuma mensagem enviada.",
	"Chat failed":                                 "Falha no chat",
	"Message sent to the room.":                   "Mensagem enviada para a sala.",
	"Event text: ":                                "Texto do evento: ",
	"No event broadcast.":                         "Nenhum evento foi anunciado.",
	"Broadcast failed":                            "Falha ao anunciar",
	"Room event broadcast.":                       "Evento da sala anunciado.",
	"Nobody nearby to whisper.":                   "Ninguém por perto para sussurrar.",
	"Choose a recipient:":                         "Escolha um destinatário:",
	"Recipient: ":                                 "Destinatário: ",
	"Invalid recipient.":                          "Destinatário inválido.",
	"Whisper: ":                                   "Sussurro: ",
	"No whisper sent.":                            "Nenhum sussurro enviado.",
	"Whisper failed":                              "Falha no sussurro",
	"You whisper to %s.":                          "Você sussurra para %s.",
	"Ask about a place, item, or mystery: ":       "Pergunte sobre um lugar, item ou mistério: ",
	"The wise man blinks slowly at your silence.": "O sábio pisca lentamente diante do seu silêncio.",
	"The wise man sighs, 'The tide caves beneath the harbor hide what you seek, but not today. The ruins, moonwell, and observatory also keep their own secrets.'": "O sábio suspira: 'As cavernas das marés sob o porto escondem o que você procura, mas não hoje. As ruínas, o poço da lua e o observatório também guardam seus segredos.'",
	"The wise man says, 'I'm too tired to talk.'":                         "O sábio diz: 'Estou cansado demais para falar.'",
	"Could not create party":                                              "Não foi possível criar a party",
	"Party created. Invite nearby players to join.":                       "Party criada. Convide jogadores próximos para entrar.",
	"Could not load party":                                                "Não foi possível carregar a party",
	"You are not in a party.":                                             "Você não está em uma party.",
	"Could not load party quest log":                                      "Não foi possível carregar o registro da quest da party",
	"Your party has no active quest log.":                                 "Sua party não tem uma quest ativa.",
	"Could not follow leader":                                             "Não foi possível seguir o líder",
	"You follow the party leader.":                                        "Você segue o líder da party.",
	"You lead the party into the beast trail.":                            "Você conduz a party para a trilha da besta.",
	"You guide the party into the void expedition.":                       "Você conduz a party para a expedição do vazio.",
	"You guide the party into the glacier pact.":                          "Você conduz a party para o pacto glacial.",
	"You guide the party into the sun covenant.":                          "Você conduz a party para a aliança solar.",
	"Your party quest does not have a dungeon entry yet.":                 "Sua quest da party ainda não tem entrada de masmorra.",
	"Could not share a party heal":                                        "Não foi possível compartilhar a cura da party",
	"You shared a healing surge with the party.":                          "Você compartilhou uma onda de cura com a party.",
	"Could not rally the party":                                           "Não foi possível reunir a party",
	"You rallied the party for battle.":                                   "Você reuniu a party para a batalha.",
	"Could not raise a party guard":                                       "Não foi possível levantar a guarda da party",
	"You raised a party guard.":                                           "Você ergueu a guarda da party.",
	"Could not load party quest board":                                    "Não foi possível carregar o quadro de quests da party",
	"Current phase:":                                                      "Fase atual:",
	"Rewards:":                                                            "Recompensas:",
	"Party members:":                                                      "Membros da party:",
	"Your party quest board is recovering for another %s.":                "Seu quadro de quests da party ainda está se recuperando por mais %s.",
	"Could not start quest":                                               "Não foi possível iniciar a quest",
	"Party quest started: Monster Hunt.":                                  "Quest da party iniciada: Monster Hunt.",
	"Party quest started: Void Expedition.":                               "Quest da party iniciada: Void Expedition.",
	"Party quest started: Frost Pact.":                                    "Quest da party iniciada: Frost Pact.",
	"Party quest started: Sun Covenant.":                                  "Quest da party iniciada: Sun Covenant.",
	"Could not leave party":                                               "Não foi possível sair da party",
	"You left the party.":                                                 "Você saiu da party.",
	"Nobody nearby to invite.":                                            "Ninguém por perto para convidar.",
	"Invite failed":                                                       "Falha no convite",
	"Party invite sent to %s.":                                            "Convite de party enviado para %s.",
	"No room chat yet.":                                                   "Ainda não há chat da sala.",
	"No whispers yet.":                                                    "Ainda não há sussurros.",
	"Your inbox is empty.":                                                "Sua caixa de entrada está vazia.",
	"Inbox dismissed without marking notifications as read.":              "Caixa de entrada fechada sem marcar notificações como lidas.",
	"Reply failed":                                                        "Falha ao responder",
	"Replied to %s.":                                                      "Respondido para %s.",
	"Party invite failed":                                                 "Falha no convite da party",
	"Joined party %s.":                                                    "Entrou na party %s.",
	"Party heal from %s is already applied.":                              "A cura da party de %s já está aplicada.",
	"Party rally from %s is already active.":                              "A reunião da party de %s já está ativa.",
	"Party guard from %s is already active.":                              "A guarda da party de %s já está ativa.",
	"Party quest update from %s is already reflected on the quest board.": "A atualização da quest da party de %s já está refletida no quadro.",
	"Thank-you failed":                                                    "Falha ao agradecer",
	"Sent a thank-you to %s.":                                             "Agradecimento enviado para %s.",
	"You already carry the %s.":                                           "Você já carrega o(a) %s.",
	"You found the %s.":                                                   "Você encontrou o(a) %s.",
	"You have no relics to trade.":                                        "Você não tem relicos para trocar.",
	"No nearby players to trade with.":                                    "Nenhum jogador próximo para trocar.",
	"Invalid relic.":                                                      "Relíquia inválida.",
	"Trade failed":                                                        "Falha na troca",
	"You traded %s for %s with %s.":                                       "Você trocou %s por %s com %s.",
	"Wise man AI configured with Gemini.":                                 "IA do sábio configurada com Gemini.",
	"Wise man AI configured with Cloudflare Workers AI.":                  "IA do sábio configurada com Cloudflare Workers AI.",
	"Wise man AI setup skipped.":                                          "Configuração da IA do sábio ignorada.",
	"Welcome!":                                                            "Bem-vindo!",
	"See ya!":                                                             "Até logo!",
	"No gold!":                                                            "Sem ouro!",
	"Bought Sword!":                                                       "Espada comprada!",
	"Bought Meat!":                                                        "Carne comprada!",
	"Bought Kit!":                                                         "Kit comprado!",
	"Bought Water!":                                                       "Água comprada!",
	"Bought Fur Coat!":                                                    "Casaco de pele comprado!",
	"You don't have enough gold!":                                         "Você não tem ouro suficiente!",
	"Level Up!":                                                           "Subiu de nível!",
	"Heading home.":                                                       "Voltando para casa.",
	"Water sound ahead.":                                                  "Há som de água à frente.",
	"Roars ahead.":                                                        "Rugidos à frente.",
	"The air gets hot.":                                                   "O ar fica quente.",
	"Old stones await beyond the trees.":                                  "Pedras antigas esperam além das árvores.",
	"Back to woods.":                                                      "De volta à mata.",
	"Crossing bridge...":                                                  "Atravessando a ponte...",
	"You follow the river until the salt air reaches you.":                "Você segue o rio até o ar salgado alcançá-lo.",
	"The path grows quiet and thoughtful.":                                "O caminho fica silencioso e reflexivo.",
	"From the beacon room you see every road you have walked. The wise man's words feel clearer now.": "Da sala do farol, você vê todas as estradas que percorreu. As palavras do sábio parecem mais claras agora.",
	"The tide guardian collapses and reveals the Moon Pearl.":                                         "O guardião da maré cai e revela a Pérola da Lua.",
	"You recover a Moon Charm and a little wisdom.":                                                   "Você recupera um Amuleto da Lua e um pouco de sabedoria.",
	"The wise man says he cannot accept an empty palm.":                                               "O sábio diz que não pode aceitar uma mão vazia.",
	"The wise man already marked your triptych as complete.":                                          "O sábio já marcou seu tríptico como concluído.",
	"The wise man says the three signs are not yet all in your hands.":                                "O sábio diz que os três sinais ainda não estão em suas mãos.",
	"You gain the Triptych Blessing, 150 gold, and a sharper sense of hidden paths.":                  "Você recebe a Bênção do Tríptico, 150 de ouro e uma percepção mais aguçada dos caminhos ocultos.",
	"Invalid slot. Please enter 1-5.":                                                                 "Slot inválido. Digite 1-5.",
	"No actions are available right now.":                                                             "Nenhuma ação está disponível agora.",
	"Invalid selection.":                                                                              "Seleção inválida.",
	"Game Saved. Goodbye!":                                                                            "Jogo salvo. Até logo!",
	"Press Enter to start your adventure...":                                                          "Pressione Enter para começar sua aventura...",
	"Select an option: ":                                                                              "Selecione uma opção: ",
	"Select a slot (1-5): ":                                                                           "Selecione um slot (1-5): ",
	"Continue":                                                                                        "Continuar",
	"Overwrite (New Game)":                                                                            "Sobrescrever (Novo Jogo)",
	"Select: ":                                                                                        "Selecione: ",
	"Enter your character's name: ":                                                                   "Digite o nome do seu personagem: ",
	"Would you like to wire an AI bot for the wise man now? (y/n): ":                                  "Deseja conectar um bot de IA ao sábio agora? (s/n): ",
	"Welcome to Magic Adventure 6.0!":                                                                 "Bem-vindo(a) ao Magic Adventure 6.0!",
	"Created by Dante Gomes with assistance from Gemini and Codex.":                                   "Criado por Dante Gomes com assistência de Gemini e Codex.",
	"--- CHARACTER SLOTS ---":                                                                         "--- SLOTS DE PERSONAGEM ---",
	"=== HOW TO PLAY ===":                                                                             "=== COMO JOGAR ===",
	"Welcome, adventurer! Magic Adventure is a multiplayer text RPG.":                                 "Bem-vindo, aventureiro! Magic Adventure é um RPG de texto multijogador.",
	"Navigation": "Navegação",
	"Choose numbered options to move between locations.": "Escolha opções numeradas para se mover entre os locais.",
	"Combat": "Combate",
	"Fight monsters to earn XP and Gold. Your stats increase automatically.": "Lute contra monstros para ganhar XP e ouro. Seus atributos aumentam automaticamente.",
	"You can see other online players and their locations.":                  "Você pode ver outros jogadores online e seus locais.",
	"Party Play": "Jogo em party",
	"Leaders can heal, rally, and guard the party for real combat advantages.": "Líderes podem curar, reunir e proteger a party para obter vantagens reais no combate.",
	"Party Quests": "Quests da party",
	"Parties can start cooperative quests that reward stronger shared buffs.":                                                  "Parties podem iniciar quests cooperativas que concedem bônus compartilhados mais fortes.",
	"The village, forest, river, harbor, lighthouse, tide caves, ruins, observatory, and sage hut all hide different rewards.": "A vila, floresta, rio, porto, farol, cavernas das marés, ruínas, observatório e cabana do sábio escondem recompensas diferentes.",
	"Wise Man": "Sábio",
	"You can wire him to Gemini or Cloudflare and reconfigure him later from the Settings menu.": "Você pode conectá-lo ao Gemini ou Cloudflare e reconfigurá-lo depois no menu Configurações.",
	"Progression": "Progressão",
	"Buy items in the Shop and complete biomes to find the Dragon.": "Compre itens na Loja e complete biomas para encontrar o Dragão.",
	"Tip: Keep an eye on your Health! Use Potions to stay alive.":   "Dica: fique de olho na sua vida! Use Poções para sobreviver.",
	"G A M E   O V E R":          "F I M   D E   J O G O",
	"Try again? (yes/no)":        "Tentar novamente? (sim/não)",
	"Choose a player to invite:": "Escolha um jogador para convidar:",
	"Notifications inbox:":       "Caixa de notificações:",
	"Select a notification, type r<number> to reply, x to dismiss, or press Enter to mark all as read: ": "Selecione uma notificação, digite r<número> para responder, x para fechar, ou Enter para marcar tudo como lido: ",
	"Reply to %s? (y/n): ":                        "Responder para %s? (s/n): ",
	"Accept party invite from %s? (y/n): ":        "Aceitar convite de party de %s? (s/n): ",
	"Send a thank-you to %s? (y/n): ":             "Enviar um agradecimento para %s? (s/n): ",
	"These notifications are now marked as read.": "Essas notificações agora estão marcadas como lidas.",
	"%s has no relics to trade.":                  "%s não tem relíquias para trocar.",
	"Choose one of your relics to offer:":         "Escolha uma de suas relíquias para oferecer:",
	"Your relic: ":                                "Sua relíquia: ",
	"Choose one of their relics to request:":      "Escolha uma das relíquias dele(a) para pedir:",
	"Their relic: ":                               "Relíquia dele(a): ",
	"Party quest board:":                          "Quadro de quests da party:",
	"No active quest.":                            "Nenhuma quest ativa.",
	"Choose an option:":                           "Escolha uma opção:",
	"Start Monster Hunt quest":                    "Iniciar quest Monster Hunt",
	"Start Void Expedition quest":                 "Iniciar quest Void Expedition",
	"Start Frost Pact quest":                      "Iniciar quest Frost Pact",
	"Start Sun Covenant quest":                    "Iniciar quest Sun Covenant",
	"Refresh board":                               "Atualizar quadro",
	"The wise man invites you to speak. Type your question, or type `back` to end the conversation.": "O sábio o convida a falar. Digite sua pergunta, ou digite `back` para encerrar a conversa.",
	"You":                           "Você",
	"You: ":                         "Você: ",
	"The wise man waits patiently.": "O sábio espera pacientemente.",
	"You end your conversation with the wise man.":                      "Você encerra sua conversa com o sábio.",
	"The wise man rubs his eyes and asks you to try again in a moment.": "O sábio esfrega os olhos e pede para você tentar novamente em instantes.",
	"The Wise Man appears with a challenge!":                             "O Homem Sábio aparece com um desafio!",
	"Choose the correct answer: ":                                        "Escolha a resposta correta: ",
	"Correct! You earned 50 gold.":                                       "Correto! Você ganhou 50 de ouro.",
	"Wrong answer! You lost 50 gold.":                                    "Resposta errada! Você perdeu 50 de ouro.",
	"The correct answer was":                                             "A resposta correta era",
	"That's 3 wrong in a row! You lost 10 Health.":                       "Já são 3 erros seguidos! Você perdeu 10 de Vida.",
	"You passed the challenge and entered!":                              "Você passou no desafio e entrou!",
	"You failed the Wise Man's challenge.":                               "Você falhou no desafio do Homem Sábio.",
	"You meditated and gained 50 Skill Points!":                          "Você meditou e ganhou 50 Pontos de Habilidade!",
	"You guessed correctly! You won 30 gold!":                            "Você adivinhou corretamente! Você ganhou 30 de ouro!",
	"Wrong cup! Better luck next time.":                                  "Copo errado! Mais sorte na próxima vez.",
	"Lift Cup 1 (10 gold)":                "Levantar Copo 1 (10 ouro)",
	"Lift Cup 2 (10 gold)":                "Levantar Copo 2 (10 ouro)",
	"Lift Cup 3 (10 gold)":                "Levantar Copo 3 (10 ouro)",
	"The pearl was in cup %s.":            "A pérola estava no copo %s.",
	"Correct":                             "Correto",
	"Visit the Wise Man":                  "Visitar o Homem Sábio",
	"Go to Knowledge Hall":                                               "Ir para o Salão do Conhecimento",
	"Enter Hall of Wisdom (Challenge)":                                   "Entrar no Salão da Sabedoria (Desafio)",
	"Go to Guessing Table":                                               "Ir para a Mesa de Adivinhação",
	"Back to Village":                     "Voltar para a Vila",
	"Leave to Knowledge Hall":             "Sair para o Salão do Conhecimento",
	"Back to Knowledge Hall":              "Voltar para o Salão do Conhecimento",

	"Meditate (+50 SP)":                                                  "Meditar (+50 SP)",
	"Millionaire Help Request":                                           "Pedido de Ajuda do Milionário",
	"Help sent!":                                                         "Ajuda enviada!",
	"Go to Millionaire Hall":                                              "Ir para o Salão do Milionário",
	"The lights get brighter and you hear dramatic music.":               "As luzes ficam mais brilhantes e você ouve música dramática.",
	"Start Millionaire Game":                                             "Começar Jogo do Milionário",
	"The Hall of Millionaires. Spotlights shine on a central stage with two chairs. A giant scoreboard hangs on the wall.": "O Salão dos Milionários. Holofotes brilham num palco central com duas cadeiras. Um placar gigante está na parede.",
	"Leaderboard reset! Rewards sent to winners.":                        "Placar reiniciado! Recompensas enviadas aos vencedores.",
}

var roomTranslations = map[string]string{
	"village":          "Você está na Vila de Oakhaven. É um lugar pacífico. Você vê uma loja e um jardim.",
	"knowledge_hall":   "Você está no Salão do Conhecimento. Grandes estantes de livros revestem as paredes. Ao norte está o Salão da Sabedoria, mas ele está guardado.",
	"millionaire_hall": "O Salão dos Milionários. Holofotes brilham num palco central com duas cadeiras. Um placar gigante está na parede.",
	"wisdom_hall":      "O Salão da Sabedoria. Um pedestal dourado está no centro. Você já se sente muito mais inteligente!",
	"guessing_game":    "Uma pequena mesa com três copos. Um velho desafia você a adivinhar onde está a pérola.",
	"shop":             "A loja geral. Tudo o que você precisa para sobreviver.",
	"garden":           "O tranquilo Jardim Zen.",
	"forest":           "Um cruzamento na floresta. Vila a oeste, Rio ao norte, Zoo ao sul, Deserto a leste, e ruínas antigas por uma trilha escondida.",
	"river":            "Um rio frio. Floresta ao sul, Caverna atrás da cachoeira, Árctico ao norte pela ponte.",
	"sage_hut":         "Uma cabana iluminada por velas, cheia de pergaminhos, conchas e mil respostas silenciosas.",
	"harbor":           "Um porto movimentado de vento salgado, cordas e barcos rangendo. Um velho farol observa a água.",
	"lighthouse":       "Um farol alto cujo feixe corta as nuvens de tempestade.",
	"moonwell":         "Um Poço da Lua alimentado pelas marés, escondido atrás do porto. A água brilha quando você fala baixo.",
	"tide_caves":       "Cavernas baixas esculpidas pela maré. Algo brilhante reluz sob a espuma.",
	"old_ruins":        "Arcos de pedra caídos e mosaicos quebrados. A floresta esconde este lugar dos viajantes casuais.",
	"desert":           "O Deserto Escaldante. Dunas infinitas e calor.",
	"sphinx":           "A Esfinge fala: 'Não tenho voz, mas posso gritar. Não tenho asas, mas posso voar. O que sou?'",
	"arctic":           "A Entrada do Árctico. Um deserto congelado. Ao norte ficam os Penhascos Gelados, e a leste há uma Ponte de Gelo.",
	"frost_cliffs":     "Os altos Penhascos Gelados. O vento uiva aqui. Um Gigante do Gelo guarda esta área.",
	"ice_bridge":       "Uma ponte estreita de gelo. Você vê um Portal Glitchado cintilando no fim dela.",
	"mountain":         "O Pico da Montanha. O Grande Verme aguarda, e um observatório fica mais acima na trilha.",
	"observatory":      "Um observatório em ruínas onde o céu parece próximo o suficiente para tocar.",
	"zoo":              "A besta-trilha se abre numa arena antiga, com rugidos ecoando ao redor.",
	"void":             "Os Ermos do Glitch. A realidade aqui pisca. Você sente uma presença estranha.",
	"binary_sea":       "Um mar binário de bits e ruído digital.",
	"party_sun_gate":   "Um portal de areia e sol.",
	"party_sun_depths": "As profundezas do deserto solar.",
	"party_sun_crown":  "A coroa do templo solar.",
}

var actionRules = []translationRule{
	{prefix: "Go ", apply: func(text string) string { return strings.Replace(text, "Go ", "Ir ", 1) }},
	{prefix: "Enter ", apply: func(text string) string { return strings.Replace(text, "Enter ", "Entrar em ", 1) }},
	{prefix: "Visit ", apply: func(text string) string { return strings.Replace(text, "Visit ", "Visitar ", 1) }},
	{prefix: "Leave ", apply: func(text string) string { return strings.Replace(text, "Leave ", "Sair de ", 1) }},
	{prefix: "Return ", apply: func(text string) string { return strings.Replace(text, "Return ", "Retornar para ", 1) }},
	{prefix: "Back to ", apply: func(text string) string { return strings.Replace(text, "Back to ", "Voltar para ", 1) }},
	{prefix: "Buy ", apply: func(text string) string { return strings.Replace(text, "Buy ", "Comprar ", 1) }},
	{prefix: "Attack ", apply: func(text string) string { return strings.Replace(text, "Attack ", "Atacar ", 1) }},
	{prefix: "Fight ", apply: func(text string) string { return strings.Replace(text, "Fight ", "Lutar contra ", 1) }},
	{prefix: "Open ", apply: func(text string) string { return strings.Replace(text, "Open ", "Abrir ", 1) }},
	{prefix: "Ask ", apply: func(text string) string { return strings.Replace(text, "Ask ", "Perguntar ", 1) }},
	{prefix: "Search ", apply: func(text string) string { return strings.Replace(text, "Search ", "Procurar ", 1) }},
	{prefix: "Share ", apply: func(text string) string { return strings.Replace(text, "Share ", "Compartilhar ", 1) }},
	{prefix: "Create ", apply: func(text string) string { return strings.Replace(text, "Create ", "Criar ", 1) }},
	{prefix: "View ", apply: func(text string) string { return strings.Replace(text, "View ", "Ver ", 1) }},
	{prefix: "Follow ", apply: func(text string) string { return strings.Replace(text, "Follow ", "Seguir ", 1) }},
}

var translationRules = []translationRule{
	{prefix: "You head ", apply: func(text string) string {
		return strings.ReplaceAll(strings.ReplaceAll(text, "You head east.", "Você segue para leste."), "You head west.", "Você segue para oeste.")
	}},
	{prefix: "You drink", apply: func(text string) string {
		return strings.Replace(text, "You drank a Minor Potion and recovered ", "Você bebeu uma Poção Menor e recuperou ", 1)
	}},
	{prefix: "You recover a ", apply: func(text string) string { return strings.Replace(text, "You recover a ", "Você recupera ", 1) }},
	{prefix: "You recover ", apply: func(text string) string { return strings.Replace(text, "You recover ", "Você recupera ", 1) }},
	{prefix: "You found the ", apply: func(text string) string { return strings.Replace(text, "You found the ", "Você encontrou a ", 1) }},
	{prefix: "You are", apply: func(text string) string { return strings.Replace(text, "You are", "Você está", 1) }},
	{prefix: "The wise man ", apply: func(text string) string { return text }},
}

func localizeActionDesc(lang string, desc string) string {
	return TranslateActionDescription(lang, desc)
}

func localizeResult(lang string, msg string) string {
	return TranslateText(lang, msg)
}

func localizeRoom(lang, roomID, desc string) string {
	return TranslateRoomDescription(lang, roomID, desc)
}

func localizeProviderLabel(lang, provider string) string {
	return FormatWiseManProvider(lang, provider)
}

func localizef(lang, format string, args ...any) string {
	return TranslateText(lang, fmt.Sprintf(format, args...))
}
