<template>
  <div class="notifications-page">
    <header class="page-header">
      <div class="header-container">
        <div class="header-left">
          <div class="title-with-badge">
            <h1>Notificações</h1>
            <span class="badge-count">{{ totalTypes }} tipos</span>
            <span class="badge-count badge-messages">{{ totalMessages }} variações</span>
          </div>
          <p class="subtitle">Catálogo centralizado e mapeamento do ciclo de vida das mensagens</p>
        </div>

        <div class="header-actions">
          <div class="view-mode-switch">
            <button 
              :class="['mode-btn', { active: currentViewMode === 'catalog' }]" 
              @click="currentViewMode = 'catalog'"
            >
              📋 Catálogo
            </button>
            <button 
              :class="['mode-btn', { active: currentViewMode === 'lifecycle' }]" 
              @click="currentViewMode = 'lifecycle'"
            >
              🧭 Ciclo de Vida
            </button>
          </div>

          <div v-if="currentViewMode === 'catalog'" class="lang-selector">
            <button 
              :class="['lang-btn', { active: currentLangView === 'ptbr' }]" 
              @click="currentLangView = 'ptbr'"
            >
              🇧🇷 PT-BR
            </button>
            <button 
              :class="['lang-btn', { active: currentLangView === 'en' }]" 
              @click="currentLangView = 'en'"
            >
              🇺🇸 EN
            </button>
            <button 
              :class="['lang-btn', { active: currentLangView === 'both' }]" 
              @click="currentLangView = 'both'"
            >
              ⚖️ Lado a Lado
            </button>
          </div>

          <button 
            class="btn btn-secondary" 
            @click="handleTranslateAll" 
            :disabled="translating"
          >
            {{ translating ? '🌐 Traduzindo...' : '🌐 Traduzir Tudo' }}
          </button>

          <button class="btn btn-primary" @click="openCreateModal()">
            + Nova Notificação
          </button>
        </div>
      </div>
    </header>

    <div v-if="toastMessage" :class="['toast-notification', `toast-${toastType}`]">
      {{ toastMessage }}
    </div>

    <main class="main-content">
      <div v-if="currentViewMode === 'lifecycle'" class="lifecycle-view-wrapper">
        <LifecycleDiagram 
          :groups="groups" 
          @edit-notification="onEditFromDiagram" 
          @refresh="loadData"
        />
      </div>

      <div v-else class="catalog-view-wrapper">
        <div class="toolbar">
          <div class="search-box">
            <span class="search-icon">🔍</span>
            <input 
              type="text" 
              v-model="searchQuery" 
              placeholder="Buscar por ID, título ou conteúdo de mensagem..."
            />
            <button v-if="searchQuery" class="clear-search" @click="searchQuery = ''">✕</button>
          </div>

          <div class="toolbar-stats">
            Exibindo <strong>{{ filteredGroupsCount }}</strong> grupos e <strong>{{ filteredTypesCount }}</strong> tipos
          </div>
        </div>

        <div v-if="loading && !groups.length" class="loading-state">
          <div class="spinner"></div>
          <p>Carregando notificações...</p>
        </div>

        <div v-else-if="!filteredGroups.length" class="empty-state">
          <div class="empty-icon">🔔</div>
          <h3>Nenhuma notificação encontrada</h3>
          <p>Tente ajustar a busca ou adicione uma nova notificação.</p>
        </div>

        <div v-else class="groups-container">
          <div 
            v-for="group in filteredGroups" 
            :key="group.group" 
            class="group-card"
          >
            <div class="group-header" @click="toggleGroup(group.group)">
              <div class="group-title-col">
                <span class="group-expand-icon">
                  {{ isGroupExpanded(group.group) ? '▼' : '▶' }}
                </span>
                <span class="group-icon">📦</span>
                <div class="group-info">
                  <h3>{{ getGroupTitle(group) }}</h3>
                  <span class="group-key">{{ group.group }}</span>
                </div>
              </div>

              <div class="group-badges">
                <span class="badge-pill">{{ group.items.length }} tipos</span>
                <span class="badge-pill badge-pill-dim">{{ getGroupTotalMessages(group) }} msgs</span>
              </div>

              <div class="group-actions" @click.stop>
                <button 
                  class="btn btn-secondary btn-sm" 
                  @click="openCreateModal(group.group)"
                  title="Adicionar item neste grupo"
                >
                  + Adicionar
                </button>
                <button 
                  class="btn btn-secondary btn-sm" 
                  @click="handleTranslateGroup(group.group)"
                  title="Traduzir este grupo"
                  :disabled="translating"
                >
                  🌐 Traduzir
                </button>
              </div>
            </div>

            <div v-show="isGroupExpanded(group.group)" class="group-items-list">
              <div 
                v-for="item in group.items" 
                :key="item.id" 
                class="item-card"
              >
                <div class="item-header">
                  <div class="item-title-col">
                    <span class="item-id-badge">{{ item.id }}</span>
                    <h4 class="item-title">{{ getItemTitle(item) }}</h4>
                  </div>

                  <div class="item-actions">
                    <span class="item-count-badge">
                      {{ (item.messages_pt?.length || item.messages?.length || 0) }} variações
                    </span>
                    <button 
                      class="btn-action btn-edit" 
                      @click="openEditModal(item, group.group)"
                      title="Editar notificação"
                    >
                      ✏️
                    </button>
                    <button 
                      class="btn-action btn-delete" 
                      @click="handleDeleteItem(group.group, item.id)"
                      title="Excluir notificação"
                    >
                      🗑️
                    </button>
                  </div>
                </div>

                <div v-if="currentLangView === 'both'" class="dual-messages-grid">
                  <div class="lang-column">
                    <div class="lang-col-header">
                      <span>🇧🇷 Português (PT-BR)</span>
                      <span class="col-count">{{ item.messages_pt?.length || 0 }}</span>
                    </div>
                    <ul class="messages-sublist">
                      <li v-for="(msg, mIdx) in (item.messages_pt || [])" :key="'pt-' + mIdx">
                        <span class="msg-index">{{ mIdx + 1 }}.</span>
                        <span class="msg-text">{{ msg }}</span>
                      </li>
                    </ul>
                  </div>

                  <div class="lang-column">
                    <div class="lang-col-header">
                      <span>🇺🇸 Inglês (EN)</span>
                      <span class="col-count">{{ item.messages_en?.length || 0 }}</span>
                    </div>
                    <ul class="messages-sublist">
                      <li v-for="(msg, mIdx) in (item.messages_en || [])" :key="'en-' + mIdx">
                        <span class="msg-index">{{ mIdx + 1 }}.</span>
                        <span class="msg-text">{{ msg }}</span>
                      </li>
                    </ul>
                  </div>
                </div>

                <div v-else class="single-messages-list">
                  <ul class="messages-sublist">
                    <li 
                      v-for="(msg, mIdx) in getMessagesForCurrentLang(item)" 
                      :key="mIdx"
                    >
                      <span class="msg-index">{{ mIdx + 1 }}.</span>
                      <span class="msg-text">{{ msg }}</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <NotificationModal 
      v-if="showModal" 
      :notification="selectedItem" 
      :group="selectedGroup" 
      :groups="groups" 
      @close="showModal = false" 
      @save="onSaveItem" 
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  fetchNotifications, 
  saveNotificationItem, 
  deleteNotificationItem, 
  translateNotifications 
} from '../api'
import NotificationModal from '../components/NotificationModal.vue'
import LifecycleDiagram from '../components/LifecycleDiagram.vue'

