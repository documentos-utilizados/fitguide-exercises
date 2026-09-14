const API_BASE = '/api'

export async function fetchExercises(lang = 'pt') {
  const res = await fetch(`${API_BASE}/exercises?lang=${encodeURIComponent(lang)}`)
  if (!res.ok) throw new Error('Falha ao buscar exercícios')
  return res.json()
}

export async function fetchExerciseById(id, lang = 'pt') {
  const res = await fetch(`${API_BASE}/exercises?id=${encodeURIComponent(id)}&lang=${encodeURIComponent(lang)}`)
  if (!res.ok) throw new Error(`Exercício ${id} não encontrado`)
  return res.json()
}

export async function saveExercise(exercise) {
  const res = await fetch(`${API_BASE}/exercises`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(exercise)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao salvar exercício')
  }
  return res.json()
}

export async function deleteExercise(id) {
  if (!id || id === 'undefined' || id === 'null') {
    throw new Error('ID do exercício é obrigatório para exclusão')
  }
  const res = await fetch(`${API_BASE}/exercises?id=${encodeURIComponent(id)}`, {
    method: 'DELETE'
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao excluir exercício')
  }
  return res.json()
}

export async function fetchMetadata() {
  const res = await fetch(`${API_BASE}/metadata`)
  if (!res.ok) throw new Error('Falha ao carregar metadados')
  return res.json()
}

export async function fetchWorkouts() {
  const res = await fetch(`${API_BASE}/workouts`)
  if (!res.ok) throw new Error('Falha ao buscar treinos')
  return res.json()
}

export async function saveWorkout(workout) {
  const res = await fetch(`${API_BASE}/workouts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(workout)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao salvar treino')
  }
  return res.json()
}

export async function deleteWorkout(type, id) {
  const res = await fetch(`${API_BASE}/workouts?type=${encodeURIComponent(type)}&id=${encodeURIComponent(id)}`, {
    method: 'DELETE'
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao excluir treino')
  }
  return res.json()
}

export async function uploadImage(file, relativePath) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('relative_path', relativePath)

  const res = await fetch(`${API_BASE}/upload-image`, {
    method: 'POST',
    body: formData
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao fazer upload da imagem')
  }

  return res.json()
}

export async function fetchTranslationAudit() {
  const res = await fetch(`${API_BASE}/translations/audit`)
  if (!res.ok) throw new Error('Falha ao obter auditoria de traduções')
  return res.json()
}

export async function fetchTranslationCompare(id) {
  const res = await fetch(`${API_BASE}/translations/compare?id=${encodeURIComponent(id)}`)
  if (!res.ok) throw new Error(`Falha ao comparar traduções do exercício ${id}`)
  return res.json()
}

export async function triggerBuild() {
  const res = await fetch(`${API_BASE}/build`, {
    method: 'POST'
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok || !data.success) {
    throw new Error(data.error || 'Falha ao executar compilação')
  }
  return data
}

export async function saveCategory(category) {
  const res = await fetch(`${API_BASE}/metadata/category`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      key: category.key,
      label_pt: category.labelPt || category.label_pt,
      label_en: category.labelEn || category.label_en
    })
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao salvar categoria')
  }
  return res.json()
}

export async function deleteCategory(key) {
  const res = await fetch(`${API_BASE}/metadata/category?key=${encodeURIComponent(key)}`, {
    method: 'DELETE'
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao excluir categoria')
  }
  return res.json()
}

export async function translateText(text, source = 'pt', target = 'en') {
  const res = await fetch(`${API_BASE}/translate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text, source, target })
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao traduzir texto')
  }
  return res.json()
}

export async function extractPrimaryMuscleCategories() {
  const res = await fetch(`${API_BASE}/metadata/categories/extract`, {
    method: 'POST'
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao extrair categorias')
  }
  return res.json()
}

export async function translateCategories(categoryType = 'primary_muscle', source = 'pt', target = 'en') {
  const res = await fetch(`${API_BASE}/metadata/categories/translate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ category_type: categoryType, source, target })
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Falha ao traduzir categorias')
  }
  return res.json()
}

