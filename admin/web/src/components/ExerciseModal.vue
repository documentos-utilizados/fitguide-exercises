<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-card">
      <div class="modal-header">
        <h3>{{ isEdit ? 'Editar Exercício' : 'Novo Exercício' }}</h3>
        <button class="btn-close" @click="$emit('close')">&times;</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-grid">
          <div class="form-group">
            <label>ID do Exercício (Canônico) *</label>
            <input 
              type="text" 
              v-model="form.id" 
              placeholder="Ex: Supino_Inclinado_Halteres" 
              :disabled="isEdit"
              required
            />
          </div>

          <div class="form-group">
            <label>Nome do Exercício (PT-BR) *</label>
            <input 
              type="text" 
              v-model="form.name" 
              placeholder="Ex: Supino Inclinado com Halteres" 
              required
            />
          </div>
        </div>

        <div class="form-grid-3">
          <div class="form-group">
            <label>Categoria *</label>
            <div v-if="!isCreatingNewCategory">
              <select v-model="form.category" @change="onCategoryChange" required>
                <option v-for="(label, key) in metadata.categories?.pt || {}" :key="key" :value="key">
                  {{ label }} ({{ key }})
                </option>
                <option value="__new__">+ Criar Nova Categoria...</option>
              </select>
            </div>
            <div v-else class="input-with-action">
              <input 
                type="text" 
                v-model="customCategoryInput" 
                placeholder="Ex: pilates, funcional, crossfit..." 
                required
              />
              <button type="button" class="btn btn-secondary btn-sm" @click="cancelNewCategory">Voltar</button>
            </div>
          </div>

          <div class="form-group">
            <label>Nível de Dificuldade *</label>
            <select v-model="form.level" required>
              <option value="iniciante">Iniciante</option>
              <option value="intermediário">Intermediário</option>
              <option value="avançado">Avançado</option>
            </select>
          </div>

          <div class="form-group">
            <label>Mecânica</label>
            <select v-model="form.mechanic">
              <option :value="null">Nenhuma / Padrão</option>
              <option value="composto">Composto</option>
              <option value="isolado">Isolado</option>
            </select>
          </div>
        </div>

        <div class="form-grid-3">
          <div class="form-group">
            <label>Equipamento</label>
            <select v-model="form.equipment">
              <option :value="null">Nenhum / Peso Corporal</option>
              <option v-for="item in availableEquipments" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Tipo de Força</label>
            <select v-model="form.force">
              <option :value="null">Nenhum / Padrão</option>
              <option value="empurrar">Empurrar (push)</option>
              <option value="tração">Tração (pull)</option>
              <option value="estático">Estático (static)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Músculo Primário *</label>
            <select v-model="primaryMuscleModel" required>
              <option value="">Selecione o músculo</option>
              <option v-for="item in availableMuscles" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </div>
        </div>

        <div class="form-section">
          <div class="section-title">
            <label>Instruções de Execução (Passo a Passo) *</label>
            <button type="button" class="btn btn-secondary btn-sm" @click="addInstruction">
              + Adicionar Passo
            </button>
          </div>

          <div v-for="(step, idx) in form.instructions" :key="idx" class="instruction-row">
            <span class="step-num">{{ idx + 1 }}</span>
            <textarea 
              v-model="form.instructions[idx]" 
              rows="2" 
              placeholder="Descreva a execução deste passo..."
              required
            ></textarea>
            <button 
              type="button" 
              class="btn-icon-danger" 
              @click="removeInstruction(idx)"
              :disabled="form.instructions.length <= 1"
            >
              &times;
            </button>
          </div>
        </div>

        <div class="form-section">
          <div class="section-title">
            <label>Imagens do Exercício</label>
            <div class="actions-group">
              <button type="button" class="btn btn-secondary btn-sm" @click="triggerFileInput">
                📁 Escolher Foto do Computador
              </button>
              <button type="button" class="btn btn-outline btn-sm" @click="addImageManual">
                + Caminho Manual
              </button>
              <input 
                type="file" 
                ref="fileInputRef" 
                accept="image/*" 
                style="display: none" 
                @change="handleFileSelected"
              />
            </div>
          </div>

          <div v-if="pendingUpload" class="upload-dialog">
            <div class="upload-preview-box">
              <img :src="pendingUpload.previewUrl" class="upload-thumb" />
              <div class="upload-details">
                <span class="upload-orig-name">Arquivo: {{ pendingUpload.file.name }}</span>
                <label>Nome / Caminho dentro do Projeto:</label>
                <div class="input-with-action">
                  <input 
                    type="text" 
                    v-model="pendingUpload.targetPath" 
                    placeholder="Ex: natacao/0.jpg"
                  />
                  <button 
                    type="button" 
                    class="btn btn-primary btn-sm" 
                    :disabled="isUploading"
                    @click="confirmUpload"
                  >
                    {{ isUploading ? 'Salvando...' : 'Salvar no Projeto' }}
                  </button>
                  <button 
                    type="button" 
                    class="btn btn-secondary btn-sm" 
                    @click="cancelUpload"
                  >
                    Cancelar
                  </button>
                </div>
              </div>
            </div>
            <p v-if="uploadError" class="upload-error">{{ uploadError }}</p>
          </div>

          <div v-if="form.images.length === 0 && !pendingUpload" class="empty-images">
            <span>Nenhuma imagem adicionada. Escolha um arquivo do seu computador ou adicione manualmente.</span>
          </div>

          <div v-for="(img, idx) in form.images" :key="idx" class="image-item-row">
            <div class="img-preview-wrap">
              <img 
                :src="'/images/' + img" 
                @error="onImgLoadError($event)" 
                class="exercise-img-thumb"
              />
            </div>
            <div class="image-input-wrap">
              <input 
                type="text" 
                v-model="form.images[idx]" 
                placeholder="Ex: Barbell_Bench_Press/0.jpg"
              />
            </div>
            <button type="button" class="btn-icon-danger" @click="removeImage(idx)">&times;</button>
          </div>
        </div>

        <div v-if="errorMessage" class="error-box">
          {{ errorMessage }}
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
          <button type="submit" class="btn btn-primary" :disabled="saving || isUploading">
            {{ saving ? 'Salvando...' : 'Salvar Exercício' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { saveExercise, uploadImage } from '../api'

const props = defineProps({
  exercise: {
    type: Object,
    default: null
  },
  metadata: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['close', 'saved'])

const isEdit = computed(() => !!props.exercise?.id)
const saving = ref(false)
const errorMessage = ref('')
const fileInputRef = ref(null)

const isCreatingNewCategory = ref(false)
const customCategoryInput = ref('')

const pendingUpload = ref(null)
const isUploading = ref(false)
const uploadError = ref('')

const form = reactive({
  id: props.exercise?.id || '',
  name: props.exercise?.name || '',
  force: props.exercise?.force || null,
  level: props.exercise?.level || 'iniciante',
  mechanic: props.exercise?.mechanic || null,
  equipment: props.exercise?.equipment || null,
  primaryMuscles: props.exercise?.primaryMuscles ? [...props.exercise.primaryMuscles] : ['peitoral'],
  secondaryMuscles: props.exercise?.secondaryMuscles ? [...props.exercise.secondaryMuscles] : [],
  instructions: props.exercise?.instructions?.length ? [...props.exercise.instructions] : [''],
  category: props.exercise?.category || 'strength',
  images: props.exercise?.images?.length ? [...props.exercise.images] : []
})

const availableMuscles = computed(() => {
  const list = Object.values(props.metadata?.muscles?.pt || {})
  return Array.from(new Set(list)).sort((a, b) => a.localeCompare(b))
})

const availableEquipments = computed(() => {
  const list = Object.values(props.metadata?.equipments?.pt || {})
  return Array.from(new Set(list)).sort((a, b) => a.localeCompare(b))
})

const primaryMuscleModel = computed({
  get: () => form.primaryMuscles[0] || '',
  set: (val) => {
    form.primaryMuscles = val ? [val] : []
  }
})

function onCategoryChange(e) {
  if (e.target.value === '__new__') {
    isCreatingNewCategory.value = true
    customCategoryInput.value = ''
  }
}

function cancelNewCategory() {
  isCreatingNewCategory.value = false
  form.category = 'strength'
}

function addInstruction() {
  form.instructions.push('')
}

function removeInstruction(index) {
  if (form.instructions.length > 1) {
    form.instructions.splice(index, 1)
  }
}

function addImageManual() {
  const folderName = form.id ? form.id.trim().replace(/\s+/g, '_') : 'novo_exercicio'
  const nextIndex = form.images.length
  const defaultPath = `${folderName}/${nextIndex}.jpg`
  form.images.push(defaultPath)
}

function removeImage(index) {
  form.images.splice(index, 1)
}

function triggerFileInput() {
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
    fileInputRef.value.click()
  }
}

function handleFileSelected(event) {
  const file = event.target.files?.[0]
  if (!file) return

  uploadError.value = ''
  const folderName = form.id ? form.id.trim().replace(/\s+/g, '_') : 'novo_exercicio'
  const nextIndex = form.images.length
  const suggestedPath = `${folderName}/${nextIndex}.jpg`

  pendingUpload.value = {
    file,
    previewUrl: URL.createObjectURL(file),
    targetPath: suggestedPath
  }
}

async function confirmUpload() {
  if (!pendingUpload.value) return
  let targetPath = pendingUpload.value.targetPath.trim()
  if (!targetPath) {
    uploadError.value = 'Informe o caminho/nome de destino.'
    return
  }

  if (!targetPath.toLowerCase().endsWith('.jpg')) {
    targetPath = targetPath.replace(/\.[a-zA-Z0-9]+$/, '') + '.jpg'
    if (!targetPath.endsWith('.jpg')) {
      targetPath += '.jpg'
    }
  }

  isUploading.value = true
  uploadError.value = ''

  try {
    const res = await uploadImage(pendingUpload.value.file, targetPath)
    form.images.push(res.path)
    cancelUpload()
  } catch (err) {
    uploadError.value = err.message || 'Erro ao enviar imagem'
  } finally {
    isUploading.value = false
  }
}

function cancelUpload() {
  if (pendingUpload.value?.previewUrl) {
    URL.revokeObjectURL(pendingUpload.value.previewUrl)
  }
  pendingUpload.value = null
  uploadError.value = ''
}

function onImgLoadError(event) {
  event.target.style.opacity = '0.3'
}

async function handleSubmit() {
  saving.value = true
  errorMessage.value = ''

  try {
    let finalCategory = form.category
    if (isCreatingNewCategory.value) {
      const sanitized = customCategoryInput.value
        .trim()
        .toLowerCase()
        .replace(/[\s-]+/g, '_')
        .replace(/[^a-z0-9_]/g, '')
      if (!sanitized) {
        throw new Error('Informe um nome válido para a nova categoria')
      }
      finalCategory = sanitized
    }

    const payload = {
      ...form,
      id: form.id.trim(),
      name: form.name.trim(),
      category: finalCategory,
      instructions: form.instructions.filter(i => i.trim() !== ''),
      images: form.images.filter(i => i.trim() !== '')
    }

    await saveExercise(payload)
    emit('saved', payload)
    emit('close')
  } catch (err) {
    errorMessage.value = err.message || 'Erro ao salvar exercício'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
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
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
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
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px 8px;
}

.btn-close:hover {
  color: var(--text-main);
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--text-muted);
}

.form-group input,
.form-group select {
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  padding: 8px 12px;
  font-size: 0.9rem;
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-section {
  border-top: 1px solid var(--border);
  padding-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title label {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-main);
}

.actions-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-outline {
  background: transparent;
  border: 1px dashed var(--border);
  color: var(--text-muted);
}

.btn-outline:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.upload-dialog {
  background: #0f172a;
  border: 1px solid var(--primary);
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.upload-preview-box {
  display: flex;
  align-items: center;
  gap: 16px;
}

.upload-thumb {
  width: 72px;
  height: 72px;
  border-radius: var(--radius-sm);
  object-fit: cover;
  border: 1px solid var(--border);
  background: #000;
}

.upload-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.upload-orig-name {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.upload-details label {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-main);
}

.input-with-action {
  display: flex;
  align-items: center;
  gap: 8px;
}

.input-with-action input {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  padding: 6px 10px;
  font-size: 0.85rem;
}

.upload-error {
  font-size: 0.8rem;
  color: var(--danger);
}

.empty-images {
  font-size: 0.85rem;
  color: var(--text-muted);
  padding: 12px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: var(--radius-sm);
  text-align: center;
}

.image-item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.03);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.img-preview-wrap {
  width: 44px;
  height: 44px;
  border-radius: 4px;
  overflow: hidden;
  background: #000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.exercise-img-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-input-wrap {
  flex: 1;
}

.image-input-wrap input {
  width: 100%;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  padding: 6px 10px;
  font-size: 0.85rem;
}

.instruction-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.step-num {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-muted);
  width: 24px;
  text-align: center;
  padding-top: 8px;
}

.instruction-row textarea {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  padding: 8px 12px;
  font-size: 0.9rem;
  resize: vertical;
}

.btn-icon-danger {
  background: transparent;
  border: none;
  color: var(--danger);
  font-size: 1.25rem;
  cursor: pointer;
  padding: 4px 8px;
}

.btn-icon-danger:hover {
  opacity: 0.8;
}

.error-box {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: var(--danger);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

@media (max-width: 640px) {
  .form-grid,
  .form-grid-3 {
    grid-template-columns: 1fr;
  }
  .input-with-action {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
