import Dashboard from './views/dashboard.js'
import OSView from './views/os.js'
import HistoryView from './views/history.js'
import VideoView from './views/video.js'
import RadioView from './views/radio.js'
import MusicView from './views/music.js'

export default [
  { path: '/', icon: 'dashboard', title: 'Dashboard', component: Dashboard },
  { path: '/os', icon: 'dashboard', title: 'OS', component: OSView },
  { path: '/history', icon: 'dashboard', title: 'History', component: HistoryView },
  { path: '/video', icon: 'dashboard', title: 'Video', component: VideoView },
  { path: '/radio', icon: 'dashboard', title: 'Radio', component: RadioView },
  { path: '/music', icon: 'dashboard', title: 'Music', component: MusicView },
]