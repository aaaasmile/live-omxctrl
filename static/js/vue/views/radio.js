import API from '../apicaller.js'

export default {
  data() {
    return {
      radioloading: false,
      selected_item: {},
      dialogItemSelected: false,
      dialogInsertEdit: false,
      dialog_title: '',
      pagesize: 20,
      pageix: 0,
      transition: 'scale-transition',
      rules: {
        name: [val => (val || '').length > 0 || 'This field is required'],
        URI: [val => (val || '').length > 0 || 'This field is required'],
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
      this.selected_item.action_name = 'Play'
      this.selected_item.continue_action = this.playSelectedItem
      if (item.title === '') {
        this.selected_item.itemquestion = item.uri
      }
      this.dialogItemSelected = true
    },
    askForDeleteItem(item) {
      console.log('ask to delete radio item: ', item)
      this.selected_item = item
      this.selected_item.itemquestion = item.title
      this.selected_item.action_name = 'Delete'
      this.selected_item.continue_action = this.deleteSelectedItem
      if (item.title === '') {
        this.selected_item.itemquestion = item.uri
      }
      this.dialogItemSelected = true
    },
    continueSelectedItem() {
      console.log('continueSelectedItem : ', this.selected_item)
      this.selected_item.action_name()
    },
    playSelectedItem() {
      console.log('playSelectedItem is: ', this.selected_item)
      this.dialogItemSelected = false

      let req = { uri: this.selected_item.uri, force_type: 'radio' }
      API.PlayUri(this, req)

      this.$router.push('/')
    },
    prepareInsert() {
      this.dialog_title = 'Insert New'
      this.selected_item = {
        id: '',
        name: '',
        URI: '',
        descr: '',
      }
      this.dialogInsertEdit = true
      this.selected_item.action_name = this.insertNewtem
    },
    prepareEdit(item) {
      console.log('prepare Edit radio')
      this.dialog_title = 'Edit Radio'
      this.selected_item = {}
      this.selected_item.id = item.id
      this.selected_item.name = item.title
      this.selected_item.URI = item.uri
      this.selected_item.descr = item.description
      this.dialogInsertEdit = true
      this.selected_item.action_name = this.editItem
    },
    editItem(){
      console.log('Edit radio')
      this.handleRadioReq('Edit')  
    },
    insertNewtem() {
      console.log('Insert new radio')
      this.handleRadioReq('Insert')  
    },
    deleteSelectedItem() {
      console.log('deleteSelectedItem : ', this.selected_item)
      this.handleRadioReq('Delete')  
    },
    handleRadioReq(req_name){
      let req = {
        radio_name: this.selected_item.name,
        id: this.selected_item.id,
        uri: this.selected_item.URI,
        descr: this.selected_item.descr,
        pageix: this.pageix, pagesize: this.pagesize
      }
      req.name = req_name
      API.HandleRadio(this, req, (ok, result) => {
        this.dialogInsertEdit = false
        if (ok) {
          this.$store.commit('radiofetch', result.data)
        }
      })
    },
    loadMore() {
      console.log('Load more')
      this.pageix += 1
      let req = { pageix: this.pageix, pagesize: this.pagesize }
      API.FetchRadio(this, req)
    },
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
              <v-btn icon @click="prepareInsert" v-on="on">
                <v-icon>mdi-plus</v-icon>
              </v-btn>
            </template>
            <span>New Radio</span>
          </v-tooltip>
        </v-toolbar>
        <v-container>
          <v-list dense nav>
            <template v-for="plitem in radio">
              <v-list-item :key="plitem.id">
                <v-list-item-content>
                  <v-list-item-title>{{ plitem.title }}</v-list-item-title>
                  <v-list-item-title>{{
                    plitem.description
                  }}</v-list-item-title>
                  <v-list-item-title>{{ plitem.genre }}</v-list-item-title>
                  <v-list-item-title>{{ plitem.uri }}</v-list-item-title>
                  <v-row>
                    <v-btn
                      icon
                      text
                      :key="plitem.id"
                      @click="askForPlayItem(plitem)"
                      ><v-icon>library_music</v-icon>
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn icon text @click="prepareEdit(plitem)"
                      ><v-icon>mdi-circle-edit-outline</v-icon>
                    </v-btn>
                    <v-btn icon text @click="askForDeleteItem(plitem)"
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
      <v-dialog v-model="dialogItemSelected" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text
            >Do you want to {{ selected_item.action_name }} the radio "{{
              selected_item.itemquestion
            }}"?</v-card-text
          >
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="continueSelectedItem"
              >OK</v-btn
            >
            <v-btn
              color="green darken-1"
              text
              @click="dialogItemSelected = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="dialogInsertEdit" persistent max-width="290">
        <v-card>
          <v-container>
            <v-col cols="12">
              <v-row justify="space-around">
                <v-card-title class="headline">{{dialog_title}}</v-card-title>
                <v-text-field
                  label="Name"
                  v-model="selected_item.name"
                  :rules="rules.name"
                  required
                ></v-text-field>
                <v-text-field
                  label="URI"
                  v-model="selected_item.URI"
                  :rules="rules.URI"
                  required
                ></v-text-field>
                <v-text-field
                  label="Description"
                  v-model="selected_item.descr"
                ></v-text-field>
              </v-row>
            </v-col>
          </v-container>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="continueSelectedItem"
              >OK</v-btn
            >
            <v-btn color="green darken-1" text @click="dialogInsertEdit = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>
`
}