const currentViewMode = ref('catalog')
const groups = ref([])
const loading = ref(false)
const translating = ref(false)
const searchQuery = ref('')
const currentLangView = ref('ptbr')

const expandedGroups = ref({})

const showModal = ref(false)
const selectedItem = ref(null)
const selectedGroup = ref('')

const toastMessage = ref('')
const toastType = ref('success')
let toastTimer = null

function showToast(msg, type = 'success') {
  toastMessage.value = msg
  toastType.value = type
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 4000)
}

async function loadData() {
  loading.value = true
  try {
    const raw = await fetchNotifications()
    groups.value = raw.map(g => {
      const ptItems = g.pt?.items || []
      const enItems = g.en?.items || []
      const enMap = new Map(enItems.map(item => [item.id, item]))

      const mergedItems = ptItems.map(ptItem => {
        const enItem = enMap.get(ptItem.id) || { title: '', messages: [] }
        return {
          id: ptItem.id,
          title_pt: ptItem.title,
          title_en: enItem.title,
          messages_pt: ptItem.messages || [],
          messages_en: enItem.messages || []
        }
      })

      return {
        group: g.group,
        name_pt: g.name_pt || g.pt?.name || g.group,
        name_en: g.name_en || g.en?.name || g.group,
        items: mergedItems
      }
    })

    groups.value.forEach(g => {
      if (expandedGroups.value[g.group] === undefined) {
        expandedGroups.value[g.group] = true
      }
    })
  } catch (err) {
    showToast(`Erro ao carregar notificações: ${err.message}`, 'error')
  } finally {
    loading.value = false
  }
}

