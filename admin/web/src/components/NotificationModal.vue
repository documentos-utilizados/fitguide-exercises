<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-card">
      <div class="modal-header">
        <h3>{{ isEdit ? 'Editar Notificação' : 'Nova Notificação' }}</h3>
        <button class="btn-close" @click="$emit('close')">&times;</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-row">
          <div class="form-group">
            <label>Grupo *</label>
            <select v-model="form.group" :disabled="isEdit" required>
              <option v-for="g in groups" :key="g.group" :value="g.group">
                {{ g.name_pt || g.group }} ({{ g.group }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>ID do Tipo (Constante) *</label>
            <input 
              type="text" 
              v-model="form.id" 
              placeholder="Ex: WORKOUT_CUSTOM_ALERT"
              :disabled="isEdit"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label>Título (PT-BR) *</label>
          <input 
            type="text" 
            v-model="form.title_pt" 
            placeholder="Ex: Hora de Treinar!" 
            required
          />
        </div>

        <div class="form-group">
          <div class="section-title-row">
            <label>Mensagens / Variações (PT-BR) *</label>
            <button type="button" class="btn btn-secondary btn-xs" @click="addMessagePt">
              + Adicionar Variação
            </button>
          </div>

          <div class="messages-list">
            <div 
              v-for="(msg, idx) in form.messages_pt" 
              :key="idx" 
              class="message-input-row"
            >
              <span class="msg-num">{{ idx + 1 }}</span>
              <input 
                type="text" 
                v-model="form.messages_pt[idx]" 
                placeholder="Ex: Bora treinar que o shape não espera!"
                required
              />
              <button 
                type="button" 
                class="btn-icon btn-danger" 
                @click="removeMessagePt(idx)"
                :disabled="form.messages_pt.length <= 1"
              >
                ✕
              </button>
            </div>
          </div>
        </div>

        <div class="auto-translate-box">
          <label class="checkbox-label">
            <input type="checkbox" v-model="form.auto_translate" />
            <span>Traduzir automaticamente para Inglês (EN) pelo script ao salvar</span>
          </label>
        </div>

        <div v-if="!form.auto_translate" class="en-manual-section">
          <div class="form-group">
            <label>Título (EN)</label>
            <input 
              type="text" 
              v-model="form.title_en" 
              placeholder="Ex: Time to Train!" 
            />
          </div>

          <div class="form-group">
            <div class="section-title-row">
              <label>Mensagens (EN)</label>
              <button type="button" class="btn btn-secondary btn-xs" @click="addMessageEn">
                + Adicionar Variação EN
              </button>
            </div>

            <div class="messages-list">
              <div 
                v-for="(msg, idx) in form.messages_en" 
                :key="idx" 
                class="message-input-row"
              >
                <span class="msg-num">{{ idx + 1 }}</span>
                <input 
                  type="text" 
                  v-model="form.messages_en[idx]" 
                  placeholder="Ex: Time to work out, gains are waiting!"
                />
                <button 
                  type="button" 
                  class="btn-icon btn-danger" 
                  @click="removeMessageEn(idx)"
                >
                  ✕
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            Cancelar
          </button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            {{ saving ? 'Salvando...' : (isEdit ? 'Salvar Alterações' : 'Criar Notificação') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  notification: {
    type: Object,
    default: null
  },
  group: {
    type: String,
    default: ''
  },
  groups: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'save'])

const isEdit = computed(() => !!props.notification)
const saving = ref(false)

const form = ref({
  group: props.group || (props.groups[0]?.group || 'workout_call'),
  id: props.notification?.id || '',
  original_id: props.notification?.id || '',
  title_pt: props.notification?.title_pt || props.notification?.title || '',
  title_en: props.notification?.title_en || '',
  messages_pt: props.notification?.messages_pt?.length 
    ? [...props.notification.messages_pt] 
    : (props.notification?.messages?.length ? [...props.notification.messages] : ['']),
  messages_en: props.notification?.messages_en?.length 
    ? [...props.notification.messages_en] 
    : [''],
  auto_translate: true
})

function addMessagePt() {
  form.value.messages_pt.push('')
}

function removeMessagePt(idx) {
  if (form.value.messages_pt.length > 1) {
    form.value.messages_pt.splice(idx, 1)
  }
}

function addMessageEn() {
  form.value.messages_en.push('')
}

function removeMessageEn(idx) {
  form.value.messages_en.splice(idx, 1)
}

function handleSubmit() {
  form.value.id = form.value.id.trim().toUpperCase().replace(/[^A-Z0-9_]/g, '_')
  form.value.messages_pt = form.value.messages_pt.map(m => m.trim()).filter(Boolean)
  form.value.messages_en = form.value.messages_en.map(m => m.trim()).filter(Boolean)

  if (!form.value.messages_pt.length) {
    form.value.messages_pt = ['']
    return
  }

  saving.value = true
  emit('save', { ...form.value })
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 680px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5);
}

.modal-header {
  padding: 18px 24px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h3 {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--text-main);
}

.btn-close {
  font-size: 1.6rem;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
}

.btn-close:hover {
  color: var(--text-main);
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-muted);
}

.section-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-xs {
  font-size: 0.75rem;
  padding: 4px 8px;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 240px;
  overflow-y: auto;
  padding-right: 4px;
}

.message-input-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.msg-num {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-dim);
  width: 22px;
  text-align: right;
  flex-shrink: 0;
}

.message-input-row input {
  flex: 1;
}

.btn-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
  cursor: pointer;
}

.auto-translate-box {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-size: 0.88rem;
  font-weight: 500;
  color: var(--text-main);
}

.checkbox-label input {
  width: 18px;
  height: 18px;
  accent-color: var(--primary);
  cursor: pointer;
}

.en-manual-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 16px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
