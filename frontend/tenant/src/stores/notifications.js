import { reactive } from 'vue'
import api from '@/services/api'

const state = reactive({
  unreadCount: 0,
  recentItems: [],
  loaded: false
})

let pollTimer = null

export function useNotifications() {
  async function fetchUnreadCount() {
    try {
      const res = await api.get('/api/v1/tenant/notifications/unread-count')
      state.unreadCount = res.data?.data?.unread_count ?? 0
    } catch (err) {
      console.error('Failed to fetch unread notification count:', err)
    }
  }

  // NOTE: the list endpoint wraps its payload as
  // { success, data: { data: [...], total, page, per_page } } — a nested
  // "data" one level deeper than most other list endpoints in this app
  // (which put data/total directly under the top-level envelope).
  async function fetchRecent() {
    try {
      const res = await api.get('/api/v1/tenant/notifications', { params: { page: 1, per_page: 10 } })
      state.recentItems = res.data?.data?.data || []
    } catch (err) {
      console.error('Failed to fetch recent notifications:', err)
    }
  }

  async function refresh() {
    await Promise.all([fetchUnreadCount(), fetchRecent()])
    state.loaded = true
  }

  async function markAsRead(id) {
    try {
      await api.patch(`/api/v1/tenant/notifications/${id}/read`)
      await refresh()
    } catch (err) {
      console.error('Failed to mark notification as read:', err)
    }
  }

  async function markAllAsRead() {
    try {
      await api.post('/api/v1/tenant/notifications/read-all')
      await refresh()
    } catch (err) {
      console.error('Failed to mark all notifications as read:', err)
    }
  }

  // startPolling — refreshes the unread badge periodically since there's no
  // push/websocket infra. A single global interval (guarded by pollTimer),
  // not one per component instance, so mounting HeaderBar more than once
  // never stacks up parallel timers.
  function startPolling(intervalMs = 60000) {
    if (pollTimer) return
    pollTimer = setInterval(fetchUnreadCount, intervalMs)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  function reset() {
    state.unreadCount = 0
    state.recentItems = []
    state.loaded = false
    stopPolling()
  }

  return {
    state,
    fetchUnreadCount,
    fetchRecent,
    refresh,
    markAsRead,
    markAllAsRead,
    startPolling,
    stopPolling,
    reset
  }
}
