const API_BASE = '/api'

export async function fetchExercises() {
  const res = await fetch(`${API_BASE}/exercises`)
  if (!res.ok) throw new Error('Falha ao buscar exercícios')
  return res.json()
}

export async function fetchExerciseById(id) {
  const res = await fetch(`${API_BASE}/exercises?id=${encodeURIComponent(id)}`)
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

