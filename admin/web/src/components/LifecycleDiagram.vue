<template>
  <div class="lifecycle-container">
    <div class="lifecycle-toolbar">
      <div class="phase-filters">
        <button 
          v-for="phase in phases" 
          :key="phase.id" 
          :class="['phase-btn', { active: selectedPhase === phase.id }]" 
          @click="filterByPhase(phase.id)"
        >
          {{ phase.label }}
        </button>
      </div>

      <div class="toolbar-hint">
        <span class="hint-icon">💡</span>
        <span>Clique em qualquer <strong>Identificador (borda verde)</strong> para ver mensagens e regras</span>
      </div>

      <div class="zoom-controls">
        <button class="zoom-btn" @click="zoomOut" title="Diminuir zoom">－</button>
        <span class="zoom-value">{{ Math.round(zoomScale * 100) }}%</span>
        <button class="zoom-btn" @click="zoomIn" title="Aumentar zoom">＋</button>
        <button class="zoom-btn zoom-reset" @click="resetZoom" title="Resetar zoom">↺</button>
      </div>
    </div>

    <div class="lifecycle-main-layout">
      <div 
        class="diagram-viewport" 
        ref="viewportRef"
        @mousedown="startPan"
        @mousemove="onPan"
        @mouseup="endPan"
        @mouseleave="endPan"
      >
        <div 
          class="diagram-canvas" 
          :style="{
            transform: `translate(${panX}px, ${panY}px) scale(${zoomScale})`,
            transformOrigin: '0 0'
          }"
        >
          <div ref="diagramContainerRef" class="mermaid-target"></div>
        </div>
      </div>

      <transition name="drawer">
        <aside v-if="selectedNotification" class="detail-drawer">
          <div class="drawer-header">
            <div class="drawer-header-info">
              <span class="drawer-group-badge">{{ selectedGroupTitle }}</span>
              <h3 class="drawer-id">{{ selectedNotification.id }}</h3>
              <p class="drawer-title">{{ selectedNotification.title_pt }}</p>
              <p v-if="selectedNotification.title_en" class="drawer-subtitle-en">
                🇺🇸 {{ selectedNotification.title_en }}
              </p>
            </div>
            <button class="btn-close-drawer" @click="selectedNotification = null">&times;</button>
          </div>

          <div class="drawer-body">
            <div class="rule-card">
              <div class="rule-header">
                <span class="rule-icon">⚙️</span>
                <h4>Regra de Negócio & Ciclo de Vida</h4>
              </div>
              <div class="rule-grid">
                <div class="rule-item">
                  <span class="rule-label">Fase</span>
                  <span class="rule-val">{{ currentRule.phase || 'N/A' }}</span>
                </div>
                <div class="rule-item">
                  <span class="rule-label">Local</span>
                  <span class="rule-val">{{ currentRule.location || 'Push Android' }}</span>
                </div>
                <div class="rule-item">
                  <span class="rule-label">Gatilho</span>
                  <span class="rule-val">{{ currentRule.trigger || 'N/A' }}</span>
                </div>
                <div class="rule-item">
                  <span class="rule-label">Cooldown</span>
                  <span class="rule-val">{{ currentRule.cooldown || 'N/A' }}</span>
                </div>
              </div>
            </div>

            <div class="drawer-tabs">
              <button 
                :class="['drawer-tab-btn', { active: activeDrawerTab === 'view' }]" 
                @click="activeDrawerTab = 'view'"
              >
                📋 Visualizar
              </button>
              <button 
                :class="['drawer-tab-btn', { active: activeDrawerTab === 'crud' }]" 
                @click="openCrudTab"
              >
                ✏️ Gerenciar Variações
              </button>
            </div>

            <div v-if="activeDrawerTab === 'view'" class="messages-section">
              <div class="messages-header-row">
                <div class="messages-title-col">
                  <h4>Variações de Notificação</h4>
                  <span class="count-badge">
                    {{ currentMessages.length }} mensagens
                  </span>
                </div>

                <div class="lang-mini-switch">
                  <button 
                    :class="['mini-btn', { active: detailLang === 'ptbr' }]" 
                    @click="detailLang = 'ptbr'"
                  >
                    PT-BR
                  </button>
                  <button 
                    :class="['mini-btn', { active: detailLang === 'en' }]" 
                    @click="detailLang = 'en'"
                  >
                    EN
                  </button>
                </div>
              </div>

              <div class="drawer-messages-list">
                <div 
                  v-for="(msg, idx) in currentMessages" 
                  :key="idx" 
                  class="drawer-msg-item"
                >
                  <span class="drawer-msg-num">{{ idx + 1 }}</span>
                  <span class="drawer-msg-text">{{ msg }}</span>
                  <button 
                    class="btn-copy" 
                    @click="copyText(msg, idx)"
                    :title="copiedIndex === idx ? 'Copiado!' : 'Copiar texto'"
                  >
                    {{ copiedIndex === idx ? '✓' : '📋' }}
                  </button>
                </div>
              </div>
            </div>

            <div v-else class="crud-section">
              <div class="crud-header-row">
                <div class="crud-title-col">
                  <h4>Variações de Notificação</h4>
                  <span class="count-badge">{{ editableMessages.length }} variações</span>
                </div>
              </div>

              <div class="crud-field-group">
                <label class="crud-label">Título da Notificação (PT-BR)</label>
                <input 
                  v-model="editableTitle" 
                  type="text" 
                  class="crud-input-title" 
                  placeholder="Ex: Você Foi Gigante Hoje!" 
                />
              </div>

              <div class="crud-messages-container">
                <label class="crud-label">Lista de Mensagens (PT-BR)</label>
                <div 
                  v-for="(msg, idx) in editableMessages" 
                  :key="idx" 
                  class="crud-message-item"
                >
                  <span class="crud-item-num">{{ idx + 1 }}</span>
                  <textarea 
                    v-model="editableMessages[idx]" 
                    rows="2" 
                    class="crud-textarea" 
                    placeholder="Texto da mensagem..."
                  ></textarea>
                  <button 
                    type="button" 
                    class="btn-delete-msg" 
                    @click="removeMessage(idx)" 
                    title="Remover variação"
                  >
                    🗑️
                  </button>
                </div>

                <button 
                  type="button" 
                  class="btn btn-secondary btn-sm btn-add-msg" 
                  @click="addNewMessage"
                >
                  + Adicionar Variação
                </button>
              </div>

              <div class="crud-options-card">
                <label class="crud-toggle-label">
                  <input type="checkbox" v-model="autoTranslate" />
                  <span>Traduzir automaticamente para inglês (EN)</span>
                </label>
                <small class="crud-help-text">
                  Gera a versão em inglês e atualiza o arquivo <code>config_notifications.json</code> com a data, hora e a lista de grupos modificados.
                </small>
              </div>

              <div v-if="crudStatusMessage" :class="['crud-status-banner', `status-${crudStatusType}`]">
                {{ crudStatusMessage }}
              </div>

              <div class="crud-save-row">
                <button 
                  class="btn btn-primary btn-full" 
                  @click="saveChanges" 
                  :disabled="isSavingCrud"
                >
                  {{ isSavingCrud ? '💾 Salvando e Traduzindo...' : '💾 Salvar Alterações' }}
                </button>
              </div>
            </div>
          </div>

          <div class="drawer-footer">
            <button 
              v-if="activeDrawerTab === 'view'"
              class="btn btn-secondary btn-full" 
              @click="openCrudTab"
            >
              ✏️ Gerenciar Variações no Ciclo de Vida
            </button>
            <button 
              v-else
              class="btn btn-secondary btn-full" 
              @click="activeDrawerTab = 'view'"
            >
              👁️ Voltar para Visualização
            </button>
          </div>
        </aside>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import mermaid from 'mermaid'
