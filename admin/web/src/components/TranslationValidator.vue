<template>
  <div class="translation-validator">
    <div class="validator-nav">
      <div class="tabs-group">
        <button 
          :class="['tab-btn', { active: activeSubTab === 'compare' }]"
          @click="activeSubTab = 'compare'"
        >
          🔍 Exercício Lado a Lado
        </button>
        <button 
          :class="['tab-btn', { active: activeSubTab === 'categories' }]"
          @click="activeSubTab = 'categories'"
        >
          🏷️ Comparador de Categorias ({{ categoryComparisonList.length }})
        </button>
        <button 
          :class="['tab-btn', { active: activeSubTab === 'audit' }]"
          @click="activeSubTab = 'audit'"
        >
          📊 Auditoria Geral ({{ auditReport?.total_discrepancies || 0 }})
        </button>
        <button 
          :class="['tab-btn', { active: activeSubTab === 'payload' }]"
          @click="activeSubTab = 'payload'"
        >
          📱 Simulador de Payload Mobile
        </button>
      </div>

      <button class="btn btn-secondary btn-sm" @click="loadData" :disabled="loading">
        🔄 Atualizar Auditoria
      </button>
    </div>

    <div v-if="loading && !auditReport" class="loading-box">
      <div class="spinner"></div>
      <p>Analisando consistência e traduções entre idiomas...</p>
    </div>

    <div v-else>
      <section v-show="activeSubTab === 'compare'" class="sub-tab-content">
        <div class="selector-card">
          <div class="selector-controls">
            <div class="form-group flex-1">
              <label>Selecionar Exercício para Análise:</label>
              <select v-model="selectedExerciseId" @change="onExerciseSelected">
                <option v-for="ex in filteredExerciseOptions" :key="ex.id" :value="ex.id">
                  {{ ex.name }} ({{ ex.id }})
                </option>
              </select>
            </div>

            <div class="form-group">
              <label>Filtro Rápido:</label>
              <select v-model="filterDiscrepancy">
                <option value="all">Todos os Exercícios ({{ exercises.length }})</option>
                <option value="steps">Com Discrepância de Passos ({{ auditReport?.step_mismatches?.length || 0 }})</option>
                <option value="images">Com Discrepância de Fotos ({{ auditReport?.image_mismatches?.length || 0 }})</option>
              </select>
            </div>
          </div>
        </div>

        <div v-if="loadingCompare" class="loading-box">
          <div class="spinner"></div>
          <p>Carregando dados comparativos do exercício...</p>
        </div>

        <div v-else-if="compareData" class="compare-container">
          <div class="compare-header-summary">
            <div class="canonical-badge">
              <span class="badge-label">ID Canônico:</span>
              <code>{{ compareData.id }}</code>
              <span class="badge badge-success">✅ Identificador Universal</span>
            </div>

            <div class="status-tags">
              <span v-if="stepMismatchWarning" class="badge badge-warning">
                ⚠️ {{ stepMismatchWarning }}
              </span>
              <span v-else class="badge badge-success">
                ✅ Paridade de Passos ({{ ptExercise?.instructions?.length || 0 }})
              </span>

              <span v-if="imageMismatchWarning" class="badge badge-warning">
                ⚠️ {{ imageMismatchWarning }}
              </span>
              <span v-else class="badge badge-success">
                ✅ Paridade de Fotos ({{ ptExercise?.images?.length || 0 }})
              </span>
            </div>
          </div>

          <div class="compare-grid">
            <div class="lang-column">
              <div class="lang-column-header header-pt">
                <span class="flag">🇧🇷</span>
                <h4>Português (PT-BR)</h4>
                <span class="source-tag">Fonte da Verdade</span>
              </div>

              <div v-if="ptExercise" class="card-details">
                <div class="field-item">
                  <span class="field-label">Nome do Exercício:</span>
                  <div class="field-value highlight">{{ ptExercise.name }}</div>
                </div>

                <div class="field-row-2">
                  <div class="field-item">
                    <span class="field-label">Categoria:</span>
                    <div class="field-value">
                      <span class="badge badge-primary">{{ ptExercise.category }}</span>
                    </div>
                  </div>
                  <div class="field-item">
                    <span class="field-label">Nível:</span>
                    <div class="field-value">
                      <span class="badge badge-gray">{{ ptExercise.level }}</span>
                    </div>
                  </div>
                </div>

                <div class="field-row-2">
                  <div class="field-item">
                    <span class="field-label">Mecânica:</span>
                    <div class="field-value">{{ ptExercise.mechanic || 'N/A' }}</div>
                  </div>
                  <div class="field-item">
                    <span class="field-label">Tipo de Força:</span>
                    <div class="field-value">{{ ptExercise.force || 'N/A' }}</div>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Equipamento:</span>
                  <div class="field-value">{{ ptExercise.equipment || 'Nenhum / Peso Corporal' }}</div>
                </div>

                <div class="field-item">
                  <span class="field-label">Músculos Primários:</span>
                  <div class="tags-group">
                    <span v-for="m in ptExercise.primaryMuscles" :key="m" class="badge badge-primary">{{ m }}</span>
                    <span v-if="!ptExercise.primaryMuscles?.length" class="text-muted">Nenhum</span>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Músculos Secundários:</span>
                  <div class="tags-group">
                    <span v-for="m in ptExercise.secondaryMuscles" :key="m" class="badge badge-gray">{{ m }}</span>
                    <span v-if="!ptExercise.secondaryMuscles?.length" class="text-muted">Nenhum</span>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Instruções de Execução ({{ ptExercise.instructions?.length || 0 }} passos):</span>
                  <div class="steps-list">
                    <div v-for="(step, sIdx) in ptExercise.instructions" :key="sIdx" class="step-card">
                      <span class="step-badge">{{ sIdx + 1 }}</span>
                      <p>{{ step }}</p>
                    </div>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Imagens ({{ ptExercise.images?.length || 0 }}):</span>
                  <div class="images-preview-grid">
                    <div v-for="(img, iIdx) in ptExercise.images" :key="iIdx" class="img-card">
                      <img :src="'/images/' + img" class="thumb" @error="onImgErr" />
                      <span class="img-path">{{ img }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="empty-lang">Exercício não encontrado em PT-BR</div>
            </div>

            <div class="lang-column">
              <div class="lang-column-header header-en">
                <span class="flag">🇺🇸</span>
                <h4>Inglês (EN)</h4>
                <span class="source-tag">Dataset Internacional</span>
              </div>

              <div v-if="enExercise" class="card-details">
                <div class="field-item">
                  <span class="field-label">Nome do Exercício:</span>
                  <div class="field-value highlight">{{ enExercise.name }}</div>
                </div>

                <div class="field-row-2">
                  <div class="field-item">
                    <span class="field-label">Categoria:</span>
                    <div class="field-value">
                      <span class="badge badge-primary">{{ enExercise.category }}</span>
                    </div>
                  </div>
                  <div class="field-item">
                    <span class="field-label">Nível:</span>
                    <div class="field-value">
                      <span class="badge badge-gray">{{ enExercise.level }}</span>
                    </div>
                  </div>
                </div>

                <div class="field-row-2">
                  <div class="field-item">
                    <span class="field-label">Mecânica:</span>
                    <div class="field-value">{{ enExercise.mechanic || 'N/A' }}</div>
                  </div>
                  <div class="field-item">
                    <span class="field-label">Tipo de Força:</span>
                    <div class="field-value">{{ enExercise.force || 'N/A' }}</div>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Equipamento:</span>
                  <div class="field-value">{{ enExercise.equipment || 'None / Body Only' }}</div>
                </div>

                <div class="field-item">
                  <span class="field-label">Músculos Primários:</span>
                  <div class="tags-group">
                    <span v-for="m in enExercise.primaryMuscles" :key="m" class="badge badge-primary">{{ m }}</span>
                    <span v-if="!enExercise.primaryMuscles?.length" class="text-muted">None</span>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Músculos Secundários:</span>
                  <div class="tags-group">
                    <span v-for="m in enExercise.secondaryMuscles" :key="m" class="badge badge-gray">{{ m }}</span>
                    <span v-if="!enExercise.secondaryMuscles?.length" class="text-muted">None</span>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Instruções de Execução ({{ enExercise.instructions?.length || 0 }} passos):</span>
                  <div class="steps-list">
                    <div v-for="(step, sIdx) in enExercise.instructions" :key="sIdx" class="step-card">
                      <span class="step-badge">{{ sIdx + 1 }}</span>
                      <p>{{ step }}</p>
                    </div>
                  </div>
                </div>

                <div class="field-item">
                  <span class="field-label">Imagens ({{ enExercise.images?.length || 0 }}):</span>
                  <div class="images-preview-grid">
                    <div v-for="(img, iIdx) in enExercise.images" :key="iIdx" class="img-card">
                      <img :src="'/images/' + img" class="thumb" @error="onImgErr" />
                      <span class="img-path">{{ img }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="empty-lang">Exercício não encontrado em EN</div>
            </div>
          </div>
        </div>
      </section>

      <section v-show="activeSubTab === 'audit'" class="sub-tab-content">
        <div class="audit-kpi-grid">
          <div class="kpi-card">
            <span class="kpi-title">Catálogo em Português</span>
            <span class="kpi-val">{{ auditReport?.counts?.pt || 0 }}</span>
            <span class="kpi-sub">database/exercises/pt/</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Catálogo em Inglês</span>
            <span class="kpi-val">{{ auditReport?.counts?.en || 0 }}</span>
            <span class="kpi-sub">database/exercises/en/</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Paridade de Catálogo</span>
            <span class="kpi-val text-success">
              {{ isParityPerfect ? '100%' : 'Divergência' }}
            </span>
            <span class="kpi-sub">Todos os IDs mapeados</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Alertas de Discrepância</span>
            <span :class="['kpi-val', (auditReport?.total_discrepancies || 0) > 0 ? 'text-warning' : 'text-success']">
              {{ auditReport?.total_discrepancies || 0 }}
            </span>
            <span class="kpi-sub">Passos ou fotos divergentes</span>
          </div>
        </div>

        <div class="audit-table-card">
          <div class="table-header">
            <h4>📋 Detalhamento dos Alertas de Consistência</h4>
            <div class="filter-pills">
              <button 
                :class="['pill', { active: auditFilter === 'all' }]"
                @click="auditFilter = 'all'"
              >
                Todos ({{ (auditReport?.step_mismatches?.length || 0) + (auditReport?.image_mismatches?.length || 0) }})
              </button>
              <button 
                :class="['pill', { active: auditFilter === 'steps' }]"
                @click="auditFilter = 'steps'"
              >
                Passos ({{ auditReport?.step_mismatches?.length || 0 }})
              </button>
              <button 
                :class="['pill', { active: auditFilter === 'images' }]"
                @click="auditFilter = 'images'"
              >
                Fotos ({{ auditReport?.image_mismatches?.length || 0 }})
              </button>
            </div>
          </div>

          <div v-if="filteredAuditList.length === 0" class="empty-audit-box">
            <span class="success-icon">🎉</span>
            <p>Nenhuma inconsistência encontrada para o filtro selecionado!</p>
          </div>

          <div v-else class="audit-table-wrap">
            <table class="audit-table">
              <thead>
                <tr>
                  <th>ID do Exercício</th>
                  <th>Nome em PT-BR</th>
                  <th>Tipo de Discrepância</th>
                  <th>Valor em PT-BR</th>
                  <th>Valor em EN</th>
                  <th>Ações</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in filteredAuditList" :key="item.id + item.field">
                  <td><code>{{ item.id }}</code></td>
                  <td><strong>{{ item.name_pt }}</strong></td>
                  <td>
                    <span :class="['badge', item.severity === 'warning' ? 'badge-warning' : 'badge-danger']">
                      {{ item.field === 'instructions' ? 'Passos' : 'Fotos' }}
                    </span>
                  </td>
                  <td>{{ item.value_pt }}</td>
                  <td>{{ item.value_en }}</td>
                  <td>
                    <button class="btn btn-secondary btn-sm" @click="inspectExercise(item.id)">
                      Inspecionar 🔍
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section v-show="activeSubTab === 'payload'" class="sub-tab-content">
        <div class="payload-card">
          <div class="payload-toolbar">
            <div class="controls-row">
              <div class="form-group">
                <label>Idioma do Payload:</label>
                <select v-model="payloadLang" @change="updatePayloadPreview">
                  <option value="pt">Português (PT-BR) - dist/exercises_pt.json</option>
                  <option value="en">Inglês (EN) - dist/exercises_en.json</option>
                </select>
              </div>

              <div class="form-group flex-1">
                <label>Exercício:</label>
                <select v-model="payloadExerciseId" @change="updatePayloadPreview">
                  <option v-for="ex in exercises" :key="ex.id" :value="ex.id">
                    {{ ex.name }} ({{ ex.id }})
                  </option>
                </select>
              </div>
            </div>

            <div class="payload-actions">
              <button class="btn btn-primary btn-sm" @click="copyPayloadJson">
                {{ copyFeedback || '📋 Copiar JSON do Payload' }}
              </button>
            </div>
          </div>

          <div class="contract-checklist">
            <h5>✅ Validação de Contrato do Payload (Android / iOS)</h5>
            <div class="checks-grid">
              <div class="check-item"><span class="check-icon">✓</span> <code>id</code> (String não-nula)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>name</code> (String localizada)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>category</code> (String padronizada)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>level</code> (String padronizada)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>primaryMuscles</code> (Array de Strings)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>instructions</code> (Array de Strings ordenadas)</div>
              <div class="check-item"><span class="check-icon">✓</span> <code>images</code> (Array de caminhos relativos)</div>
            </div>
          </div>

          <div class="json-code-box">
            <pre><code>{{ payloadFormattedJson }}</code></pre>
          </div>
        </div>
      </section>

      <section v-show="activeSubTab === 'categories'" class="sub-tab-content">
        <div class="audit-kpi-grid">
          <div class="kpi-card">
            <span class="kpi-title">Categorias em PT-BR</span>
            <span class="kpi-val">{{ categoryCountPt }}</span>
            <span class="kpi-sub">database/metadata/categories/pt.json</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Categorias em EN</span>
            <span class="kpi-val">{{ categoryCountEn }}</span>
            <span class="kpi-sub">database/metadata/categories/en.json</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Paridade de Categorias</span>
            <span :class="['kpi-val', isCategoryParityPerfect ? 'text-success' : 'text-warning']">
              {{ isCategoryParityPerfect ? '100% Sincronizado' : `${categoryMismatchCount} Divergência(s)` }}
            </span>
            <span class="kpi-sub">Chaves em snake_case</span>
          </div>

          <div class="kpi-card">
            <span class="kpi-title">Total de Exercícios Mapeados</span>
            <span class="kpi-val text-primary">{{ totalCategorizedExercises }}</span>
            <span class="kpi-sub">Distribuição nos catálogos</span>
          </div>
        </div>

        <div class="category-compare-card">
          <div class="category-toolbar">
            <div class="form-group flex-1">
              <label>Buscar Categoria / Chave:</label>
              <input 
                type="text" 
                v-model="categorySearch" 
                placeholder="Ex: olympic_weightlifting, cardio, Força..." 
              />
            </div>

            <div class="form-group">
              <label>Filtro de Status:</label>
              <select v-model="categoryFilterStatus">
                <option value="all">Todas as Categorias ({{ categoryComparisonList.length }})</option>
                <option value="synced">Apenas Sincronizadas ({{ syncedCategoriesCount }})</option>
                <option value="mismatch">Com Divergência ({{ categoryMismatchCount }})</option>
              </select>
            </div>

            <button class="btn btn-primary btn-add-cat" @click="openCreateCategory">
              ✨ Nova Categoria
            </button>
          </div>

          <div class="category-table-wrap">
            <table class="category-table">
              <thead>
                <tr>
                  <th>Chave Canônica (snake_case)</th>
                  <th>Rótulo PT-BR (🇧🇷)</th>
                  <th>Rótulo EN (🇺🇸)</th>
                  <th>Recurso Android / strings.xml</th>
                  <th>Exercícios (PT / EN)</th>
                  <th>Status</th>
                  <th>Ações</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="cat in filteredCategoryList" :key="cat.key">
                  <td>
                    <div class="key-chip">
                      <code>{{ cat.key }}</code>
                      <button class="btn-copy-chip" title="Copiar chave" @click="copyText(cat.key, 'Chave copiada!')">📋</button>
                    </div>
                  </td>
                  <td>
                    <span v-if="cat.labelPt" class="label-badge label-pt">{{ cat.labelPt }}</span>
                    <span v-else class="badge badge-warning">Não definido em PT</span>
                  </td>
                  <td>
                    <span v-if="cat.labelEn" class="label-badge label-en">{{ cat.labelEn }}</span>
                    <span v-else class="badge badge-warning">Não definido em EN</span>
                  </td>
                  <td>
                    <div class="key-chip android-chip">
                      <code>@string/cat_{{ cat.key }}</code>
                      <button class="btn-copy-chip" title="Copiar recurso Android" @click="copyText(`@string/cat_${cat.key}`, 'Recurso copiado!')">📋</button>
                    </div>
                  </td>
                  <td>
                    <span class="count-badge">
                      🇧🇷 {{ cat.countPt }} ex &bull; 🇺🇸 {{ cat.countEn }} ex
                    </span>
                  </td>
                  <td>
                    <span v-if="cat.status === 'synced'" class="badge badge-success">✅ Sincronizado</span>
                    <span v-else-if="cat.status === 'missing_en'" class="badge badge-warning">⚠️ Falta em EN</span>
                    <span v-else class="badge badge-danger">⚠️ Falta em PT</span>
                  </td>
                  <td>
                    <div class="table-actions-group">
                      <button 
                        class="btn btn-outline btn-xs" 
                        title="Editar Categoria"
                        @click="openEditCategory(cat)"
                      >
                        ✏️ Editar
                      </button>
                      <button 
                        class="btn btn-secondary btn-xs" 
                        @click="toggleInspectCategory(cat.key)"
                      >
                        {{ inspectedCategoryKey === cat.key ? 'Fechar ✕' : 'Ver 🔍' }}
                      </button>
                      <button 
                        class="btn btn-danger btn-xs" 
                        title="Excluir Categoria"
                        @click="handleDeleteCategory(cat)"
                      >
                        🗑️
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="showCategoryModal" class="modal-backdrop" @click.self="showCategoryModal = false">
            <div class="category-modal-card">
              <div class="modal-header">
                <h3>{{ categoryModalMode === 'create' ? '✨ Nova Categoria' : '✏️ Editar Categoria' }}</h3>
                <button class="btn-close" @click="showCategoryModal = false">&times;</button>
              </div>

              <div class="modal-body">
                <div v-if="categoryModalError" class="modal-error-box">
                  {{ categoryModalError }}
                </div>

                <div class="form-group">
                  <label>Nome em Português (PT-BR) *</label>
                  <input 
                    type="text" 
                    v-model="categoryForm.labelPt" 
                    @input="onCategoryLabelPtInput"
                    placeholder="Ex: Mobilidade Articular"
                    required
                  />
                </div>

                <div class="form-group">
                  <label>Chave Canônica / Slug (snake_case) *</label>
                  <div class="input-with-action">
                    <input 
                      type="text" 
                      v-model="categoryForm.key" 
                      @input="isAutoSlug = false"
                      :disabled="categoryModalMode === 'edit'"
                      placeholder="Ex: mobilidade_articular"
                      required
                    />
                    <span v-if="categoryModalMode === 'create'" class="input-hint-chip">
                      {{ isAutoSlug ? 'Auto' : 'Manual' }}
                    </span>
                  </div>
                  <small class="field-hint" v-if="categoryModalMode === 'create'">
                    Identificador no app Android: <code>@string/cat_{{ categoryForm.key || 'chave' }}</code>
                  </small>
                </div>

                <div class="form-group">
                  <label>Nome em Inglês (EN)</label>
                  <div class="input-with-btn">
                    <input 
                      type="text" 
                      v-model="categoryForm.labelEn" 
                      placeholder="Ex: Joint Mobility"
                    />
                    <button 
                      type="button" 
                      class="btn btn-secondary btn-sm" 
                      :disabled="translatingEn || !categoryForm.labelPt"
                      @click="autoTranslateCategoryEn"
                      title="Traduzir automaticamente do Português"
                    >
                      <span v-if="translatingEn" class="spinner-inline"></span>
                      <span v-else>✨ Traduzir</span>
                    </button>
                  </div>
                </div>
              </div>

              <div class="modal-footer">
                <button class="btn btn-secondary" @click="showCategoryModal = false" :disabled="categorySaving">
                  Cancelar
                </button>
                <button class="btn btn-primary" @click="handleSaveCategory" :disabled="categorySaving">
                  <span v-if="categorySaving" class="spinner-inline"></span>
                  <span v-else>{{ categoryModalMode === 'create' ? 'Salvar Categoria' : 'Atualizar Categoria' }}</span>
                </button>
              </div>
            </div>
          </div>

          <div v-if="inspectedCategory" class="category-drilldown-panel">
            <div class="drilldown-header">
              <div class="drilldown-title">
                <h5>Exercícios na categoria: <code>{{ inspectedCategory.key }}</code></h5>
                <span class="badge badge-primary">{{ inspectedCategoryExercises.length }} encontrados</span>
              </div>
              <button class="btn btn-secondary btn-sm" @click="inspectedCategoryKey = null">Fechar Painel</button>
            </div>

            <div v-if="inspectedCategoryExercises.length === 0" class="empty-drilldown">
              Nenhum exercício encontrado com esta categoria.
            </div>

            <div v-else class="drilldown-grid">
              <div 
                v-for="item in inspectedCategoryExercises" 
                :key="item.id" 
                class="drilldown-item-card"
              >
                <div class="drilldown-item-header">
                  <code>{{ item.id }}</code>
                  <button class="btn btn-outline btn-xs" @click="inspectExercise(item.id)">
                    Ver Comparativo 🔍
                  </button>
                </div>
                <div class="drilldown-item-names">
                  <div class="item-name-pt">🇧🇷 {{ item.namePt || '—' }}</div>
                  <div class="item-name-en">🇺🇸 {{ item.nameEn || '—' }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  fetchTranslationAudit, 
  fetchTranslationCompare, 
  fetchExerciseById, 
  fetchExercises,
  saveCategory,
  deleteCategory,
  translateText
} from '../api'

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

const emit = defineEmits(['metadata-updated'])

const activeSubTab = ref('compare')
const loading = ref(false)
const loadingCompare = ref(false)

const auditReport = ref(null)
const auditFilter = ref('all')

const selectedExerciseId = ref('')
const compareData = ref(null)
const filterDiscrepancy = ref('all')

const payloadLang = ref('pt')
const payloadExerciseId = ref('')
const payloadExerciseData = ref(null)
const copyFeedback = ref('')

const exercisesPtList = ref([])
const exercisesEnList = ref([])
const categorySearch = ref('')
const categoryFilterStatus = ref('all')
const inspectedCategoryKey = ref(null)

const isParityPerfect = computed(() => {
  const missingEn = auditReport.value?.missing_in_en?.length || 0
  const missingPt = auditReport.value?.missing_in_pt?.length || 0
  return missingEn === 0 && missingPt === 0
})

const ptExercise = computed(() => compareData.value?.languages?.pt || null)
const enExercise = computed(() => compareData.value?.languages?.en || null)

const stepMismatchWarning = computed(() => {
  if (!ptExercise.value || !enExercise.value) return null
  const ptCount = ptExercise.value.instructions?.length || 0
  const enCount = enExercise.value.instructions?.length || 0
  if (ptCount !== enCount) {
    return `Discrepância de Passos: PT tem ${ptCount} e EN tem ${enCount}`
  }
  return null
})

const imageMismatchWarning = computed(() => {
  if (!ptExercise.value || !enExercise.value) return null
  const ptCount = ptExercise.value.images?.length || 0
  const enCount = enExercise.value.images?.length || 0
  if (ptCount !== enCount) {
    return `Discrepância de Fotos: PT tem ${ptCount} e EN tem ${enCount}`
  }
  return null
})

const filteredExerciseOptions = computed(() => {
  if (filterDiscrepancy.value === 'steps') {
    const ids = new Set((auditReport.value?.step_mismatches || []).map(m => m.id))
    return props.exercises.filter(ex => ids.has(ex.id))
  }
  if (filterDiscrepancy.value === 'images') {
    const ids = new Set((auditReport.value?.image_mismatches || []).map(m => m.id))
    return props.exercises.filter(ex => ids.has(ex.id))
  }
  return props.exercises
})

const filteredAuditList = computed(() => {
  if (!auditReport.value) return []
  const steps = auditReport.value.step_mismatches || []
  const images = auditReport.value.image_mismatches || []

  if (auditFilter.value === 'steps') return steps
  if (auditFilter.value === 'images') return images
  return [...steps, ...images]
})

const payloadFormattedJson = computed(() => {
  if (!payloadExerciseData.value) return '// Selecione um exercício para visualizar o payload'
  return JSON.stringify(payloadExerciseData.value, null, 2)
})

const categoryComparisonList = computed(() => {
  const ptCats = props.metadata.categories?.pt || {}
  const enCats = props.metadata.categories?.en || {}

  const allKeys = new Set([
    ...Object.keys(ptCats),
    ...Object.keys(enCats)
  ])

  for (const ex of exercisesPtList.value) {
    if (ex.category) {
      allKeys.add(ex.category.trim().toLowerCase().replace(/[\s-]+/g, '_'))
    }
  }
  for (const ex of exercisesEnList.value) {
    if (ex.category) {
      allKeys.add(ex.category.trim().toLowerCase().replace(/[\s-]+/g, '_'))
    }
  }

  const list = []
  for (const key of Array.from(allKeys).sort()) {
    const labelPt = ptCats[key] || ''
    const labelEn = enCats[key] || ''

    const countPt = exercisesPtList.value.filter(ex => {
      const c = (ex.category || '').trim().toLowerCase().replace(/[\s-]+/g, '_')
      return c === key
    }).length

    const countEn = exercisesEnList.value.filter(ex => {
      const c = (ex.category || '').trim().toLowerCase().replace(/[\s-]+/g, '_')
      return c === key
    }).length

    let status = 'synced'
    if (labelPt && !labelEn) {
      status = 'missing_en'
    } else if (!labelPt && labelEn) {
      status = 'missing_pt'
    }

    list.push({
      key,
      labelPt,
      labelEn,
      countPt,
      countEn,
      status
    })
  }

  return list
})

const categoryCountPt = computed(() => Object.keys(props.metadata.categories?.pt || {}).length)
const categoryCountEn = computed(() => Object.keys(props.metadata.categories?.en || {}).length)

const categoryMismatchCount = computed(() => {
  return categoryComparisonList.value.filter(c => c.status !== 'synced').length
})

const isCategoryParityPerfect = computed(() => categoryMismatchCount.value === 0)
const syncedCategoriesCount = computed(() => categoryComparisonList.value.filter(c => c.status === 'synced').length)

const totalCategorizedExercises = computed(() => {
  return exercisesPtList.value.length || props.exercises.length
})

const filteredCategoryList = computed(() => {
  let list = categoryComparisonList.value

  if (categoryFilterStatus.value === 'synced') {
    list = list.filter(c => c.status === 'synced')
  } else if (categoryFilterStatus.value === 'mismatch') {
    list = list.filter(c => c.status !== 'synced')
  }

  if (categorySearch.value.trim()) {
    const q = categorySearch.value.trim().toLowerCase()
    list = list.filter(c => 
      c.key.toLowerCase().includes(q) ||
      c.labelPt.toLowerCase().includes(q) ||
      c.labelEn.toLowerCase().includes(q)
    )
  }

  return list
})

const inspectedCategory = computed(() => {
  if (!inspectedCategoryKey.value) return null
  return categoryComparisonList.value.find(c => c.key === inspectedCategoryKey.value) || null
})

const inspectedCategoryExercises = computed(() => {
  if (!inspectedCategoryKey.value) return []
  const key = inspectedCategoryKey.value

  const ptMap = new Map()
  for (const ex of exercisesPtList.value) {
    const c = (ex.category || '').trim().toLowerCase().replace(/[\s-]+/g, '_')
    if (c === key) {
      ptMap.set(ex.id, ex.name)
    }
  }

  const enMap = new Map()
  for (const ex of exercisesEnList.value) {
    const c = (ex.category || '').trim().toLowerCase().replace(/[\s-]+/g, '_')
    if (c === key) {
      enMap.set(ex.id, ex.name)
    }
  }

  const allIds = new Set([...ptMap.keys(), ...enMap.keys()])
  const result = []
  for (const id of Array.from(allIds).sort()) {
    result.push({
      id,
      namePt: ptMap.get(id) || '',
      nameEn: enMap.get(id) || ''
    })
  }

  return result
})

function toggleInspectCategory(key) {
  if (inspectedCategoryKey.value === key) {
    inspectedCategoryKey.value = null
  } else {
    inspectedCategoryKey.value = key
  }
}

function copyText(text, feedbackMsg) {
  if (!text) return
  navigator.clipboard.writeText(text)
  copyFeedback.value = feedbackMsg || 'Copiado!'
  setTimeout(() => {
    copyFeedback.value = ''
  }, 2000)
}

async function loadData() {
  loading.value = true
  try {
    const [audit, ptList, enList] = await Promise.all([
      fetchTranslationAudit(),
      fetchExercises('pt').catch(() => props.exercises || []),
      fetchExercises('en').catch(() => [])
    ])

    auditReport.value = audit
    exercisesPtList.value = ptList || []
    exercisesEnList.value = enList || []

    if (!selectedExerciseId.value && props.exercises.length > 0) {
      selectedExerciseId.value = props.exercises[0].id
      await loadCompare(selectedExerciseId.value)
    }

    if (!payloadExerciseId.value && props.exercises.length > 0) {
      payloadExerciseId.value = props.exercises[0].id
      await updatePayloadPreview()
    }
  } catch (err) {
    console.error('Erro ao carregar auditoria:', err)
  } finally {
    loading.value = false
  }
}

async function loadCompare(id) {
  if (!id) return
  loadingCompare.value = true
  try {
    compareData.value = await fetchTranslationCompare(id)
  } catch (err) {
    console.error('Erro ao carregar comparativo:', err)
  } finally {
    loadingCompare.value = false
  }
}

function onExerciseSelected() {
  loadCompare(selectedExerciseId.value)
}

function inspectExercise(id) {
  selectedExerciseId.value = id
  activeSubTab.value = 'compare'
  loadCompare(id)
}

async function updatePayloadPreview() {
  if (!payloadExerciseId.value) return
  try {
    payloadExerciseData.value = await fetchExerciseById(payloadExerciseId.value, payloadLang.value)
  } catch (err) {
    console.error('Erro ao carregar payload:', err)
  }
}

function copyPayloadJson() {
  if (!payloadFormattedJson.value) return
  navigator.clipboard.writeText(payloadFormattedJson.value)
  copyFeedback.value = '✅ Copiado com sucesso!'
  setTimeout(() => {
    copyFeedback.value = ''
  }, 2500)
}

function onImgErr(event) {
  event.target.style.opacity = '0.3'
}

const showCategoryModal = ref(false)
const categoryModalMode = ref('create')
const categoryForm = ref({
  key: '',
  labelPt: '',
  labelEn: ''
})
const isAutoSlug = ref(true)
const categoryModalError = ref('')
const categorySaving = ref(false)
const translatingEn = ref(false)

function slugifyCategory(text) {
  return (text || '')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[\s-]+/g, '_')
    .replace(/[^a-z0-9_]/g, '')
    .replace(/__+/g, '_')
}

function openCreateCategory() {
  categoryForm.value = {
    key: '',
    labelPt: '',
    labelEn: ''
  }
  isAutoSlug.value = true
  categoryModalError.value = ''
  categoryModalMode.value = 'create'
  showCategoryModal.value = true
}

function openEditCategory(cat) {
  categoryForm.value = {
    key: cat.key,
    labelPt: cat.labelPt || '',
    labelEn: cat.labelEn || ''
  }
  isAutoSlug.value = false
  categoryModalError.value = ''
  categoryModalMode.value = 'edit'
  showCategoryModal.value = true
}

function onCategoryLabelPtInput() {
  if (isAutoSlug.value && categoryModalMode.value === 'create') {
    categoryForm.value.key = slugifyCategory(categoryForm.value.labelPt)
  }
}

async function autoTranslateCategoryEn() {
  if (!categoryForm.value.labelPt.trim()) return
  translatingEn.value = true
  categoryModalError.value = ''
  try {
    const res = await translateText(categoryForm.value.labelPt.trim(), 'pt', 'en')
    if (res.translated) {
      categoryForm.value.labelEn = res.translated
    }
  } catch (err) {
    categoryModalError.value = `Erro ao traduzir: ${err.message}`
  } finally {
    translatingEn.value = false
  }
}

async function handleSaveCategory() {
  const pt = categoryForm.value.labelPt.trim()
  if (!pt) {
    categoryModalError.value = 'O nome da categoria em português é obrigatório.'
    return
  }

  let key = categoryForm.value.key.trim()
  if (!key) {
    key = slugifyCategory(pt)
    categoryForm.value.key = key
  }

  if (!key) {
    categoryModalError.value = 'Chave canônica inválida.'
    return
  }

  categorySaving.value = true
  categoryModalError.value = ''

  try {
    await saveCategory({
      key,
      labelPt: pt,
      labelEn: categoryForm.value.labelEn.trim()
    })
    showCategoryModal.value = false
    emit('metadata-updated')
    await loadData()
  } catch (err) {
    categoryModalError.value = err.message || 'Erro ao salvar categoria'
  } finally {
    categorySaving.value = false
  }
}

async function handleDeleteCategory(cat) {
  let msg = `Deseja realmente excluir a categoria "${cat.labelPt || cat.key}" (${cat.key})?`
  if (cat.countPt > 0 || cat.countEn > 0) {
    msg += `\n\n⚠️ ATENÇÃO: Esta categoria possui ${cat.countPt} exercício(s) em PT e ${cat.countEn} em EN associados a ela!`
  }

  if (!confirm(msg)) {
    return
  }

  try {
    await deleteCategory(cat.key)
    emit('metadata-updated')
    await loadData()
  } catch (err) {
    alert(`Erro ao excluir categoria: ${err.message}`)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.translation-validator {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.validator-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.tabs-group {
  display: flex;
  background: var(--bg-card);
  padding: 4px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  gap: 4px;
}

.tab-btn {
  background: transparent;
  border: none;
  color: var(--text-muted);
  padding: 8px 16px;
  font-size: 0.88rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: var(--text-main);
}

.tab-btn.active {
  background: var(--primary);
  color: #064e3b;
}

.selector-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 16px 20px;
  margin-bottom: 20px;
}

.selector-controls {
  display: flex;
  gap: 20px;
  align-items: flex-end;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.flex-1 {
  flex: 1;
}

.form-group label {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-muted);
}

.form-group select {
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-main);
  padding: 8px 12px;
  font-size: 0.9rem;
}

.compare-header-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 14px 20px;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.canonical-badge {
  display: flex;
  align-items: center;
  gap: 10px;
}

.badge-label {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.canonical-badge code {
  font-size: 1rem;
  font-weight: 700;
  color: var(--primary);
  background: rgba(16, 185, 129, 0.1);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}

.status-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.compare-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.lang-column {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.lang-column-header {
  padding: 14px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  border-bottom: 1px solid var(--border);
}

.header-pt {
  background: rgba(16, 185, 129, 0.08);
}

.header-en {
  background: rgba(59, 130, 246, 0.08);
}

.flag {
  font-size: 1.3rem;
}

.lang-column-header h4 {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
  flex: 1;
}

.source-tag {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.card-details {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.field-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.field-value {
  font-size: 0.92rem;
  color: var(--text-main);
}

.field-value.highlight {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--primary);
}

.tags-group {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.steps-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.step-card {
  display: flex;
  gap: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--border);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
}

.step-badge {
  background: var(--border);
  color: var(--text-main);
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
  flex-shrink: 0;
}

.step-card p {
  font-size: 0.88rem;
  line-height: 1.45;
  color: var(--text-main);
}

.images-preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 10px;
}

.img-card {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.img-card .thumb {
  width: 100%;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.img-path {
  font-size: 0.7rem;
  color: var(--text-muted);
  text-align: center;
  word-break: break-all;
}

.audit-kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.kpi-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kpi-title {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.kpi-val {
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--text-main);
}

.kpi-sub {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.text-success {
  color: var(--primary) !important;
}

.text-warning {
  color: #f59e0b !important;
}

.badge-warning {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.badge-success {
  background: rgba(16, 185, 129, 0.15);
  color: var(--primary);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.audit-table-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.table-header {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border);
}

.table-header h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-main);
}

.filter-pills {
  display: flex;
  gap: 8px;
}

.pill {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
}

.pill.active {
  background: rgba(16, 185, 129, 0.15);
  border-color: var(--primary);
  color: var(--primary);
}

.audit-table-wrap {
  overflow-x: auto;
}

.audit-table {
  width: 100%;
  border-collapse: collapse;
}

.audit-table th {
  background: rgba(0, 0, 0, 0.2);
  padding: 12px 16px;
  text-align: left;
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 1px solid var(--border);
}

.audit-table td {
  padding: 12px 16px;
  font-size: 0.88rem;
  color: var(--text-main);
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.payload-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.payload-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 20px;
  flex-wrap: wrap;
}

.controls-row {
  display: flex;
  gap: 16px;
  flex: 1;
}

.contract-checklist {
  background: rgba(16, 185, 129, 0.05);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: var(--radius-sm);
  padding: 14px 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.contract-checklist h5 {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--primary);
}

.checks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 8px;
}

.check-item {
  font-size: 0.82rem;
  color: var(--text-main);
  display: flex;
  align-items: center;
  gap: 6px;
}

.check-icon {
  color: var(--primary);
  font-weight: bold;
}

.json-code-box {
  background: #060911;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 16px;
  max-height: 450px;
  overflow-y: auto;
}

.json-code-box pre {
  margin: 0;
}

.json-code-box code {
  color: #38bdf8;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
  line-height: 1.5;
}

.loading-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 16px;
  color: var(--text-muted);
}

.empty-audit-box {
  padding: 48px 20px;
  text-align: center;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.category-compare-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.category-toolbar {
  padding: 16px 20px;
  display: flex;
  gap: 16px;
  align-items: flex-end;
  border-bottom: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.02);
}

.category-table-wrap {
  overflow-x: auto;
}

.category-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.category-table th {
  padding: 14px 18px;
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
  background: rgba(0, 0, 0, 0.2);
}

.category-table td {
  padding: 14px 18px;
  font-size: 0.88rem;
  border-bottom: 1px solid var(--border);
  color: var(--text-main);
  vertical-align: middle;
}

.category-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.03);
}

.key-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 4px 8px;
}

.key-chip code {
  color: #38bdf8;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.85rem;
}

.android-chip code {
  color: #a78bfa;
}

.btn-copy-chip {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 0.85rem;
  padding: 2px 4px;
  border-radius: 4px;
  opacity: 0.6;
  transition: opacity 0.2s, background-color 0.2s;
}

.btn-copy-chip:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.1);
}

