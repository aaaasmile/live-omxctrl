import API from '../apicaller.js'

export default {
  data() {
    return {
        videoloading: false,
        selected_item: {},
        dialogPlaySelected: false,
        pagesize: 5,
        pageix: 0,
        transition: 'scale-transition',
    }
  },
  created() {
    this.pageix = 0
    let req = { pageix: this.pageix, pagesize: this.pagesize }
    API.FetchVideo(this, req)
  },
  computed: {
    ...Vuex.mapState({
        video: state => {
          return state.fs.video
        },
        last_video_fetch: state => {
            return state.fs.last_video_fetch
          }
      })
  },
  methods: {
    scan_video() {
      console.log('scan for video')
     
    },
    askForPlayItem(item) {
        console.log('ask to play video item: ', item)
        this.selected_item = item
        this.selected_item.itemquestion = item.title
        if (item.title === ''){
          this.selected_item.itemquestion = item.uri
        }
        this.dialogPlaySelected = true
      },
      playSelectedItem() {
        console.log('playSelectedItem is: ', this.selected_item)
        this.dialogPlaySelected = false
  
        let req = { uri: this.selected_item.uri }
        API.PlayUri(this, req)
  
        this.$router.push('/')
      },
      loadMore() {
        console.log('Load more')
        this.pageix += 1
        let req = { pageix: this.pageix, pagesize: this.pagesize }
        API.FetchHistory(this, req)
      }
  },
  template: `
  <v-container pa-1>
    <v-skeleton-loader
      :loading="videoloading"
      :transition="transition"
      height="94"
      type="list-item-three-line"
    >
      <v-card color="grey lighten-4" flat tile>
        <v-card-title> Video</v-card-title>
        <v-container>
          <v-list dense nav>
            <v-list-item
              v-for="plitem in video"
              :key="plitem.id"
              @click="askForPlayItem(plitem)"
            >
              <v-list-item-icon>
                <v-icon>{{ plitem.icon }}</v-icon>
              </v-list-item-icon>

              <v-list-item-content>
                <v-list-item-title>{{ plitem.title }}</v-list-item-title>
                <v-list-item-title>{{ plitem.uri }}</v-list-item-title>
                <v-list-item-title>{{ plitem.playedAt }}</v-list-item-title>
                <v-list-item-title
                  >Duration: {{ plitem.duration }}</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-divider></v-divider>
          <v-row justify="center">
            <v-btn icon text @click="loadMore" :disabled="last_video_fetch"
              >More<v-icon>more_horiz</v-icon>
            </v-btn>
          </v-row>
        </v-container>
      </v-card>
    </v-skeleton-loader>
    <v-container>
      <v-dialog v-model="dialogPlaySelected" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text
            >Do you want to play the video "{{ selected_item.itemquestion }}"?</v-card-text
          >
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="playSelectedItem"
              >OK</v-btn
            >
            <v-btn
              color="green darken-1"
              text
              @click="dialogPlaySelected = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>`
}