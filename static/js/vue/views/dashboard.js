import API from '../apicaller.js'
import Playerbar from '../components/Playerbar.js'

export default {
  components: { Playerbar },
  data() {
    return {
      loadingyoutube: false,
      uriToPlay: '',
    }
  },
  computed: {
    ...Vuex.mapState({
      PlayingURI: state => {
        return state.ps.uri
      },
      PlayingInfo: state => {
        return state.ps.info
      },
      PlayingType: state => {
        return state.ps.itemtype
      },
      ListName: state => {
        return state.ps.listname
      },
      PlayingPrev: state => {
        return state.ps.previous
      },
      PlayingNext: state => {
        return state.ps.next
      },
    })
  },
  methods: {
    playYoutube() {
      console.log('play youtube url')
      let req = { uri: this.uriToPlay }
      this.loadingyoutube = true
      API.PlayYoutube(this, req)
    }
  },
  template: `
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Omx Player</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="playYoutube" :loading="loadingyoutube" v-on="on">
            <v-icon>airplay</v-icon>
          </v-btn>
        </template>
        <span>Play youtube url</span>
      </v-tooltip>
    </v-toolbar>
    <v-col cols="12">
      <v-row>
        <v-col cols="12">
          <v-text-field v-model="uriToPlay" label="Select an URL"></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card flat tile>
            <v-card-title>Playlist {{ListName}}</v-card-title>
            <div class="mx-4">
              <div class="subtitle-2 text--secondary">URI: {{PlayingURI}}</div>
              <div class="subtitle-2 text--secondary">Info: {{PlayingInfo}}</div>
              <div class="subtitle-2 text--secondary">Type: {{PlayingType}}</div>
              <div class="subtitle-2 text--secondary">Previous: {{PlayingPrev}}</div>
              <div class="subtitle-2 text--secondary">Next: {{PlayingNext}}</div>
            </div>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <Playerbar />
        </v-col>
      </v-row>
    </v-col>
  </v-card>`
}