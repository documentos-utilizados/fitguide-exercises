<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <span class="brand-icon">⚡</span>
        <div class="brand-info">
          <h2>FitGuide</h2>
          <span class="brand-badge">Admin</span>
        </div>
      </div>

      <nav class="sidebar-menu">
        <router-link to="/treinos" class="sidebar-link">
          <span class="link-icon">🏋️</span>
          <span class="link-text">Treinos</span>
        </router-link>
        <router-link to="/notificacoes" class="sidebar-link">
          <span class="link-icon">🔔</span>
          <span class="link-text">Notificações</span>
        </router-link>
      </nav>
    </aside>

    <div class="content-area">
      <router-view />
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background-color: var(--bg-main);
}

.sidebar {
  width: 240px;
  background-color: var(--bg-card);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  height: 100vh;
  z-index: 200;
}

.sidebar-header {
  padding: 18px 20px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  font-size: 1.5rem;
}

.brand-info h2 {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.2;
}

.brand-badge {
  font-size: 0.7rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px 12px;
  flex: 1;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  text-decoration: none;
  font-size: 0.95rem;
  font-weight: 600;
  transition: all 0.2s ease;
}

.sidebar-link:hover {
  color: var(--text-main);
  background-color: rgba(255, 255, 255, 0.05);
}

.sidebar-link.router-link-active {
  color: var(--primary);
  background-color: rgba(16, 185, 129, 0.12);
}

.link-icon {
  font-size: 1.2rem;
  display: flex;
  align-items: center;
}

.content-area {
  flex: 1;
  min-width: 0;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .app-layout {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    height: auto;
    position: relative;
    border-right: none;
    border-bottom: 1px solid var(--border);
  }

  .sidebar-menu {
    flex-direction: row;
    padding: 8px 12px;
  }

  .sidebar-link {
    flex: 1;
    justify-content: center;
  }
}
</style>
