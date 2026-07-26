import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: 'Login', public: true }
  },
  {
    path: '',
    component: () => import('@/layouts/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { 
          title: 'Dashboard', 
          icon: 'pi pi-home',
          description: 'Platform overview & real-time statistics'
        }
      },
      {
        path: 'companies',
        name: 'Companies',
        component: () => import('@/views/Companies.vue'),
        meta: { 
          title: 'Companies', 
          icon: 'pi pi-building',
          description: 'Manage tenant companies & their lifecycle'
        }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users.vue'),
        meta: { 
          title: 'Users', 
          icon: 'pi pi-users',
          description: 'Manage platform admin accounts'
        }
      },
      {
        path: 'modules',
        name: 'Modules',
        component: () => import('@/views/Modules.vue'),
        meta: { 
          title: 'Modules', 
          icon: 'pi pi-cog',
          description: 'Module management, including module installation and uninstallation.'
        }
      },
      {
        path: 'licenses',
        name: 'Licenses',
        component: () => import('@/views/Licenses.vue'),
        meta: { 
          title: 'Licenses', 
          icon: 'pi pi-id-card',
          description: 'Manage software licenses for tenant companies.'
        }
      },
      {
        path: 'monitoring',
        name: 'Monitoring',
        component: () => import('@/views/Monitoring.vue'),
        meta: {
          title: 'Monitoring', 
          icon: 'pi pi-chart-bar',
          description: 'Platform health, database connectivity & tenant status.'
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Auth guard
router.beforeEach((to, from, next) => {
  const { state } = useAuth()
  if (to.meta.requiresAuth && !state.isAuthenticated) {
    next('/login')
  } else if (to.name === 'Login' && state.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