.label-badge {
  display: inline-block;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
}

.label-pt {
  background: rgba(16, 185, 129, 0.12);
  color: var(--primary);
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.label-en {
  background: rgba(56, 189, 248, 0.12);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.25);
}

.count-badge {
  font-size: 0.8rem;
  color: var(--text-muted);
  font-weight: 500;
}

.category-drilldown-panel {
  padding: 20px;
  background: rgba(0, 0, 0, 0.25);
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.drilldown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.drilldown-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.drilldown-title h5 {
  margin: 0;
  font-size: 0.95rem;
  color: var(--text-main);
}

.drilldown-title code {
  color: var(--primary);
  background: var(--bg-input);
  padding: 2px 6px;
  border-radius: 4px;
}

.empty-drilldown {
  padding: 24px;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.drilldown-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.drilldown-item-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.drilldown-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.drilldown-item-header code {
  font-size: 0.78rem;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.drilldown-item-names {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.85rem;
}

.item-name-pt {
  font-weight: 600;
  color: var(--text-main);
}

.item-name-en {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.btn-add-cat {
  white-space: nowrap;
}

.table-actions-group {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.btn-xs {
  padding: 4px 8px;
  font-size: 0.75rem;
}

.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.category-modal-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 500px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: modalScale 0.2s ease-out;
}

@keyframes modalScale {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.02);
}

.modal-header h3 {
  margin: 0;
  font-size: 1.15rem;
  color: var(--text-main);
  font-weight: 700;
}

.btn-close {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
}

.btn-close:hover {
  color: var(--danger);
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-error-box {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: var(--danger);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
  font-weight: 500;
}

.input-with-btn {
  display: flex;
  gap: 8px;
  align-items: center;
}

.input-with-btn input {
  flex: 1;
}

.input-with-action {
  display: flex;
  align-items: center;
  position: relative;
}

.input-with-action input {
  width: 100%;
  padding-right: 60px;
}

.input-hint-chip {
  position: absolute;
  right: 10px;
  font-size: 0.7rem;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.05em;
  padding: 2px 6px;
  border-radius: 4px;
  background: var(--bg-input);
  color: var(--text-muted);
  border: 1px solid var(--border);
}

.field-hint {
  display: block;
  margin-top: 4px;
  font-size: 0.75rem;
  color: var(--text-dim);
}

.field-hint code {
  color: var(--primary);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.02);
}

.spinner-inline {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 900px) {
  .compare-grid,
  .audit-kpi-grid,
  .category-toolbar {
    grid-template-columns: 1fr;
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