const totalTypes = computed(() => {
  return groups.value.reduce((acc, g) => acc + g.items.length, 0)
})

const totalMessages = computed(() => {
  return groups.value.reduce((acc, g) => {
    return acc + g.items.reduce((sum, item) => sum + (item.messages_pt?.length || 0), 0)
  }, 0)
})

const filteredGroups = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return groups.value

  return groups.value.map(g => {
    const matchedItems = g.items.filter(item => {
      if (item.id.toLowerCase().includes(query)) return true
      if (item.title_pt?.toLowerCase().includes(query)) return true
      if (item.title_en?.toLowerCase().includes(query)) return true
      const hasPtMsg = item.messages_pt?.some(m => m.toLowerCase().includes(query))
      const hasEnMsg = item.messages_en?.some(m => m.toLowerCase().includes(query))
      return hasPtMsg || hasEnMsg
    })

    return {
      ...g,
      items: matchedItems
    }
  }).filter(g => g.items.length > 0)
})

const filteredGroupsCount = computed(() => filteredGroups.value.length)

const filteredTypesCount = computed(() => {
  return filteredGroups.value.reduce((acc, g) => acc + g.items.length, 0)
})

function isGroupExpanded(groupKey) {
  return !!expandedGroups.value[groupKey]
}

function toggleGroup(groupKey) {
  expandedGroups.value[groupKey] = !expandedGroups.value[groupKey]
}

function getGroupTitle(group) {
  if (currentLangView.value === 'en') return group.name_en || group.group
  if (currentLangView.value === 'both') return `${group.name_pt} / ${group.name_en}`
  return group.name_pt || group.group
}

function getGroupTotalMessages(group) {
  return group.items.reduce((acc, item) => acc + (item.messages_pt?.length || 0), 0)
}

function getItemTitle(item) {
  if (currentLangView.value === 'en') return item.title_en || item.title_pt
  if (currentLangView.value === 'both') return `${item.title_pt} — ${item.title_en || 'Sem tradução'}`
  return item.title_pt
}

function getMessagesForCurrentLang(item) {
  if (currentLangView.value === 'en') return item.messages_en || []
  return item.messages_pt || []
}

function openCreateModal(groupKey = '') {
  selectedItem.value = null
  selectedGroup.value = groupKey || (groups.value[0]?.group || 'workout_call')
  showModal.value = true
}

function openEditModal(item, groupKey) {
  selectedItem.value = item
  selectedGroup.value = groupKey
  showModal.value = true
}

function onEditFromDiagram({ item, group }) {
  selectedItem.value = item
  selectedGroup.value = group
  showModal.value = true
}

async function onSaveItem(payload) {
  try {
    await saveNotificationItem(payload)
    showModal.value = false
    showToast(`Notificação "${payload.id}" salva com sucesso!`)
    await loadData()
  } catch (err) {
    showToast(`Erro ao salvar: ${err.message}`, 'error')
  }
}

async function handleDeleteItem(groupKey, id) {
  if (!confirm(`Tem certeza que deseja remover a notificação "${id}"?`)) {
    return
  }
  try {
    await deleteNotificationItem(groupKey, id)
    showToast(`Notificação "${id}" removida com sucesso.`)
    await loadData()
  } catch (err) {
    showToast(`Erro ao remover: ${err.message}`, 'error')
  }
}

async function handleTranslateGroup(groupKey) {
  translating.value = true
  try {
    await translateNotifications(groupKey)
    showToast(`Grupo "${groupKey}" traduzido com sucesso!`)
    await loadData()
  } catch (err) {
    showToast(`Erro ao traduzir: ${err.message}`, 'error')
  } finally {
    translating.value = false
  }
}

