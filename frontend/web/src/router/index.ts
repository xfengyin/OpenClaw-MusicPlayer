import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/home',
      redirect: '/'
    },
    {
      path: '/playlist',
      name: 'playlist',
      component: () => import('../views/Playlist.vue')
    },
    {
      path: '/favorites',
      name: 'favorites',
      component: () => import('../views/Favorites.vue')
    },
    {
      path: '/history',
      name: 'history',
      component: () => import('../views/History.vue')
    },
    {
      path: '/downloads',
      name: 'downloads',
      component: () => import('../views/Downloads.vue')
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/Search.vue')
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/Settings.vue')
    }
  ]
})

export default router