import { saveNotificationItem } from '../api'

const props = defineProps({
  groups: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['edit-notification', 'refresh'])

const viewportRef = ref(null)
const diagramContainerRef = ref(null)

const selectedPhase = ref('all')
const selectedNotification = ref(null)
const detailLang = ref('ptbr')
const copiedIndex = ref(null)

const activeDrawerTab = ref('view')
const editableTitle = ref('')
const editableMessages = ref([])
const autoTranslate = ref(true)
const isSavingCrud = ref(false)
const crudStatusMessage = ref('')
const crudStatusType = ref('success')

const zoomScale = ref(0.95)
const panX = ref(40)
const panY = ref(20)
let isPanning = false
let startX = 0
let startY = 0

const phases = [
  { id: 'all', label: 'Todas as Fases' },
  { id: 'FASE_1', label: '1. Fora do App' },
  { id: 'FASE_2', label: '2. Entrada / Wi-Fi' },
  { id: 'FASE_3', label: '3. Em Treino' },
  { id: 'FASE_4', label: '4. Pós-Treino' }
]

const LIFECYCLE_RULES = {
  WORKOUT_DAY_ONE: {
    phase: '1. Fora do App / Boas-Vindas',
    trigger: 'Usuário sem nenhum treino no histórico',
    cooldown: 'Exibido até o primeiro treino',
    location: 'Push Android e Card na Home',
    rule: 'Ativado no dia planejado de treino quando o usuário ainda não completou nenhuma sessão no app.'
  },
  COMEBACK_SHORT_ABSENCE: {
    phase: '1. Fora do App / Recuperação Curta',
    trigger: '3 a 6 dias sem treinar',
    cooldown: '20 horas',
    location: 'Push Android e Card na Home',
    rule: 'Incentiva o retorno após os primeiros dias de falta antes de virar perda de consistência.'
  },
  COMEBACK_WEEK_ABSENCE: {
    phase: '1. Fora do App / Recuperação Média',
    trigger: '7 a 13 dias sem treinar',
    cooldown: '48 horas',
    location: 'Push Android e Card na Home',
    rule: 'Passa o nome do último treino concluído para reativar o usuário após uma semana inteira fora.'
  },
  COMEBACK_LONG_ABSENCE: {
    phase: '1. Fora do App / Hibernação',
    trigger: '≥ 14 dias sem treinar',
    cooldown: '168 horas (7 dias)',
    location: 'Push Android e Card na Home',
    rule: 'Mensagem de retorno gradual e leve com cooldown semanal estendido para não saturar o usuário inativo.'
  },
  ONE_HOUR_BEFORE_HABITUAL: {
    phase: '1. Fora do App / Lembrete Horário',
    trigger: 'Exatamente 1h antes do horário habitual',
    cooldown: 'Diário (alarme exato)',
    location: 'Push Android',
    rule: 'Calculado pelo alarme com base no horário médio usual de treino. Não dispara se já treinou hoje.'
  },
  WORKOUT_PREFERRED_TIME: {
    phase: '1. Fora do App / Horário Habitual',
    trigger: 'Slot habitual atingido (Manhã, Tarde ou Noite)',
    cooldown: 'Diário (20h)',
    location: 'Push Android e Card na Home',
    rule: 'Dispara no momento exato em que o horário do dia coincide com o slot preferido de treino.'
  },
  WORKOUT_SCHEDULED_TODAY: {
    phase: '1. Fora do App / Rotina Planejada',
    trigger: 'Dia presente na rotina planejada',
    cooldown: 'Diário (20h)',
    location: 'Push Android e Card na Home',
    rule: 'Ativado nos dias cadastrados na agenda do usuário (ex: Seg/Qua/Sex) caso ainda não tenha treinado.'
  },
  MISSED_MORNING_AFTERNOON: {
    phase: '1. Fora do App / Chamada da Tarde',
    trigger: 'Alarme das 15h para usuários matutinos ausentes',
    cooldown: 'Diário',
    location: 'Push Android',
    rule: 'Dispara se o usuário costuma treinar entre 4h e 11h, não treinou e não conectou no Wi-Fi da academia hoje.'
  },
  STREAK_AT_RISK: {
    phase: '1. Fora do App / Risco de Sequência',
    trigger: 'Horário ≥ 19h com sequência ativa ≥ 3 dias sem treino',
    cooldown: 'Diário',
    location: 'Push Android e Card na Home',
    rule: 'Gera alerta urgente no final do dia para salvar os dias consecutivos de treino antes de zerar.'
  },
  WORKOUT_QUICK_CALL: {
    phase: '1. Fora do App / Fim do Dia',
    trigger: 'Horário ≥ 19h em dia planejado',
    cooldown: 'Diário',
    location: 'Push Android e Card na Home',
    rule: 'Sugere uma sessão rápida e expressa de poucos minutos para fechar o dia com dever cumprido.'
  },
  MISSED_DAY_NIGHT: {
    phase: '1. Fora do App / Última Chamada',
    trigger: 'Alarme das 20h para quem não treinou no dia',
    cooldown: 'Diário',
    location: 'Push Android',
    rule: 'Última chamada noturna para resgatar a rotina antes de dormir.'
  },
  EVOLUTION_NEGLECTED_MUSCLE: {
    phase: '1. Fora do App / Equilíbrio Muscular',
    trigger: 'Músculo sem treino há ≥ 10 dias',
    cooldown: 'Diário',
    location: 'Push Android e Card na Home',
    rule: 'Identifica o grupo muscular mais esquecido nas últimas semanas e sugere incluir na sessão.'
  },
  GYM_WIFI_ARRIVAL: {
    phase: '2. Entrada / Presença Física',
    trigger: 'Conectou ao Wi-Fi cadastrado da academia',
    cooldown: '4 horas',
    location: 'Push Android',
    rule: 'Dispara no momento em que o celular reconhece o SSID da academia e o usuário ainda não iniciou o treino.'
  },
  WORKOUT_OPENED_NO_WORKOUT: {
    phase: '2. Interação / Abandono Inicial',
    trigger: 'Abriu o FitGuide e saiu sem iniciar treino',
    cooldown: '72 horas',
    location: 'Push Android (atraso 30 min)',
    rule: 'Alarme checado 30 minutos após sair do app para resgatar quem apenas espiou e fechou.'
  },
  IN_WORKOUT_PR: {
    phase: '3. Em Treino / Recorde Pessoal',
    trigger: 'Carga superou o recorde histórico no exercício',
    cooldown: 'Instantâneo',
    location: 'Banner no Topo do Treino (10s)',
    rule: 'Exibido imediatamente ao bater nova carga máxima histórica na série.'
  },
  IN_WORKOUT_REST_OVERRUN: {
    phase: '3. Em Treino / Descanso Excedido',
    trigger: 'Descanso ultrapassou o planejado em +30s',
    cooldown: 'Por série',
    location: 'Banner no Topo do Treino (10s)',
    rule: 'Chama a atenção com tom bem-humorado avisando que o descanso acabou e o aparelho está ocupado.'
  },
  IN_WORKOUT_HALFWAY: {
    phase: '3. Em Treino / Metade da Sessão',
    trigger: 'Concluiu exatamente 50% dos exercícios (fichas ≥ 4)',
    cooldown: '1 vez por treino',
    location: 'Banner no Topo do Treino (10s)',
    rule: 'Marca que a metade do treino já foi vencida e incentiva a segunda etapa.'
  },
  IN_WORKOUT_LOAD_PROGRESSION: {
    phase: '3. Em Treino / Progressão de Carga',
    trigger: 'Carga dominada com facilidade em 2 sessões seguidas',
    cooldown: '1 vez por exercício',
    location: 'Banner no Topo do Treino (10s)',
    rule: 'Sugere subir a carga no exercício para manter a sobrecarga progressiva.'
  },
  IN_WORKOUT_LAST_EXERCISE: {
    phase: '3. Em Treino / Reta Final',
    trigger: 'Resta apenas 1 exercício para encerrar a ficha',
    cooldown: '1 vez por treino',
    location: 'Banner no Topo do Treino (10s)',
    rule: 'Alerta de motivação final para entregar tudo no último exercício.'
  },
  WORKOUT_ABANDONED: {
    phase: '3. Em Treino / Treino Abandonado',
    trigger: 'Treino iniciado e fechado com 0 exercícios marcados',
    cooldown: '24 horas',
    location: 'Push Android (atraso 30 min)',
    rule: 'Dispara 30 minutos após encerramento com zero séries feitas chamando para voltar.'
  },
  POST_WORKOUT_CONGRATS: {
    phase: '4. Pós-Treino / Reconhecimento',
    trigger: 'Alarme exato de 1 hora pós-treino concluído',
    cooldown: 'Por treino',
    location: 'Push Android',
    rule: 'Dispara parabéns destacando se completou 100% ou parcial dos exercícios planejados.'
  },
  STREAK_ACTIVE: {
    phase: '4. Pós-Treino / Ofensiva Ativa',
    trigger: 'Sequência consecutiva de treinos ≥ 3 dias',
    cooldown: 'Diário',
    location: 'Push Android e Card na Home',
    rule: 'Reforça o hábito e o fogo da ofensiva ativa por manter a rotina diária.'
  },
  STREAK_NEW_RECORD: {
    phase: '4. Pós-Treino / Recorde Histórico',
    trigger: 'Sequência atual superou o maior recorde histórico',
    cooldown: 'Por novo recorde',
    location: 'Push Android e Card na Home',
    rule: 'Comemoração de marco histórico de consistência máxima já alcançada.'
  },
  WEEKLY_GOAL_COMPLETED: {
    phase: '4. Pós-Treino / Meta Semanal',
    trigger: '100% dos treinos da semana entregues',
    cooldown: 'Semanal',
    location: 'Push Android e Card na Home',
    rule: 'Disparado ao cumprir a meta da semana sem dever nenhum treino.'
  },
  EVOLUTION_VOLUME_MILESTONE: {
    phase: '4. Pós-Treino / Marco de Volume',
    trigger: 'Volume vitalício cruzou múltiplo de 10.000 kg',
    cooldown: 'A cada 10 toneladas',
    location: 'Push Android e Card na Home',
    rule: 'Celebra marcos históricos de toneladas totais levantadas no app.'
  },
  EVOLUTION_HEALTH_BENEFIT: {
    phase: '4. Pós-Treino / Saúde e Função',
    trigger: 'Treino concluído com benefício anatômico mapeado',
    cooldown: 'Diário',
    location: 'Push Android e Card na Home',
    rule: 'Conecta o treino de um músculo primário à postura, densidade óssea ou estabilidade articular.'
  }
}

const NOTIF_IDS = Object.keys(LIFECYCLE_RULES)

const currentRule = computed(() => {
  if (!selectedNotification.value) return {}
  return LIFECYCLE_RULES[selectedNotification.value.id] || {}
})

const selectedGroupTitle = computed(() => {
  if (!selectedNotification.value) return ''
  for (const g of props.groups) {
    if (g.items.some(item => item.id === selectedNotification.value.id)) {
      return g.name_pt || g.group
    }
  }
  return ''
})

const currentMessages = computed(() => {
  if (!selectedNotification.value) return []
  if (detailLang.value === 'en') {
    return selectedNotification.value.messages_en || []
  }
  return selectedNotification.value.messages_pt || []
})

function getDiagramDefinition() {
  return `flowchart TD
    classDef notifNode fill:#152238,stroke:#10b981,stroke-width:2.5px,color:#f8fafc;
    classDef decisionNode fill:#1e293b,stroke:#64748b,stroke-width:1.5px,color:#94a3b8;
    classDef actionNode fill:#0f172a,stroke:#3b82f6,stroke-width:1.5px,color:#cbd5e1;

    subgraph FASE_1["1. FORA DO APP: Rotina, Ausência e Lembretes"]
        START([Dia Iniciado]) --> CHK_TRAINED{Já treinou hoje?}
        
        CHK_TRAINED -- "Não" --> CHK_ABSENCE{Dias sem treinar?}
        CHK_ABSENCE -- "0 treinos no app" --> WORKOUT_DAY_ONE["⚡ WORKOUT_DAY_ONE<br/>Primeiro Passo Começa Hoje"]
        CHK_ABSENCE -- "3 a 6 dias" --> COMEBACK_SHORT_ABSENCE["⚡ COMEBACK_SHORT_ABSENCE<br/>Hora de Retomar o Ritmo"]
        CHK_ABSENCE -- "7 a 13 dias" --> COMEBACK_WEEK_ABSENCE["⚡ COMEBACK_WEEK_ABSENCE<br/>Sentimos Sua Falta no Treino"]
        CHK_ABSENCE -- "≥ 14 dias" --> COMEBACK_LONG_ABSENCE["⚡ COMEBACK_LONG_ABSENCE<br/>Todo Recomeço Importa"]
        
        CHK_ABSENCE -- "0 a 2 dias" --> CHK_ROUTINE{Verificar Horário e Rotina}
        CHK_ROUTINE -- "1h antes do habitual" --> ONE_HOUR_BEFORE_HABITUAL["⚡ ONE_HOUR_BEFORE_HABITUAL<br/>Hora de se Preparar"]
        CHK_ROUTINE -- "Horário habitual atingido" --> WORKOUT_PREFERRED_TIME["⚡ WORKOUT_PREFERRED_TIME<br/>Hora do Seu Treino"]
        CHK_ROUTINE -- "Dia da rotina planejada" --> WORKOUT_SCHEDULED_TODAY["⚡ WORKOUT_SCHEDULED_TODAY<br/>Treino do Dia Planejado"]
        CHK_ROUTINE -- "15h e faltou de manhã" --> MISSED_MORNING_AFTERNOON["⚡ MISSED_MORNING_AFTERNOON<br/>Ainda Dá Tempo - Tarde"]
        CHK_ROUTINE -- "≥ 19h e ofensiva ≥ 3d" --> STREAK_AT_RISK["⚡ STREAK_AT_RISK<br/>Proteja Sua Sequência"]
        CHK_ROUTINE -- "≥ 19h e dia corrido" --> WORKOUT_QUICK_CALL["⚡ WORKOUT_QUICK_CALL<br/>Ainda Dá Tempo de Treinar"]
        CHK_ROUTINE -- "20h e não treinou" --> MISSED_DAY_NIGHT["⚡ MISSED_DAY_NIGHT<br/>Última Chamada Hoje"]
        CHK_ROUTINE -- "Músculo parado ≥ 10d" --> EVOLUTION_NEGLECTED_MUSCLE["⚡ EVOLUTION_NEGLECTED_MUSCLE<br/>Equilíbrio Muscular"]
    end

    subgraph FASE_2["2. INTERAÇÃO E ENTRADA NO APP"]
        WIFI_DETECT[Chegou na Academia] --> CHK_WIFI{Wi-Fi da academia?}
        CHK_WIFI -- "Conectou SSID (cooldown 4h)" --> GYM_WIFI_ARRIVAL["⚡ GYM_WIFI_ARRIVAL<br/>Chegou no Treino!"]
        
        APP_OPEN[Abriu o FitGuide] --> CHK_START{Iniciou treino em 30 min?}
        CHK_START -- "Saiu sem iniciar (cooldown 72h)" --> WORKOUT_OPENED_NO_WORKOUT["⚡ WORKOUT_OPENED_NO_WORKOUT<br/>A Gente Viu Você por Aqui"]
        CHK_START -- "Iniciou treino" --> WORKOUT_ACTIVE[Treino em Andamento]
    end

    subgraph FASE_3["3. DURANTE O TREINO ATIVO"]
        WORKOUT_ACTIVE --> EVT_SET[Conclusão de Série / Descanso]
        
        EVT_SET -- "Superou peso máximo histórico" --> IN_WORKOUT_PR["⚡ IN_WORKOUT_PR<br/>Novo Recorde Pessoal!"]
        EVT_SET -- "Descanso > planejado + 30s" --> IN_WORKOUT_REST_OVERRUN["⚡ IN_WORKOUT_REST_OVERRUN<br/>Descanso Finalizado"]
        EVT_SET -- "Concluiu 50% dos exercícios" --> IN_WORKOUT_HALFWAY["⚡ IN_WORKOUT_HALFWAY<br/>Metade Concluída"]
        EVT_SET -- "Carga dominada 2 sessões" --> IN_WORKOUT_LOAD_PROGRESSION["⚡ IN_WORKOUT_LOAD_PROGRESSION<br/>Sugestão de Progressão"]
        EVT_SET -- "Resta 1 exercício para o fim" --> IN_WORKOUT_LAST_EXERCISE["⚡ IN_WORKOUT_LAST_EXERCISE<br/>Reta Final do Treino"]
        
        WORKOUT_ACTIVE --> CHK_FINISH{Como encerrou?}
        CHK_FINISH -- "0 exercícios marcados (30 min)" --> WORKOUT_ABANDONED["⚡ WORKOUT_ABANDONED<br/>Treino Incompleto!"]
        CHK_FINISH -- "Treino entregue com sucesso" --> WORKOUT_DONE[Treino Finalizado com Sucesso]
    end

    subgraph FASE_4["4. PÓS-TREINO E RECONHECIMENTO"]
        WORKOUT_DONE --> ALARM_CONGRATS[Alarme de 1 hora pós-treino]
        ALARM_CONGRATS --> POST_WORKOUT_CONGRATS["⚡ POST_WORKOUT_CONGRATS<br/>Você Foi Gigante Hoje!"]
        
        WORKOUT_DONE --> CHK_ACHIEVEMENTS{Conquistas do Dia}
        CHK_ACHIEVEMENTS -- "Ofensiva ≥ 3 dias seguidos" --> STREAK_ACTIVE["⚡ STREAK_ACTIVE<br/>Ofensiva em Chamas"]
        CHK_ACHIEVEMENTS -- "Superou recorde de ofensiva" --> STREAK_NEW_RECORD["⚡ STREAK_NEW_RECORD<br/>Novo Recorde de Consistência"]
        CHK_ACHIEVEMENTS -- "Semana 100% concluída" --> WEEKLY_GOAL_COMPLETED["⚡ WEEKLY_GOAL_COMPLETED<br/>Meta Semanal Concluída"]
        CHK_ACHIEVEMENTS -- "Cruzou múltiplo 10.000 kg" --> EVOLUTION_VOLUME_MILESTONE["⚡ EVOLUTION_VOLUME_MILESTONE<br/>Marco Histórico de Volume"]
        CHK_ACHIEVEMENTS -- "Benefício anatômico mapeado" --> EVOLUTION_HEALTH_BENEFIT["⚡ EVOLUTION_HEALTH_BENEFIT<br/>Evolução e Saúde"]
    end

    class WORKOUT_DAY_ONE,COMEBACK_SHORT_ABSENCE,COMEBACK_WEEK_ABSENCE,COMEBACK_LONG_ABSENCE,ONE_HOUR_BEFORE_HABITUAL,WORKOUT_PREFERRED_TIME,WORKOUT_SCHEDULED_TODAY,MISSED_MORNING_AFTERNOON,STREAK_AT_RISK,WORKOUT_QUICK_CALL,MISSED_DAY_NIGHT,EVOLUTION_NEGLECTED_MUSCLE,GYM_WIFI_ARRIVAL,WORKOUT_OPENED_NO_WORKOUT,IN_WORKOUT_PR,IN_WORKOUT_REST_OVERRUN,IN_WORKOUT_HALFWAY,IN_WORKOUT_LOAD_PROGRESSION,IN_WORKOUT_LAST_EXERCISE,WORKOUT_ABANDONED,POST_WORKOUT_CONGRATS,STREAK_ACTIVE,STREAK_NEW_RECORD,WEEKLY_GOAL_COMPLETED,EVOLUTION_VOLUME_MILESTONE,EVOLUTION_HEALTH_BENEFIT notifNode;
    class CHK_TRAINED,CHK_ABSENCE,CHK_ROUTINE,CHK_WIFI,CHK_START,CHK_FINISH,CHK_ACHIEVEMENTS decisionNode;
    class START,WIFI_DETECT,APP_OPEN,WORKOUT_ACTIVE,EVT_SET,WORKOUT_DONE,ALARM_CONGRATS actionNode;
`
}

async function renderDiagram() {
  if (!diagramContainerRef.value) return

  try {
    mermaid.initialize({
      startOnLoad: false,
      theme: 'dark',
      securityLevel: 'loose',
      flowchart: {
        useMaxWidth: false,
        htmlLabels: true,
        curve: 'basis'
      }
    })

    const id = 'mermaid-lifecycle-' + Math.random().toString(36).substring(2, 9)
    const { svg } = await mermaid.render(id, getDiagramDefinition())
    diagramContainerRef.value.innerHTML = svg

    nextTick(() => {
      attachNodeListeners()
    })
  } catch (err) {
    diagramContainerRef.value.innerHTML = `<div class="render-error">Erro ao renderizar diagrama: ${err.message}</div>`
  }
}

function attachNodeListeners() {
  if (!diagramContainerRef.value) return

  const nodes = diagramContainerRef.value.querySelectorAll('.node')
  nodes.forEach(nodeEl => {
    const text = nodeEl.textContent || ''
    const matchedId = NOTIF_IDS.find(id => nodeEl.id.includes(id) || text.includes(id))

    if (matchedId) {
      nodeEl.classList.add('clickable-notif-node')
      nodeEl.style.cursor = 'pointer'

      nodeEl.addEventListener('click', (e) => {
        e.stopPropagation()
        selectNotificationById(matchedId)
      })
    }
  })
}

function selectNotificationById(notifId) {
  let found = null
  for (const g of props.groups) {
    const itm = g.items.find(i => i.id === notifId)
    if (itm) {
      found = { ...itm, group: g.group }
      break
    }
  }

  if (!found) {
    found = {
      id: notifId,
      title_pt: notifId.replace(/_/g, ' '),
      title_en: notifId.replace(/_/g, ' '),
      messages_pt: [],
      messages_en: [],
      group: 'workout_call'
    }
  }

  selectedNotification.value = found
  detailLang.value = 'ptbr'
  activeDrawerTab.value = 'view'
  initCrudState()

  if (diagramContainerRef.value) {
    diagramContainerRef.value.querySelectorAll('.clickable-notif-node').forEach(el => {
      const text = el.textContent || ''
      if (el.id.includes(notifId) || text.includes(notifId)) {
        el.classList.add('selected-node')
      } else {
        el.classList.remove('selected-node')
      }
    })
  }
}

function initCrudState() {
  if (!selectedNotification.value) return
  editableTitle.value = selectedNotification.value.title_pt || ''
  editableMessages.value = [...(selectedNotification.value.messages_pt || [])]
  crudStatusMessage.value = ''
  isSavingCrud.value = false
}

function openCrudTab() {
  initCrudState()
  activeDrawerTab.value = 'crud'
}

function addNewMessage() {
  editableMessages.value.push('')
}

function removeMessage(idx) {
  editableMessages.value.splice(idx, 1)
}

async function saveChanges() {
  if (!selectedNotification.value) return
  const validMessages = editableMessages.value.map(m => m.trim()).filter(Boolean)
  if (validMessages.length === 0) {
    crudStatusMessage.value = 'Adicione ao menos uma mensagem.'
    crudStatusType.value = 'error'
    return
  }
  isSavingCrud.value = true
  crudStatusMessage.value = ''
  try {
    const payload = {
      group: selectedNotification.value.group || 'workout_call',
      id: selectedNotification.value.id,
      title_pt: editableTitle.value.trim() || selectedNotification.value.title_pt,
      title_en: selectedNotification.value.title_en || '',
      messages_pt: validMessages,
      messages_en: autoTranslate.value ? [] : (selectedNotification.value.messages_en || []),
      auto_translate: autoTranslate.value
    }
    await saveNotificationItem(payload)
    
    selectedNotification.value.title_pt = payload.title_pt
    selectedNotification.value.messages_pt = validMessages
    editableMessages.value = [...validMessages]

    emit('refresh')

    crudStatusMessage.value = 'Alterações salvas e sincronizadas com sucesso!'
    crudStatusType.value = 'success'
  } catch (err) {
    crudStatusMessage.value = `Erro ao salvar: ${err.message}`
    crudStatusType.value = 'error'
  } finally {
    isSavingCrud.value = false
  }
}

function filterByPhase(phaseId) {
  selectedPhase.value = phaseId
  if (!diagramContainerRef.value) return

  const subgraphs = diagramContainerRef.value.querySelectorAll('.cluster')
  subgraphs.forEach(sub => {
    const text = sub.textContent || ''
    if (phaseId === 'all') {
      sub.style.opacity = '1'
      sub.style.filter = 'none'
    } else {
      const match = (phaseId === 'FASE_1' && text.includes('1. FORA DO APP')) ||
                    (phaseId === 'FASE_2' && text.includes('2. INTERAÇÃO')) ||
                    (phaseId === 'FASE_3' && text.includes('3. DURANTE O TREINO')) ||
                    (phaseId === 'FASE_4' && text.includes('4. PÓS-TREINO'))
      if (match) {
        sub.style.opacity = '1'
        sub.style.filter = 'drop-shadow(0 0 16px rgba(16, 185, 129, 0.4))'
      } else {
        sub.style.opacity = '0.25'
        sub.style.filter = 'grayscale(0.8)'
      }
    }
  })
}

function zoomIn() {
  zoomScale.value = Math.min(zoomScale.value + 0.15, 2.0)
}

function zoomOut() {
  zoomScale.value = Math.max(zoomScale.value - 0.15, 0.4)
}

function resetZoom() {
  zoomScale.value = 0.95
  panX.value = 40
  panY.value = 20
}

function startPan(e) {
  if (e.target.closest('.detail-drawer') || e.target.closest('.lifecycle-toolbar')) return
  isPanning = true
  startX = e.clientX - panX.value
  startY = e.clientY - panY.value
}

function onPan(e) {
  if (!isPanning) return
  panX.value = e.clientX - startX
  panY.value = e.clientY - startY
}

function endPan() {
  isPanning = false
}

function copyText(text, idx) {
  navigator.clipboard.writeText(text).then(() => {
    copiedIndex.value = idx
    setTimeout(() => {
      copiedIndex.value = null
    }, 2000)
  })
}

function emitEdit() {
  if (!selectedNotification.value) return
  emit('edit-notification', {
    item: selectedNotification.value,
    group: selectedNotification.value.group
  })
}

watch(() => props.groups, () => {
  renderDiagram()
  if (selectedNotification.value) {
    for (const g of props.groups) {
      const itm = g.items.find(i => i.id === selectedNotification.value.id)
      if (itm) {
        selectedNotification.value = { ...itm, group: g.group }
        if (activeDrawerTab.value === 'view') {
          initCrudState()
        }
        break
      }
    }
  }
}, { deep: true })

onMounted(() => {
  renderDiagram()
})
</script>

<style scoped>
.lifecycle-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.lifecycle-toolbar {
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.phase-filters {
  display: flex;
  gap: 6px;
  background: rgba(0, 0, 0, 0.25);
  padding: 4px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.phase-btn {
  padding: 6px 12px;
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
  border-radius: 4px;
  transition: all 0.2s;
}

.phase-btn:hover {
  color: var(--text-main);
}

.phase-btn.active {
  background: rgba(16, 185, 129, 0.15);
  color: var(--primary);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.toolbar-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;
  color: var(--text-muted);
}

.toolbar-hint strong {
  color: var(--primary);
}

.zoom-controls {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(0, 0, 0, 0.25);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.zoom-btn {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  color: var(--text-main);
  font-size: 0.9rem;
  font-weight: 700;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.05);
}

.zoom-btn:hover {
  background: var(--bg-card-hover);
}

.zoom-value {
  font-size: 0.78rem;
  font-family: monospace;
  color: var(--text-muted);
  min-width: 44px;
  text-align: center;
}

.lifecycle-main-layout {
  flex: 1;
  display: flex;
  position: relative;
  overflow: hidden;
}

.diagram-viewport {
  flex: 1;
  height: 100%;
  overflow: hidden;
  cursor: grab;
  position: relative;
  background: radial-gradient(circle at 50% 50%, rgba(30, 41, 59, 0.4) 0%, rgba(11, 15, 25, 0.8) 100%);
}

.diagram-viewport:active {
  cursor: grabbing;
}

.diagram-canvas {
  width: max-content;
  height: max-content;
  padding: 40px;
  transition: transform 0.05s ease-out;
}

.mermaid-target {
  user-select: none;
}

:deep(.clickable-notif-node) {
  cursor: pointer !important;
  transition: all 0.25s ease;
}

:deep(.clickable-notif-node rect),
:deep(.clickable-notif-node polygon),
:deep(.clickable-notif-node path) {
  cursor: pointer !important;
  stroke: #10b981 !important;
  stroke-width: 2.5px !important;
  fill: #132238 !important;
  transition: all 0.25s ease;
}

:deep(.clickable-notif-node:hover rect),
:deep(.clickable-notif-node:hover polygon) {
  stroke: #34d399 !important;
  stroke-width: 3.5px !important;
  filter: drop-shadow(0 0 10px rgba(16, 185, 129, 0.75)) !important;
}

:deep(.selected-node rect),
:deep(.selected-node polygon) {
  stroke: #3b82f6 !important;
  stroke-width: 4px !important;
  fill: #1e293b !important;
  filter: drop-shadow(0 0 14px rgba(59, 130, 246, 0.9)) !important;
}

:deep(.cluster rect) {
  fill: rgba(15, 23, 42, 0.6) !important;
  stroke: #24324d !important;
  rx: 12px;
}

:deep(.edgePath path) {
  stroke: #475569 !important;
  stroke-width: 1.5px !important;
}

.detail-drawer {
  width: 440px;
  background: var(--bg-card);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  z-index: 50;
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.5);
  animation: slideDrawer 0.25s ease-out;
}

.drawer-header {
  padding: 18px 20px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.drawer-header-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.drawer-group-badge {
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--primary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.drawer-id {
  font-size: 1.05rem;
  font-family: monospace;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.2;
}

.drawer-title {
  font-size: 0.95rem;
  color: #60a5fa;
  font-weight: 600;
}

.drawer-subtitle-en {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.btn-close-drawer {
  font-size: 1.6rem;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
}

.btn-close-drawer:hover {
  color: var(--text-main);
}

.drawer-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rule-card {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rule-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rule-header h4 {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-main);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.rule-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.rule-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.rule-label {
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--text-dim);
  text-transform: uppercase;
}

.rule-val {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-main);
}

.drawer-tabs {
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 3px;
  gap: 4px;
}

.drawer-tab-btn {
  flex: 1;
  padding: 8px 12px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-muted);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.drawer-tab-btn.active {
  background: var(--primary);
  color: #ffffff;
}

.crud-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.crud-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.crud-title-col {
  display: flex;
  align-items: center;
  gap: 8px;
}

.crud-title-col h4 {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-main);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.crud-field-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.crud-label {
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.crud-input-title {
  width: 100%;
  padding: 8px 10px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  font-size: 0.85rem;
  font-weight: 600;
  outline: none;
  transition: border-color 0.2s;
}

.crud-input-title:focus {
  border-color: var(--primary);
}

.crud-messages-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.crud-message-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 6px 8px;
}

.crud-item-num {
  font-size: 0.74rem;
  font-weight: 700;
  color: var(--primary);
  padding-top: 6px;
  min-width: 16px;
}

.crud-textarea {
  flex: 1;
  resize: vertical;
  min-height: 48px;
  font-size: 0.83rem;
  line-height: 1.4;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-main);
  padding: 6px 8px;
  font-family: inherit;
  outline: none;
}

.crud-textarea:focus {
  border-color: var(--primary);
}

.btn-delete-msg {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px;
  font-size: 0.85rem;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.btn-delete-msg:hover {
  opacity: 1;
}

.btn-add-msg {
  align-self: flex-start;
  font-size: 0.78rem;
  padding: 6px 12px;
}

.crud-options-card {
  background: rgba(0, 0, 0, 0.25);
  border: 1px dashed var(--border);
  border-radius: var(--radius-sm);
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.crud-toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-main);
  cursor: pointer;
}

.crud-help-text {
  font-size: 0.72rem;
  color: var(--text-muted);
  line-height: 1.35;
}

.crud-status-banner {
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 500;
}

.crud-status-banner.status-success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid var(--primary);
  color: #34d399;
}

.crud-status-banner.status-error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid #ef4444;
  color: #f87171;
}

.crud-save-row {
  margin-top: 4px;
}

.messages-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.messages-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.messages-title-col {
  display: flex;
  align-items: center;
  gap: 8px;
}

.messages-title-col h4 {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-main);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.count-badge {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  font-size: 0.72rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 999px;
}

.lang-mini-switch {
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 2px;
}

.mini-btn {
  padding: 2px 8px;
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--text-muted);
  border-radius: 3px;
}

.mini-btn.active {
  background: var(--bg-card-hover);
  color: var(--primary);
}

.drawer-messages-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 380px;
  overflow-y: auto;
  padding-right: 4px;
}

.drawer-msg-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 10px 12px;
}

.drawer-msg-num {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-dim);
  padding-top: 1px;
}

.drawer-msg-text {
  flex: 1;
  font-size: 0.86rem;
  color: var(--text-main);
  line-height: 1.4;
}

.btn-copy {
  font-size: 0.85rem;
  color: var(--text-muted);
  cursor: pointer;
  padding: 2px;
  transition: color 0.2s;
}

.btn-copy:hover {
  color: var(--primary);
}

.drawer-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.2);
}

.btn-full {
  width: 100%;
}

@keyframes slideDrawer {
  from {
    transform: translateX(100%);
  }
  to {
    transform: translateX(0);
  }
}

.drawer-enter-active,
.drawer-leave-active {
  transition: transform 0.25s ease-out;
}

.drawer-enter-from,
.drawer-leave-to {
  transform: translateX(100%);
}

@media (max-width: 900px) {
  .detail-drawer {
    width: 100%;
    position: absolute;
    inset: 0;
  }
}
</style>

