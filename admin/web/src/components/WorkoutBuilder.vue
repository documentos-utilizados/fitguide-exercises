<template>
  <div class="workout-section">
    <div class="workout-nav">
      <div class="tabs">
        <button 
          :class="['tab-btn', activeTab === 'templates' ? 'active' : '']"
          @click="activeTab = 'templates'"
        >
          Fichas / Templates ({{ workouts.templates?.length || 0 }})
        </button>
        <button 
          :class="['tab-btn', activeTab === 'seeds' ? 'active' : '']"
          @click="activeTab = 'seeds'"
        >
          Treinos Iniciais / Seeds ({{ workouts.seeds?.length || 0 }})
        </button>
      </div>

      <button class="btn btn-primary" @click="openCreateModal">
        + Novo Treino ({{ activeTab === 'seeds' ? 'Seed' : 'Template' }})
      </button>
    </div>

    <div class="workout-grid">
      <div 
        v-for="w in currentList" 
        :key="w.id" 
        class="workout-card"
      >
        <div class="workout-card-header">
          <div>
            <h4>{{ formatWorkoutTitle(w) }}</h4>
            <span class="workout-id">{{ w.id }}</span>
          </div>
          <div class="card-actions">
            <button class="btn-icon" title="Editar" @click="openEditModal(w)">✏️</button>
            <button class="btn-icon btn-icon-delete" title="Excluir" @click="confirmDelete(w)">🗑️</button>
          </div>
        </div>

        <div class="workout-badges">
          <span class="badge badge-primary">{{ w.category_key }}</span>
          <span class="badge badge-blue">⏱️ {{ w.estimated_duration_min }} min</span>
          <span v-if="w.days?.length" class="badge badge-gray">
            📅 {{ w.days.join(', ') }}
          </span>
        </div>

        <div class="exercises-preview-list">
          <div 
            v-for="(ex, idx) in w.exercises" 
            :key="idx" 
            class="ex-item"
          >
            <span class="ex-name">{{ getExerciseName(ex.exercise_id) }}</span>
            <span class="ex-meta">
              {{ ex.is_time_based ? `${ex.sets}x ${ex.time_seconds}s` : `${ex.sets}x ${ex.reps}` }}
              <template v-if="ex.weight_kg"> ({{ ex.weight_kg }}kg)</template>
            </span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-backdrop" @click.self="showModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <h3>{{ isEditMode ? 'Editar Treino' : 'Novo Treino' }}</h3>
          <button class="btn-close" @click="showModal = false">&times;</button>
        </div>

        <form @submit.prevent="handleSaveWorkout" class="modal-body">
          <div class="form-grid">
            <div class="form-group">
              <label>ID do Treino (Único)</label>
              <input 
                type="text" 
                v-model="modalForm.id" 
                :disabled="isEditMode" 
                placeholder="Ex: template_ombros"
                required
              />
            </div>

            <div class="form-group">
              <label>Chave de Tradução do Nome (name_key)</label>
              <input 
                type="text" 
                v-model="modalForm.name_key" 
                placeholder="Ex: template_shoulders"
                required
              />
            </div>

            <div class="form-group">
              <label>Categoria do Treino</label>
              <select v-model="modalForm.category_key" required>
                <option value="cat_hypertrophy">Hipertrofia (cat_hypertrophy)</option>
                <option value="cat_strength">Força (cat_strength)</option>
                <option value="cat_endurance">Resistência (cat_endurance)</option>
                <option value="cat_high_intensity">Alta Intensidade (cat_high_intensity)</option>
              </select>
            </div>

            <div class="form-group">
              <label>Duração Estimada (min)</label>
              <input 
                type="number" 
                v-model.number="modalForm.estimated_duration_min" 
                min="5" 
                max="240"
                required
              />
            </div>
          </div>

          <div v-if="activeTab === 'seeds'" class="form-group">
            <label>Dias da Semana</label>
            <div class="days-selector">
              <label v-for="d in ['MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT', 'SUN']" :key="d" class="day-chip">
                <input type="checkbox" :value="d" v-model="modalForm.days" />
                <span>{{ d }}</span>
              </label>
            </div>
          </div>

          <div class="exercises-builder">
            <div class="builder-header">
              <label>Exercícios do Treino ({{ modalForm.exercises.length }})</label>
              
              <div class="add-exercise-bar">
                <select v-model="selectedExToAdd">
                  <option value="" disabled>Selecione um exercício para adicionar...</option>
                  <option v-for="ex in exercisesList" :key="ex.id" :value="ex.id">
                    {{ ex.name }} ({{ ex.id }})
                  </option>
                </select>
                <button 
                  type="button" 
                  class="btn btn-primary btn-sm" 
                  :disabled="!selectedExToAdd"
                  @click="addExerciseToWorkout"
                >
                  + Adicionar
                </button>
              </div>
            </div>

            <div v-if="modalForm.exercises.length === 0" class="empty-builder">
              Selecione exercícios acima para montar esta ficha de treino.
            </div>

            <div v-for="(item, idx) in modalForm.exercises" :key="idx" class="exercise-row-card">
              <div class="row-top">
                <span class="row-num">{{ idx + 1 }}.</span>
                <span class="row-title">{{ getExerciseName(item.exercise_id) }}</span>
                <div class="row-actions">
                  <button type="button" class="btn-icon" :disabled="idx === 0" @click="moveEx(idx, -1)">⬆️</button>
                  <button type="button" class="btn-icon" :disabled="idx === modalForm.exercises.length - 1" @click="moveEx(idx, 1)">⬇️</button>
                  <button type="button" class="btn-icon btn-icon-delete" @click="removeEx(idx)">&times;</button>
                </div>
              </div>

              <div class="row-inputs">
                <div class="param-input">
                  <label>Séries</label>
                  <input type="number" v-model.number="item.sets" min="1" max="20" required />
                </div>

                <div class="param-input" v-if="!item.is_time_based">
                  <label>Reps</label>
                  <input type="text" v-model="item.reps" placeholder="Ex: 10-12" required />
                </div>

                <div class="param-input" v-if="item.is_time_based">
                  <label>Tempo (s)</label>
                  <input type="number" v-model.number="item.time_seconds" min="5" max="3600" required />
                </div>

                <div class="param-input">
                  <label>Carga (kg)</label>
                  <input type="number" step="0.5" v-model.number="item.weight_kg" min="0" />
                </div>

                <div class="param-input">
                  <label>Descanso (s)</label>
                  <input type="number" v-model.number="item.rest_time_seconds" min="0" step="5" required />
                </div>

                <div class="param-input-checkbox">
                  <label>
                    <input type="checkbox" v-model="item.is_time_based" />
                    Tempo?
                  </label>
                </div>
              </div>
            </div>
          </div>

          <div v-if="workoutError" class="error-box">
            {{ workoutError }}
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showModal = false">Cancelar</button>
            <button type="submit" class="btn btn-primary" :disabled="isSaving">
              {{ isSaving ? 'Salvando...' : 'Salvar Treino' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { saveWorkout, deleteWorkout } from '../api'

const props = defineProps({
  workouts: {
    type: Object,
    default: () => ({ seeds: [], templates: [] })
  },
  exercisesList: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['workout-updated'])

const activeTab = ref('templates')
const showModal = ref(false)
const isEditMode = ref(false)
const isSaving = ref(false)
const workoutError = ref('')
const selectedExToAdd = ref('')

const currentList = computed(() => {
  return props.workouts[activeTab.value] || []
})

const exerciseMap = computed(() => {
  const map = {}
  for (const ex of props.exercisesList) {
    map[ex.id] = ex.name
  }
  return map
})

function getExerciseName(id) {
  return exerciseMap.value[id] || id
}

function formatWorkoutTitle(w) {
  return w.name_key.replace('workout_', '').replace('template_', '').replace(/_/g, ' ').toUpperCase()
}

const modalForm = reactive({
  id: '',
  name_key: '',
  category_key: 'cat_hypertrophy',
  days: [],
  estimated_duration_min: 45,
  cover_image_url: '',
  type: 'templates',
  exercises: []
})

function openCreateModal() {
  isEditMode.value = false
  workoutError.value = ''
  selectedExToAdd.value = ''
  
  modalForm.id = ''
  modalForm.name_key = ''
  modalForm.category_key = 'cat_hypertrophy'
  modalForm.days = []
  modalForm.estimated_duration_min = 45
  modalForm.cover_image_url = ''
  modalForm.type = activeTab.value
  modalForm.exercises = []

  showModal.value = true
}

function openEditModal(w) {
  isEditMode.value = true
  workoutError.value = ''
  selectedExToAdd.value = ''

  modalForm.id = w.id
  modalForm.name_key = w.name_key
  modalForm.category_key = w.category_key
  modalForm.days = w.days ? [...w.days] : []
  modalForm.estimated_duration_min = w.estimated_duration_min || 45
  modalForm.cover_image_url = w.cover_image_url || ''
  modalForm.type = activeTab.value
  modalForm.exercises = (w.exercises || []).map(e => ({ ...e }))

  showModal.value = true
}

function addExerciseToWorkout() {
  if (!selectedExToAdd.value) return
  modalForm.exercises.push({
    exercise_id: selectedExToAdd.value,
    sets: 4,
    reps: '10-12',
    weight_kg: 0.0,
    rest_time_seconds: 60,
    is_time_based: false,
    time_seconds: 0
  })
  selectedExToAdd.value = ''
}

function removeEx(idx) {
  modalForm.exercises.splice(idx, 1)
}

function moveEx(idx, offset) {
  const target = idx + offset
  if (target < 0 || target >= modalForm.exercises.length) return
  const temp = modalForm.exercises[idx]
  modalForm.exercises[idx] = modalForm.exercises[target]
  modalForm.exercises[target] = temp
}

async function handleSaveWorkout() {
  if (modalForm.exercises.length === 0) {
    workoutError.value = 'Adicione ao menos 1 exercício ao treino.'
    return
  }

  isSaving.value = true
  workoutError.value = ''

  try {
    const payload = {
      ...modalForm,
      id: modalForm.id.trim(),
      name_key: modalForm.name_key.trim(),
      type: activeTab.value
    }

    await saveWorkout(payload)
    showModal.value = false
    emit('workout-updated')
  } catch (err) {
    workoutError.value = err.message || 'Erro ao salvar treino'
  } finally {
    isSaving.value = false
  }
}

async function confirmDelete(w) {
  if (confirm(`Tem certeza que deseja excluir o treino "${w.id}"?`)) {
    try {
      await deleteWorkout(activeTab.value, w.id)
      emit('workout-updated')
    } catch (err) {
      alert(`Erro: ${err.message}`)
    }
  }
}
</script>

<style scoped>
.workout-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.workout-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.tabs {
  display: flex;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 4px;
}

.tab-btn {
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-muted);
}
.tab-btn.active {
  background: var(--primary);
  color: #041f18;
  font-weight: 600;
}

.workout-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.workout-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.workout-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.workout-card-header h4 {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--text-main);
}

.workout-id {
  font-size: 0.75rem;
  color: var(--text-dim);
  font-family: monospace;
}

.workout-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.exercises-preview-list {
  background: #0d1322;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 200px;
  overflow-y: auto;
}

.ex-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.05);
  padding-bottom: 4px;
}
.ex-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.ex-name {
  color: var(--text-main);
  font-weight: 500;
}

