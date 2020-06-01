

import Navbar from './components/Navbar.js'
import Playerbar from './components/Playerbar.js'
import store from './store/index.js'
import routes from './routes.js'


export const app = new Vue({
  el: '#app',
  router: new VueRouter({ routes }),
  components: { Navbar, Playerbar },
  vuetify: new Vuetify(),
  store,
  data() {
    return {
      Buildnr: "",
      links: routes,
      AppTitle: "Omx Control",
      drawer: false,
    }
  },
  computed: {
    ...Vuex.mapState({
     
    })
  },
  created() {
    // keep in mind that all that is comming from index.html is a string. Boolean or numerics need to be parsed.
    this.Buildnr = window.myapp.buildnr    
  },
  methods: {

  },
  template: `
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer absolute>
      <v-container>
        <v-row>
          <v-col class="d-flex text-center caption" cols="12">
            <Playerbar />
          </v-col>
        </v-row>
        <v-row>
          <v-col class="d-flex text-center caption" cols="12">
            <div>
              {{ new Date().getFullYear() }} —
              <span>Buildnr: {{Buildnr}}</span>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>
  </v-app>`
})

console.log('Main is here!')