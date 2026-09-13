<template>
  <div class="build-card">
    <div class="build-header">
      <div>
        <h2>Compilação de Datasets & Metadados</h2>
        <p>Gera os arquivos compilados em <code>dist/</code> e sincroniza metadados e treinos hidratados.</p>
      </div>
      <button 
        class="btn btn-primary" 
        :disabled="loading" 
        @click="runBuild"
      >
        <span v-if="loading" class="spinner"></span>
        <span v-else>⚡ Compilar Agora</span>
      </button>
    </div>

    <div v-if="statusMessage" :class="['status-box', isSuccess ? 'status-success' : 'status-error']">
      {{ statusMessage }}
    </div>

    <div v-if="logs" class="logs-container">
      <div class="logs-header">
        <span>Console de Saída</span>
        <button class="btn-clear" @click="logs = ''">Limpar</button>
      </div>
      <pre><code>{{ logs }}</code></pre>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { triggerBuild } from '../api'

const emit = defineEmits(['build-completed'])

const loading = ref(false)
const logs = ref('')
const statusMessage = ref('')
const isSuccess = ref(true)

async function runBuild() {
  loading.value = true
  statusMessage.value = 'Executando pipeline de compilação...'
  isSuccess.value = true
  logs.value = ''

  try {
    const res = await triggerBuild()
    logs.value = res.output || 'Compilação concluída com sucesso!'
    statusMessage.value = 'Datasets, metadados e treinos compilados com sucesso!'
    isSuccess.value = true
    emit('build-completed')
  } catch (err) {
    statusMessage.value = `Erro: ${err.message}`
    isSuccess.value = false
    if (err.message) {
      logs.value = err.message
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.build-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 24px;
}

.build-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.build-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
  margin-bottom: 4px;
}

.build-header p {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.build-header code {
  color: var(--primary);
  background: rgba(16, 185, 129, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.status-box {
  margin-top: 16px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  font-size: 0.9rem;
  font-weight: 500;
}

.status-success {
  background: rgba(16, 185, 129, 0.15);
  color: var(--primary);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.status-error {
  background: rgba(239, 68, 68, 0.15);
  color: var(--danger);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.logs-container {
  margin-top: 16px;
  background: #060911;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #0f172a;
  font-size: 0.75rem;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.btn-clear {
  color: var(--text-muted);
  font-size: 0.75rem;
}
.btn-clear:hover {
  color: var(--text-main);
}

pre {
  padding: 12px;
  font-family: 'JetBrains Mono', Consolas, Monaco, monospace;
  font-size: 0.8rem;
  color: #38bdf8;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 250px;
  overflow-y: auto;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-top-color: #041f18;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
