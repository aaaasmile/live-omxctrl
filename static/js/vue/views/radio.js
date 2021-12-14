import API from '../apicaller.js'

export default {
  data() {
    return {
      radioloading: false,
      selected_item: {},
      dialogPlaySelected: false,
      dialogEditSelected: false,
      dialogInsert: false,
      pagesize: 20,
      pageix: 0,
      transition: 'scale-transition',
      radio_name: '',
      radio_URI: '',
      radio_descr: '',
      rules: {
        radio_name: [val => (val || '').length > 0 || 'This field is required'],
        radio_URI: [val => (val || '').length > 0 || 'This field is required'],
      },
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
    },
    insertNewtem() {
      console.log('Insert new radio')
      let req = {  radio_name: this.radio_name, uri: this.radio_URI, descr: this.radio_descr, pageix: this.pageix, pagesize: this.pagesize  }
      req.name = 'InsertRadio'
      API.HandleRadio(this, req, (ok,result) => {
        this.dialogInsert = false
        if (ok){
          this.$store.commit('radiofetch', result.data)
        }
      })
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
          <v-toolbar-title class="subheading grey--text">Radio</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon @click="dialogInsert = true" v-on="on">
                <v-icon>mdi-plus</v-icon>
              </v-btn>
            </template>
            <span>New Radio</span>
          </v-tooltip>
        </v-toolbar>
        <v-container>
          <v-list dense nav>
            <template v-for="plitem in radio" >
              <v-list-item :key="plitem.id">
              <v-list-item-content>
                <v-list-item-title>{{ plitem.title }}</v-list-item-title>
                <v-list-item-title>{{ plitem.description }}</v-list-item-title>
                <v-list-item-title>{{ plitem.genre }}</v-list-item-title>
                <v-list-item-title>{{ plitem.uri }}</v-list-item-title>
                <v-row >
                  <v-btn icon text :key="plitem.id"
                    @click="askForPlayItem(plitem)"
                  ><v-icon>library_music</v-icon> </v-btn>
                  <v-spacer></v-spacer>
                  <v-btn icon text
                    ><v-icon>mdi-circle-edit-outline</v-icon>
                  </v-btn>
                  <v-btn icon text
                    ><v-icon>mdi-delete-forever-outline</v-icon>
                  </v-btn>
                </v-row>
              </v-list-item-content>
            </v-list-item>
            </template>
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
      <v-dialog v-model="dialogInsert" persistent max-width="290">
        <v-card>
          <v-container>
            <v-col cols="12">
              <v-row justify="space-around">
                <v-card-title class="headline">Insert New</v-card-title>
                <v-text-field
                  label="Name"
                  v-model="radio_name"
                  :rules="rules.radio_name"
                  required
                ></v-text-field>
                <v-text-field
                  label="URI"
                  v-model="radio_URI"
                  :rules="rules.radio_URI"
                  required
                ></v-text-field>
                <v-text-field
                  label="Description"
                  v-model="radio_descr"
                ></v-text-field>
              </v-row>
            </v-col>
          </v-container>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="insertNewtem">OK</v-btn>
            <v-btn color="green darken-1" text @click="dialogInsert = false"
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
  </v-container>`
}