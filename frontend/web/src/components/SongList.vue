<template>
  <n-data-table
    :columns="columns"
    :data="songs"
    :bordered="false"
    :single-line="false"
    size="small"
    @row-click="handleRowClick"
  />
</template>

<script setup lang="ts">
import { h } from 'vue'
import { NDataTable, NButton, NIcon, NSpace } from 'naive-ui'
import { Play, Heart, Download } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'

interface Song {
  id: string
  title: string
  artist: string
  album: string
  duration: number
}

const props = defineProps<{
  songs: Song[]
}>()

const emit = defineEmits<{
  play: [song: Song]
  download: [song: Song]
}>()

const formatDuration = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const columns: DataTableColumns<Song> = [
  {
    title: '',
    key: 'action',
    width: 60,
    render(row) {
      return h(
        NButton,
        {
          quaternary: true,
          circle: true,
          size: 'small',
          onClick: (e) => {
            e.stopPropagation()
            emit('play', row)
          }
        },
        { icon: () => h(NIcon, null, { default: () => h(Play) }) }
      )
    }
  },
  {
    title: '歌曲',
    key: 'title',
    ellipsis: true
  },
  {
    title: '歌手',
    key: 'artist',
    ellipsis: true
  },
  {
    title: '专辑',
    key: 'album',
    ellipsis: true
  },
  {
    title: '时长',
    key: 'duration',
    width: 80,
    render(row) {
      return formatDuration(row.duration)
    }
  },
  {
    title: '操作',
    key: 'operations',
    width: 100,
    render(row) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              quaternary: true,
              circle: true,
              size: 'small'
            },
            { icon: () => h(NIcon, null, { default: () => h(Heart) }) }
          ),
          h(
            NButton,
            {
              quaternary: true,
              circle: true,
              size: 'small',
              onClick: (e) => {
                e.stopPropagation()
                emit('download', row)
              }
            },
            { icon: () => h(NIcon, null, { default: () => h(Download) }) }
          )
        ]
      })
    }
  }
]

const handleRowClick = (row: Song) => {
  emit('play', row)
}
</script>
