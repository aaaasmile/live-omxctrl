import API from '../apicaller.js'

export default {
    data() {
        return {
            musicloading: false,
            selected_item: {},
            dialogScan: false,
            pagesize: 10,
            pageix: 0,
            transition: 'scale-transition',
        }
    },
    created() {
     
    },
    computed: {
        ...Vuex.mapState({
     
        })
    },
    methods: {
        scanForMusic() {
            console.log('scan for music')
            this.pageix = 0
            let req = { pageix: this.pageix, pagesize: this.pagesize }
            API.ScanMusic(this, req)
        }
    },
    template: `
  <v-container pa-1>
    <v-skeleton-loader
      :loading="musicloading"
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
        </v-toolbar>
        <v-card-title>Music available</v-card-title>
        <v-container>
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
  </v-container>`
}