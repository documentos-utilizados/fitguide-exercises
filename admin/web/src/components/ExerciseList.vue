<template>
  <div class="exercise-section">
    <div class="toolbar">
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Buscar exercício por nome ou ID..."
        />
        <span v-if="searchQuery" class="clear-search" @click="searchQuery = ''">&times;</span>
      </div>

      <div class="filters">
        <select v-model="selectedMuscle">
          <option value="">Todos os Músculos</option>
          <option v-for="(label, key) in metadata.muscles?.pt || {}" :key="key" :value="label">
            {{ label }}
          </option>
        </select>

        <select v-model="selectedCategory">
          <option value="">Todas as Categorias</option>
          <option value="strength">Força</option>
          <option value="cardio">Cardio</option>
          <option value="stretching">Alongamento</option>
          <option value="plyometrics">Pliometria</option>
          <option value="powerlifting">Powerlifting</option>
          <option value="olympic weightlifting">LPO</option>
          <option value="strongman">Strongman</option>
        </select>

        <select v-model="selectedLevel">
          <option value="">Todos os Níveis</option>
          <option value="iniciante">Iniciante</option>
          <option value="intermediário">Intermediário</option>
          <option value="avançado">Avançado</option>
        </select>

        <select v-model="selectedEquipment">
          <option value="">Todos os Equipamentos</option>
          <option v-for="(label, key) in metadata.equipments?.pt || {}" :key="key" :value="label">
            {{ label }}
          </option>
        </select>

        <button class="btn btn-primary" @click="$emit('create')">
          + Novo Exercício
        </button>
      </div>
    </div>

    <div class="list-meta">
      <span>Exibindo <strong>{{ filteredExercises.length }}</strong> de <strong>{{ exercises.length }}</strong> exercícios em <code>exercises_pt/</code></span>
    </div>

    <div v-if="paginatedExercises.length === 0" class="empty-state">
      <p>Nenhum exercício encontrado com os filtros selecionados.</p>
    </div>

    <div v-else class="exercise-grid">
      <div v-for="ex in paginatedExercises" :key="ex.id" class="exercise-card">
        <div class="card-header">
          <div class="title-group">
            <h4>{{ ex.name }}</h4>
            <span class="ex-id">{{ ex.id }}</span>
          </div>
          <div class="card-actions">
            <button class="btn-icon" title="Editar" @click="$emit('edit', ex)">✏️</button>
            <button class="btn-icon btn-icon-delete" title="Excluir" @click="confirmDelete(ex)">🗑️</button>
          </div>
        </div>

        <div class="tags-row">
          <span class="badge badge-primary">{{ ex.primaryMuscles?.[0] || 'N/A' }}</span>
          <span v-if="ex.equipment" class="badge badge-blue">{{ ex.equipment }}</span>
          <span class="badge badge-gray">{{ ex.level }}</span>
          <span v-if="ex.force" class="badge badge-gray">{{ ex.force }}</span>
        </div>

        <div class="instructions-preview">
          <p>{{ ex.instructions?.[0] || 'Sem instruções cadastradas.' }}</p>
        </div>

        <div class="card-footer">
          <span class="img-count">📷 {{ ex.images?.length || 0 }} fotos</span>
          <span class="steps-count">📝 {{ ex.instructions?.length || 0 }} passos</span>
        </div>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button 
        class="btn btn-secondary btn-page" 
        :disabled="currentPage === 1" 
        @click="currentPage--"
      >
        &laquo; Anterior
      </button>

      <span class="page-info">Página {{ currentPage }} de {{ totalPages }}</span>

      <button 
        class="btn btn-secondary btn-page" 
        :disabled="currentPage === totalPages" 
        @click="currentPage++"
      >
        Próxima &raquo;
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  exercises: {
    type: Array,
    default: () => []
  },
  metadata: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['create', 'edit', 'delete'])

const searchQuery = ref('')
const selectedMuscle = ref('')
const selectedCategory = ref('')
const selectedLevel = ref('')
const selectedEquipment = ref('')
const currentPage = ref(1)
const pageSize = 24

watch([searchQuery, selectedMuscle, selectedCategory, selectedLevel, selectedEquipment], () => {
  currentPage.value = 1
})

const filteredExercises = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  return props.exercises.filter(ex => {
    if (q) {
      const matchName = ex.name?.toLowerCase().includes(q)
      const matchId = ex.id?.toLowerCase().includes(q)
      if (!matchName && !matchId) return false
    }

    if (selectedMuscle.value) {
      const prim = (ex.primaryMuscles || []).map(m => m.toLowerCase())
      if (!prim.includes(selectedMuscle.value.toLowerCase())) return false
    }

    if (selectedCategory.value && ex.category?.toLowerCase() !== selectedCategory.value.toLowerCase()) {
      return false
    }

    if (selectedLevel.value && ex.level?.toLowerCase() !== selectedLevel.value.toLowerCase()) {
      return false
    }

    if (selectedEquipment.value && ex.equipment?.toLowerCase() !== selectedEquipment.value.toLowerCase()) {
      return false
    }

    return true
  })
})

const totalPages = computed(() => Math.ceil(filteredExercises.value.length / pageSize) || 1)

const paginatedExercises = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredExercises.value.slice(start, start + pageSize)
})

function confirmDelete(ex) {
  if (confirm(`Tem certeza que deseja excluir o exercício "${ex.name}" (${ex.id})?`)) {
    emit('delete', ex.id)
  }
}
</script>

<style scoped>
.exercise-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 16px;
}

@media (min-width: 900px) {
  .toolbar {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

.search-box {
  position: relative;
  flex: 1;
}

.search-box input {
  width: 100%;
  padding-right: 30px;
}

.clear-search {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  color: var(--text-muted);
  font-size: 1.2rem;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.list-meta {
  font-size: 0.85rem;
  color: var(--text-muted);
}
.list-meta code {
  color: var(--primary);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: var(--radius-lg);
}

.exercise-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.exercise-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 12px;
  transition: all 0.2s ease;
}

.exercise-card:hover {
  border-color: var(--border-focus);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.title-group h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-main);
  line-height: 1.3;
}

.ex-id {
  font-size: 0.75rem;
  color: var(--text-dim);
  font-family: monospace;
}

.card-actions {
  display: flex;
  gap: 4px;
}

.btn-icon {
  font-size: 0.9rem;
  padding: 4px;
  border-radius: 4px;
  opacity: 0.8;
}
.btn-icon:hover {
  opacity: 1;
  background: var(--bg-card-hover);
}
.btn-icon-delete:hover {
  background: rgba(239, 68, 68, 0.2);
}

.tags-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.instructions-preview {
  font-size: 0.82rem;
  color: var(--text-muted);
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--text-dim);
  border-top: 1px solid var(--border);
  padding-top: 8px;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 16px;
  padding: 16px 0;
}

.page-info {
  font-size: 0.85rem;
  color: var(--text-muted);
}
</style>
