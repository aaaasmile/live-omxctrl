import API from '../apicaller.js'

export default {
    data() {
        return {
            radioloading: false,
            selected_item: {},
            dialogPlaySelected: false,
            dialogEditSelected: false,
            dialogScan: false,
            pagesize: 20,
            pageix: 0,
            transition: 'scale-transition',
        }
    },
    created() {
        this.pageix = 0
        let req = { pageix: this.pageix, pagesize: this.pagesize }
        API.FetchRadio(this, req)
    },
    computed: {
        ...Vuex.mapState({
            radio: state => {
                return state.fs.radio
            },
            last_radio_fetch: state => {
                return state.fs.last_radio_fetch
            }
        })
    },
    methods: {
        askForPlayItem(item) {
            console.log('ask to play radio item: ', item)
            this.selected_item = item
            this.selected_item.itemquestion = item.title
            if (item.title === '') {
                this.selected_item.itemquestion = item.uri
            }
            this.dialogPlaySelected = true
        },
        playSelectedItem() {
            console.log('playSelectedItem is: ', this.selected_item)
            this.dialogPlaySelected = false

            let req = { uri: this.selected_item.uri, force_type: 'radio' }
            API.PlayUri(this, req)

            this.$router.push('/')
        },
        loadMore() {
            console.log('Load more')
            this.pageix += 1
            let req = { pageix: this.pageix, pagesize: this.pagesize }
            API.FetchRadio(this, req)
        }
    },
    template: `
  <v-container pa-1>
    <v-skeleton-loader
      :loading="radioloading"
      :transition="transition"
      height="94"
      type="list-item-three-line"
    >
      <v-card color="grey lighten-4" flat tile>
        <v-toolbar flat dense>
          <v-toolbar-title class="subheading grey--text"
            >Radio commands</v-toolbar-title
          >
          <v-spacer></v-spacer>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon @click="dialogEditSelected = true" v-on="on">
                <v-icon>mdi-edit</v-icon>
              </v-btn>
            </template>
            <span>Edit Radio</span>
          </v-tooltip>
        </v-toolbar>
        <v-card-title>Radio available</v-card-title>
        <v-container>
          <v-list dense nav>
            <v-list-item
              v-for="plitem in radio"
              :key="plitem.id"
              @click="askForPlayItem(plitem)"
            >
              
              <v-list-item-content>
                <v-list-item-title>{{ plitem.title }}</v-list-item-title>
                <v-list-item-title>{{ plitem.description }}</v-list-item-title>
                <v-list-item-title>{{ plitem.genre }}</v-list-item-title>
                <v-list-item-title>{{ plitem.uri }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-divider></v-divider>
          <v-row justify="center">
            <v-btn icon text @click="loadMore" :disabled="last_radio_fetch"
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
            >Do you want to play the radio "{{
              selected_item.itemquestion
            }}"?</v-card-text
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
      <v-dialog v-model="dialogEditSelected" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Edit</v-card-title>
         
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="dialogEditSelected"
              >OK</v-btn
            >
            <v-btn
              color="green darken-1"
              text
              @click="dialogEditSelected = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>`
}