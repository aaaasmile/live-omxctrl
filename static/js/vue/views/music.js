import API from '../apicaller.js'

import FolderList from '../components/folderlist.js'
import Playerbar from '../components/playerbar.js'

export default {
  components: { FolderList, Playerbar },
  data() {
    return {
      loadingData: false,
      loadingSync: false,
      selected_item: {},
      dialogScan: false,
      transition: 'scale-transition',
    }
  },
  created() {
    if (this.parent_folder === '') {
      let req = { parent: this.parent_folder }
      API.FetchMusic(this, req)
    }
  },
  computed: {
    ...Vuex.mapState({
      music: state => {
        return state.fs.music
      },
      parent_folder: state => {
        return state.fs.parent_folder
      },
    })
  },
  methods: {
    scanForMusic() {
      console.log('scan for music')
      let req = { parent: this.parent_folder }
      API.ScanMusic(this, req)
    },
    fetchRoot() {
      console.log('fetchRoot music')
      let req = { parent: '' }
      API.FetchMusic(this, req)
    },
    refreshFolder(){
      console.log('refreshFolder ')
      this.loadingSync = true
      let req = { parent: this.parent_folder }
      API.FetchMusic(this, req, () => this.loadingSync = false)
    }
  },
  template: `
  <v-container pa-1>
    <v-skeleton-loader
      :loading="loadingData"
      :transition="transition"
      height="94"
      type="list-item-three-line"
    >
      <v-card color="grey lighten-4" flat tile>
        <v-toolbar flat dense>
          <v-toolbar-title class="subheading grey--text"
            >Music commands</v-toolbar-title
          >
          <v-spacer></v-spacer>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon @click="dialogScan = true" v-on="on">
                <v-icon>mdi-magnify-scan</v-icon>
              </v-btn>
            </template>
            <span>Scan for music</span>
          </v-tooltip>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon @click="fetchRoot()" v-on="on">
                <v-icon>mdi-home</v-icon>
              </v-btn>
            </template>
            <span>Fetch root</span>
          </v-tooltip>
          <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="refreshFolder" :loading="loadingSync" v-on="on">
              <v-icon>mdi-sync</v-icon>
            </v-btn>
          </template>
          <span>Update view</span>
        </v-tooltip>

        </v-toolbar>
        <v-container>
          <FolderList></FolderList>
           <Playerbar />
        </v-container>
      </v-card>
    </v-skeleton-loader>
    <v-container>
      <v-dialog v-model="dialogScan" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text>Do you want to scan the Pi for music and rebuild the list?</v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="scanForMusic">OK</v-btn>
            <v-btn color="green darken-1" text @click="dialogScan = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>
`
}