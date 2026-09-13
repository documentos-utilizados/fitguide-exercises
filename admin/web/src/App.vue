<template>
  <div class="app-wrapper">
    <header class="navbar">
      <div class="navbar-container">
        <div class="brand">
          <span class="brand-icon">⚡</span>
          <div class="brand-text">
            <h1>FitGuide Admin</h1>
            <span class="brand-subtitle">Catálogo & Treinos</span>
          </div>
        </div>

        <nav class="nav-tabs">
          <button 
            :class="['nav-item', { active: currentTab === 'exercises' }]" 
            @click="currentTab = 'exercises'"
          >
            📋 Exercícios ({{ exercises.length }})
          </button>
          <button 
            :class="['nav-item', { active: currentTab === 'workouts' }]" 
            @click="currentTab = 'workouts'"
          >
            🏋️ Treinos & Fichas ({{ totalWorkouts }})
          </button>
          <button 
            :class="['nav-item', { active: currentTab === 'validator' }]" 
            @click="currentTab = 'validator'"
          >
            🌐 Validador i18n
          </button>
          <button 
            :class="['nav-item', { active: currentTab === 'build' }]" 
            @click="currentTab = 'build'"
          >
            🚀 Compilação
          </button>
        </nav>

        <div class="nav-actions">
          <button class="btn btn-secondary btn-sm" @click="refreshAll" :disabled="loading">
            🔄 Recarregar
          </button>
        </div>
      </div>
    </header>

    <main class="main-content">
      <div v-if="toastMessage" :class="['toast-notification', `toast-${toastType}`]">
        {{ toastMessage }}
      </div>

      <div v-if="loading && !exercises.length" class="loading-state">
        <div class="spinner"></div>
        <p>Carregando base de dados...</p>
      </div>

      <div v-else>
        <section v-show="currentTab === 'exercises'">
          <ExerciseList 
            :exercises="exercises" 
            :metadata="metadata" 
            @create="openCreateModal"
            @edit="openEditModal"
            @delete="handleDeleteExercise"
          />
        </section>

        <section v-show="currentTab === 'workouts'">
          <WorkoutBuilder 
            :workouts="workouts" 
            :exercisesList="exercises" 
            @workout-updated="loadWorkouts"
          />
        </section>

        <section v-show="currentTab === 'validator'">
          <TranslationValidator 
            :exercises="exercises"
            :metadata="metadata"
            @metadata-updated="loadMetadata"
          />
        </section>

        <section v-show="currentTab === 'build'">
          <BuildController @build-completed="onBuildCompleted" />
        </section>
      </div>
    </main>

    <ExerciseModal 
      v-if="showExerciseModal" 
      :exercise="selectedExercise" 
      :metadata="metadata" 
      @close="showExerciseModal = false" 
      @saved="onExerciseSaved"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchExercises, fetchMetadata, fetchWorkouts, deleteExercise } from './api'
import ExerciseList from './components/ExerciseList.vue'
import ExerciseModal from './components/ExerciseModal.vue'
import WorkoutBuilder from './components/WorkoutBuilder.vue'
import TranslationValidator from './components/TranslationValidator.vue'
import BuildController from './components/BuildController.vue'

const currentTab = ref('exercises')
const exercises = ref([])
const metadata = ref({})
const workouts = ref({ seeds: [], templates: [] })
const loading = ref(false)

const showExerciseModal = ref(false)
const selectedExercise = ref(null)

const toastMessage = ref('')
const toastType = ref('success')
let toastTimer = null

const totalWorkouts = computed(() => {
  const seedsCount = workouts.value.seeds?.length || 0
  const templatesCount = workouts.value.templates?.length || 0
  return seedsCount + templatesCount
})

function showToast(msg, type = 'success') {
  toastMessage.value = msg
  toastType.value = type
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 4000)
}

async function loadExercises() {
  try {
    exercises.value = await fetchExercises('pt')
  } catch (err) {
    showToast(`Erro ao carregar exercícios: ${err.message}`, 'error')
  }
}

async function loadMetadata() {
  try {
    metadata.value = await fetchMetadata()
  } catch (err) {
    showToast(`Erro ao carregar metadados: ${err.message}`, 'error')
  }
}

async function loadWorkouts() {
  try {
    workouts.value = await fetchWorkouts()
  } catch (err) {
    showToast(`Erro ao carregar treinos: ${err.message}`, 'error')
  }
}

async function refreshAll() {
  loading.value = true
  try {
    await Promise.all([loadExercises(), loadMetadata(), loadWorkouts()])
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  selectedExercise.value = null
  showExerciseModal.value = true
}

function openEditModal(exercise) {
  selectedExercise.value = exercise
  showExerciseModal.value = true
}

async function handleDeleteExercise(exercise) {
  if (!confirm(`Tem certeza que deseja excluir o exercício "${exercise.name}" (${exercise.id})?`)) {
    return
  }

  try {
    await deleteExercise(exercise.id)
    showToast(`Exercício "${exercise.id}" excluído com sucesso.`)
    await loadExercises()
  } catch (err) {
    showToast(`Erro ao excluir: ${err.message}`, 'error')
  }
}

async function onExerciseSaved(savedExercise) {
  showExerciseModal.value = false
  showToast(`Exercício "${savedExercise.name || savedExercise.id}" salvo com sucesso!`)
  await loadExercises()
}

async function onBuildCompleted() {
  showToast('Compilação finalizada! Dados sincronizados.')
  await refreshAll()
}

onMounted(() => {
  refreshAll()
})
</script>

<style scoped>
.app-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  font-size: 1.5rem;
}

.brand h1 {
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.2;
}

.brand-subtitle {
  font-size: 0.72rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.nav-tabs {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 100%;
}

.nav-item {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 0.9rem;
  font-weight: 600;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-item:hover {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.05);
}

.nav-item.active {
  color: var(--primary);
  background: rgba(16, 185, 129, 0.12);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.main-content {
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
  flex: 1;
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
  animation: slideIn 0.3s ease-out;
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

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  gap: 16px;
  color: var(--text-muted);
}

@keyframes slideIn {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .navbar-container {
    height: auto;
    padding: 12px 16px;
    flex-wrap: wrap;
  }

  .nav-tabs {
    order: 3;
    width: 100%;
    overflow-x: auto;
    padding-bottom: 4px;
  }
}
</style>