async function handleTranslateAll() {
  translating.value = true
  try {
    await translateNotifications('')
    showToast('Todas as notificações foram traduzidas e sincronizadas com sucesso!')
    await loadData()
  } catch (err) {
    showToast(`Erro ao traduzir: ${err.message}`, 'error')
  } finally {
    translating.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.notifications-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 18px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title-with-badge {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-with-badge h1 {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--text-main);
}

.badge-count {
  background: rgba(16, 185, 129, 0.15);
  color: var(--primary);
  border: 1px solid rgba(16, 185, 129, 0.3);
  font-size: 0.75rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
}

.badge-messages {
  background: rgba(59, 130, 246, 0.15);
  color: var(--accent);
  border-color: rgba(59, 130, 246, 0.3);
}

.subtitle {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.view-mode-switch {
  display: flex;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 3px;
}

.mode-btn {
  padding: 6px 14px;
  font-size: 0.84rem;
  font-weight: 700;
  color: var(--text-muted);
  border-radius: 4px;
  transition: all 0.2s;
}

.mode-btn:hover {
  color: var(--text-main);
}

.mode-btn.active {
  background: var(--bg-card-hover);
  color: var(--primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.lang-selector {
  display: flex;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 2px;
}

.lang-btn {
  padding: 6px 12px;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-muted);
  border-radius: 4px;
  transition: all 0.2s;
}

.lang-btn:hover {
  color: var(--text-main);
}

.lang-btn.active {
  background: var(--bg-card-hover);
  color: var(--primary);
  border: 1px solid var(--border);
}

.main-content {
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.lifecycle-view-wrapper {
  width: 100%;
  height: 100%;
}

.catalog-view-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 6px 12px;
  width: 100%;
  max-width: 480px;
}

.search-box input {
  background: transparent;
  border: none;
  padding: 0;
  flex: 1;
  font-size: 0.9rem;
}

.search-icon {
  font-size: 0.9rem;
  color: var(--text-dim);
}

.clear-search {
  font-size: 0.85rem;
  color: var(--text-muted);
  cursor: pointer;
}

.toolbar-stats {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.groups-container {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.group-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.group-header {
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.01);
  transition: background 0.2s;
}

.group-header:hover {
  background: rgba(255, 255, 255, 0.03);
}

.group-title-col {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-expand-icon {
  font-size: 0.75rem;
  color: var(--text-dim);
  width: 14px;
}

.group-icon {
  font-size: 1.3rem;
}

.group-info h3 {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.2;
}

.group-key {
  font-size: 0.75rem;
  font-family: monospace;
  color: var(--text-dim);
}

.group-badges {
  display: flex;
  align-items: center;
  gap: 8px;
}

.badge-pill {
  font-size: 0.75rem;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border);
  color: var(--text-main);
  padding: 3px 10px;
  border-radius: 999px;
}

.badge-pill-dim {
  color: var(--text-muted);
}

.group-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-items-list {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  border-top: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.15);
}

.item-card {
  background: var(--bg-main);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.item-title-col {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.item-id-badge {
  font-size: 0.75rem;
  font-family: monospace;
  font-weight: 700;
  background: rgba(59, 130, 246, 0.12);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.25);
  padding: 2px 8px;
  border-radius: 4px;
}

.item-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-main);
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-count-badge {
  font-size: 0.75rem;
  color: var(--text-dim);
  background: rgba(255, 255, 255, 0.03);
  padding: 2px 8px;
  border-radius: 4px;
}

.btn-action {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-action:hover {
  background: var(--bg-card-hover);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: var(--danger);
}

.dual-messages-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  padding-top: 6px;
  border-top: 1px dashed rgba(255, 255, 255, 0.07);
}

.lang-column {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.lang-col-header {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-dim);
  display: flex;
  align-items: center;
  justify-content: space-between;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.col-count {
  background: rgba(255, 255, 255, 0.05);
  padding: 1px 6px;
  border-radius: 999px;
}

.single-messages-list {
  padding-top: 4px;
  border-top: 1px dashed rgba(255, 255, 255, 0.07);
}

.messages-sublist {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.messages-sublist li {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 0.86rem;
  color: var(--text-main);
  line-height: 1.4;
  background: rgba(255, 255, 255, 0.015);
  padding: 6px 10px;
  border-radius: 4px;
}

.msg-index {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-dim);
  flex-shrink: 0;
  padding-top: 1px;
}

.msg-text {
  flex: 1;
}

.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 12px;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 2.5rem;
  opacity: 0.6;
}

.toast-notification {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 14px 20px;
  border-radius: var(--radius-sm);
  font-size: 0.9rem;
  font-weight: 600;
  z-index: 1000;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
}

.toast-success {
  background: #064e3b;
  color: #6ee7b7;
  border: 1px solid #059669;
}

.toast-error {
  background: #7f1d1d;
  color: #fca5a5;
  border: 1px solid #dc2626;
}

@media (max-width: 768px) {
  .dual-messages-grid {
    grid-template-columns: 1fr;
  }
}
</style>