.ex-meta {
  color: var(--text-muted);
  font-size: 0.75rem;
}

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  z-index: 100;
}

.modal-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 800px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border);
}

.modal-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
}

.btn-close {
  font-size: 1.5rem;
  color: var(--text-muted);
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
}

.days-selector {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.day-chip {
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  cursor: pointer;
}

.exercises-builder {
  border-top: 1px solid var(--border);
  padding-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.builder-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.add-exercise-bar {
  display: flex;
  gap: 8px;
  flex: 1;
  max-width: 450px;
}

.add-exercise-bar select {
  flex: 1;
}

.empty-builder {
  text-align: center;
  padding: 24px;
  color: var(--text-dim);
  border: 1px dashed var(--border);
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
}

.exercise-row-card {
  background: #0d1322;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.row-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.row-num {
  color: var(--primary);
  font-weight: 700;
  margin-right: 6px;
}

.row-title {
  font-weight: 600;
  color: var(--text-main);
  flex: 1;
}

.row-inputs {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 10px;
  align-items: flex-end;
}

.param-input {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.param-input label {
  font-size: 0.75rem;
  color: var(--text-dim);
}

.param-input input {
  padding: 6px 8px;
  font-size: 0.85rem;
}

.param-input-checkbox {
  display: flex;
  align-items: center;
  padding-bottom: 6px;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.error-box {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: var(--danger);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border);
  background: #0d1322;
}

.btn-icon-delete:hover {
  background: rgba(239, 68, 68, 0.2);
}
</style>
